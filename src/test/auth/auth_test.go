// go:build integration

package auth

import (
	"fmt"
	"github.com/africarealty/server/src/http/auth"
	"github.com/africarealty/server/src/kit"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/sdk"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type authTestSuite struct {
	kitTestSuite.Suite
	cfg *service.CfgSdk
}

func (s *authTestSuite) SetupSuite() {
	s.Init(service.LF())

	os.Setenv("ARROOT", "/home/mikhailb/work/africarealty/dev")

	// load config
	cfg, err := service.LoadConfig()
	if err != nil {
		s.Fatal(err)
	}
	s.cfg = cfg.Sdk
}

func (s *authTestSuite) SetupTest() {}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(authTestSuite))
}

func (s *authTestSuite) TearDownSuite(t *testing.T) {
}

func (s *authTestSuite) Test_Register_Ok() {

	ownerSdk := sdk.New(s.cfg)
	defer ownerSdk.Close(s.Ctx)

	user, err := ownerSdk.RegisterUser(s.Ctx, &auth.RegistrationRequest{
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
