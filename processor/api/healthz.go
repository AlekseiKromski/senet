package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}
