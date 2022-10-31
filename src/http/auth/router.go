package auth

import (
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/kit/http"
	_ "github.com/africarealty/server/src/swagger"
)

type Router struct {
	ctrl         Controller
	routeBuilder *http.RouteBuilder
}

func NewRouter(c Controller, routeBuilder *http.RouteBuilder) http.RouteSetter {
	return &Router{
		ctrl:         c,
		routeBuilder: routeBuilder,
	}
}

func (r *Router) Set() error {
	return r.routeBuilder.Build(
		// authorization
		http.R("/api/auth/users/login", r.ctrl.Login).POST().NoAuth(),
		http.R("/api/auth/token/refresh", r.ctrl.TokenRefresh).POST().NoAuth(),
		http.R("/api/auth/users/logout", r.ctrl.Logout).POST(),
		http.R("/api/auth/users/registration", r.ctrl.Registration).POST().NoAuth(),
		http.R("/api/auth/users/activation", r.ctrl.Activation).POST().NoAuth(),
		http.R("/api/auth/users/password", r.ctrl.SetPassword).POST(),
		http.R("/api/auth/users", r.ctrl.CreateUser).POST().Authorize(auth.Resource(domain.AuthResUserProfileAll, "rw").B()),
	)
}
