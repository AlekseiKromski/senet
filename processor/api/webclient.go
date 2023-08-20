package api

import (
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Webclient(frontend embed.FS) func(c *gin.Context) {
	return func(c *gin.Context) {
		file, err := frontend.ReadFile("webclient/build/index.html")
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		if _, err = c.Writer.Write(file); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}
}
