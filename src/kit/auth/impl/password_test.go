package impl

import (
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/kit/log"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/stretchr/testify/suite"
	"testing"
)

type passwordTestSuite struct {
	kitTestSuite.Suite
	pwdService auth.PasswordService
}

func (s *passwordTestSuite) SetupSuite() {
	logger := func() log.CLogger {
		return log.L(log.Init(&log.Config{Level: log.InfoLevel}))
	}
	s.Suite.Init(logger)
	s.pwdService = NewPasswordService(logger)
}

func (s *passwordTestSuite) SetupTest() {
}

func TestPasswordSuite(t *testing.T) {
	suite.Run(t, new(passwordTestSuite))
}

func (s *passwordTestSuite) Test_GetHash() {
	// empty pwd
	_, err := s.pwdService.GetHash(s.Ctx, "")
	s.AssertAppErr(err, auth.ErrCodeAuthPwdEmpty)
	// valid hash
	hash, err := s.pwdService.GetHash(s.Ctx, "123")
	s.Nil(err)
	s.NotEmpty(hash)
}

func (s *passwordTestSuite) Test_CheckPolicy() {
	// empty pwd
	_, err := s.pwdService.CheckPolicy(s.Ctx, "", nil)
	s.AssertAppErr(err, auth.ErrCodeAuthPwdEmpty)
	// empty policy
	res, err := s.pwdService.CheckPolicy(s.Ctx, "1234", nil)
	s.Nil(err)
	s.True(res)
	// min len
	var ln uint = 6
	res, err = s.pwdService.CheckPolicy(s.Ctx, "1234", &auth.PasswordPolicy{MinLen: &ln})
	s.AssertAppErr(err, auth.ErrCodeAuthPwdPolicy)
	s.False(res)
	// valid
	ln = 4
	res, err = s.pwdService.CheckPolicy(s.Ctx, "12342", &auth.PasswordPolicy{MinLen: &ln})
	s.Nil(err)
	s.True(res)
}

func (s *passwordTestSuite) Test_Verify() {
	hash, err := s.pwdService.GetHash(s.Ctx, "123")
	s.Nil(err)
	s.NotEmpty(hash)
	// invalid
	r, err := s.pwdService.Verify(s.Ctx, "1234", hash)
	s.Nil(err)
	s.False(r)
	// valid
	r, err = s.pwdService.Verify(s.Ctx, "123", hash)
	s.Nil(err)
	s.True(r)
}
