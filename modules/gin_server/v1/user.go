package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (v *V1) GetUser(c *gin.Context) {
	username := c.Param("username")
	if len(username) == 0 {
		c.Status(204)
		return
	}

	users, err := v.storage.FindUsersByUsername(username)
	if err != nil {
		v.log("cannot find users", err.Error())
		c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("cannot get users by %s username", username)))
		return
	}

	c.JSON(200, users)
}
