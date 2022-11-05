package ads

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
		http.R("/api/ads", r.ctrl.CreateAds).POST().Authorize(auth.Resource(domain.AuthResUserProfileAll, "rw").B()),
		http.R("/api/ads/{adsId}", r.ctrl.UpdateAds).PUT().Authorize(auth.Resource(domain.AuthResUserProfileAll, "rw").B()),
		http.R("/api/ads/{adsId}", r.ctrl.GetAds).GET().Authorize(auth.Resource(domain.AuthResUserProfileAll, "rw").B()),
		http.R("/api/ads/{adsId}/activate", r.ctrl.ActivateAds).POST().Authorize(auth.Resource(domain.AuthResUserProfileAll, "rw").B()),
		http.R("/api/ads/{adsId}/close", r.ctrl.CloseAds).POST().Authorize(auth.Resource(domain.AuthResUserProfileAll, "rw").B()),
		http.R("/api/ads/search/query", r.ctrl.SearchAds).GET().Authorize(auth.Resource(domain.AuthResUserProfileAll, "rw").B()),

		// anonymous access
		http.R("/api/guest/ads/{adsId}", r.ctrl.GetAdsGuest).GET().Authorize(auth.Resource(domain.AuthResUserProfileAll, "rw").B()),
		http.R("/api/guest/ads/search/query", r.ctrl.SearchAdsGuest).GET().Authorize(auth.Resource(domain.AuthResUserProfileAll, "rw").B()),
	)
}
