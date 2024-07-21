package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (v1 *V1) GetMessages(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusOK, []string{})
		return
	}

	currentUserId, ok := uid.(string)
	if !ok {
		v1.log("wrong user id")
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse("wrong user id"))
		return
	}

	chatid := c.Param("chatid")
	if len(chatid) == 0 {
		v1.log("request does not have chatid")
		c.JSON(http.StatusBadRequest, NewErrorResponse("request does not have chatid"))
		return
	}

	offsetParam := c.Param("offset")
	if len(chatid) == 0 {
		v1.log("request does not have offset")
		c.JSON(http.StatusBadRequest, NewErrorResponse("request does not have offset"))
		return
	}
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		v1.log("cannot convert offset from string to int")
		c.JSON(http.StatusBadRequest, NewErrorResponse("offset should be int"))
		return
	}

	chat, err := v1.storage.GetChat(chatid)
	if err != nil {
		v1.log("cannot get chat by id: ", chatid, err.Error())
		c.JSON(http.StatusBadRequest, NewErrorResponse("cannot get chat by id: "+chatid))
		return
	}

	accesable := false
	for _, u := range chat.Users {
		if u.Id != currentUserId {
			continue
		}
		accesable = true
	}

	if !accesable {
		c.JSON(http.StatusBadRequest, NewErrorResponse("you don't have access to chat: "+chatid))
		return
	}

	ms, err := v1.storage.GetMessagesByChatId(chatid, offset, 15)
	if err != nil {
		v1.log("cannot get messages by chatid: ", chatid, err.Error())
		c.JSON(http.StatusBadRequest, NewErrorResponse("cannot get messages by id: "+chatid))
		return
	}

	// TODO: make chat security level check
	for _, m := range ms {
		decoded, err := v1.serverKeyStorage.Decrypt(m.Message)
		if err != nil {
			v1.log("cannot decode message", chatid, err.Error())
			continue
		}
		m.Message = decoded
	}

	c.JSON(http.StatusOK, ms)
}
