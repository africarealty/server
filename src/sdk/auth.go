package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/africarealty/server/src/http/auth"
	"github.com/africarealty/server/src/kit/log"
)

func (s *Sdk) RegisterUser(ctx context.Context, rq *auth.RegistrationRequest) (*auth.User, error) {
	l := s.l().C(ctx).Mth("reg-user").F(log.FF{"email": rq.Email}).Trc()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}

	rs, err := s.POST(ctx, fmt.Sprintf("%s/api/auth/users/registration", s.cfg.Url), "", rqJs)
	if err != nil {
		return nil, err
	}

	var user *auth.User
	err = json.Unmarshal(rs, &user)
	if err != nil {
		return nil, err
	}
	l.F(log.FF{"userId": user.Id}).Dbg("ok")

	return user, nil
}

func (s *Sdk) CreateActiveUser(ctx context.Context, rq *auth.RegistrationRequest) (*auth.User, error) {
	l := s.l().C(ctx).Mth("active-user").F(log.FF{"email": rq.Email}).Trc()

	rqJs, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}

	rs, err := s.POST(ctx, fmt.Sprintf("%s/api/auth/users", s.cfg.Url), s.accessToken, rqJs)
	if err != nil {
		return nil, err
	}

	var user *auth.User
	err = json.Unmarshal(rs, &user)
	if err != nil {
		return nil, err
	}
	l.F(log.FF{"userId": user.Id}).Dbg("ok")

	return user, nil
}

func (s *Sdk) LoginUser(ctx context.Context, email, pwd string) (*auth.SessionToken, error) {
	l := s.l().C(ctx).Mth("login").F(log.FF{"email": email}).Dbg()

	// if logged in already, logged out
	if s.accessToken != "" {
		_ = s.Logout(ctx)
	}

	rq := &auth.LoginRequest{
		Email:    email,
		Password: pwd,
	}

	rqJ, _ := json.Marshal(rq)
	r, err := s.POST(ctx, fmt.Sprintf("%s/api/auth/users/login", s.cfg.Url), "", rqJ)
	if err != nil {
		return nil, err
	}

	var rs *auth.LoginResponse
	err = json.Unmarshal(r, &rs)
	if err != nil {
		return nil, err
	}

	l.Dbg("logged in")

	s.accessToken = rs.Token.AccessToken

	return rs.Token, nil
}

func (s *Sdk) Logout(ctx context.Context) error {
	s.l().C(ctx).Mth("logout").Dbg()
	_, err := s.POST(ctx, fmt.Sprintf("%s/api/auth/users/logout", s.cfg.Url), s.accessToken, []byte{})
	if err != nil {
		return err
	}
	s.accessToken = ""
	return err
}
