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

	rs, err := s.POST(ctx, fmt.Sprintf("%s/api/auth/registration", s.cfg.Url), "", rqJs)
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
