package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"senet/processor/errors"
)

func (api *Api) GetUsers(c *gin.Context) {
	search := c.Param("search")

	if len(search) == 0 {
		c.JSON(http.StatusBadRequest, errors.NewApiErrorMessage(
			fmt.Errorf("not search string"),
		))
		return
	}

	users, err := api.storage.GetUser(search, true)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, errors.NewApiErrorMessage(
			fmt.Errorf("cannot find user by name"),
		))
		return
	}

	for _, user := range users {
		user.Password = nil
	}

	c.JSON(http.StatusOK, users)
}
