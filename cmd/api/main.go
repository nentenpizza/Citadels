package main

import (
	"log"
	"os"

	"github.com/nentenpizza/citadels/internal/service"

	"github.com/nentenpizza/citadels/internal/repository"
	"github.com/nentenpizza/citadels/internal/repository/postgres"
	httpapi "github.com/nentenpizza/citadels/internal/transport/http"
)

func main() {
	db, err := postgres.Open(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	repos := repository.New(db)

	usersService := service.Users{Repos: repos}

	server := httpapi.App{UsersService: usersService}

	log.Fatal(server.Run(":8080"))
}
