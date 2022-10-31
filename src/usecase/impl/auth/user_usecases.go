package auth

import (
	"context"
	"fmt"
	"github.com/africarealty/server/src/domain"
	errors "github.com/africarealty/server/src/errors/auth"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
	"github.com/africarealty/server/src/usecase"
)

type userUseCasesImpl struct {
	userService  domain.UserService
	emailService domain.EmailService
	pwdService   auth.PasswordService
	cfg          *service.Config
}

func NewUserUseCases(userService domain.UserService, emailService domain.EmailService, pwdService auth.PasswordService, cfg *service.Config) usecase.UserUseCases {
	return &userUseCasesImpl{
		userService:  userService,
		emailService: emailService,
		pwdService:   pwdService,
		cfg:          cfg,
	}
}

func (u *userUseCasesImpl) l() log.CLogger {
	return service.L().Cmp("user-reg-uc")
}

// checkPassword checks password and if it's ok builds a hash
func (u *userUseCasesImpl) checkPassword(ctx context.Context, password, confirmation string) (string, error) {
	// check password and confirmation
	if password == "" {
		return "", errors.ErrUserRegPasswordNotSpecified(ctx)
	}
	if password != confirmation {
		return "", errors.ErrUserRegPasswordConfirmationNotEqual(ctx)
	}
	// check password policy
	check, err := u.pwdService.CheckPolicy(ctx, password, &auth.PasswordPolicy{MinLen: &u.cfg.Auth.Password.MinLen})
	if err != nil {
		return "", err
	}
	if !check {
		return "", errors.ErrUserRegPasswordTooSimple(ctx)
	}
	// get pwd hash
	pwdHash, err := u.pwdService.GetHash(ctx, password)
	if err != nil {
		return "", err
	}
	return pwdHash, nil
}

func (u *userUseCasesImpl) Register(ctx context.Context, rq *usecase.UserRegistrationRq) (*domain.User, error) {
	l := u.l().C(ctx).Mth("registration").Trc()

	// validate input
	if rq == nil {
		return nil, errors.ErrUserRegEmptyRq(ctx)
	}

	// check password and get hash
	pwdHash, err := u.checkPassword(ctx, rq.Password, rq.PasswordConfirmation)
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

func (u *userUseCasesImpl) CreateActiveUser(ctx context.Context, rq *usecase.UserRegistrationRq) (*domain.User, error) {
	l := u.l().C(ctx).Mth("create-active").Trc()

	// validate input
	if rq == nil {
		return nil, errors.ErrUserRegEmptyRq(ctx)
	}

	// check password and get hash
	pwdHash, err := u.checkPassword(ctx, rq.Password, rq.PasswordConfirmation)
	if err != nil {
		return nil, err
	}

	now := kit.Now()

	// create new user
	user := &domain.User{
		User: auth.User{
			Username:    rq.Email,
			Password:    pwdHash,
			Type:        rq.UserType,
			LastName:    rq.LastName,
			FirstName:   rq.FirstName,
			ActivatedAt: &now,
		},
	}
	user, err = u.userService.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	l.F(log.FF{"userId": user.Id}).Dbg("created")

	return user, nil
}
