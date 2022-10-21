package auth

import (
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/errors/auth"
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/kit/context"
	kitHttp "github.com/africarealty/server/src/kit/http"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
	"net/http"
	"strings"
)

type Controller interface {

	// Ready returns OK if service is ready
	Ready(http.ResponseWriter, *http.Request)

	// auth
	Login(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	Registration(http.ResponseWriter, *http.Request)
	TokenRefresh(http.ResponseWriter, *http.Request)
	SetPassword(http.ResponseWriter, *http.Request)
}

type controllerIml struct {
	kitHttp.BaseController
	userService    domain.UserService
	sessionService auth.SessionsService
}

func NewController(sessionService auth.SessionsService, userService domain.UserService) Controller {
	return &controllerIml{
		BaseController: kitHttp.BaseController{
			Logger: service.LF(),
		},
		sessionService: sessionService,
		userService:    userService,
	}
}

func (c *controllerIml) l() log.CLogger {
	return service.L().Cmp("auth-controller")
}

// Ready godoc
// @Summary check system is ready
// @Router /ready [get]
// @Success 200
// @tags system
func (c *controllerIml) Ready(w http.ResponseWriter, r *http.Request) {
	c.RespondWithStatus(w, http.StatusOK, "OK")
}

// Registration godoc
// @Summary registers a new client
// @Accept json
// @produce json
// @Param regRequest body ClientRegistrationRequest true "registration request"
// @Success 200 {object} ClientUser
// @Failure 500 {object} http.Error
// @Router /auth/registration [post]
// @tags auth
func (c *controllerIml) Registration(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	c.l().C(ctx).Mth("registration").Trc()

	rq := &ClientRegistrationRequest{}
	if err := c.DecodeRequest(r, ctx, rq); err != nil {
		c.RespondError(w, err)
		return
	}

	user, err := c.userService.Create(ctx, c.toClientRegRequestDomain(rq))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toClientUserApi(user))
}

// Login godoc
// @Summary logins user by email/password
// @Accept json
// @produce json
// @Param loginRequest body LoginRequest true "auth request"
// @Success 200 {object} LoginResponse
// @Failure 500 {object} http.Error
// @Router /auth/login [post]
// @tags auth
func (c *controllerIml) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.l().C(ctx).Mth("login").Trc()

	rq := &LoginRequest{}
	if err := c.DecodeRequest(r, ctx, rq); err != nil {
		c.RespondError(w, err)
		return
	}

	session, token, err := c.sessionService.LoginPassword(ctx, c.toLoginRequest(rq))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toLoginResponseApi(session, token))
}

// Logout godoc
// @Summary logouts user
// @Accept json
// @Produce json
// @Router /auth/logout [post]
// @Success 200
// @Failure 400 {object} http.Error
// @Failure 500 {object} http.Error
// @tags auth
func (c *controllerIml) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.l().C(ctx).Mth("logout").Trc()

	// take userId from the currently logged session
	rq, err := context.MustRequest(ctx)
	if err != nil || rq.Sid == "" {
		c.RespondError(w, errors.ErrLogoutNoSID(ctx))
		return
	}

	err = c.sessionService.Logout(ctx, rq.Sid)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// TokenRefresh godoc
// @Summary refreshes auth token
// @Accept json
// @Produce json
// @Router /auth/token/refresh [post]
// @Success 200 {object} SessionToken
// @Failure 400 {object} http.Error
// @Failure 500 {object} http.Error
// @tags auth
func (c *controllerIml) TokenRefresh(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	c.l().C(ctx).Mth("token-refresh").Trc()

	_, err := context.MustRequest(ctx)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		c.RespondError(w, errors.ErrNoAuthHeader(ctx))
		return
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) < 2 {
		c.RespondError(w, errors.ErrAuthHeaderInvalid(ctx))
		return
	}
	token := splitToken[1]

	sessionToken, err := c.sessionService.RefreshToken(ctx, token)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toTokenApi(sessionToken))
}

// SetPassword godoc
// @Summary sets a new password for the user
// @Param request body SetPasswordRequest true "set password request"
// @Accept json
// @Produce json
// @Router /auth/password [post]
// @Success 200
// @Failure 400 {object} http.Error
// @Failure 500 {object} http.Error
// @tags auth
func (c *controllerIml) SetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.l().C(ctx).Mth("set-password").Trc()

	// take userId from the currently logged session
	rqCtx, err := context.MustRequest(ctx)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if rqCtx.Uid == "" {
		c.RespondError(w, errors.ErrNoUID(ctx))
		return
	}

	rq := &SetPasswordRequest{}
	if err := c.DecodeRequest(r, ctx, rq); err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.userService.SetPassword(ctx, rqCtx.Uid, rq.PrevPassword, rq.NewPassword)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}
