package impl

import (
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/communications"
	"github.com/africarealty/server/src/kit"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/mocks"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type templateGeneratorTestSuite struct {
	kitTestSuite.Suite
	templGen domain.TemplateGenerator
	storage  *mocks.TemplateStorage
}

func (s *templateGeneratorTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
}

func (s *templateGeneratorTestSuite) SetupTest() {
	s.storage = &mocks.TemplateStorage{}
	s.templGen = NewTemplateGenerator(s.storage)
}

func TestPushTemplateSuite(t *testing.T) {
	suite.Run(t, new(templateGeneratorTestSuite))
}

func (s *templateGeneratorTestSuite) Test_Generate_Success() {
	t := &domain.Template{
		Id:    kit.NewRandString(),
		Title: "test",
		Body:  "Some body {{Placeholder}}",
	}
	s.storage.On("Get", s.Ctx, t.Id).Return(t, nil)
	res, err := s.templGen.Generate(s.Ctx, &domain.TemplateRequest{Id: t.Id, Data: map[string]interface{}{"Placeholder": "value"}})
	s.NoError(err)
	s.Equal("Some body value", res.Body)
	s.Equal("test", res.Title)
}

func (s *templateGeneratorTestSuite) Test_Generate_WithTitlePlaceholder_Success() {
	t := &domain.Template{
		Id:    kit.NewRandString(),
		Title: "test {{TitlePlaceholder}}",
		Body:  "Some body {{Placeholder}}",
	}
	s.storage.On("Get", s.Ctx, t.Id).Return(t, nil)
	res, err := s.templGen.Generate(s.Ctx, &domain.TemplateRequest{
		Id: t.Id,
		Data: map[string]interface{}{
			"Placeholder":      "value",
			"TitlePlaceholder": "value2",
		},
	})
	s.NoError(err)
	s.Equal("Some body value", res.Body)
	s.Equal("test value2", res.Title)
}

func (s *templateGeneratorTestSuite) Test_Generate_WhenDataNil_Success() {
	t := &domain.Template{
		Id:    kit.NewRandString(),
		Title: "test",
		Body:  "Some body",
	}
	s.storage.On("Get", s.Ctx, t.Id).Return(t, nil)
	res, err := s.templGen.Generate(s.Ctx, &domain.TemplateRequest{Id: t.Id, Data: nil})
	s.NoError(err)
	s.Equal("Some body", res.Body)
	s.Equal("test", res.Title)
}

func (s *templateGeneratorTestSuite) Test_Generate_Error() {
	t := &domain.Template{
		Id:    kit.NewRandString(),
		Title: "test",
		Body:  "Some body {{{Placeholder}}",
	}
	s.storage.On("Get", s.Ctx, t.Id).Return(t, nil)
	_, err := s.templGen.Generate(s.Ctx, &domain.TemplateRequest{Id: t.Id, Data: nil})
	s.T().Log(err)
	s.Error(err)
	s.AssertAppErr(err, errors.ErrCodeTemplateGenerator)
}
