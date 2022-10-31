package smtp

import (
	"context"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/communications"
	"github.com/africarealty/server/src/kit/log"
	kitService "github.com/africarealty/server/src/kit/service"
	"github.com/africarealty/server/src/service"
	"net/smtp"
)

type Adapter interface {
	kitService.Adapter
	domain.EmailRepository
}

type adapterImpl struct {
	serviceImpl  *emailImpl
	cfg          *service.Config
	msgBuilder   *messageBuilder
	emailStorage domain.EmailStorage
}

func (a *adapterImpl) l() log.CLogger {
	return service.L().Cmp("smtp-adapter")
}

func NewAdapter(emailStorage domain.EmailStorage) Adapter {
	return &adapterImpl{
		serviceImpl:  newEmailClient(),
		msgBuilder:   newMessageBuilder(),
		emailStorage: emailStorage,
	}
}

func (a *adapterImpl) Init(ctx context.Context, cfg interface{}) error {
	a.l().Mth("init").Trc()

	var ok bool
	a.cfg, ok = cfg.(*service.Config)
	if !ok {
		return errors.ErrEmailSmtpInvalidConfig(ctx)
	}

	smtpClient := mailClient{sendMail: func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
		err := smtp.SendMail(addr, auth, from, to, msg)

		return err
	}}

	if err := a.serviceImpl.Init(a.cfg.Communications.Email, a.emailStorage, smtpClient); err != nil {
		return err
	}

	return nil
}

func (a *adapterImpl) Send(ctx context.Context, email *domain.Email) error {
	return a.serviceImpl.send(ctx, email, a.msgBuilder.build(email))
}

func (a *adapterImpl) Close(ctx context.Context) error {
	return nil
}
