package smtp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/africarealty/server/src/domain"
	"mime/multipart"
	"strings"
)

type messageBuilder struct {
	boundary string
}

func newMessageBuilder() *messageBuilder {
	return &messageBuilder{}
}

func (b *messageBuilder) build(email *domain.Email) []byte {
	buf := bytes.NewBuffer(nil)
	b.boundary = multipart.NewWriter(buf).Boundary()

	buf.WriteString("MIME-version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", email.Email))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))

	withAttachments := len(email.Attachments) > 0

	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", b.boundary))
		buf.WriteString(fmt.Sprintf("--%s\r\n", b.boundary))
		buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		buf.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	} else {
		buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
	}

	buf.WriteString(fmt.Sprintf("%s\r\n\r\n\r\n", email.Text))

	if withAttachments {
		for _, attachment := range email.Attachments {
			buf.WriteString(fmt.Sprintf("--%s\r\n", b.boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", attachment.ContentType, attachment.Filename))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", attachment.Filename))

			contentBuf := new(bytes.Buffer)
			_, _ = contentBuf.ReadFrom(attachment.Content)
			body := make([]byte, base64.StdEncoding.EncodedLen(contentBuf.Len()))
			base64.StdEncoding.Encode(body, contentBuf.Bytes())

			// to avoid error "552 message line is too long" coming from some smtp server
			// we must have lines with less than 998 symbols
			// so add new line after each 998 symbols
			b := strings.Builder{}
			b.Grow(len(body))
			for i, r := range []rune(string(body)) {
				b.WriteRune(r)
				if i > 0 && i%998 == 0 {
					b.WriteString("\r\n")
				}
			}
			buf.Write([]byte(b.String()))
			buf.WriteString("\r\n\r\n")
		}
		buf.WriteString(fmt.Sprintf("--%s--", b.boundary))
	}

	return buf.Bytes()
}
