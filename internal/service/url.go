package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/d4niells/shorten/internal/entity"
	"github.com/d4niells/shorten/internal/repository"
)

type URLService interface {
	Shorten(ctx context.Context, longURL string) (*entity.URL, error)
}

type URLServiceImpl struct {
	cache repository.CacheRepository
}

func NewURLService(cache repository.CacheRepository) *URLServiceImpl {
	return &URLServiceImpl{cache}
}

func (s *URLServiceImpl) Shorten(ctx context.Context, longURL string) (*entity.URL, error) {
	key := genSHA256Hash(longURL, entity.KEY_SIZE)

	cachedURL, err := s.cache.Get(ctx, key)
	if err == nil {
		return cachedURL, nil
	}

	if err != repository.ErrKeyDoesNotExist {
		return nil, err
	}

	newURL := entity.NewURL(key, longURL)
	if err := newURL.Validate(); err != nil {
		return nil, err
	}

	// TODO: Add a great expiration time for shortened URLs
	if err := s.cache.Set(ctx, newURL, time.Duration(0)); err != nil {
		return nil, err
	}

	return newURL, nil
}

func genSHA256Hash(url string, length int) string {
	hash := sha256.Sum256([]byte(url))
	encodedHash := base64.URLEncoding.EncodeToString(hash[:])
	return encodedHash[:length]
}
