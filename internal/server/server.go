package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/imide/debashare-go/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "debashare-go",
			AppName:      "debashare-go",
		}),

		db: database.New(),
	}

	return server
}
