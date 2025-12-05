package main

import (
	"embed"
	"log"
	"os"
	"strconv"

	"alekseikromski.com/senet/modules/gin_server"
	server_key_storage "alekseikromski.com/senet/modules/server-key-storage"
	"alekseikromski.com/senet/modules/storage/postgres"
	"github.com/AlekseiKromski/server-core/core"
	"github.com/joho/godotenv"
)

var (
	//go:embed front-end/build
	resources embed.FS
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Main: cannot load env form .env file: %v", err)
		return
	}

	ginCookieDomain := os.Getenv("GIN_COOKIE_DOMAIN")

	ginAddress := os.Getenv("GIN_ADDRESS")
	ginSecret := os.Getenv("GIN_SECRET")

	dbDatabase := os.Getenv("DB_DATABASE")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Printf("Main: cannot load DB_PORT: %v", err)
		return
	}

	c := core.NewCore()
	c.Init([]core.Module{
		gin_server.NewServer(
			gin_server.NewServerConfig(ginSecret, ginAddress, ginCookieDomain),
			resources,
		),
		postgres.NewPostgres(
			postgres.NewConfig(
				dbHost,
				dbDatabase,
				dbUsername,
				dbPassword,
				dbPort,
			),
		),
		server_key_storage.NewStorage(
			server_key_storage.NewConfig(
				"./keys/public-key.asc",
				"./keys/private-key.asc",
			),
		),
	})
}
