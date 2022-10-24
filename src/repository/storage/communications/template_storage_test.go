//+build integration

package communications

import (
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type templateStorageTestSuite struct {
	kitTestSuite.Suite
	storage domain.TemplateStorage
	cfg     *service.Config
	adapter Adapter
}

// SetupSuite is called once for a suite
func (s *templateStorageTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())

	// load config
	var err error
	s.cfg, err = service.LoadConfig()
	if err != nil {
		s.Fatal(err)
	}
	// disable applying migrations
	s.cfg.Storages.Pg.MigPath = ""
}

// SetupTest is called for each test
func (s *templateStorageTestSuite) SetupTest() {
	// initialize adapter for each test,
	// because in some of them we can close internal services
	s.adapter = NewAdapter()
	err := s.adapter.Init(s.Ctx, s.cfg)
	if err != nil {
		s.Fatal(err)
	}
	s.storage = s.adapter.GetTemplateStorage()
}

func (s *templateStorageTestSuite) TearDownTest() {
	_ = s.adapter.Close(s.Ctx)
}

func TestTemplateSuite(t *testing.T) {
	suite.Run(t, new(templateStorageTestSuite))
}

func (s *templateStorageTestSuite) Test_CreateUpdateGetSearchTemplate_Success() {
	t := &domain.Template{
		Id:    kit.NewRandString(),
		Title: kit.NewRandString(),
		Body:  kit.NewRandString(),
	}

	//create
	err := s.storage.Create(s.Ctx, t)
	s.NoError(err)

	//get
	tStored, err := s.storage.Get(s.Ctx, t.Id)
	s.NoError(err)
	s.Equal(t.Id, tStored.Id)
	s.NotEmpty(t.Title, tStored.Title)
	s.NotEmpty(t.Body, tStored.Body)

	//update
	t.Title = kit.NewRandString()
	err = s.storage.Update(s.Ctx, t)
	s.NoError(err)

	//get
	tStored, err = s.storage.Get(s.Ctx, t.Id)
	s.NoError(err)
	s.NotEmpty(t.Title, tStored.Title)

	//search
	tStoredList, err := s.storage.Search(s.Ctx, t.Title[:2])
	s.NoError(err)
	s.NotEmpty(t.Title, tStoredList)
	s.True(len(tStoredList) > 0)

	//delete
	err = s.storage.Delete(s.Ctx, t.Id)
	s.NoError(err)
}

func (s *templateStorageTestSuite) Test_GetTemplateBody_Error() {
	t, err := s.storage.Get(s.Ctx, "not_exist")
	s.Nil(err)
	s.Nil(t)
}

func (s *templateStorageTestSuite) Test_GetTemplateBody_WhenIdIsEmpty() {
	t, err := s.storage.Get(s.Ctx, "")
	s.Nil(err)
	s.Nil(t)
}
