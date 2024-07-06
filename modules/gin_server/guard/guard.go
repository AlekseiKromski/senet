package guard

import (
	"alekseikromski.com/senet/modules/storage"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
)

type errorMessage struct {
	Message string `json:"message"`
}

func newErrorMessage(message string) *errorMessage {
	return &errorMessage{
		Message: message,
	}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type Guard struct {
	secret       []byte
	permissions  map[string][]*storage.Endpoint
	store        storage.Storage
	cookieDomain string
	log          func(messages ...string)
}

func NewGuard(log func(messages ...string), secret []byte, store storage.Storage, cookieDomain string) *Guard {
	permissions, err := store.GetPermissions()
	if err != nil {
		log("cannot get permissions from database", err.Error())
	}

	return &Guard{
		secret:       secret,
		store:        store,
		permissions:  permissions,
		cookieDomain: cookieDomain,
		log:          log,
	}
}

func (g *Guard) Auth(c *gin.Context) {
	defer c.Request.Body.Close()

	ar := &authRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&ar); err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	user, err := g.store.GetUserByUsername(ar.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, newErrorMessage("user or password is not correct"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ar.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, newErrorMessage("user or password is not correct"))
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["authorized"] = true
	claims["id"] = user.Id

	tokenString, err := token.SignedString(g.secret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, newErrorMessage("cannot sign token"))
		return
	}

	if ar.Type == "cookie" {
		c.SetCookie("token", tokenString, 3600, "/", g.cookieDomain, true, true)
	} else {
		c.JSON(http.StatusOK, struct {
			Token string `json:"token"`
		}{
			Token: tokenString,
		})
	}

	return
}

func (g *Guard) Check(c *gin.Context) {
	req := c.Request
	tokenRequest := ""
	t, _ := c.Cookie("token")

	if t != "" {
		tokenRequest = t
	} else {
		tokenRequest = req.Header.Get("Authorization")
		if len(tokenRequest) != 0 && len(tokenRequest) > 10 {
			tokenRequest = tokenRequest[7:len(tokenRequest)]
		} else {
			g.log("there is not token in request", req.URL.String())
			c.AbortWithStatusJSON(http.StatusUnauthorized, newErrorMessage("cannot find token in request"))
			return
		}
	}

	userID := ""
	if tokenRequest == "" || len(tokenRequest) < 10 {
		g.log("there is not token in request", req.URL.String())
		c.AbortWithStatusJSON(http.StatusUnauthorized, newErrorMessage("cannot find token in request"))
		return
	}

	token, err := jwt.Parse(tokenRequest, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong sign method")
		}
		claims := token.Claims.(jwt.MapClaims)
		if claims["id"] == nil {
			return nil, fmt.Errorf("wrong format of JWT")
		}

		userid, ok := claims["id"].(string)
		if !ok {
			return nil, fmt.Errorf("wrong format of userid")
		}

		userID = userid

		return g.secret, nil
	})

	if err != nil {
		g.log("token verify failed", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, newErrorMessage("token verification failed"))

		return
	}

	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, newErrorMessage("token is not valid"))
		return
	}

	user, err := g.store.GetUserById(userID)
	if err != nil {
		g.log("cannot find user by id:", userID, err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, newErrorMessage("token payload error"))
		return
	}

	denied := true
	for _, path := range g.permissions[user.Role] {
		matched, err := regexp.MatchString(path.Urn, req.URL.Path)
		if err != nil {
			g.log("regexp match string error", err.Error())
			continue
		}

		if matched {
			denied = false
		}
	}

	if denied {
		g.log("access denied by permission restrictions role/resource", user.Role, req.URL.String())
		c.AbortWithStatusJSON(http.StatusForbidden, newErrorMessage("access denied by permission restrictions"))
		return
	}

	// Set for at-socket-server
	c.Set("uid", userID)
	c.Next()
}

func (g *Guard) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", g.cookieDomain, true, true)
	c.Status(http.StatusOK)
}
