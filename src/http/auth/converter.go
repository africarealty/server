package auth

import (
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/usecase"
)

func (c *controllerIml) toLoginRequest(rq *LoginRequest) *auth.LoginRequest {
	if rq == nil {
		return nil
	}
	return &auth.LoginRequest{
		Username: rq.Email,
		Password: rq.Password,
	}
}

func (c *controllerIml) toTokenApi(t *auth.SessionToken) *SessionToken {
	if t == nil {
		return nil
	}
	return &SessionToken{
		SessionId:             t.SessionId,
		AccessToken:           t.AccessToken,
		AccessTokenExpiresAt:  t.AccessTokenExpiresAt,
		RefreshToken:          t.RefreshToken,
		RefreshTokenExpiresAt: t.RefreshTokenExpiresAt,
	}
}

func (c *controllerIml) toLoginResponseApi(s *auth.Session, t *auth.SessionToken) *LoginResponse {
	return &LoginResponse{
		UserId: s.UserId,
		Token:  c.toTokenApi(t),
	}
}

func (c *controllerIml) toRegRequestDomain(rq *RegistrationRequest) *usecase.UserRegistrationRq {
	return &usecase.UserRegistrationRq{
		Email:                rq.Email,
		Password:             rq.Password,
		UserType:             rq.UserType,
		FirstName:            rq.FirstName,
		LastName:             rq.LastName,
		PasswordConfirmation: rq.Confirmation,
	}
}

func (c *controllerIml) toOwnerApi(p *domain.OwnerProfile) *OwnerProfile {
	if p == nil {
		return nil
	}
	return &OwnerProfile{
		Avatar: p.Avatar,
	}
}

func (c *controllerIml) toAgentApi(p *domain.AgentProfile) *AgentProfile {
	if p == nil {
		return nil
	}
	return &AgentProfile{
		Avatar: p.Avatar,
	}
}

func (c *controllerIml) toUserApi(usr *domain.User) *User {
	return &User{
		Id:          usr.Id,
		Email:       usr.Username,
		FirstName:   usr.FirstName,
		LastName:    usr.LastName,
		ActivatedAt: usr.ActivatedAt,
		LockedAt:    usr.LockedAt,
		Owner:       c.toOwnerApi(usr.Owner),
		Agent:       c.toAgentApi(usr.Agent),
	}
}
