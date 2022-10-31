package communications

import (
	"context"
	"github.com/africarealty/server/src/domain"
	kitService "github.com/africarealty/server/src/kit/service"
	kitAero "github.com/africarealty/server/src/kit/storages/aerospike"
	"github.com/africarealty/server/src/kit/storages/pg"
	"github.com/africarealty/server/src/service"
)

// Adapter provides a contract to access a remote service
type Adapter interface {
	kitService.Adapter
	// GetEmailStorage - return an email storage interface
	GetEmailStorage() domain.EmailStorage
	// GetTemplateStorage - return a template storage interface
	GetTemplateStorage() domain.TemplateStorage
}

// adapterImpl implements storage adapter
type adapterImpl struct {
	*emailStorageImpl
	*templateStorageImpl
	aero kitAero.Aerospike
	pg   *pg.Storage
}

// NewAdapter creates a new instance of the adapter
func NewAdapter() Adapter {
	return &adapterImpl{
		aero: kitAero.New(),
	}
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
	a.emailStorageImpl = newEmailStorage(a.pg)
	a.templateStorageImpl = newTemplateStorage(a.pg, a.aero)
	return nil
}

func (a *adapterImpl) GetEmailStorage() domain.EmailStorage {
	return a
}

func (a *adapterImpl) GetTemplateStorage() domain.TemplateStorage {
	return a
}

func (a *adapterImpl) Close(ctx context.Context) error {
	a.pg.Close()
	_ = a.aero.Close(ctx)
	return nil
}
