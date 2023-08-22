package api

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

type user struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (api *Api) Signup(c *gin.Context) {

	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	user := &user{}
	if err := json.Unmarshal(body, user); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if user.Password == "" || user.PasswordConfirmation == "" || user.Username == "" {
		c.Status(http.StatusInternalServerError)
		return
	}

	if user.Password != user.PasswordConfirmation {
		c.Status(http.StatusInternalServerError)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if err := api.storage.CreateUser(uuid.New(), user.Username, string(password)); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
