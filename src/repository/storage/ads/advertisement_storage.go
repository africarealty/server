package ads

import (
	"context"
	"fmt"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/aerospike/aerospike-client-go/v6/types"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/ads"
	"github.com/africarealty/server/src/kit/goroutine"
	"github.com/africarealty/server/src/kit/log"
	kitAero "github.com/africarealty/server/src/kit/storages/aerospike"
	"github.com/africarealty/server/src/kit/storages/es"
	"github.com/africarealty/server/src/kit/storages/pg"
	"github.com/africarealty/server/src/repository/storage"
	"github.com/africarealty/server/src/service"
	"time"
)

const (
	AeroSetAdsCache = "ads_cache"
	EsIndexAds      = "ads"
)

type advertisement struct {
	pg.GormDto
	Id          string     `gorm:"column:id"`
	Code        string     `gorm:"column:code"`
	UserId      string     `gorm:"column:user_id"`
	Status      string     `gorm:"column:status"`
	SubStatus   string     `gorm:"column:sub_status"`
	Type        string     `gorm:"column:type"`
	SubType     string     `gorm:"column:sub_type"`
	Details     string     `gorm:"column:details"`
	ActivatedAt *time.Time `gorm:"column:activated_at"`
	ClosedAt    *time.Time `gorm:"column:closed_at"`
}

type esIdxAds struct {
	Code      string `json:"code" es:"type:keyword"`
	UserId    string `json:"userId" es:"type:keyword"`
	Status    string `json:"status" es:"type:keyword"`
	SubStatus string `json:"subStatus" es:"type:keyword"`
	Type      string `json:"type" es:"type:keyword"`
	SubType   string `json:"subType" es:"type:keyword"`
}

type AdvertisementStorageImpl struct {
	pg   *pg.Storage
	aero kitAero.Aerospike
	es   es.Search
	cfg  *kitAero.Config
}

func (s *AdvertisementStorageImpl) l() log.CLogger {
	return service.L().Cmp("ads-storage")
}

func NewAdvertisementStorage(pg *pg.Storage, aero kitAero.Aerospike, es es.Search, cfg *kitAero.Config) *AdvertisementStorageImpl {
	return &AdvertisementStorageImpl{
		pg:   pg,
		aero: aero,
		es:   es,
		cfg:  cfg,
	}
}

func (s *AdvertisementStorageImpl) init(ctx context.Context) error {
	// build ES index
	return s.es.NewBuilder().WithAlias(EsIndexAds).WithMappingModel(&esIdxAds{}).Build()
}

func (s *AdvertisementStorageImpl) clearCache(ctx context.Context, adsId string) error {
	s.l().Mth("clear-cache").C(ctx).F(log.FF{"adsId": adsId}).Trc()
	key, err := aero.NewKey(storage.AeroNsCache, AeroSetAdsCache, adsId)
	if err != nil {
		return errors.ErrAdsStorageAeroKey(err, ctx)
	}
	_, err = s.aero.Instance().Delete(nil, key)
	if err != nil {
		return errors.ErrAdsStorageClearCache(err, ctx)
	}
	return nil
}

func (s *AdvertisementStorageImpl) getFromCacheById(ctx context.Context, adsId string) (*domain.Advertisement, error) {
	s.l().Mth("get-cache").C(ctx).F(log.FF{"adsId": adsId}).Trc()
	key, err := aero.NewKey(storage.AeroNsCache, AeroSetAdsCache, adsId)
	if err != nil {
		return nil, errors.ErrAdsStorageAeroKey(err, ctx)
	}
	policy := aero.NewPolicy()
	policy.SendKey = true
	rec, err := s.aero.Instance().Get(policy, key)
	if err != nil && !err.Matches(types.KEY_NOT_FOUND_ERROR) {
		return nil, errors.ErrAdsStorageGetCache(err, ctx)
	}
	return s.toAdsCacheDomain(rec), nil
}

