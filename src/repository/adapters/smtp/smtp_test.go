package smtp

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/africarealty/server/src/domain"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/mocks"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/smtp"
	"testing"
)

type smtpClientTestSuite struct {
	kitTestSuite.Suite
	storage *mocks.EmailStorage
	email   *mocks.EmailService
}

var cfg = &service.CfgEmail{
	SmtpServer:     "localhost",
	SmtpServerPort: "25",
	SmtpUser:       "none",
	SmtpPassword:   "none",
	SmtpFrom:       "from@test.com",
}

// SetupSuite is called once for a suite
func (s *smtpClientTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
	s.email = &mocks.EmailService{}
}

// SetupTest is called once for a test
func (s *smtpClientTestSuite) SetupTest() {
	s.storage = &mocks.EmailStorage{}
}

func TestSmtpClientSuite(t *testing.T) {
	suite.Run(t, new(smtpClientTestSuite))
}

func (s *smtpClientTestSuite) Test_Send_SentSuccessfully() {
	reader1 := bytes.NewReader([]byte("filecontent1"))
	reader2 := bytes.NewReader([]byte("filecontent2"))
	tests := []struct {
		name          string
		requestEmail  *domain.Email
		expectedEmail *domain.Email
		expectedMsg   []byte
	}{
		{
			name: "without attachments",
			requestEmail: &domain.Email{
				Id:      "1",
				Email:   "test@mail.com",
				From:    "from@test.mail",
				Text:    "bodytext",
				Subject: "subj",
				Template: &domain.TemplateRequest{
					Id: "verifcode",
					Data: map[string]interface{}{
						"UserFullName": "Иван",
						"verifcode":    "123456",
					},
				},
			},
			expectedEmail: &domain.Email{
				Id:      "1",
				Email:   "test@mail.com",
				From:    "from@test.mail",
				Text:    "bodytext",
				Subject: "subj",
				Template: &domain.TemplateRequest{
					Id: "verifcode",
					Data: map[string]interface{}{
						"UserFullName": "Иван",
						"verifcode":    "123456",
					},
				},
			},
		},
		{
			name: "with attachments",
			requestEmail: &domain.Email{
				Id:      "2",
				Email:   "test@mail.com",
				From:    "from@test.mail",
				Text:    "bodytext",
				Subject: "subj",
				Attachments: []*domain.FileContent{
					{FileID: "13", Filename: "filename13", Content: reader1, ContentType: "pdf"},
					{FileID: "42", Filename: "filename42", Content: reader2, ContentType: "jpg"},
				},
			},
			expectedEmail: &domain.Email{
				Id:      "2",
				Email:   "test@mail.com",
				From:    "from@test.mail",
				Text:    "bodytext",
				Subject: "subj",
				Attachments: []*domain.FileContent{
					{FileID: "13", Filename: "filename13", Content: reader1, ContentType: "pdf"},
					{FileID: "42", Filename: "filename42", Content: reader2, ContentType: "jpg"},
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.storage.On("CreateEmail", mock.Anything, mock.MatchedBy(func(email *domain.Email) bool {
				s.Equal(tt.expectedEmail, email)
				return true
			})).Return(nil)

			s.storage.On("UpdateEmail", mock.Anything, mock.MatchedBy(func(email *domain.Email) bool {
				tt.expectedEmail.SendStatus = domain.EmailRqStatusSent
				tt.expectedEmail.ErrorDescription = ""
				s.Equal(tt.expectedEmail, email)
				return true
			})).Return(nil)

			emailClient := newEmailClient()

			mockSendEmail := func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
				wantAuth := smtp.CRAMMD5Auth(cfg.SmtpUser, cfg.SmtpPassword)

				s.Equal(fmt.Sprintf("%s:%s", cfg.SmtpServer, cfg.SmtpServerPort), addr)
				s.Equal(wantAuth, auth)
				s.Equal([]string{tt.requestEmail.Email}, to)
				s.Equal(tt.requestEmail.From, from)

				return nil
			}

			smtpClient := mailClient{sendMail: mockSendEmail}

			err := emailClient.Init(cfg, s.storage, smtpClient)
			s.NoError(err)

			err = emailClient.send(s.Ctx, tt.requestEmail, newMessageBuilder().build(tt.requestEmail))
			s.NoError(err)
		})
	}
}

func (s *smtpClientTestSuite) Test_Send_WithErrorCreate() {
	errCreate := errors.New("errCreate")
	s.storage.On("CreateEmail", s.Ctx, mock.Anything).Return(errCreate)

	emailClient := newEmailClient()
	err := emailClient.Init(cfg, s.storage, mailClient{})
	s.NoError(err)

	email := &domain.Email{}
	err = emailClient.send(s.Ctx, email, newMessageBuilder().build(email))
	s.ErrorIs(err, errCreate)
}

func (s *smtpClientTestSuite) Test_Send_WithErrorUpdate() {
	errUpdate := errors.New("errUpdate")
	s.storage.On("CreateEmail", s.Ctx, mock.Anything).Return(nil)
	s.storage.On("UpdateEmail", s.Ctx, mock.Anything).Return(errUpdate)

	emailClient := newEmailClient()
	err := emailClient.Init(cfg, s.storage, mailClient{
		sendMail: func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error { return nil },
	})
	s.NoError(err)

	email := &domain.Email{}
	err = emailClient.send(s.Ctx, email, newMessageBuilder().build(email))
	s.ErrorIs(err, errUpdate)
}

func (s *smtpClientTestSuite) Test_Send_WithErrorSend() {
	errSend := errors.New("errSend")

	requestEmail := &domain.Email{
		Email: "test@mail.com",
	}
	expectedEmail := &domain.Email{
		Email:            "test@mail.com",
		SendStatus:       domain.EmailRqStatusSmtpError,
		ErrorDescription: errSend.Error(),
	}

	s.storage.On("CreateEmail", s.Ctx, mock.Anything).Return(nil)
	s.storage.On("UpdateEmail", s.Ctx, mock.MatchedBy(func(email *domain.Email) bool {
		s.Equal(expectedEmail.Email, email.Email)
		return true
	})).Return(nil)
	emailClient := newEmailClient()
	err := emailClient.Init(cfg, s.storage, mailClient{
		sendMail: func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error { return errSend },
	})
	s.NoError(err)

	err = emailClient.send(s.Ctx, requestEmail, newMessageBuilder().build(requestEmail))
	s.NoError(err)
}
