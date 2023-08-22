package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"senet/processor/lb"
	"senet/processor/storage"
)

//TODO: added JWT check middleware -> https://github.com/AlekseiKromski/alekseikromski-blog/blob/master/api/guard/jwt/jwt.go

type Api struct {
	lb        *lb.LoadBalancer
	storage   storage.Storage
	pfs       *ProcessorFs
	engine    *gin.Engine
	jwtSecret []byte
}

func NewApi(lb *lb.LoadBalancer, pfs *ProcessorFs, engine *gin.Engine, store storage.Storage, jwtSecret []byte) *Api {
	return &Api{
		lb:        lb,
		pfs:       pfs,
		engine:    engine,
		storage:   store,
		jwtSecret: jwtSecret,
	}
}

func (api *Api) Register() {

	authGroup := api.engine.Group("/auth/")
	{
		authGroup.POST("/register", api.Signup)
		authGroup.POST("/login", api.Login)
	}

	apiGroup := api.engine.Group("/api/").Use(api.JwtCheck)
	{
		apiGroup.GET("/healthz", api.Healthz)

	}

	api.engine.GET("/", Webclient(api.pfs.Content))
}

func (api *Api) RegisterStaticFiles() {
	api.engine.StaticFS("/static", http.FS(api.pfs))
}
