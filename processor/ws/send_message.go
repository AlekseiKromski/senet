package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/AlekseiKromski/at-socket-server/core"
)

const SEND_MESSAGE core.HandlerName = "SEND_MESSAGE"

type senderPayload struct {
	ReceiverID string `json:"receiverID"`
	Message    string `json:"message"`
}

type SenderHandler struct{}

func (sh *SenderHandler) Handle(payload string, client *core.Client, clients core.Clients) {
	sp := senderPayload{}

	if err := json.Unmarshal([]byte(payload), &sp); err != nil {
		if err := client.Conn.WriteJSON(core.ActionModel{
			Action:  core.HandlerName(core.ERR_DECODE),
			Payload: fmt.Sprintf("cannot decode payload: %v", err),
		}); err != nil {
			fmt.Printf("cannot send error back: %v", err)
		}
	}

	receiver := clients[sp.ReceiverID]
	if receiver == nil {
		if err := client.Conn.WriteJSON(core.ActionModel{
			Action:  core.ERR_HANDLER,
			Payload: fmt.Sprintf("cannot find client: %s", sp.ReceiverID),
		}); err != nil {
			fmt.Printf("cannot send error back: %v", err)
		}
		return
	}

	if err := receiver.Conn.WriteJSON(core.ActionModel{
		Action:  "NEW_MESSAGE",
		Payload: sp.Message,
	}); err != nil {
		if err := client.Conn.WriteJSON(core.ActionModel{
			Action:  core.ERR_DECODE,
			Payload: fmt.Sprintf("cannot decode payload: %v", err),
		}); err != nil {
			fmt.Printf("cannot send error back: %v", err)
		}
	}
}
