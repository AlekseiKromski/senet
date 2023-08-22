package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (api *Api) Users(c *gin.Context) {
	users, err := api.lb.GetUsers()
	if err != nil {
		fmt.Printf("lb error: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, users)
}
