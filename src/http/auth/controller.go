package auth

import (
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/errors/auth"
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/kit/context"
	kitHttp "github.com/africarealty/server/src/kit/http"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
	"github.com/africarealty/server/src/usecase"
	"net/http"
	"strings"
)

type Controller interface {
	Login(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	Registration(http.ResponseWriter, *http.Request)
	Activation(http.ResponseWriter, *http.Request)
	TokenRefresh(http.ResponseWriter, *http.Request)
	SetPassword(http.ResponseWriter, *http.Request)
}

type controllerIml struct {
	kitHttp.BaseController
	userRegUc      usecase.UserRegistrationUseCase
	userService    domain.UserService
	sessionService auth.SessionsService
}

func NewController(sessionService auth.SessionsService, userService domain.UserService, userRegUc usecase.UserRegistrationUseCase) Controller {
	return &controllerIml{
		BaseController: kitHttp.BaseController{
			Logger: service.LF(),
		},
		sessionService: sessionService,
		userService:    userService,
		userRegUc:      userRegUc,
	}
}

func (c *controllerIml) l() log.CLogger {
	return service.L().Cmp("auth-controller")
}

// Registration godoc
// @Summary registers a new client
// @Accept json
// @produce json
// @Param regRequest body RegistrationRequest true "registration request"
// @Success 200 {object} User
// @Failure 500 {object} http.Error
// @Router /auth/registration [post]
// @tags auth
func (c *controllerIml) Registration(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	c.l().C(ctx).Mth("registration").Trc()

	rq := &RegistrationRequest{}
	if err := c.DecodeRequest(r, ctx, rq); err != nil {
		c.RespondError(w, err)
		return
	}

	user, err := c.userRegUc.Register(ctx, c.toRegRequestDomain(rq))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toUserApi(user))
}

// Activation godoc
// @Summary activates a user by token
// @produce json
// @Param userId query string true "user id"
// @Param token query string true "activation token"
// @Success 200 {object} User
// @Failure 500 {object} http.Error
// @Router /auth/activation [post]
// @tags auth
func (c *controllerIml) Activation(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	c.l().C(ctx).Mth("activation").Trc()

	userId, err := c.FormVal(r, ctx, "userId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	token, err := c.FormVal(r, ctx, "token", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	user, err := c.userService.ActivateByToken(ctx, userId, token)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toUserApi(user))
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
	//rqCtx, err := context.MustRequest(ctx)
	//if err != nil {
	//	c.RespondError(w, err)
	//	return
	//}
	//if rqCtx.Uid == "" {
	//	c.RespondError(w, errors.ErrNoUID(ctx))
	//	return
	//}
	//
	//rq := &SetPasswordRequest{}
	//if err := c.DecodeRequest(r, ctx, rq); err != nil {
	//	c.RespondError(w, err)
	//	return
	//}
	//
	//err = c.userService.SetPassword(ctx, rqCtx.Uid, rq.PrevPassword, rq.NewPassword)
	//if err != nil {
	//	c.RespondError(w, err)
	//	return
	//}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}
