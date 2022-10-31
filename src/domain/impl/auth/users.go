package auth

import (
	"context"
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/errors/auth"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
)

type userSvcImpl struct {
	storage domain.UserStorage
}

func NewUserService(storage domain.UserStorage) domain.UserService {
	return &userSvcImpl{
		storage: storage,
	}
}

func (u *userSvcImpl) l() log.CLogger {
	return service.L().Cmp("users-domain-svc")
}

func (u *userSvcImpl) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	l := u.l().C(ctx).Mth("create").Trc()

	// check email
	if user.Username == "" {
		return nil, errors.ErrUserEmailEmpty(ctx)
	}
	if !kit.IsEmailValid(user.Username) {
		return nil, errors.ErrUserNoValidEmail(ctx)
	}

	user.Id = kit.NewId()
	user.LockedAt = nil

	// check username uniqueness
	another, err := u.storage.GetByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}
	if another != nil {
		return nil, errors.ErrUserNameNotUnique(ctx, user.Username)
	}

	// save to storage
	err = u.storage.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	l.F(log.FF{"id": user.Id}).Dbg("created")

	return user, nil
}

func (u *userSvcImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u.l().C(ctx).Mth("get-by-username").F(log.FF{"email": email}).Trc()
	if email == "" {
		return nil, errors.ErrUserEmailEmpty(ctx)
	}
	return u.storage.GetByEmail(ctx, email)
}

func (u *userSvcImpl) Get(ctx context.Context, userId string) (*domain.User, error) {
	u.l().C(ctx).Mth("get").F(log.FF{"userId": userId}).Trc()
	if userId == "" {
		return nil, errors.ErrUserIdEmpty(ctx)
	}
	return u.storage.GetUser(ctx, userId)
}

func (u *userSvcImpl) GetByIds(ctx context.Context, userIds []string) ([]*domain.User, error) {
	u.l().C(ctx).Mth("get-ids").Trc()
	return u.storage.GetUserByIds(ctx, userIds)
}

func (u *userSvcImpl) SetPassword(ctx context.Context, userId, newPasswordHash string) error {
	l := u.l().C(ctx).Mth("reset-password").F(log.FF{"userId": userId}).Trc()

	if userId == "" {
		return errors.ErrUserIdEmpty(ctx)
	}

	// find user
	user, err := u.storage.GetUser(ctx, userId)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound(ctx, userId)
	}

	//check user status
	if user.ActivatedAt == nil {
		return errors.ErrUserNotActive(ctx, userId)
	}
	if user.LockedAt != nil {
		return errors.ErrUserLocked(ctx, userId)
	}

	user.Password = newPasswordHash

	if err := u.storage.UpdateUser(ctx, user); err != nil {
		return err
	}
	l.Trc("updated")
	return nil
}

func (u *userSvcImpl) SetActivationToken(ctx context.Context, userId, token string, ttl uint32) error {
	u.l().C(ctx).Mth("set-activation-token").F(log.FF{"userId": userId}).Trc()
	if userId == "" {
		return errors.ErrUserIdEmpty(ctx)
	}
	return u.storage.SetActivationToken(ctx, userId, token, ttl)
}

func (u *userSvcImpl) ActivateByToken(ctx context.Context, userId, token string) (*domain.User, error) {
	u.l().C(ctx).Mth("activate-by-token").F(log.FF{"userId": userId}).Trc()
	if userId == "" {
		return nil, errors.ErrUserIdEmpty(ctx)
	}
	if token == "" {
		return nil, errors.ErrUserActivationTokenEmpty(ctx, userId)
	}
	// get stored token by userId
	storedToken, err := u.storage.GetActivationToken(ctx, userId)
	if err != nil {
		return nil, err
	}
	if storedToken == "" || token != storedToken {
		return nil, errors.ErrUserActivationNotExistedOnInvalidToken(ctx, userId)
	}
	// get user
	user, err := u.storage.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user.ActivatedAt != nil {
		return nil, errors.ErrUserActivationInvalidOperation(ctx, userId)
	}
	if user.LockedAt != nil {
		return nil, errors.ErrUserActivationInvalidOperation(ctx, userId)
	}
	// set activation date
	now := kit.Now()
	user.ActivatedAt = &now
	// update storage
	err = u.storage.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
