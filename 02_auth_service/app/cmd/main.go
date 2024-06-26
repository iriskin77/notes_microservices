package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/iriskin77/notes_microservices/app/internal/auth"
	"github.com/iriskin77/notes_microservices/app/internal/config"
	"github.com/iriskin77/notes_microservices/app/pkg/db"
	"github.com/iriskin77/notes_microservices/app/pkg/jwt"
	"github.com/iriskin77/notes_microservices/app/pkg/logging"
	"github.com/iriskin77/notes_microservices/app/pkg/redis_client"
	"github.com/joho/godotenv"
)

func main() {

	// initializing config
	conf := config.LoadConfig("./app/config/config.yml")

	// initializing logger
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	// load .env variables

	if err := godotenv.Load(); err != nil {
		logger.Fatal(err)
	}

	// initializin db
	confDB := db.NewPostgresConfig(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("SSLMODE"),
	)

	ctx := context.Background()
	client, err := db.NewPostgresDB(ctx, confDB)

	if err != nil {
		logger.Fatal(err)
	}

	// initializing Redis

	redisConfig, err := redis_client.NewRedisConfig(
		os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("REDIS_DB"),
	)
	if err != nil {
		logger.Fatal(err)
	}

	clientRedis, err := redis_client.NewRedisClient(redisConfig)
	if err != nil {
		logger.Fatal(err)
	}

	storeRedis := auth.NewRedisStore(clientRedis)

	// initializing tokenManager to generate JWT

	tokenManager, err := jwt.NewManager(conf.JWTSecret, time.Duration(conf.AccessTokenTTL)*time.Minute, time.Duration(conf.RefreshTokenTTL)*time.Minute)
	if err != nil {
		logger.Fatal(err)
	}

	// initializing server
	repo := auth.NewRepository(client, &logger)
	service := auth.NewService(repo, tokenManager, storeRedis, &logger,
		time.Duration(conf.AccessTokenTTL)*time.Minute,
		time.Duration(conf.RefreshTokenTTL)*time.Minute)

	h := auth.NewHandler(service, tokenManager, &logger)

	router := mux.NewRouter()

	h.RegisterHandlers(router)

	srv, err := NewHttpServer(router, conf)

	if err != nil {
		logger.Fatal(err)
	}

	srv.ListenAndServe()

}

func NewHttpServer(router *mux.Router, cfg *config.Config) (*http.Server, error) {

	return &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}, nil
}
