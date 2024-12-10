package server

import (
	bucket "debashare-go/internal/minio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"debashare-go/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int

	db    database.Service
	minio bucket.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db:    database.New(),
		minio: bucket.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
