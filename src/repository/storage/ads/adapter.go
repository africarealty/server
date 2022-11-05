package ads

import (
	"context"
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit/goroutine"
	kitService "github.com/africarealty/server/src/kit/service"
	kitAero "github.com/africarealty/server/src/kit/storages/aerospike"
	"github.com/africarealty/server/src/kit/storages/es"
	"github.com/africarealty/server/src/kit/storages/pg"
	"github.com/africarealty/server/src/service"
)

type Adapter interface {
	kitService.Adapter
	domain.AdvertisementStorage
}

type adapterImpl struct {
	*AdvertisementStorageImpl
	aero kitAero.Aerospike
	pg   *pg.Storage
	es   es.Search
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		aero: kitAero.New(),
	}
	return a
}

func (a *adapterImpl) Init(ctx context.Context, cfg interface{}) error {
	config := cfg.(*service.Config)

	grp := goroutine.
		NewGroup(context.Background()).
		WithLoggerFn(service.LF()).
		Cmp("ads-storage").
		Mth("init")

	// init postgres
	grp.Go(func() error {
		var err error
		a.pg, err = pg.Open(config.Storages.Pg.Master, service.LF())
		if err != nil {
			return err
		}

		// applying migrations
		if config.Storages.Pg.MigPath != "" {
			db, _ := a.pg.Instance.DB()
			m := pg.NewMigration(db, config.Storages.Pg.MigPath, service.LF())
			if err := m.Up(); err != nil {
				return err
			}
		}
		return nil
	})

	// init aero
	grp.Go(func() error {
		return a.aero.Open(ctx, config.Storages.Aero, service.LF())
	})

	// Index search
	grp.Go(func() error {
		var err error
		a.es, err = es.NewEs(config.Storages.Es, service.LF())
		return err
	})

	if err := grp.Wait(); err != nil {
		return err
	}

	// init storages
	a.AdvertisementStorageImpl = NewAdvertisementStorage(a.pg, a.aero, a.es, config.Storages.Aero)
	err := a.init(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapterImpl) Close(ctx context.Context) error {
	if a.aero != nil {
		_ = a.aero.Close(ctx)
	}
	if a.pg != nil {
		a.pg.Close()
	}
	return nil
}
