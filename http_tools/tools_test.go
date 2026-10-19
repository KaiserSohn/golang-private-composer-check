package http_tools

import "testing"

func TestValidateUrl(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"https url", "https://packages.org/repository/composer-proxy-packagist/symfony/console/v7.0.0/symfony-console-v7.0.0.zip", false},
		{"http url", "http://example.com/package.zip", false},
		{"url with credentials", "https://user:pass@example.com/package.zip", false},
		{"url with port", "https://example.com:8080/package.zip", false},
		{"empty url", "", true},
		{"whitespace only", "   ", true},
		{"no scheme", "packages.org/repository/package.zip", true},
		{"ftp scheme", "ftp://example.com/package.zip", true},
		{"file scheme", "file:///etc/passwd", true},
		{"no host", "https:///package.zip", true},
		{"scheme only", "https://", true},
		{"garbage with spaces", "https://exa mple.com/package.zip", true},
		{"control character", "https://example.com/\x7f", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUrl(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUrl(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}
