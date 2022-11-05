package ads

import (
	"github.com/africarealty/server/src/domain"
	kitHttp "github.com/africarealty/server/src/kit/http"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
	"net/http"
)

type Controller interface {
	CreateAds(http.ResponseWriter, *http.Request)
	UpdateAds(http.ResponseWriter, *http.Request)
	GetAds(http.ResponseWriter, *http.Request)
	ActivateAds(http.ResponseWriter, *http.Request)
	CloseAds(http.ResponseWriter, *http.Request)
	SearchAds(http.ResponseWriter, *http.Request)

	GetAdsGuest(http.ResponseWriter, *http.Request)
	SearchAdsGuest(http.ResponseWriter, *http.Request)
}

type controllerIml struct {
	kitHttp.BaseController
	adsService domain.AdvertisementService
}

func NewController(adsService domain.AdvertisementService) Controller {
	return &controllerIml{
		BaseController: kitHttp.BaseController{
			Logger: service.LF(),
		},
		adsService: adsService,
	}
}

func (c *controllerIml) l() log.CLogger {
	return service.L().Cmp("ads-controller")
}

func (c *controllerIml) CreateAds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.l().C(ctx).Mth("create").Trc()
	//TODO implement me
	panic("implement me")
}

func (c *controllerIml) UpdateAds(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (c *controllerIml) GetAds(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (c *controllerIml) ActivateAds(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (c *controllerIml) CloseAds(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (c *controllerIml) SearchAds(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (c *controllerIml) GetAdsGuest(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (c *controllerIml) SearchAdsGuest(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}
