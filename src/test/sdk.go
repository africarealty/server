package test

import (
	"context"
	"github.com/africarealty/server/src/http/auth"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/sdk"
	"github.com/africarealty/server/src/service"
)

type TestSdk struct {
	sdk.Sdk
	cfg   *service.CfgTests
	suite *kitTestSuite.Suite
}

func NewSdk(cfg *service.Config, suite *kitTestSuite.Suite) *TestSdk {
	return &TestSdk{
		Sdk:   *sdk.New(cfg.Sdk),
		cfg:   cfg.Tests,
		suite: suite,
	}
}

func (s *TestSdk) LoginAdmin(ctx context.Context) *auth.SessionToken {
	sess, err := s.Sdk.LoginUser(ctx, s.cfg.User, s.cfg.Password)
	if err != nil {
		s.suite.Fatal(err)
	}
	s.suite.NotEmpty(sess)
	s.suite.NotEmpty(sess.SessionId)
	s.suite.NotEmpty(sess.AccessToken)
	return sess
}

func (s *TestSdk) Logout(ctx context.Context) {
	_ = s.Sdk.Logout(ctx)
}

func (s *TestSdk) CreateAndLogin(ctx context.Context, email, pwd string) (*auth.User, *auth.SessionToken) {
	// login admin
	s.LoginAdmin(ctx)
	// create user
	user, err := s.Sdk.CreateActiveUser(ctx, &auth.RegistrationRequest{
		Email:    email,
		Password: pwd,
		UserType: "owner",
	})
	if err != nil {
		s.suite.Fatal(err)
	}
	s.suite.NotEmpty(user)
	// login new user
	sess, err := s.Sdk.LoginUser(ctx, email, pwd)
	if err != nil {
		s.suite.Fatal(err)
	}
	s.suite.NotEmpty(sess)
	s.suite.NotEmpty(sess.SessionId)
	s.suite.NotEmpty(sess.AccessToken)
	return user, sess
}
