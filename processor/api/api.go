package api

import (
	"github.com/AlekseiKromski/at-socket-server/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"senet/processor/lb"
	"senet/processor/storage"
)

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

	apiGroup := api.engine.Group("/api/").Use(middleware.JwtCheck(api.jwtSecret))
	{
		apiGroup.GET("/healthz", api.Healthz)
		apiGroup.GET("/user/search/:search", api.GetUsers)
		apiGroup.POST("/chat/create", api.CreateChat)
		apiGroup.GET("/chat/get/all-chats", api.GetChats)
	}

	api.engine.GET("/", Webclient(api.pfs.Content))
}

func (api *Api) RegisterStaticFiles() {
	api.engine.StaticFS("/static", http.FS(api.pfs))
}
