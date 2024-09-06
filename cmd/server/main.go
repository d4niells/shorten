package main

import (
	"log"
	"net/http"
	"time"

	"github.com/d4niells/shorten/internal/handler"
	"github.com/d4niells/shorten/internal/repository"
	"github.com/d4niells/shorten/internal/service"
	"github.com/go-redis/redis"
)

const (
	PORT       = ":8080"
	REDIS_ADDR = "localhost:6379"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{Addr: REDIS_ADDR})
	if err := redisClient.Ping().Err(); err != nil {
		log.Fatalf("couldn't connect to Redis: %v\n", err)
	}

	cache := repository.NewCache(redisClient)
	urlService := service.NewURLService(cache)
	urlHandler := handler.NewURLHandler(urlService)

	r := http.NewServeMux()
	r.HandleFunc("POST /", urlHandler.Shorten)
	r.HandleFunc("GET /{key}", urlHandler.Resolver)

	srv := http.Server{
		Addr:         PORT,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("couldn't listen on port %v: %v\n", PORT, err)
	}
}
