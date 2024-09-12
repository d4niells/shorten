package utils

import "testing"

func TestIsValidURL(t *testing.T) {
	testCases := []struct {
		name     string
		url      string
		expected bool
	}{
		{name: "with https scheme", url: "https://example.com", expected: true},
		{name: "with http scheme", url: "http://example.com", expected: true},
		{name: "with query param and fragment", url: "https://example.com/path?query=string#fragment", expected: true},
		{name: "plain text", url: "example", expected: false},
		{name: "invalid scheme", url: "ftp://example.com", expected: false},
		{name: "without scheme", url: "://example.com", expected: false},
		{name: "empty string", url: "", expected: false},
		{name: "without hostname", url: "https://", expected: false},
		{name: "empty space between hostname", url: "https://exa mple.com", expected: false},
		{name: "triple slash", url: "https:///example.com", expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsValidURL(tc.url)
			if got != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, got)
			}
		})
	}
}

func BenchmarkIsValidURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsValidURL("https://example.com")
	}
}
