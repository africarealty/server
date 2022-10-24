package smtp

import (
	"context"
	"fmt"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/communications"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
	"net/smtp"
)

type mailClient struct {
	sendMail func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error
}

type emailImpl struct {
	cfg     *service.CfgEmail
	addr    string
	storage domain.EmailStorage
	client  mailClient
}

func newEmailClient() *emailImpl {
	return &emailImpl{}
}

func (e *emailImpl) Init(cfg *service.CfgEmail, storage domain.EmailStorage, client mailClient) error {
	e.cfg = cfg
	e.storage = storage
	e.client = client
	e.addr = fmt.Sprintf("%s:%s", e.cfg.SmtpServer, e.cfg.SmtpServerPort)
	return nil
}

func (e *emailImpl) l() log.CLogger {
	return service.L().Cmp("smtp")
}

func (e *emailImpl) send(ctx context.Context, email *domain.Email, body []byte) error {
	l := e.l().C(ctx).Mth("send").F(log.FF{"emailId": email.Id}).Dbg("sending")

	if err := e.storage.CreateEmail(ctx, email); err != nil {
		return err
	}

	email.SendStatus = domain.EmailRqStatusSent
	email.ErrorDescription = ""

	auth := smtp.CRAMMD5Auth(e.cfg.SmtpUser, e.cfg.SmtpPassword)
	err := e.client.sendMail(e.addr, auth, email.From, []string{email.Email}, body)
	if err != nil {
		// log error
		l.E(errors.ErrEmailSmtpSend(err, ctx)).St().Err()

		email.SendStatus = domain.EmailRqStatusSmtpError
		email.ErrorDescription = err.Error()
	}

	return e.storage.UpdateEmail(ctx, email)
}
