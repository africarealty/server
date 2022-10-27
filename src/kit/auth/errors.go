package auth

import (
	"context"
	"github.com/africarealty/server/src/kit/er"
	"net/http"
)

const (
	ErrCodeAccessTokenCreation           = "KITAUTH-001"
	ErrCodeSessionPasswordValidation     = "KITAUTH-002"
	ErrCodeUserNotFound                  = "KITAUTH-003"
	ErrCodeUserNotActive                 = "KITAUTH-004"
	ErrCodeUserLocked                    = "KITAUTH-005"
	ErrCodeSessionLoggedOut              = "KITAUTH-006"
	ErrCodeSessionAuthWrongSigningMethod = "KITAUTH-007"
	ErrCodeSessionAuthTokenExpired       = "KITAUTH-008"
	ErrCodeSessionAuthTokenInvalid       = "KITAUTH-009"
	ErrCodeSessionAuthTokenClaimsInvalid = "KITAUTH-010"
	ErrCodeSessionTokenInvalid           = "KITAUTH-011"
	ErrCodeSessionNoSessionFound         = "KITAUTH-012"
	ErrCodeAuthPwdHashGenerate           = "KITAUTH-013"
	ErrCodeAuthPwdEmpty                  = "KITAUTH-014"
	ErrCodeAuthPwdPolicy                 = "KITAUTH-015"
)

var (
	ErrAccessTokenCreation = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeAccessTokenCreation, "").C(ctx).Err()
	}
	ErrSessionPasswordValidation = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionPasswordValidation, "invalid password").Business().HttpSt(http.StatusUnauthorized).C(ctx).Err()
	}
	ErrUserNotFound = func(ctx context.Context, userId string) error {
		return er.WithBuilder(ErrCodeUserNotFound, "user not found").Business().F(er.FF{"userId": userId}).C(ctx).Err()
	}
	ErrUserNotActive = func(ctx context.Context, userId string) error {
		return er.WithBuilder(ErrCodeUserNotActive, "user not active").Business().F(er.FF{"userId": userId}).C(ctx).Err()
	}
	ErrUserLocked = func(ctx context.Context, userId string) error {
		return er.WithBuilder(ErrCodeUserLocked, "user locked").Business().F(er.FF{"userId": userId}).C(ctx).Err()
	}
	ErrSessionLoggedOut = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionLoggedOut, "session is logged out").C(ctx).Business().HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSessionAuthWrongSigningMethod = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionAuthWrongSigningMethod, "wrong signing method").C(ctx).HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSessionAuthTokenExpired = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionAuthTokenExpired, "token expired").C(ctx).Business().HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSessionAuthTokenInvalid = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionAuthTokenInvalid, "invalid token").C(ctx).Business().HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSessionAuthTokenClaimsInvalid = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionAuthTokenClaimsInvalid, "invalid token claims").Business().HttpSt(http.StatusUnauthorized).C(ctx).Err()
	}
	ErrSessionTokenInvalid = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionTokenInvalid, "session token is invalid").C(ctx).Business().HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSessionNoSessionFound = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionNoSessionFound, "no session found").C(ctx).Business().HttpSt(http.StatusUnauthorized).Err()
	}
	ErrAuthPwdHashGenerate = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeAuthPwdHashGenerate, "").C(ctx).Err()
	}
	ErrAuthPwdEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeAuthPwdEmpty, "password isn't specified").Business().C(ctx).Err()
	}
	ErrAuthPwdPolicy = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeAuthPwdPolicy, "password is too simple").Business().C(ctx).Err()
	}
)
