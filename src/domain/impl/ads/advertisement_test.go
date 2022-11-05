package ads

import (
	"github.com/africarealty/server/src/domain"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/mocks"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type adsTestSuite struct {
	kitTestSuite.Suite
	storage *mocks.AdvertisementStorage
	svc     domain.AdvertisementService
}

func (s *adsTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
}

func (s *adsTestSuite) SetupTest() {
	s.storage = &mocks.AdvertisementStorage{}
	s.svc = NewAdvertisementService(s.storage)
}

func TestAdsSuite(t *testing.T) {
	suite.Run(t, new(adsTestSuite))
}
