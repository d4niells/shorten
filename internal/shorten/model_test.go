package shorten

import (
	"errors"
	"testing"
)

func TestNewURL(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		key, longURL, shortURL := "exmpl", "https://example.com", "https://shorten.io/exmpl"

		url := NewURL(key, longURL)

		if url.Key != key {
			t.Errorf("expected %s, got %s", url.Key, key)
		}
		if url.LongURL != longURL {
			t.Errorf("expected %s, got %s", longURL, url.LongURL)
		}
		if url.ShortURL != shortURL {
			t.Errorf("expected %s, got %s", shortURL, url.ShortURL)
		}
	})
}

func TestValidate(t *testing.T) {
	testcases := []struct {
		name string
		url  *URL
		err  error
	}{
		{
			name: "valid URL",
			url: &URL{
				Key:     "exmpl",
				LongURL: "https://example.com",
			},
			err: nil,
		},
		{
			name: "missing key",
			url: &URL{
				LongURL: "https://example.com",
			},
			err: ErrEmptyKey,
		},
		{
			name: "missing long URL",
			url: &URL{
				Key: "exmpl",
			},
			err: ErrEmptyLongURL,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.url.Validate()

			if !errors.Is(err, tc.err) {
				t.Errorf("expected %v, got %v", tc.err, err)
			}
		})
	}
}
