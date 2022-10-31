package auth

import (
	"context"
	"fmt"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/auth"
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
	"github.com/africarealty/server/src/usecase"
)

type userRegistrationImpl struct {
	userService  domain.UserService
	emailService domain.EmailService
	pwdService   auth.PasswordService
	cfg          *service.Config
}

func NewUserRegistrationImpl(userService domain.UserService, emailService domain.EmailService, pwdService auth.PasswordService, cfg *service.Config) usecase.UserRegistrationUseCase {
	return &userRegistrationImpl{
		userService:  userService,
		emailService: emailService,
		pwdService:   pwdService,
		cfg:          cfg,
	}
}

func (u *userRegistrationImpl) l() log.CLogger {
	return service.L().Cmp("user-reg-uc")
}

func (u *userRegistrationImpl) Register(ctx context.Context, rq *usecase.UserRegistrationRq) (*domain.User, error) {
	l := u.l().C(ctx).Mth("reg").Trc()

	// validate input
	if rq == nil {
		return nil, errors.ErrUserRegEmptyRq(ctx)
	}

	// check password policy
	check, err := u.pwdService.CheckPolicy(ctx, rq.Password, &auth.PasswordPolicy{MinLen: &u.cfg.Auth.Password.MinLen})
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, errors.ErrUserRegPasswordTooSimple(ctx)
	}

	// get pwd hash
	pwdHash, err := u.pwdService.GetHash(ctx, rq.Password)
	if err != nil {
		return nil, err
	}

	// create new user
	user := &domain.User{
		User: auth.User{
			Username:  rq.Email,
			Password:  pwdHash,
			Type:      rq.UserType,
			LastName:  rq.LastName,
			FirstName: rq.FirstName,
		},
	}
	user, err = u.userService.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	l.F(log.FF{"userId": user.Id}).Dbg("created")

	// generate activation token
	activationToken := auth.GenerateToken()

	// set activation token with the configured ttl
	err = u.userService.SetActivationToken(ctx, user.Id, activationToken, u.cfg.Auth.Activation.Ttl)
	if err != nil {
		return nil, err
	}

	// send activation email
	templateRq := &domain.TemplateRequest{
		Id: domain.EmailTemplateUserActivation,
		Data: map[string]interface{}{
			"Name":             fmt.Sprintf("%s %s", rq.FirstName, rq.LastName),
			"RegistrationLink": fmt.Sprintf("%s?userId=%s&token=%s", u.cfg.Auth.Activation.Url, user.Id, activationToken),
		},
	}
	_, err = u.emailService.Send(ctx, &domain.EmailRequest{
		UserId:   user.Id,
		Email:    user.Username,
		Template: templateRq,
		From:     u.cfg.Communications.Email.SmtpFrom,
	})
	if err != nil {
		return nil, err
	}
	l.Dbg("email sent")

	return user, nil
}
