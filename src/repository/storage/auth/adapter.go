package auth

import (
	"context"
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit/auth"
	kitService "github.com/africarealty/server/src/kit/service"
	kitAero "github.com/africarealty/server/src/kit/storages/aerospike"
	"github.com/africarealty/server/src/kit/storages/pg"
	"github.com/africarealty/server/src/service"
)

type Adapter interface {
	kitService.Adapter
	domain.UserStorage
	auth.SessionStorage
}

type adapterImpl struct {
	*UserStorageImpl
	*SessionStorageImpl
	aero kitAero.Aerospike
	pg   *pg.Storage
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		aero: kitAero.New(),
	}
	return a
}

func (a *adapterImpl) Init(ctx context.Context, cfg interface{}) error {
	config := cfg.(*service.Config)

	// init postgres
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

	// init aero
	err = a.aero.Open(ctx, config.Storages.Aero, service.LF())
	if err != nil {
		return err
	}

	// init storages
	a.UserStorageImpl = NewUserStorage(a.pg, a.aero, config.Storages.Aero)
	err = a.UserStorageImpl.init(ctx)
	if err != nil {
		return err
	}
	a.SessionStorageImpl = NewSessionStorage(a.pg, a.aero, config.Storages.Aero)
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
