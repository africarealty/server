package system

import (
	kitHttp "github.com/africarealty/server/src/kit/http"
	"github.com/africarealty/server/src/service"
	"net/http"
)

type Controller interface {
	// Ready returns OK if service is ready
	Ready(http.ResponseWriter, *http.Request)
}

type controllerIml struct {
	kitHttp.BaseController
}

func NewController() Controller {
	return &controllerIml{
		BaseController: kitHttp.BaseController{
			Logger: service.LF(),
		},
	}
}

// Ready godoc
// @Summary check system is ready
// @Router /ready [get]
// @Success 200
// @tags system
func (c *controllerIml) Ready(w http.ResponseWriter, r *http.Request) {
	c.RespondWithStatus(w, http.StatusOK, "OK")
}
