package ws

import (
	v1 "alekseikromski.com/senet/modules/gin_server/v1"
	server_key_storage "alekseikromski.com/senet/modules/server-key-storage"
	"alekseikromski.com/senet/modules/storage"
	"encoding/json"
	"github.com/AlekseiKromski/ws-gin-upgrader/core"
	"github.com/go-playground/validator/v10"
)

const SENT_MESSAGE = "SENT_MESSAGE"
const INCOMING_MESSAGE = "INCOMING_MESSAGE"
const STATE_MESSAGE_OK = "STATE_MESSAGE_OK"
const STATE_MESSAGE_FAIL = "STATE_MESSAGE_FAIL"

type SentMessageHandler struct {
	serverKeyStorage server_key_storage.ServerKeyStorage
	log              func(m ...string)
	store            storage.Storage
}

type data struct {
	CID     string `json:"cid" validate:"required"`
	To      string `json:"to" validate:"required"`
	Message string `json:"message" validate:"required"`
}

func (sm *SentMessageHandler) Handle(payload string, session *core.Session, clients core.Clients) {
	validate := validator.New()

	d := data{}
	if err := json.Unmarshal([]byte(payload), &d); err != nil {
		sm.log("cannot decode payload in ws event", err.Error())
		sm.sendResponse(STATE_MESSAGE_FAIL, "incorrect payload format", session)
		return
	}

	err := validate.Struct(d)
	if err != nil {
		sm.log("data is not valid", err.Error())
		sm.sendResponse(STATE_MESSAGE_FAIL, "incorrect payload format", session)
		return
	}

	chatId, err := sm.store.IsChatBetweenUsersExists(session.ID, d.To)
	if err != nil {
		sm.log("cannot check chat between two users", err.Error())
		sm.sendResponse(STATE_MESSAGE_FAIL, "server error", session)
		return
	}
	if len(chatId) == 0 {
		sm.log("no chat between users")
		sm.sendResponse(STATE_MESSAGE_FAIL, "no chat between users", session)
		return
	}

	chat, err := sm.store.GetChat(chatId)
	if err != nil {
		sm.log("cannot get chat", err.Error())
		sm.sendResponse(STATE_MESSAGE_FAIL, "server error", session)
		return
	}

	if chat.SecurityLevel == v1.SERVER_PRIVATE_KEY {
		encryptedMessage, err := sm.serverKeyStorage.Encrypt(d.Message)
		if err != nil {
			sm.log("error during encryption", err.Error())
			sm.sendResponse(STATE_MESSAGE_FAIL, "server encryption errors", session)
			return
		}
		message, err := sm.store.CreateMessage(d.CID, session.ID, encryptedMessage)
		if err != nil {
			sm.log("cannot save message", err.Error())
			sm.sendResponse(STATE_MESSAGE_FAIL, "server error", session)
			return
		}

		//encode message to json
		message.Message = d.Message
		incomingMessagePayload, err := json.Marshal(message)
		if err != nil {
			sm.log("cannot encode message", err.Error())
			sm.sendResponse(STATE_MESSAGE_FAIL, "server error", session)
			return
		}

		if err := clients.Send(d.To, string(incomingMessagePayload), INCOMING_MESSAGE); err != nil {
			sm.log("cannot sent message to user: ", d.To, err.Error())
			sm.sendResponse(STATE_MESSAGE_FAIL, "cannot sent message to use", session)
			return
		}

		sm.sendResponse(STATE_MESSAGE_OK, string(incomingMessagePayload), session)
	} else {
		sm.log("unsupported security level: ", err.Error())
		sm.sendResponse(STATE_MESSAGE_FAIL, "unsupported security level", session)
		return
	}
}

func (sm *SentMessageHandler) sendResponse(state core.HandlerName, message string, session *core.Session) {
	if err := session.Send(message, state); err != nil {
		sm.log("cannot sent state message to: ", session.ID, err.Error())
		return
	}
}
