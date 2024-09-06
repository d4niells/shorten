package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/d4niells/shorten/internal/repository"
	"github.com/d4niells/shorten/internal/service"
	"github.com/go-redis/redis"
)

func setupTestEnvironment() *http.ServeMux {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	cacheRepo := repository.NewCache(redisClient)
	urlService := service.NewURLService(cacheRepo)
	urlHandler := NewURLHandler(urlService)

	r := http.NewServeMux()
	r.HandleFunc("/", urlHandler.Shorten)
	r.HandleFunc("/{key}", urlHandler.Resolver)

	return r
}

// TODO: add test cleanup to delete keys from Redis to avoid misunderstandings
func TestShortenHandler(t *testing.T) {
	t.Run("return a short url", func(t *testing.T) {
		urlHandler := setupTestEnvironment()

		server := httptest.NewServer(urlHandler)
		defer server.Close()

		longURL := "https://github.com/d4niells"
		reqBody := map[string]any{"long_url": longURL}
		reqBytes, _ := json.Marshal(reqBody)

		req, err := http.NewRequest(http.MethodPost, server.URL+"/", bytes.NewBuffer(reqBytes))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("expected request to be successfuly, got %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Fatalf("expected status code 201 Created, got %d", res.StatusCode)
		}

		var resBody map[string]string
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatalf("expected response body to be decoded, got %v", err)
		}

		if resBody["key"] == "" {
			t.Fatal("expected a valid key value, got a empty string")
		}

		if len(resBody["key"]) != 8 {
			t.Fatalf("expected key length to be 8, got %d", len(resBody["key"]))
		}

		if resBody["short_url"] == "" {
			t.Fatalf("expected a valid short_url value, got a empty string")
		}

		if !strings.Contains(resBody["short_url"], resBody["key"]) {
			t.Fatalf("expected a valid short_url, got %s", resBody["short_url"])
		}

		if resBody["long_url"] != longURL {
			t.Fatalf("expected %s, got %s", longURL, resBody["long_url"])
		}
	})

	t.Run("invalid request payload", func(t *testing.T) {
		urlHandler := setupTestEnvironment()

		server := httptest.NewServer(urlHandler)
		defer server.Close()

		req, err := http.NewRequest(http.MethodPost, server.URL+"/", nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("expected request to be successful, got %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status code 400 Bad Request, got %d", res.StatusCode)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected response body to be read, got %v", err)
		}

		if strings.TrimSpace(string(resBody)) != "invalid request payload" {
			t.Errorf(`expected error message "invalid request payload", got %s`, resBody)
		}
	})

	t.Run("missing field", func(t *testing.T) {
		urlHandler := setupTestEnvironment()

		server := httptest.NewServer(urlHandler)
		defer server.Close()

		reqBody := map[string]any{"long_url": ""}
		reqBytes, _ := json.Marshal(reqBody)

		req, err := http.NewRequest(http.MethodPost, server.URL+"/", bytes.NewBuffer(reqBytes))
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("expected request to be successful, got %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status code 400 Bad Request, got %d", res.StatusCode)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected response body to be read, got %v", err)
		}

		if strings.TrimSpace(string(resBody)) != "missing field: long_url cannot be empty" {
			t.Errorf(`expected error message "invalid request payload", got %s`, resBody)
		}
	})

	t.Run("invalid URL format", func(t *testing.T) {
		urlHandler := setupTestEnvironment()

		server := httptest.NewServer(urlHandler)
		defer server.Close()

		reqBody := map[string]any{"long_url": "x:///github.com"}
		reqBytes, _ := json.Marshal(reqBody)

		req, err := http.NewRequest(http.MethodPost, server.URL+"/", bytes.NewBuffer(reqBytes))
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("expected request to be successful, got %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status code 400 Bad Request, got %d", res.StatusCode)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("expected response body to be read, got %v", err)
		}

		if strings.TrimSpace(string(resBody)) != "invalid URL format" {
			t.Errorf(`expected error message "invalid URL format", got %s`, resBody)
		}
	})
}
