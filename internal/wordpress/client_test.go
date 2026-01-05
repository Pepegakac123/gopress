package wordpress

import (
	"testing"
)

func TestNewClient_URLFormatting(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		wantURL  string
	}{
		{
			name:     "Standard URL",
			inputURL: "https://example.com",
			wantURL:  "https://example.com/wp-json",
		},
		{
			name:     "URL with trailing slash",
			inputURL: "https://example.com/",
			wantURL:  "https://example.com/wp-json",
		},
		{
			name:     "URL with wp-json suffix",
			inputURL: "https://example.com/wp-json",
			wantURL:  "https://example.com/wp-json",
		},
		{
			name:     "URL with wp-json and trailing slash",
			inputURL: "https://example.com/wp-json/",
			wantURL:  "https://example.com/wp-json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.inputURL, "user", "pass", "token")
			if client.baseURL != tt.wantURL {
				t.Errorf("NewClient() baseURL = %v, want %v", client.baseURL, tt.wantURL)
			}
		})
	}
}
