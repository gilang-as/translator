package gt

import (
	"net/http"

	"gopkg.gilang.dev/translator/googletranslate"
)

// Re-export types from googletranslate for backward compatibility.
type (
	GoogleTranslate = googletranslate.GoogleTranslate
	Option          = googletranslate.Option
	Translated      = googletranslate.Translated
)

// NewGoogleTranslate creates a new GoogleTranslate client.
// Deprecated: Use googletranslate.New() instead.
func NewGoogleTranslate(opts ...Option) *GoogleTranslate {
	return googletranslate.New(opts...)
}

// WithHost sets a custom host.
func WithHost(host string) Option {
	return googletranslate.WithHost(host)
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return googletranslate.WithHTTPClient(client)
}
