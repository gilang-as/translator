package googletranslate

import (
	"net/http"
	"sync"
)

const (
	DefaultHost = "google.com"
)

// GoogleTranslate is a concurrency-safe client for the Google Translate API.
type GoogleTranslate struct {
	mu       sync.RWMutex
	host     string
	client   *http.Client
	proxyURL string
}

// Option is a functional option for configuring GoogleTranslate.
type Option func(*GoogleTranslate)

// WithHost sets a custom host for the Google Translate API.
func WithHost(host string) Option {
	return func(gt *GoogleTranslate) {
		gt.host = host
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(gt *GoogleTranslate) {
		gt.client = client
	}
}

// WithProxyURL sets a proxy URL for requests.
func WithProxyURL(proxyURL string) Option {
	return func(gt *GoogleTranslate) {
		gt.proxyURL = proxyURL
	}
}

// New creates a new GoogleTranslate client with the given options.
func New(opts ...Option) *GoogleTranslate {
	gt := &GoogleTranslate{
		host:   DefaultHost,
		client: &http.Client{},
	}
	for _, opt := range opts {
		opt(gt)
	}
	return gt
}

// Host returns the current host.
func (gt *GoogleTranslate) Host() string {
	gt.mu.RLock()
	defer gt.mu.RUnlock()
	return gt.host
}

// SetHost sets the host.
func (gt *GoogleTranslate) SetHost(host string) {
	gt.mu.Lock()
	defer gt.mu.Unlock()
	gt.host = host
}

// Client returns the current HTTP client.
func (gt *GoogleTranslate) Client() *http.Client {
	gt.mu.RLock()
	defer gt.mu.RUnlock()
	return gt.client
}

// SetClient sets the HTTP client.
func (gt *GoogleTranslate) SetClient(client *http.Client) {
	gt.mu.Lock()
	defer gt.mu.Unlock()
	gt.client = client
}

// ProxyURL returns the current proxy URL.
func (gt *GoogleTranslate) ProxyURL() string {
	gt.mu.RLock()
	defer gt.mu.RUnlock()
	return gt.proxyURL
}

// SetProxyURL sets the proxy URL.
func (gt *GoogleTranslate) SetProxyURL(proxyURL string) {
	gt.mu.Lock()
	defer gt.mu.Unlock()
	gt.proxyURL = proxyURL
}
