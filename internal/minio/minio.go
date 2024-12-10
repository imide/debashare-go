package minio

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

type Service interface {
	// Health returns a map of health status information
	// The keys and values in the map are service-specific
	Health() map[string]string

	// Close terminates the bucket connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	client *minio.Client
}

var (
	endpoint        = os.Getenv("MINIO_ENDPOINT")
	accessKeyID     = os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey = os.Getenv("MINIO_SECRET_KEY")

	minioInstance *service
)

func New() Service {
	if minioInstance != nil {
		return minioInstance
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize minio client")
	}

	minioInstance = &service{
		client: client,
	}

	log.Info().Msg("successfully connected to the bucket")

	return minioInstance
}

func (s *service) Health() map[string]string {
	// TODO implement me
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stats := make(map[string]string)

	// Ping the bucket (list buckets)
	_, err := s.client.ListBuckets(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("failed to list buckets: %v", err)
		log.Fatal().Err(err).Msg("failed to connect to minio")
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "minio healthy"
}

func (s service) Close() error {
	// TODO implement me
	panic("implement me")
}
