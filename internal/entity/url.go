package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	utils "github.com/d4niells/shorten/pkg"
)

const KEY_SIZE = 8

var (
	ErrEmptyKey             = errors.New("key cannot be empty")
	ErrEmptyLongURL         = errors.New("long url cannot be empty")
	ErrInvalidKeySize       = errors.New("key must be 8 characters long")
	ErrInvalidLongURLFormat = errors.New("invalid URL format")
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
		ShortURL:  fmt.Sprintf("http://localhost:8080/%s", key),
		CreatedAt: time.Now(),
	}
}

func (u *URL) Validate() error {
	if u.Key == "" {
		return ErrEmptyKey
	}
	if len(u.Key) != KEY_SIZE {
		return ErrInvalidKeySize
	}
	if u.LongURL == "" {
		return ErrEmptyLongURL
	}
	if !utils.IsValidURL(u.LongURL) {
		return ErrInvalidLongURLFormat
	}
	return nil
}

func (u *URL) ToJSON() (string, error) {
	bytes, err := json.Marshal(&u)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (u *URL) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), u)
}
