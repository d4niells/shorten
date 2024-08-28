package shorten

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrEmptyKey     = errors.New("key cannot be empty")
	ErrEmptyLongURL = errors.New("long url cannot be empty")
)

type URL struct {
	Key       string    `json:"key"`
	LongURL   string    `json:"long_url"`
	ShortURL  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
}

func NewURL(key, longURL string) *URL {
	return &URL{
		Key:       key,
		LongURL:   longURL,
		ShortURL:  fmt.Sprintf("https://shorten.io/%s", key),
		CreatedAt: time.Now(),
	}
}

func (u *URL) Validate() error {
	if u.Key == "" {
		return ErrEmptyKey
	}
	if u.LongURL == "" {
		return ErrEmptyLongURL
	}
	return nil
}
