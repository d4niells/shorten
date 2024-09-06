package utils

import "testing"

func TestIsValidURL(t *testing.T) {
	testcase := []struct {
		name     string
		url      string
		expected bool
	}{
		{name: "https", url: "https://example.com", expected: true},
		{name: "http", url: "http://example.com", expected: true},
		{name: "https with query params", url: "https://example.com/path?query=string#d", expected: true},
		{name: "without schema and domain", url: "invalid-url", expected: false},
		{name: "invalid schema", url: "ftp://example.com", expected: false},
		{name: "without schema", url: "://example.com", expected: false},
		{name: "empty string", url: "", expected: false},
		{name: "ispace in domain", url: "https://exa mple.com", expected: false},
		{url: "xpto:///example.com", expected: false},
	}

	for _, c := range testcase {
		t.Run(c.name, func(t *testing.T) {
			got := IsValidURL(c.url)
			if got != c.expected {
				t.Errorf("expected %t, got %t", c.expected, got)
			}
		})
	}
}
