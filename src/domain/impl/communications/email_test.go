package communications

import (
	"bytes"
	"encoding/json"
	er "errors"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/communications"
	"github.com/africarealty/server/src/kit"
	kitCtx "github.com/africarealty/server/src/kit/context"
	"github.com/africarealty/server/src/kit/queue"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/mocks"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var cfg = &service.Config{
	Communications: &service.CfgCommunications{
		Email: &service.CfgEmail{
			SmtpServer:     "localhost",
			SmtpServerPort: "25",
			SmtpUser:       "none",
			SmtpPassword:   "none",
			SmtpFrom:       "test-default-from@test.mail",
		},
	},
}

type emailTestSuite struct {
	kitTestSuite.Suite
	emailSvc          domain.EmailService
	emailRepo         *mocks.EmailRepository
	templateGenerator *mocks.TemplateGenerator
	store             *mocks.StoreService
	queue             *mocks.Queue
}

func (s *emailTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
	s.queue = &mocks.Queue{}
	s.queue.On("Publish", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("queue.QueueType"), mock.AnythingOfType("string"), mock.AnythingOfType("*queue.Message")).Return(nil)
}

func (s *emailTestSuite) SetupTest() {
	s.store = &mocks.StoreService{}
	s.templateGenerator = &mocks.TemplateGenerator{}
	s.emailRepo = &mocks.EmailRepository{}
	s.emailSvc = NewEmailService(s.templateGenerator, s.queue, s.store, s.emailRepo)
	if err := s.emailSvc.Init(cfg); err != nil {
		s.Fatal(err)
	}
}

func TestEmailSuite(t *testing.T) {
	suite.Run(t, new(emailTestSuite))
}

func (s *emailTestSuite) Test_Send_Success() {
	eq := &domain.EmailRequest{
		UserId:      "1",
		Email:       "test@test.mail",
		From:        "from@test.mail",
		Template:    &domain.TemplateRequest{Id: "1"},
		LinkFileIds: []string{kit.NewRandString(), kit.NewRandString()},
	}
	title := "title"
	body := "body"

	s.templateGenerator.On("Generate", s.Ctx, mock.AnythingOfType("*domain.TemplateRequest")).Return(&domain.TemplateResponse{Title: title, Body: body}, nil)

	email, err := s.emailSvc.Send(s.Ctx, eq)
	s.NoError(err)
	s.Equal(title, email.Subject)
	s.Equal(body, email.Text)
	s.Equal(eq.UserId, email.UserId)
	s.Equal(eq.Email, email.Email)
	s.Equal(eq.From, email.From)
	s.Equal(eq.Template, email.Template)
	s.Equal(eq.LinkFileIds, email.LinkFileIds)
	s.Equal(domain.EmailRqStatusNotSend, email.SendStatus)
}

func (s *emailTestSuite) Test_Send_WithDefaultFrom_Success() {
	eq := &domain.EmailRequest{
		UserId:   "1",
		Email:    "test@test.mail",
		Template: &domain.TemplateRequest{Id: "1"},
	}
	title := "title"
	body := "body"

	s.templateGenerator.On("Generate", s.Ctx, mock.AnythingOfType("*domain.TemplateRequest")).Return(&domain.TemplateResponse{Title: title, Body: body}, nil)

	email, err := s.emailSvc.Send(s.Ctx, eq)
	s.NoError(err)

	s.Equal(cfg.Communications.Email.SmtpFrom, email.From)
}

