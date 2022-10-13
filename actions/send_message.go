package actions

import (
	"github.com/AlekseiKromski/at-socket-server/core"
)

type SendMessage struct {
	Data   string
	client *core.Client
}

func (sm *SendMessage) SetData(data string) {
	sm.Data = data
}

func (sm *SendMessage) Do() {
	sm.run()
}
func (sm *SendMessage) TrigType() string {
	return "to-all"
}
func (sm *SendMessage) SetClient(client *core.Client) {
	sm.client = client
}

func (sm *SendMessage) run() {
	println("ok")
}
