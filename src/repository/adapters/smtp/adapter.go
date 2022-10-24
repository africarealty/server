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

type InitRq struct {
	EmailStorage domain.EmailStorage
	Cfg          *service.CfgEmail
}

type Adapter interface {
	kitService.Adapter
	domain.EmailRepository
}

type adapterImpl struct {
	serviceImpl *emailImpl
	cfg         *service.CfgEmail
	msgBuilder  *messageBuilder
}

func (a *adapterImpl) l() log.CLogger {
	return service.L().Cmp("smtp-adapter")
}

func NewAdapter() Adapter {
	return &adapterImpl{
		serviceImpl: newEmailClient(),
		msgBuilder:  newMessageBuilder(),
	}
}

func (a *adapterImpl) Init(ctx context.Context, cfg interface{}) error {
	l := a.l().Mth("init")

	rq, ok := cfg.(*InitRq)
	if !ok {
		return errors.ErrEmailSmtpInvalidRequest(ctx)
	}
	a.cfg = rq.Cfg

	smtpClient := mailClient{sendMail: func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
		err := smtp.SendMail(addr, auth, from, to, msg)

		return err
	}}

	if err := a.serviceImpl.Init(a.cfg, rq.EmailStorage, smtpClient); err != nil {
		return err
	}
	l.Dbg("smtp-client initialized")

	return nil
}

func (a *adapterImpl) Send(ctx context.Context, email *domain.Email) error {
	return a.serviceImpl.send(ctx, email, a.msgBuilder.build(email))
}

func (a *adapterImpl) Close(ctx context.Context) error {
	return nil
}
