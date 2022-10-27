package communications

import (
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/communications"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/mocks"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type templateTestSuite struct {
	kitTestSuite.Suite
	storage *mocks.TemplateStorage
	service domain.TemplateService
}

func (s *templateTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
	s.storage = &mocks.TemplateStorage{}
}

func (s *templateTestSuite) SetupTest() {
	s.storage.ExpectedCalls = nil
	s.service = NewTemplateService(s.storage)
}

func TestTemplateSuite(t *testing.T) {
	suite.Run(t, new(templateTestSuite))
}

func (s *templateTestSuite) testTemplate() *domain.Template {
	return &domain.Template{
		Id:    "new",
		Title: "My new rec",
		Body:  "Dear {user}! Have a good day!",
	}
}

func (s *templateTestSuite) Test_Create_Success() {
	t := s.testTemplate()
	s.storage.On("Get", s.Ctx, t.Id).Return(nil, nil)
	s.storage.On("Create", s.Ctx, t).Return(nil)
	t, err := s.service.CreateTemplate(s.Ctx, t)
	s.NoError(err)
	s.NotNil(t)
}

func (s *templateTestSuite) Test_Create_EmptyTitle_Fail() {
	t := s.testTemplate()
	t.Title = ""
	s.storage.On("Get", s.Ctx, t.Id).Return(nil, nil)
	_, err := s.service.CreateTemplate(s.Ctx, t)
	s.Error(err)
	s.AssertAppErr(err, errors.ErrCodeTemplateTitleEmpty)
}

func (s *templateTestSuite) Test_Create_EmptyId_Fail() {
	t := s.testTemplate()
	t.Id = ""
	s.storage.On("Get", s.Ctx, t.Id).Return(nil, nil)
	_, err := s.service.CreateTemplate(s.Ctx, t)
	s.Error(err)
	s.AssertAppErr(err, errors.ErrCodeTemplateIdEmpty)
}

func (s *templateTestSuite) Test_Create_EmptyBody_Fail() {
	t := s.testTemplate()
	t.Body = ""
	s.storage.On("Get", s.Ctx, t.Id).Return(nil, nil)
	_, err := s.service.CreateTemplate(s.Ctx, t)
	s.Error(err)
	s.AssertAppErr(err, errors.ErrCodeTemplateBodyEmpty)
}

func (s *templateTestSuite) Test_Create_AlreadyExists_Fail() {
	t := s.testTemplate()
	s.storage.On("Get", s.Ctx, t.Id).Return(t, nil)
	_, err := s.service.CreateTemplate(s.Ctx, t)
	s.Error(err)
	s.AssertAppErr(err, errors.ErrCodeTemplateAlreadyExists)
}

func (s *templateTestSuite) Test_Update_Success() {
	t := s.testTemplate()
	s.storage.On("Get", s.Ctx, t.Id).Return(t, nil)
	s.storage.On("Update", s.Ctx, t).Return(nil)
	t, err := s.service.UpdateTemplate(s.Ctx, t)
	s.NoError(err)
	s.NotNil(t)
}

func (s *templateTestSuite) Test_Update_EmptyTitle_Fail() {
	t := s.testTemplate()
	t.Title = ""
	s.storage.On("Get", s.Ctx, t.Id).Return(nil, nil)
	_, err := s.service.UpdateTemplate(s.Ctx, t)
	s.Error(err)
	s.AssertAppErr(err, errors.ErrCodeTemplateTitleEmpty)
}

func (s *templateTestSuite) Test_Update_EmptyId_Fail() {
	t := s.testTemplate()
	t.Id = ""
	s.storage.On("Get", s.Ctx, t.Id).Return(nil, nil)
	_, err := s.service.UpdateTemplate(s.Ctx, t)
	s.Error(err)
	s.AssertAppErr(err, errors.ErrCodeTemplateIdEmpty)
}

func (s *templateTestSuite) Test_Update_EmptyBody_Fail() {
	t := s.testTemplate()
	t.Body = ""
	s.storage.On("Get", s.Ctx, t.Id).Return(nil, nil)
	_, err := s.service.UpdateTemplate(s.Ctx, t)
	s.Error(err)
	s.AssertAppErr(err, errors.ErrCodeTemplateBodyEmpty)
}

func (s *templateTestSuite) Test_Update_NotExists_Fail() {
	t := s.testTemplate()
	s.storage.On("Get", s.Ctx, t.Id).Return(nil, nil)
	_, err := s.service.UpdateTemplate(s.Ctx, t)
	s.Error(err)
	s.AssertAppErr(err, errors.ErrCodeTemplateNotFound)
}

func (s *templateTestSuite) Test_Get_Success() {
	t := s.testTemplate()
	s.storage.On("Get", s.Ctx, t.Id).Return(t, nil)
	t, err := s.service.GetTemplate(s.Ctx, t.Id)
	s.NoError(err)
	s.NotNil(t)
}

func (s *templateTestSuite) Test_Get_NotFound() {
	t := s.testTemplate()
	s.storage.On("Get", s.Ctx, t.Id).Return(nil, nil)
	t, err := s.service.GetTemplate(s.Ctx, t.Id)
	s.AssertAppErr(err, errors.ErrCodeTemplateNotFound)
	s.Nil(t)
}

func (s *templateTestSuite) Test_Search_Success() {
	t := s.testTemplate()
	s.storage.On("Search", s.Ctx, "xxx").Return([]*domain.Template{t, t}, nil)
	list, err := s.service.SearchTemplates(s.Ctx, "xxx")
	s.NoError(err)
	s.NotNil(list)
	s.Len(list, 2)
}

func (s *templateTestSuite) Test_Search_NotFound() {
	s.storage.On("Search", s.Ctx, "xxx").Return(nil, nil)
	list, err := s.service.SearchTemplates(s.Ctx, "xxx")
	s.Nil(list)
	s.AssertAppErr(err, errors.ErrCodeTemplateNotFound)
}
