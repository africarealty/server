package storage

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
	kitService.StorageAdapter
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

func (c *adapterImpl) Init(ctx context.Context, cfg interface{}) error {
	config := cfg.(*service.Config)

	// init postgres
	var err error
	c.pg, err = pg.Open(config.Storages.Pg.Master, service.LF())
	if err != nil {
		return err
	}

	// applying migrations
	if config.Storages.Pg.MigPath != "" {
		db, _ := c.pg.Instance.DB()
		m := pg.NewMigration(db, config.Storages.Pg.MigPath, service.LF())
		if err := m.Up(); err != nil {
			return err
		}
	}

	// init aero
	err = c.aero.Open(ctx, config.Storages.Aero, service.LF())
	if err != nil {
		return err
	}

	// init storages
	c.UserStorageImpl = NewUserStorage(c.pg, c.aero, config.Storages.Aero)
	err = c.UserStorageImpl.init(ctx)
	if err != nil {
		return err
	}
	c.SessionStorageImpl = NewSessionStorage(c.pg, c.aero, config.Storages.Aero)
	return nil
}

func (c *adapterImpl) Close(ctx context.Context) error {
	if c.aero != nil {
		_ = c.aero.Close(ctx)
	}
	if c.pg != nil {
		c.pg.Close()
	}
	return nil
}
