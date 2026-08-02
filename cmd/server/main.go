package main

import (
	"fmt"
	"html/template"
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

	// Parse Go templates
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v\n", err)
	}

	cache := repository.NewCache(redisClient)
	urlService := service.NewURLService(cache)
	urlHandler := handler.NewURLHandler(urlService)
	webHandler := handler.NewWebHandler(tmpl, urlService)

	r := http.NewServeMux()

	// Static files (register BEFORE catch-all routes)
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// API endpoints (specific routes)
	r.HandleFunc("POST /api/shorten", urlHandler.Shorten)
	r.HandleFunc("GET /{key}", urlHandler.Resolver)

	// Frontend with Go templates (catch-all, register LAST)
	r.HandleFunc("GET /", webHandler.Home)
	r.HandleFunc("POST /", webHandler.Home)

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

	log.Printf("Server listening on http://localhost:%s\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("couldn't listen on port %v: %v\n", port, err)
	}
}
