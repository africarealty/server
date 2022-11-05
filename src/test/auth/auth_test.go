//go:build integration

package auth

import (
	"fmt"
	"github.com/africarealty/server/src/http/auth"
	"github.com/africarealty/server/src/kit"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/service"
	"github.com/africarealty/server/src/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type authTestSuite struct {
	kitTestSuite.Suite
	cfg *service.Config
}

func (s *authTestSuite) SetupSuite() {
	s.Init(service.LF())

	// load config
	cfg, err := service.LoadConfig()
	if err != nil {
		s.Fatal(err)
	}
	s.cfg = cfg
}

func (s *authTestSuite) SetupTest() {}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(authTestSuite))
}

func (s *authTestSuite) TearDownSuite(t *testing.T) {
}

func (s *authTestSuite) Test_Register_Ok() {

	userSdk := test.NewSdk(s.cfg, &s.Suite)
	defer userSdk.Close(s.Ctx)

	user, err := userSdk.RegisterUser(s.Ctx, &auth.RegistrationRequest{
		Email:     fmt.Sprintf("%s@example.test", kit.NewRandString()),
		Password:  "123456",
		FirstName: "test",
		LastName:  "test",
		UserType:  "owner",
	})
	if err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(user)
	s.NotEmpty(user.Id)
}

func (s *authTestSuite) Test_CreateActiveUser_LoginAndLogout_Ok() {

	// login as admin
	sdk := test.NewSdk(s.cfg, &s.Suite)
	defer sdk.Close(s.Ctx)
	_ = sdk.LoginAdmin(s.Ctx)
	defer sdk.Logout(s.Ctx)

	// create an active user
	user, sess := sdk.CreateAndLogin(s.Ctx, fmt.Sprintf("%s@example.test", kit.NewRandString()), "123456")
	s.NotEmpty(user)
	s.NotEmpty(user.Id)
	s.NotEmpty(sess.SessionId)
	s.NotEmpty(sess.AccessToken)
	s.NotEmpty(sess.RefreshToken)

	// logout
	err := sdk.Sdk.Logout(s.Ctx)
	if err != nil {
		s.Fatal(err)
	}
}
