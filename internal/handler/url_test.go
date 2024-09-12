package handler

import (
	"bytes"
	"context"
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

func setupTestEnvironment(t *testing.T) (*http.ServeMux, *service.URLServiceImpl) {
	t.Helper()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	cacheRepo := repository.NewCache(redisClient)
	urlService := service.NewURLService(cacheRepo)
	urlHandler := NewURLHandler(urlService)

	r := http.NewServeMux()
	r.HandleFunc("/", urlHandler.Shorten)
	r.HandleFunc("/{key}", urlHandler.Resolver)

	return r, urlService
}

func TestShorten(t *testing.T) {
	t.Run("returns valid short URL on successful input", func(t *testing.T) {
		handler, urlService := setupTestEnvironment(t)

		longURL := "https://github.com/d4niells"
		reqBody := map[string]any{"long_url": longURL}
		reqBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBytes))
		req.Header.Add("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Errorf("expected status code 201 Created, got %d", res.StatusCode)
		}

		var resBody map[string]string
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if resBody["key"] == "" {
			t.Error("expected valid key value, got empty string")
		}

		if len(resBody["key"]) != 8 {
			t.Errorf("expected key length 8, got %d", len(resBody["key"]))
		}

		if resBody["short_url"] == "" {
			t.Error("expected valid short_url value, got empty string")
		}

		if !strings.Contains(resBody["short_url"], resBody["key"]) {
			t.Errorf("expected short_url to contain '%s', got %s", resBody["key"], resBody["short_url"])
		}

		if resBody["long_url"] != longURL {
			t.Errorf("expected long URL %s, got %s", longURL, resBody["long_url"])
		}

		t.Cleanup(func() {
			err := urlService.Delete(context.Background(), resBody["key"])
			if err != nil {
				t.Fatalf("failed to delete key from Redis: %v", err)
			}
		})
	})

	t.Run("invalid request payload", func(t *testing.T) {
		handler, _ := setupTestEnvironment(t)

		testCases := []struct {
			name        string
			requestBody map[string]any
			expectedErr string
		}{
			{"returns 400 for invalid URL format", map[string]any{"long_url": "x:///github.com"}, "invalid URL format"},
			{"returns 400 for missing long_url field", map[string]any{"long_url": ""}, "missing field: long_url cannot be empty"},
			{"returns 400 for missing request payload", map[string]any{}, "missing field: long_url cannot be empty"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				reqBytes, _ := json.Marshal(tc.requestBody)
				req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBytes))
				req.Header.Add("Content-Type", "application/json")

				rec := httptest.NewRecorder()

				handler.ServeHTTP(rec, req)

				res := rec.Result()
				defer res.Body.Close()

				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("expected status code 400 Bad Request, got %d", res.StatusCode)
				}

				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}

				if strings.TrimSpace(string(resBody)) != tc.expectedErr {
					t.Errorf("expected error %s, got %s", tc.expectedErr, resBody)
				}
			})
		}
	})
}

func TestResolver(t *testing.T) {
	t.Run("redirects to long URL when short URL is found", func(t *testing.T) {
		handler, urlService := setupTestEnvironment(t)

		longURL := "https://example.com"
		newURL, err := urlService.Shorten(context.Background(), longURL)
		if err != nil {
			t.Fatalf("failed to create short URL: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/"+newURL.Key, nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusFound {
			t.Errorf("expected status code to be 302 Found, got %d", res.StatusCode)
		}

		if res.Header.Get("Location") != longURL {
			t.Errorf("expected redirect to %s, got %s", longURL, res.Header.Get("Location"))
		}

		t.Cleanup(func() {
			err := urlService.Delete(context.Background(), newURL.Key)
			if err != nil {
				t.Fatalf("failed to delete key from Redis: %v", err)
			}
		})
	})

	t.Run("returns 404 when short URL is not found", func(t *testing.T) {
		handler, _ := setupTestEnvironment(t)

		req := httptest.NewRequest(http.MethodGet, "/test-key", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("expected status code to be 404 Not Found, got %d", res.StatusCode)
		}

		if res.Header.Get("Location") != "" {
			t.Errorf("expected empty Location header, got %s", res.Header.Get("Location"))
		}
	})
}