func (s *emailTestSuite) Test_Send_ValidationError() {
	tests := []struct {
		name string
		eq   *domain.EmailRequest
		err  string
	}{
		{
			name: "with not valid email",
			eq:   &domain.EmailRequest{Email: "t@s@test.cc", Template: &domain.TemplateRequest{Id: "1"}},
			err:  errors.ErrCodeEmailValidationInvalidEmail,
		},
		{
			name: "with not valid from",
			eq:   &domain.EmailRequest{Email: "test@test.cc", From: "t@s@test.cc", Template: &domain.TemplateRequest{Id: "1"}},
			err:  errors.ErrCodeEmailValidationInvalidFrom,
		},
		{
			name: "with empty template",
			eq:   &domain.EmailRequest{Email: "test@test.cc"},
			err:  errors.ErrCodeTemplateEmpty,
		},
		{
			name: "with empty template id",
			eq:   &domain.EmailRequest{Email: "test@test.cc", Template: &domain.TemplateRequest{Id: ""}},
			err:  errors.ErrCodeTemplateEmpty,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			_, err := s.emailSvc.Send(s.Ctx, tt.eq)
			s.AssertAppErr(err, tt.err)
		})
	}
}

func (s *emailTestSuite) Test_Send_TemplateGeneratorError() {
	eq := &domain.EmailRequest{Email: "test@test.cc", Template: &domain.TemplateRequest{Id: "1"}}
	errTemplate := errors.ErrTemplateNotFound(s.Ctx)

	s.templateGenerator.On("Generate", s.Ctx, mock.AnythingOfType("*domain.TemplateRequest")).Return(nil, errTemplate)

	_, err := s.emailSvc.Send(s.Ctx, eq)
	s.AssertAppErr(err, errors.ErrCodeTemplateNotFound)
}

func (s *emailTestSuite) Test_RequestHandler_Success() {
	fileId := kit.NewRandString()
	email := domain.Email{
		Id:          kit.NewId(),
		Subject:     "Отправка пароля пользователю",
		Text:        "Уважаемый Иван, Ваш пароль доступа к moi-service.ru: 123456",
		Email:       "test@test.mail",
		From:        "from@test.mail",
		Template:    &domain.TemplateRequest{Id: "generated-password", Data: map[string]interface{}{"Password": "123456", "UserFullName": "Иван"}},
		SendStatus:  domain.EmailRqStatusNotSend,
		LinkFileIds: []string{fileId},
	}
	file := &domain.FileContent{
		Filename:    "filename",
		FileID:      fileId,
		ContentType: "text/plain",
		Extension:   "txt",
		Content:     bytes.NewReader([]byte("content")),
	}

	s.store.On("GetFile", mock.AnythingOfType("*context.valueCtx"), mock.MatchedBy(func(fileID string) bool {
		s.Equal(fileId, fileID)
		return true
	})).Return(file, nil)
	s.emailRepo.On("Send", mock.AnythingOfType("*context.valueCtx"), mock.MatchedBy(func(e *domain.Email) bool {
		s.Equal([]*domain.FileContent{file}, e.Attachments)
		s.Equal(email.Email, e.Email)
		return true
	})).Return(nil)

	msg := &queue.Message{
		Ctx:     &kitCtx.RequestContext{Rid: kit.NewRandString(), Caller: "test"},
		Payload: email,
	}
	payload, _ := json.Marshal(msg)

	err := s.emailSvc.RequestHandler()(payload)
	s.NoError(err)
}

func (s *emailTestSuite) Test_RequestHandler_GetFileError() {
	getFileError := er.New("get-file-error")

	s.store.On("GetFile", mock.AnythingOfType("*context.valueCtx"), mock.Anything).Return(nil, getFileError)

	payload, _ := json.Marshal(&queue.Message{
		Ctx: &kitCtx.RequestContext{Rid: "rid", Caller: "test"},
		Payload: domain.Email{
			LinkFileIds: []string{kit.NewRandString()},
		},
	})

	err := s.emailSvc.RequestHandler()(payload)
	s.ErrorIs(err, getFileError)
}

func (s *emailTestSuite) Test_RequestHandler_DecodeError() {
	err := s.emailSvc.RequestHandler()([]byte("notvalidjson"))
	s.AssertAppErr(err, queue.ErrCodeQueueMsgUnmarshal)
}
