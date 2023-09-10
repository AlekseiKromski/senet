package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"senet/processor/errors"
	"senet/processor/storage/models"
	"time"
)

func (api *Api) Login(c *gin.Context) {
	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot read request body: %v", err),
		))
		return
	}

	login := &login{}
	if err := json.Unmarshal(body, login); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot unmarshal incoming data: %v", err),
		))
		return
	}

	vr := login.validate()
	if !vr.Result {
		c.JSON(http.StatusBadRequest, vr)
		return
	}

	users, err := api.storage.GetUser(login.Username, false)
	if err != nil || len(users) == 0 {
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot find user: %v", err),
		))
		return
	}

	user := users[0]

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(login.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, errors.NewApiErrorMessage(
			fmt.Errorf("invalid credits: %v", err),
		))
		return
	}

	//remove password
	user.Password = nil

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["username"] = user.Username

	tokenString, err := token.SignedString(api.jwtSecret)
	if err != nil {
		log.Fatalf("cannot create token: %v", err)
	}

	c.JSON(http.StatusOK, struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
	}{
		Token: tokenString,
		User:  user,
	})
}
