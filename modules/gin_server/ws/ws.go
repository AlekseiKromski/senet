package ws

import (
	"alekseikromski.com/senet/modules/gin_server/guard"
	server_key_storage "alekseikromski.com/senet/modules/server-key-storage"
	"alekseikromski.com/senet/modules/storage"
	"fmt"
	"github.com/AlekseiKromski/ws-gin-upgrader/core"
	"github.com/gin-gonic/gin"
	"log"
)

type WebSocket struct {
	app *core.App
}

func NewWebSocket(engine *gin.Engine, secret []byte, guard *guard.Guard, store storage.Storage, serverKeyStorage server_key_storage.ServerKeyStorage, log func(m ...string)) (*WebSocket, error) {
	app, err := core.Start(
		engine,
		&core.Handlers{
			SENT_MESSAGE: &SentMessageHandler{
				serverKeyStorage: serverKeyStorage,
				log:              log,
				store:            store,
			},
		},
		guard.Check,
		&core.Config{
			JwtSecret: secret,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot start websocket server: %v", err)
	}

	go func() {
		for {
			event := <-app.Hooks
			log(event.Data)
		}
	}()

	return &WebSocket{
		app: app,
	}, nil
}

func (ws *WebSocket) SendDatapointsToAllClients(data string) error {
	for _, sessions := range ws.app.Clients.Storage {
		for _, s := range sessions {
			log.Println(s.Send(data, "INFO"))
		}
	}
	return nil
}
