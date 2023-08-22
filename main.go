package main

import (
	"embed"
	"fmt"
	"github.com/AlekseiKromski/at-socket-server/core"
	"github.com/gin-contrib/cors"
	"net/http"
	"os"
	"senet/config"
	"senet/processor"
	"strconv"
)

var (
	//go:embed webclient/build
	frontend embed.FS
)

func main() {
	conf := getConfig()

	sp, err := processor.NewProcessor(conf)
	if err != nil {
		fmt.Printf("problem with processor: %v", err)
		return
	}

	if err := sp.Start(frontend); err != nil {
		fmt.Printf("Cannot start porcessor: %v", err)
	}
}

func getConfig() *config.Config {
	host := os.Getenv("ADDRESS")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Printf("cannot transform port to int: %v", err)
		return nil
	}
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		fmt.Printf("cannot transform port to int: %v", err)
		return nil
	}

	dbHostname := os.Getenv("DB_HOSTNAME")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")

	ac := &core.Config{
		CorsOptions: cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
			},
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
		},
		Host:  host,
		Port:  port,
		Debug: debug,
	}
	dc := config.NewDbConfig(dbHostname, dbUsername, dbPassword, dbDatabase)

	return config.NewConfig(ac, dc)
}
