package main

import (
	"fmt"
	"github.com/AlekseiKromski/at-socket-server/core"
	"senet/actions"
)

var actionHandlers = []*core.ActionHandler{
	{
		ActionType: "send-message",
		Action:     &actions.SendMessage{},
	},
}

var triggerHandlers = []*core.TriggerHandler{
	{},
}

func main() {
	_, err := core.Start(actionHandlers, triggerHandlers)
	if err != nil {
		fmt.Println(err)
	}
	println("ok")
}
