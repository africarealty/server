package impl

import (
	"context"
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/kit/log"
	"golang.org/x/crypto/bcrypt"
	"unicode/utf8"
)

type passwordServiceImpl struct {
	logger log.CLoggerFunc
}

func NewPasswordService(logger log.CLoggerFunc) auth.PasswordService {
	return &passwordServiceImpl{
		logger: logger,
	}
}

func (s *passwordServiceImpl) l() log.CLogger {
	return s.logger().Cmp("password-svc")
}

func (s *passwordServiceImpl) GetHash(ctx context.Context, password string) (string, error) {
	s.l().C(ctx).Mth("get-hash").Trc()
	// check password empty
	if password == "" {
		return "", auth.ErrAuthPwdEmpty(ctx)
	}
	// hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", auth.ErrAuthPwdHashGenerate(err, ctx)
	}
	return string(bytes), nil
}

func (s *passwordServiceImpl) CheckPolicy(ctx context.Context, password string, policy *auth.PasswordPolicy) (bool, error) {
	s.l().C(ctx).Mth("check-policy").Trc()
	// check password empty
	if password == "" {
		return false, auth.ErrAuthPwdEmpty(ctx)
	}
	if policy == nil {
		return true, nil
	}
	// check password policy
	if policy.MinLen != nil && utf8.RuneCountInString(password) < int(*policy.MinLen) {
		return false, auth.ErrAuthPwdPolicy(ctx)
	}
	return true, nil
}

func (s *passwordServiceImpl) Verify(ctx context.Context, password, hash string) (bool, error) {
	s.l().C(ctx).Mth("verify").Trc()
	// check password policy
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
