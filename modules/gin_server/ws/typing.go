package ws

import (
	"alekseikromski.com/senet/modules/storage"
	"encoding/json"
	"github.com/AlekseiKromski/ws-gin-upgrader/core"
	"github.com/go-playground/validator/v10"
)

const TYPING = "TYPING"

type TypingHandler struct {
	log   func(m ...string)
	store storage.Storage
}

type dataTyping struct {
	CID string `json:"cid" validate:"required"`
	To  string `json:"to" validate:"required"`
}

func (tp *TypingHandler) Handle(payload string, session *core.Session, clients core.Clients) {
	validate := validator.New()

	d := dataTyping{}
	if err := json.Unmarshal([]byte(payload), &d); err != nil {
		tp.log("cannot decode payload in ws event", err.Error())
		return
	}

	err := validate.Struct(d)
	if err != nil {
		tp.log("data is not valid", err.Error())
		return
	}

	chatId, err := tp.store.IsChatBetweenUsersExists(session.ID, d.To)
	if err != nil {
		tp.log("cannot check chat between two users", err.Error())
		return
	}
	if len(chatId) == 0 {
		tp.log("no chat between users")
		return
	}

	if err := clients.Send(d.To, d.CID, TYPING); err != nil {
		tp.log("cannot sent message to user: ", d.To, err.Error())
		return
	}
}
