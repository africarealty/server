//+build integration

package smtp

import (
	"bytes"
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/repository/storage/communications"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type smtpTestSuite struct {
	kitTestSuite.Suite
	storageAdapter communications.Adapter
	adapter        Adapter
}

// SetupSuite is called once for a suite
func (s *smtpTestSuite) SetupSuite() {
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
	s.storageAdapter = communications.NewAdapter()
	if err := s.storageAdapter.Init(s.Ctx, cfg); err != nil {
		s.Fatal(err)
	}

	if err := s.adapter.Init(s.Ctx, &InitRq{
		EmailStorage: s.storageAdapter.GetEmailStorage(),
		Cfg:          cfg.Communications.Email,
	}); err != nil {
		s.Fatal(err)
	}
}

func (s *smtpTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
	_ = s.storageAdapter.Close(s.Ctx)
}

func TestSmtpSuite(t *testing.T) {
	suite.Run(t, new(smtpTestSuite))
}

func (s *smtpTestSuite) Test_SendEmail_Success() {
	fileId := kit.NewRandString()
	reader := bytes.NewReader([]byte("content"))
	email := &domain.Email{
		Id:          kit.NewId(),
		Subject:     "Тест отправки email",
		Text:        "Текст тестового письма",
		Email:       "test@test.mail",
		From:        "from@test.mail",
		Template:    &domain.TemplateRequest{Id: "test-template"},
		SendStatus:  domain.EmailRqStatusNotSend,
		LinkFileIds: []string{fileId},
		Attachments: []*domain.FileContent{
			{
				Filename:    "filename",
				FileID:      fileId,
				ContentType: "text/plain",
				Extension:   "txt",
				Content:     reader,
			},
		},
	}
	err := s.adapter.Send(s.Ctx, email)
	s.NoError(err)
}
