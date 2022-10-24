package communications

import (
	"context"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/aerospike/aerospike-client-go/v6/types"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/communications"
	"github.com/africarealty/server/src/kit/goroutine"
	"github.com/africarealty/server/src/kit/log"
	kitAero "github.com/africarealty/server/src/kit/storages/aerospike"
	"github.com/africarealty/server/src/kit/storages/pg"
	"github.com/africarealty/server/src/repository/storage"
	"github.com/africarealty/server/src/service"
)

const (
	AeroSetCommunicationTemplateCache = "template_cache"
)

type TemplateStorageImpl struct {
	pg   *pg.Storage
	aero kitAero.Aerospike
}

func NewTemplateStorage(pg *pg.Storage, aero kitAero.Aerospike) *TemplateStorageImpl {
	return &TemplateStorageImpl{
		pg:   pg,
		aero: aero,
	}
}

type template struct {
	pg.GormDto
	Id    string `gorm:"column:id;primaryKey"`
	Title string `gorm:"column:title"`
	Body  string `gorm:"column:body"`
}

func (t *TemplateStorageImpl) l() log.CLogger {
	return service.L().Cmp("template-storage")
}

func (t *TemplateStorageImpl) clearCacheTemplates(ctx context.Context, templatesId string) error {
	key, err := aero.NewKey(storage.AeroNsCache, AeroSetCommunicationTemplateCache, templatesId)
	if err != nil {
		return errors.ErrTemplateStorageAeroKey(err, ctx)
	}
	_, err = t.aero.Instance().Delete(nil, key)
	if err != nil {
		return errors.ErrTemplateStorageClearCache(err, ctx)
	}
	return nil
}

func (t *TemplateStorageImpl) getFromCacheById(ctx context.Context, templateId string) (*domain.Template, error) {
	t.l().Mth("get-cache").C(ctx).F(log.FF{"templateId": templateId}).Trc()
	key, err := aero.NewKey(storage.AeroNsCache, AeroSetCommunicationTemplateCache, templateId)
	if err != nil {
		return nil, errors.ErrTemplateStorageAeroKey(err, ctx)
	}
	policy := aero.NewPolicy()
	policy.SendKey = true
	rec, err := t.aero.Instance().Get(policy, key)
	if err != nil && !err.Matches(types.KEY_NOT_FOUND_ERROR) {
		return nil, errors.ErrTemplateStorageGetCache(err, ctx)
	}
	return t.toTemplateCacheDomain(rec), nil
}

func (t *TemplateStorageImpl) setCache(ctx context.Context, template *domain.Template) error {
	t.l().Mth("set-cache").C(ctx).F(log.FF{"templateId": template.Id}).Trc()
	key, err := aero.NewKey(storage.AeroNsCache, AeroSetCommunicationTemplateCache, template.Id)
	if err != nil {
		return errors.ErrTemplateStorageAeroKey(err, ctx)
	}
	writePolicy := aero.NewWritePolicy(0, 3600)
	writePolicy.SendKey = true
	err = t.aero.Instance().Put(writePolicy, key, t.toTemplateCache(template))
	if err != nil {
		return errors.ErrTemplateStoragePutCache(err, ctx)
	}
	return nil
}

func (t *TemplateStorageImpl) Get(ctx context.Context, templateId string) (*domain.Template, error) {
	l := t.l().Mth("get").C(ctx).F(log.FF{"templateId": templateId}).Trc()
	if templateId == "" {
		return nil, nil
	}
	// check cache first
	templ, err := t.getFromCacheById(ctx, templateId)
	if err != nil {
		return nil, err
	}
	if templ != nil {
		l.Trc("found in cache")
		return templ, nil
	}
	// get from db
	dto := &template{Id: templateId}
	res := t.pg.Instance.Limit(1).Find(&dto)
	if res.Error != nil {
		return nil, errors.ErrTemplateStorageGetDb(res.Error, ctx)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	templ = t.toTemplateDomain(dto)
	// set cache
	err = t.setCache(ctx, templ)
	if err != nil {
		return nil, err
	}
	return templ, nil
}

func (t *TemplateStorageImpl) Create(ctx context.Context, template *domain.Template) error {
	t.l().C(ctx).Mth("create").Dbg()
	if err := t.pg.Instance.Create(t.toTemplateDto(template)).Error; err != nil {
		return errors.ErrTemplateStorageDbCreate(err, ctx)
	}
	return nil
}

func (t *TemplateStorageImpl) Update(ctx context.Context, template *domain.Template) error {
	l := t.l().C(ctx).Mth("update").Dbg()
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	// save to store
	eg.Go(func() error {
		err := t.pg.Instance.Updates(t.toTemplateDto(template)).Error
		if err != nil {
			return errors.ErrTemplateStorageDbUpdate(err, ctx)
		}
		return nil
	})
	// delete cache key
	eg.Go(func() error {
		return t.clearCacheTemplates(ctx, template.Id)
	})
	return eg.Wait()
}

func (t *TemplateStorageImpl) Delete(ctx context.Context, id string) error {
	l := t.l().C(ctx).Mth("delete").Dbg()
	eg := goroutine.NewGroup(ctx).WithLogger(l)
	// delete from storage
	eg.Go(func() error {
		result := t.pg.Instance.Delete(&template{Id: id})
		if result.Error != nil {
			return errors.ErrTemplateStorageDbDelete(result.Error, ctx)
		}
		return nil
	})
	// delete from cache
	eg.Go(func() error {
		return t.clearCacheTemplates(ctx, id)
	})
	return eg.Wait()
}

func (t *TemplateStorageImpl) Search(ctx context.Context, query string) ([]*domain.Template, error) {
	t.l().C(ctx).Mth("search").Dbg()
	var dtos []*template
	if res := t.pg.Instance.Where("UPPER(title) like UPPER(?)", "%"+query+"%").Find(&dtos); res.Error == nil {
		return t.toTemplatesDomain(dtos), nil
	} else {
		return nil, errors.ErrTemplateStorageDbSearch(res.Error, ctx)
	}
}
