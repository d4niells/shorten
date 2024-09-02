package entity

import (
	"errors"
	"testing"
)

func TestNewURL(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		key, longURL, shortURL := "EAaArVRs", "https://example.com", "http://localhost:8080/EAaArVRs"

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
				Key:     "EAaArVRs",
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
				Key: "EAaArVRs",
			},
			err: ErrEmptyLongURL,
		},
		{
			name: "invalid key size",
			url: &URL{
				Key:     "example-of-large-key",
				LongURL: "https://example.com",
			},
			err: ErrInvalidKeySize,
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

func TestFromToJSON(t *testing.T) {
	t.Run("encode to JSON", func(t *testing.T) {
		url := NewURL("EAaArVRs", "https://example.com")
		str, err := url.ToJSON()
		if err != nil {
			t.Errorf("expected %v, got %v", nil, err)
		}
		if str == "" {
			t.Error("expected encoded url, got empty string")
		}
	})

	t.Run("decode to JSON", func(t *testing.T) {
		url := NewURL("EAaArVRs", "https://example.com")
		str, _ := url.ToJSON()

		err := url.FromJSON(str)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}
