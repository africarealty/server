package auth

import (
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
		// authentication & authorization
		http.R("/api/auth/login", r.ctrl.Login).POST().NoAuth(),
		http.R("/api/auth/token/refresh", r.ctrl.TokenRefresh).POST().NoAuth(),
		http.R("/api/auth/logout", r.ctrl.Logout).POST(),
		http.R("/api/auth/registration", r.ctrl.Registration).POST().NoAuth(),
		http.R("/api/auth/password", r.ctrl.SetPassword).POST(),
	)
}
