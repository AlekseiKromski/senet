package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (v *V1) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
