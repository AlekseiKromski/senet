package v1

import (
	"alekseikromski.com/senet/modules/gin_server/guard"
	"alekseikromski.com/senet/modules/storage"
	"embed"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"path/filepath"
	"time"
)

type V1 struct {
	router  *gin.Engine
	storage storage.Storage
	log     func(messages ...string)
	secret  []byte
	guard   *guard.Guard
}

func NewV1Api(storage storage.Storage, secret []byte, cookieDomain string, log func(messages ...string)) *V1 {
	return &V1{
		router:  gin.Default(),
		storage: storage,
		log:     log,
		secret:  secret,
		guard:   guard.NewGuard(log, secret, storage, cookieDomain),
	}
}

func (v *V1) RegisterRoutes(resources embed.FS) error {
	v.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type, Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v.router.GET("/", v.application(resources))
	v.router.POST("/api/auth", v.guard.Auth)
	v.router.Static("/static", filepath.Join("front-end", "build", "static"))

	api := v.router.Group("/api").Use(v.guard.Check)
	{
		api.GET("/healthz", v.Healthz)
		api.GET("/auth/logout", v.guard.Logout)
	}

	return nil
}

func (v *V1) GetEngine() *gin.Engine {
	return v.router
}

func (v *V1) GetGuard() *guard.Guard {
	return v.guard
}

func (v *V1) application(resources embed.FS) func(c *gin.Context) {
	return func(c *gin.Context) {
		content, err := resources.ReadFile("front-end/build/index.html")
		if err != nil {
			log.Printf("cannot return index.html: %v", err)
			c.Status(500)
			return
		}

		c.Writer.Write(content)
	}
}
