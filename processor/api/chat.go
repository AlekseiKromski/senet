package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"senet/processor/errors"
	"senet/processor/storage/creators/chat"
)

func (api *Api) CreateChat(c *gin.Context) {
	requestorUid, exists := c.Get("uid")
	if !exists {
		c.Status(http.StatusForbidden)
		return
	}

	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot read request body: %v", err),
		))
		return
	}

	cc := &createChat{}
	if err := json.Unmarshal(body, cc); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot unmarshal incoming data: %v", err),
		))
		return
	}

	if vr := cc.validate(); !vr.Result {
		c.JSON(http.StatusBadRequest, vr)
		return
	}

	cID, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot create random chat id: %v", err),
		))
		return
	}

	var creator chat.ChatCreator
	switch cc.Type {
	case "private":
		creator = chat.NewPrivateChatCreator(cID.String(), requestorUid.(string), cc.Users[0])
	default:
		c.JSON(http.StatusBadRequest, errors.NewApiErrorMessage(
			fmt.Errorf("unsupported chat type: %s", cc.Type),
		))
		return
	}

	ch, err := api.lb.CreateChat(creator)
	if err != nil {
		log.Printf("Error during creation chat: %v", err)
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot create chat"),
		))
		return
	}

	c.JSON(http.StatusOK, ch)
}

func (api *Api) GetChats(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.Status(http.StatusForbidden)
		return
	}

	chats, err := api.lb.GetChats(uid.(string))
	if err != nil {
		log.Printf("Error during getting chat/chat messages: %v", err)
		c.JSON(http.StatusInternalServerError, errors.NewApiErrorMessage(
			fmt.Errorf("cannot get chats"),
		))
		return
	}

	c.JSON(http.StatusOK, chats)
}
