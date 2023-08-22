package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"senet/processor/lb"
)

type Api struct {
	lb     *lb.LoadBalancer
	pfs    *ProcessorFs
	engine *gin.Engine
}

func NewApi(lb *lb.LoadBalancer, pfs *ProcessorFs, engine *gin.Engine) *Api {
	return &Api{
		lb:     lb,
		pfs:    pfs,
		engine: engine,
	}
}

func (api *Api) Register() {
	apiGroup := api.engine.Group("/api/")
	{
		apiGroup.GET("/healthz", api.Healthz)
		apiGroup.GET("/users", api.Users)
	}

	api.engine.GET("/", Webclient(api.pfs.Content))
}

func (api *Api) RegisterStaticFiles() {
	api.engine.StaticFS("/static", http.FS(api.pfs))
}
