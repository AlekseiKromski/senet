package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"senet/processor/lb"
	"senet/processor/storage"
)

type Api struct {
	lb      *lb.LoadBalancer
	storage storage.Storage
	pfs     *ProcessorFs
	engine  *gin.Engine
}

func NewApi(lb *lb.LoadBalancer, pfs *ProcessorFs, engine *gin.Engine, store storage.Storage) *Api {
	return &Api{
		lb:      lb,
		pfs:     pfs,
		engine:  engine,
		storage: store,
	}
}

func (api *Api) Register() {
	apiGroup := api.engine.Group("/api/")
	{
		apiGroup.GET("/healthz", api.Healthz)
		apiGroup.POST("/register", api.Signup)
	}

	api.engine.GET("/", Webclient(api.pfs.Content))
}

func (api *Api) RegisterStaticFiles() {
	api.engine.StaticFS("/static", http.FS(api.pfs))
}
