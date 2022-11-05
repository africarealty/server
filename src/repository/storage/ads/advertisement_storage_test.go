// go:build integration

package ads

import (
	"github.com/africarealty/server/src/domain"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type adsStorageTestSuite struct {
	kitTestSuite.Suite
	storage domain.AdvertisementStorage
	adapter Adapter
}

func (s *adsStorageTestSuite) SetupSuite() {

	s.Suite.Init(service.LF())

	// load config
	cfg, err := service.LoadConfig()
	if err != nil {
		s.Fatal(err)
	}

	// disable applying migrations
	cfg.Storages.Pg.MigPath = ""

	// initialize adapter
	s.adapter = NewAdapter()
	err = s.adapter.Init(s.Ctx, cfg)
	if err != nil {
		s.Fatal(err)
	}
	s.storage = s.adapter
}

func (s *adsStorageTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func (s *adsStorageTestSuite) SetupTest() {}

func TestAdsStorageSuite(t *testing.T) {
	suite.Run(t, new(adsStorageTestSuite))
}
