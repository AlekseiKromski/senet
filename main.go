package main

import (
	"fmt"
	"github.com/AlekseiKromski/at-socket-server/core"
	"github.com/rs/cors"
	"net/http"
	hs "senet/handlers"
)

func main() {
	conf := &core.Config{
		CorsOptions: cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		},
		Host:  "localhost",
		Port:  3000,
		Debug: true,
	}

	handlers := make(core.Handlers)
	handlers[hs.SEND_MESSAGE] = &hs.SenderHandler{}

	app, err := core.Start(handlers, conf)
	if err != nil {
		fmt.Println(err)
	}

	for {
		hook := <-app.Hooks
		switch hook.HookType {
		case core.CLIENT_ADDED:
			fmt.Printf("Client added: %s\n", hook.Data)
		case core.CLIENT_CLOSED_CONNECTION:
			fmt.Printf("Client closed connection: %s\n", hook.Data)
		case core.ERROR:
			fmt.Printf("Error: %s\n", hook.Data)
		}
	}
}
