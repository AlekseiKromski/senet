package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"senet/processor/errors"
)

func (api *Api) Signup(c *gin.Context) {

	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot read request body: %v", err),
		))
		return
	}

	user := &register{}
	if err := json.Unmarshal(body, user); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot unmarshal incoming data: %v", err),
		))
		return
	}

	vr := user.validate()
	if !vr.Result {
		c.JSON(http.StatusBadRequest, vr)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot generate password: %v", err),
		))
		return
	}
	if err := api.storage.CreateUser(uuid.New(), user.Username, string(password)); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot create register: %v", err),
		))
		return
	}

	c.Status(http.StatusOK)
}