func (s *AdvertisementStorageImpl) setCache(ctx context.Context, ads *domain.Advertisement) error {
	s.l().Mth("set-cache").C(ctx).F(log.FF{"adsId": ads.Id}).Trc()
	key, err := aero.NewKey(storage.AeroNsCache, AeroSetAdsCache, ads.Id)
	if err != nil {
		return errors.ErrAdsStorageAeroKey(err, ctx)
	}
	writePolicy := aero.NewWritePolicy(0, 3600)
	writePolicy.SendKey = true
	err = s.aero.Instance().Put(writePolicy, key, s.toAdsCache(ads))
	if err != nil {
		return errors.ErrAdsStoragePutCache(err, ctx)
	}
	return nil
}

func (s *AdvertisementStorageImpl) Create(ctx context.Context, ads *domain.Advertisement) error {
	l := s.l().Mth("create").C(ctx).F(log.FF{"adsId": ads.Id}).Trc()
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	// save to DB
	eg.Go(func() error {
		if err := s.pg.Instance.Create(s.toAdsDto(ads)).Error; err != nil {
			return errors.ErrAdsStorageCreate(err, ctx)
		}
		return nil
	})
	// indexing
	eg.Go(func() error {
		return s.es.Index(EsIndexAds, ads.Id, s.toAdsIndex(ads))
	})
	return eg.Wait()
}

func (s *AdvertisementStorageImpl) Update(ctx context.Context, ads *domain.Advertisement) error {
	l := s.l().Mth("update").C(ctx).F(log.FF{"adsId": ads.Id}).Trc()
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	// save to DB
	eg.Go(func() error {
		dto := s.toAdsDto(ads)
		result := s.pg.Instance.Omit("created_at").Save(dto)
		if result.Error != nil {
			return errors.ErrAdsStorageUpdate(result.Error, ctx)
		}
		return nil
	})
	// clear cache
	eg.Go(func() error {
		return s.clearCache(ctx, ads.Id)
	})
	// indexing
	eg.Go(func() error {
		return s.es.Index(EsIndexAds, ads.Id, s.toAdsIndex(ads))
	})
	return eg.Wait()
}

func (s *AdvertisementStorageImpl) Delete(ctx context.Context, adsId string) error {
	l := s.l().C(ctx).Mth("delete").F(log.FF{"adsId": adsId}).Dbg()
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	eg.Go(func() error {
		// save to DB
		if err := s.pg.Instance.Delete(&advertisement{Id: adsId}).Error; err != nil {
			return errors.ErrAdsStorageDelete(err, ctx)
		}
		return nil
	})
	// clear cache
	eg.Go(func() error {
		return s.clearCache(ctx, adsId)
	})
	// indexing
	eg.Go(func() error {
		return s.es.Delete(EsIndexAds, adsId)
	})
	return eg.Wait()
}

func (s *AdvertisementStorageImpl) Get(ctx context.Context, adsId string) (*domain.Advertisement, error) {
	l := s.l().Mth("get").C(ctx).F(log.FF{"adsId": adsId}).Trc()
	if adsId == "" {
		return nil, nil
	}
	// check cache first
	usr, err := s.getFromCacheById(ctx, adsId)
	if err != nil {
		return nil, err
	}
	if usr != nil {
		l.Trc("found in cache")
		return usr, nil
	}
	// get from db
	dto := &advertisement{Id: adsId}
	res := s.pg.Instance.Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrAdsStorageGetDb(res.Error, ctx)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	usr = s.toAdsDomain(dto)
	// set cache
	err = s.setCache(ctx, usr)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (s *AdvertisementStorageImpl) Search(ctx context.Context, rq *domain.AdsSearchRequest) (*domain.AdsSearchResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AdvertisementStorageImpl) GetCode(ctx context.Context) (string, error) {
	s.l().C(ctx).Mth("gen-code").Dbg()
	var res struct{ Val uint }
	if err := s.pg.Instance.Raw("select nextval('seq_ads_code') as val").First(&res).Error; err != nil {
		return "", errors.ErrAdsStorageDbNextSeqVal(err, ctx)
	}
	return fmt.Sprintf("%d", res.Val), nil
}
