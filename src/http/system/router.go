package system

import (
	"github.com/africarealty/server/src/kit/http"
	_ "github.com/africarealty/server/src/swagger"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	routeBuilder *http.RouteBuilder
	ctrl         Controller
}

func NewRouter(c Controller, routeBuilder *http.RouteBuilder) http.RouteSetter {
	return &Router{
		routeBuilder: routeBuilder,
		ctrl:         c,
	}
}

func (r *Router) Set() error {
	return r.routeBuilder.Build(
		// readiness
		http.R("/ready", r.ctrl.Ready).GET().NoAuth(),
		// swagger
		http.R("", nil).PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler),
	)
}
