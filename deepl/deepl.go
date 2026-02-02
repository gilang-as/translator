package deepl

import (
	"context"
	"net/http"
	"sync"
)

const (
	DefaultHost = "www2.deepl.com"
)

// DeepL is a concurrency-safe client for the DeepL API.
type DeepL struct {
	mu        sync.RWMutex
	host      string
	client    *http.Client
	proxyURL  string
	dlSession string
}

// Option is a functional option for configuring DeepL.
type Option func(*DeepL)

// WithHost sets a custom host for the DeepL API.
func WithHost(host string) Option {
	return func(d *DeepL) {
		d.host = host
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(d *DeepL) {
		d.client = client
	}
}

// WithProxyURL sets a proxy URL for requests.
func WithProxyURL(proxyURL string) Option {
	return func(d *DeepL) {
		d.proxyURL = proxyURL
	}
}

// WithDLSession sets a DeepL session cookie for Pro features.
func WithDLSession(dlSession string) Option {
	return func(d *DeepL) {
		d.dlSession = dlSession
	}
}

// New creates a new DeepL client with the given options.
func New(opts ...Option) *DeepL {
	d := &DeepL{
		host:   DefaultHost,
		client: &http.Client{},
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Host returns the current host.
func (d *DeepL) Host() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.host
}

// SetHost sets the host.
func (d *DeepL) SetHost(host string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.host = host
}

// Client returns the current HTTP client.
func (d *DeepL) Client() *http.Client {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.client
}

// SetClient sets the HTTP client.
func (d *DeepL) SetClient(client *http.Client) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.client = client
}

// ProxyURL returns the current proxy URL.
func (d *DeepL) ProxyURL() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.proxyURL
}

// SetProxyURL sets the proxy URL.
func (d *DeepL) SetProxyURL(proxyURL string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.proxyURL = proxyURL
}

// DLSession returns the current DeepL session.
func (d *DeepL) DLSession() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.dlSession
}

// SetDLSession sets the DeepL session.
func (d *DeepL) SetDLSession(dlSession string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.dlSession = dlSession
}

// TranslateFromLanguage contains detected language information.
type TranslateFromLanguage struct {
	DidYouMean bool   `json:"did_you_mean"`
	Iso        string `json:"iso"`
}

// TranslateFromText contains text correction information.
type TranslateFromText struct {
	AutoCorrected bool    `json:"auto_corrected"`
	Value         *string `json:"value"`
	DidYouMean    bool    `json:"did_you_mean"`
}

// TranslateFrom contains source language and text information.
type TranslateFrom struct {
	Language TranslateFromLanguage `json:"language"`
	Text     TranslateFromText     `json:"text"`
}

// Translated represents a translation result.
type Translated struct {
	Text          string        `json:"text"`
	Pronunciation *string       `json:"pronunciation"`
	Alternatives  []string      `json:"alternatives"`
	From          TranslateFrom `json:"from"`
	Method        string        `json:"method"`
}

// Translate translates text from one language to another.
func (d *DeepL) Translate(ctx context.Context, text string, from string, to string) (*Translated, error) {
	d.mu.RLock()
	proxyURL := d.proxyURL
	dlSession := d.dlSession
	d.mu.RUnlock()

	result, err := TranslateByDeepL(ctx, from, to, text, "", proxyURL, dlSession)
	if err != nil {
		return nil, err
	}

	if result.Code != http.StatusOK {
		return nil, &TranslationError{
			Code:    result.Code,
			Message: result.Message,
		}
	}

	return &Translated{
		Text:         result.Data,
		Alternatives: result.Alternatives,
		From: TranslateFrom{
			Language: TranslateFromLanguage{
				Iso: result.SourceLang,
			},
		},
		Method: result.Method,
	}, nil
}

// TranslationError represents a translation error.
type TranslationError struct {
	Code    int
	Message string
}

func (e *TranslationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "translation failed"
}
