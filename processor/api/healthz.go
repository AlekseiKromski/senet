package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (api *Api) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}
