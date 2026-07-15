package utils

import (
	"testing"
)

func TestValidateUrl(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "HTTPS URL",
			url:     "https://google.com",
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			url:     "invalid-url",
			wantErr: true,
		},
		{
			name:    "HTTP URL",
			url:     "http://github.com",
			wantErr: false,
		},
		{
			name:    "FTP URL",
			url:     "ftp://server.com",
			wantErr: true,
		},
		{
			name:    "No Scheme URL",
			url:     "google.com",
			wantErr: true,
		},
		{
			name:    "Javascript URL",
			url:     "javascript:alert(1)",
			wantErr: true,
		},
		{
			name:    "Empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "Missing host URL",
			url:     "https://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateURL(%q) = %v; expected %v", tt.url, err, tt.wantErr)
			}
		})
	}
}
