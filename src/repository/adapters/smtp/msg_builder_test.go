package smtp

import (
	"bytes"
	"github.com/africarealty/server/src/domain"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type msgBuilderTestSuite struct {
	kitTestSuite.Suite
}

// SetupSuite is called once for a suite
func (s *msgBuilderTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
}

func TestMsgBuilderSuite(t *testing.T) {
	suite.Run(t, new(msgBuilderTestSuite))
}

func (s *msgBuilderTestSuite) Test_MsgBuilder_Build() {
	content1 := bytes.NewReader([]byte("filecontent13"))
	content2 := bytes.NewReader([]byte("filecontent13"))
	content3 := bytes.NewReader([]byte("filecontent42"))
	tests := []struct {
		name         string
		requestEmail *domain.Email
		expectedMsg  string
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
			expectedMsg: "MIME-version: 1.0\r\n" +
				"From: from@test.mail\r\n" +
				"To: test@mail.com\r\n" +
				"Subject: subj\r\n" +
				"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
				"bodytext\r\n\r\n\r\n",
		},
		{
			name: "with only one attachment",
			requestEmail: &domain.Email{
				Id:      "2",
				Email:   "test@mail.com",
				From:    "from@test.mail",
				Text:    "bodytext",
				Subject: "subj",
				Attachments: []*domain.FileContent{
					{FileID: "13", Filename: "filename13", Content: content1, ContentType: "pdf"},
				},
			},
			expectedMsg: "MIME-version: 1.0\r\n" +
				"From: from@test.mail\r\n" +
				"To: test@mail.com\r\n" +
				"Subject: subj\r\n" +
				"Content-Type: multipart/mixed; boundary=\"%%boundary%%\"\r\n\r\n" +
				"--%%boundary%%\r\n" +
				"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
				"Content-Transfer-Encoding: 7bit\r\n\r\n" +
				"bodytext\r\n\r\n\r\n" +
				"--%%boundary%%\r\n" +
				"Content-Type: pdf; name=\"filename13\"\r\n" +
				"Content-Transfer-Encoding: base64\r\n" +
				"Content-Disposition: attachment; filename=\"filename13\"\r\n\r\n" +
				"ZmlsZWNvbnRlbnQxMw==\r\n\r\n" +
				"--%%boundary%%--",
		},
		{
			name: "with more than one attachments",
			requestEmail: &domain.Email{
				Id:      "3",
				Email:   "test@mail.com",
				From:    "from@test.mail",
				Text:    "bodytext",
				Subject: "subj",
				Attachments: []*domain.FileContent{
					{FileID: "13", Filename: "filename13", Content: content2, ContentType: "pdf"},
					{FileID: "42", Filename: "filename42", Content: content3, ContentType: "jpg"},
				},
			},
			expectedMsg: "MIME-version: 1.0\r\n" +
				"From: from@test.mail\r\n" +
				"To: test@mail.com\r\n" +
				"Subject: subj\r\n" +
				"Content-Type: multipart/mixed; boundary=\"%%boundary%%\"\r\n\r\n" +
				"--%%boundary%%\r\n" +
				"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
				"Content-Transfer-Encoding: 7bit\r\n\r\n" +
				"bodytext\r\n\r\n\r\n" +
				"--%%boundary%%\r\n" +
				"Content-Type: pdf; name=\"filename13\"\r\n" +
				"Content-Transfer-Encoding: base64\r\n" +
				"Content-Disposition: attachment; filename=\"filename13\"\r\n\r\n" +
				"ZmlsZWNvbnRlbnQxMw==\r\n\r\n" +
				"--%%boundary%%\r\n" +
				"Content-Type: jpg; name=\"filename42\"\r\n" +
				"Content-Transfer-Encoding: base64\r\n" +
				"Content-Disposition: attachment; filename=\"filename42\"\r\n\r\n" +
				"ZmlsZWNvbnRlbnQ0Mg==\r\n\r\n" +
				"--%%boundary%%--",
		},
	}

	msgBuilder := newMessageBuilder()

	for _, tt := range tests {
		s.Run(tt.name, func() {
			msg := msgBuilder.build(tt.requestEmail)
			expectedMsg := strings.Replace(tt.expectedMsg, "%%boundary%%", msgBuilder.boundary, -1)
			s.Equal([]byte(expectedMsg), msg)
		})
	}
}
