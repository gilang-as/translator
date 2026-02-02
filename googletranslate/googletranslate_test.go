package googletranslate

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	gt := New()
	if gt == nil {
		t.Error("New() returned nil")
	}
	if gt.Host() != DefaultHost {
		t.Errorf("Expected host %s, got %s", DefaultHost, gt.Host())
	}
}

func TestWithHost(t *testing.T) {
	customHost := "google.co.id"
	gt := New(WithHost(customHost))
	if gt.Host() != customHost {
		t.Errorf("Expected host %s, got %s", customHost, gt.Host())
	}
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 30 * time.Second}
	gt := New(WithHTTPClient(customClient))
	if gt.Client() != customClient {
		t.Error("Custom HTTP client not set correctly")
	}
}

func TestSetHost(t *testing.T) {
	gt := New()
	newHost := "google.co.uk"
	gt.SetHost(newHost)
	if gt.Host() != newHost {
		t.Errorf("Expected host %s, got %s", newHost, gt.Host())
	}
}

func TestWithProxyURL(t *testing.T) {
	proxyURL := "http://proxy:8080"
	gt := New(WithProxyURL(proxyURL))
	if gt.ProxyURL() != proxyURL {
		t.Errorf("Expected proxyURL %s, got %s", proxyURL, gt.ProxyURL())
	}
}

func TestSetProxyURL(t *testing.T) {
	gt := New()
	proxyURL := "http://proxy:8080"
	gt.SetProxyURL(proxyURL)
	if gt.ProxyURL() != proxyURL {
		t.Errorf("Expected proxyURL %s, got %s", proxyURL, gt.ProxyURL())
	}
}

func TestTranslate(t *testing.T) {
	gt := New()
	data, err := gt.Translate(context.Background(), "Hello World", "en", "id")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if data.Text == "" {
		t.Error("Translation returned empty text")
	}
}
