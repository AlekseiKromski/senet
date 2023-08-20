package main

import (
	"embed"
	"fmt"
	"github.com/AlekseiKromski/at-socket-server/core"
	"github.com/gin-contrib/cors"
	"net/http"
	"senet/processor"
)

var (
	//go:embed webclient/build
	frontend embed.FS
)

func main() {
	conf := &core.Config{
		CorsOptions: cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
			},
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
		},
		Host:  "localhost",
		Port:  3000,
		Debug: true,
	}

	sp := processor.NewProcessor(conf)
	if err := sp.Start(frontend); err != nil {
		fmt.Printf("Cannot strat porcessor: %v", err)
	}
}
