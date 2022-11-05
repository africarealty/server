package ads

import (
	"context"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/ads"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
)

type advertisementServiceImpl struct {
	storage domain.AdvertisementStorage
}

func NewAdvertisementService(storage domain.AdvertisementStorage) domain.AdvertisementService {
	return &advertisementServiceImpl{
		storage: storage,
	}
}

func (a *advertisementServiceImpl) l() log.CLogger {
	return service.L().Cmp("ads-svc")
}

func (a *advertisementServiceImpl) validate(ctx context.Context, ads *domain.Advertisement) error {
	return nil
}

func (a *advertisementServiceImpl) Create(ctx context.Context, ads *domain.Advertisement) (*domain.Advertisement, error) {
	a.l().C(ctx).Mth("create").Dbg()

	// set id
	ads.Id = kit.NewId()

	// set code
	var err error
	ads.Code, err = a.storage.GetCode(ctx)
	if err != nil {
		return nil, err
	}

	// status/sub-status
	if ads.Status == "" {
		ads.Status = domain.AdStatusDraft
		ads.SubStatus = domain.AdSubStatusDraft
	}

	// validate ads
	err = a.validate(ctx, ads)
	if err != nil {
		return nil, err
	}

	// save to storage
	err = a.storage.Create(ctx, ads)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

func (a *advertisementServiceImpl) Update(ctx context.Context, ads *domain.Advertisement) (*domain.Advertisement, error) {
	a.l().C(ctx).Mth("update").F(log.FF{"adsId": ads.Id}).Dbg()

	// get stored ads
	stored, err := a.Get(ctx, ads.Id)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrAdsNotFound(ctx, ads.Id)
	}

	// populate
	ads.Type = stored.Type
	ads.SubType = stored.SubType
	ads.Status = stored.Status
	ads.SubStatus = stored.SubStatus
	ads.Id = stored.Id

	// validate ads
	err = a.validate(ctx, ads)
	if err != nil {
		return nil, err
	}

	// save to storage
	err = a.storage.Update(ctx, ads)
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (a *advertisementServiceImpl) Delete(ctx context.Context, adsId string) error {
	a.l().C(ctx).Mth("delete").F(log.FF{"adsId": adsId}).Dbg()

	// get stored ads
	stored, err := a.Get(ctx, adsId)
	if err != nil {
		return err
	}
	if stored == nil {
		return errors.ErrAdsNotFound(ctx, adsId)
	}

	// delete
	err = a.storage.Delete(ctx, adsId)
	if err != nil {
		return err
	}
	return nil
}

func (a *advertisementServiceImpl) SetStatus(ctx context.Context, adsId, status, subStatus string) (*domain.Advertisement, error) {
	a.l().C(ctx).Mth("set-status").F(log.FF{"adsId": adsId, "status": status, "subStatus": subStatus}).Dbg()

	// get stored ads
	stored, err := a.Get(ctx, adsId)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrAdsNotFound(ctx, adsId)
	}

	stored.Status = status
	stored.SubStatus = subStatus

	// delete
	err = a.storage.Update(ctx, stored)
	if err != nil {
		return nil, err
	}
	return stored, nil
}

func (a *advertisementServiceImpl) Get(ctx context.Context, adsId string) (*domain.Advertisement, error) {
	a.l().C(ctx).Mth("set-status").F(log.FF{"adsId": adsId}).Dbg()
	// check id
	if adsId == "" {
		return nil, errors.ErrAdsIdEmpty(ctx)
	}
	return a.storage.Get(ctx, adsId)
}

func (a *advertisementServiceImpl) Search(ctx context.Context, rq *domain.AdsSearchRequest) (*domain.AdsSearchResponse, error) {
	//TODO implement me
	panic("implement me")
}
