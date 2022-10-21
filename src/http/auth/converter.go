package auth

import (
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit/auth"
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

func (c *controllerIml) toClientRegRequestDomain(rq *ClientRegistrationRequest) *auth.User {
	return &auth.User{
		Username:  rq.Email,
		Password:  rq.Password,
		Type:      domain.UserTypeClient,
		FirstName: rq.FirstName,
		LastName:  rq.LastName,
		Groups:    []string{domain.AuthGroupClient},
	}
}

func (c *controllerIml) toClientUserApi(usr *auth.User) *ClientUser {
	return &ClientUser{
		Id:        usr.Id,
		Email:     usr.Username,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
	}
}
