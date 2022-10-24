//+build integration

package communications

import (
	"bytes"
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type emailStorageTestSuite struct {
	kitTestSuite.Suite
	storage domain.EmailStorage
	adapter Adapter
}

// SetupSuite is called once for a suite
func (s *emailStorageTestSuite) SetupSuite() {
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
	s.storage = s.adapter.GetEmailStorage()
}

func (s *emailStorageTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestEmailSuite(t *testing.T) {
	suite.Run(t, new(emailStorageTestSuite))
}

func (s *emailStorageTestSuite) Test_CreateUpdateEmail_Success() {
	fileId := kit.NewRandString()
	email := &domain.Email{
		Id:          kit.NewId(),
		Subject:     "Тестовое письмо",
		Text:        "Текст тестового письма",
		Email:       "test@test.mail",
		From:        "from@test.mail",
		Template:    &domain.TemplateRequest{Id: "test"},
		SendStatus:  domain.EmailRqStatusNotSend,
		LinkFileIds: []string{fileId},
		Attachments: []*domain.FileContent{
			{
				Filename:    "filename",
				FileID:      fileId,
				ContentType: "text/plain",
				Extension:   "txt",
				Content:     bytes.NewReader([]byte("content")),
			},
		},
	}

	err := s.storage.CreateEmail(s.Ctx, email)
	s.NoError(err)

	email.SendStatus = domain.EmailRqStatusSent

	err = s.storage.UpdateEmail(s.Ctx, email)
	s.NoError(err)
}
