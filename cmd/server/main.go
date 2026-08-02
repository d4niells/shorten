package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/d4niells/shorten/internal/handler"
	"github.com/d4niells/shorten/internal/repository"
	"github.com/d4niells/shorten/internal/service"
	"github.com/go-redis/redis"
)

func main() {
	rHost := os.Getenv("REDIS_HOST")
	if rHost == "" {
		rHost = "localhost"
	}
	rPort := os.Getenv("REDIS_PORT")
	if rPort == "" {
		rPort = "6379"
	}

	redisClient := redis.NewClient(
		&redis.Options{
			Addr: fmt.Sprintf("%s:%s", rHost, rPort),
		},
	)
	if err := redisClient.Ping().Err(); err != nil {
		log.Fatalf("couldn't connect to Redis: %v\n", err)
	}

	cache := repository.NewCache(redisClient)
	urlService := service.NewURLService(cache)
	urlHandler := handler.NewURLHandler(urlService)

	r := http.NewServeMux()
	r.HandleFunc("POST /", urlHandler.Shorten)
	r.HandleFunc("GET /{key}", urlHandler.Resolver)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("couldn't listen on port %v: %v\n", port, err)
	}
}
