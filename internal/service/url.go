package service

import (
	"context"
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
	// TODO: Create an algorithm to be able to generate an unique key. We must
	// check for key collision and idempotency.
	url := entity.NewURL("unique-key", longURL)

	if err := url.Validate(); err != nil {
		return nil, err
	}

	// TODO: I need to think about a great expiration time for shortened URLs
	err := s.cache.Set(ctx, url, time.Duration(0))
	if err != nil {
		return nil, err
	}

	return url, nil
}
