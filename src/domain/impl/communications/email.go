package impl

import (
	"context"
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/errors/communications"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/kit/queue"
	"github.com/africarealty/server/src/kit/queue/listener"
	kitService "github.com/africarealty/server/src/kit/service"
	"github.com/africarealty/server/src/service"
)

type emailImpl struct {
	kitService.BaseService
	templateGenerator domain.TemplateGenerator
	store             domain.StoreService
	emailRepository   domain.EmailRepository
	cfg               *service.Config
}

func NewEmailService(
	templateGenerator domain.TemplateGenerator,
	queue queue.Queue,
	store domain.StoreService,
	emailRepository domain.EmailRepository,
) domain.EmailService {
	return &emailImpl{
		templateGenerator: templateGenerator,
		BaseService:       kitService.BaseService{Queue: queue},
		store:             store,
		emailRepository:   emailRepository,
	}
}

func (e *emailImpl) Init(cfg *service.Config) error {
	e.cfg = cfg
	return nil
}

func (e *emailImpl) l() log.CLogger {
	return service.L().Cmp("email-svc")
}

func (e *emailImpl) validateRequest(ctx context.Context, rq *domain.Email) error {
	// if email is not correct
	if !kit.IsEmailValid(rq.Email) {
		return errors.ErrEmailValidationInvalidEmail(ctx, rq.Email)
	}
	if !kit.IsEmailValid(rq.From) {
		return errors.ErrEmailValidationInvalidFrom(ctx, rq.From)
	}
	// check template
	if rq.Template == nil || rq.Template.Id == "" {
		return errors.ErrTemplateEmpty(ctx)
	}
	return nil
}

func (e *emailImpl) Send(ctx context.Context, rq *domain.EmailRequest) (*domain.Email, error) {
	email := domain.Email{
		Id:          kit.NewId(),
		UserId:      rq.UserId,
		Email:       rq.Email,
		Template:    rq.Template,
		SendStatus:  domain.EmailRqStatusNotSend,
		LinkFileIds: rq.LinkFileIds,
		From:        rq.From,
	}

	l := e.l().Mth("send").C(ctx).F(log.FF{"emailId": email.Id}).Dbg().Trc(kit.ToJson(rq))

	if email.From == "" {
		email.From = e.cfg.Communications.Email.SmtpFrom
	}

	// validate incoming request
	if err := e.validateRequest(ctx, &email); err != nil {
		return nil, err
	}

	//apply template
	rs, err := e.templateGenerator.Generate(ctx, email.Template)
	if err != nil {
		return nil, err
	}
	email.Text = rs.Body
	email.Subject = rs.Title

	// publish to queue to be sent
	if err := e.Publish(ctx, email, queue.QueueTypeAtLeastOnce, domain.EmailRequestTopic); err != nil {
		l.F(log.FF{"topic": domain.EmailRequestTopic}).E(err).St().Err("publish")
	}

	l.Dbg("published to queue")

	return &email, nil
}

func (e *emailImpl) RequestHandler() listener.QueueMessageHandler {
	return func(payload []byte) error {
		l := e.l().Mth("email-send-request-handler")

		var emailRq *domain.EmailRequest
		ctx, err := queue.Decode(context.Background(), payload, &emailRq)
		if err != nil {
			return err
		}
		l.C(ctx).Dbg()

		var email = &domain.Email{
			UserId:   emailRq.UserId,
			Email:    emailRq.Email,
			Template: emailRq.Template,
			From:     emailRq.From,
		}

		for _, fileId := range emailRq.LinkFileIds {
			file, err := e.store.GetFile(ctx, fileId)
			if err != nil {
				return err
			}
			email.Attachments = append(email.Attachments, file)
		}
		return e.emailRepository.Send(ctx, email)
	}
}
