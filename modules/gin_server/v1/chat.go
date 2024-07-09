package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	PRIVATE            = "PRIVATE"
	SERVER_PRIVATE_KEY = "SERVER_PRIVATE_KEY"
)

type CreateChatRequest struct {
	User1         string `json:"user1"`
	User2         string `json:"user2"`
	ChatType      string `json:"chat_type"`
	SecurityLevel string `json:"security_level"`
}

func (v *V1) CreateChat(c *gin.Context) {
	createChatRequest := CreateChatRequest{}

	if err := c.BindJSON(&createChatRequest); err != nil {
		v.log("cannot unmarshal request payload", err.Error())
		c.JSON(http.StatusBadRequest, NewErrorResponse("incorrect format of payload"))
		return
	}

	if createChatRequest.ChatType == PRIVATE {
		if createChatRequest.SecurityLevel != SERVER_PRIVATE_KEY {
			message := fmt.Sprintf("security level (%s) does not applicable for (%s) chat", createChatRequest.SecurityLevel, createChatRequest.ChatType)
			v.log(message)
			c.JSON(http.StatusBadRequest, NewErrorResponse(message))
			return
		}

		chatId, err := v.storage.IsChatBetweenUsersExists(createChatRequest.User1, createChatRequest.User2)
		if err != nil {
			v.log("cannot find chat between users", err.Error())
			c.JSON(http.StatusBadRequest, NewErrorResponse("incorrect users id"))
			return
		}
		if len(chatId) != 0 {
			v.log("chat already exists id:", chatId)
			c.JSON(http.StatusBadRequest, NewErrorResponse("chat already exists"))
			return
		}

		user1, err := v.storage.GetUserById(createChatRequest.User1)
		if err != nil {
			v.log("cannot get user by id", createChatRequest.User1, err.Error())
			c.JSON(http.StatusBadRequest, NewErrorResponse("incorrect user1 id"))
			return
		}

		user2, err := v.storage.GetUserById(createChatRequest.User2)
		if err != nil {
			v.log("cannot get user by id", createChatRequest.User2, err.Error())
			c.JSON(http.StatusBadRequest, NewErrorResponse("incorrect user2 id"))
			return
		}

		chat, err := v.storage.CreateChat(
			fmt.Sprintf("%s <-> %s", user1.Username, user2.Username),
			createChatRequest.ChatType,
			createChatRequest.SecurityLevel,
		)
		if err != nil {
			v.log("cannot create chat", err.Error())
			c.JSON(http.StatusBadRequest, NewErrorResponse("cannot create chat"))
			return
		}

		if err := v.storage.AddUserToChat(user1.Id, chat.Id); err != nil {
			v.log("cannot add user1 by id (to chat)", user1.Id, err.Error())
			c.JSON(http.StatusBadRequest, NewErrorResponse("cannot add user1 to chat"))
			return
		}

		if err := v.storage.AddUserToChat(user2.Id, chat.Id); err != nil {
			v.log("cannot add user2 by id (to chat)", user2.Id, err.Error())
			c.JSON(http.StatusBadRequest, NewErrorResponse("cannot add user2 to chat"))
			return
		}

		c.JSON(http.StatusOK, chat)
		return
	}

	c.JSON(http.StatusBadRequest, NewErrorResponse("current chat type does not exists"))
}

func (v *V1) GetAllChats(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusOK, []string{})
		return
	}

	currentUserId, ok := uid.(string)
	if !ok {
		v.log("wrong user id")
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResponse("wrong user id"))
		return
	}

	chats, err := v.storage.GetChats(currentUserId)
	if err != nil {
		v.log("cannot get chats", err.Error())
		c.JSON(http.StatusBadRequest, NewErrorResponse("cannot get chats"))
		return
	}

	c.JSON(http.StatusOK, chats)
	return
}
