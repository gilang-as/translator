package deepl

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	d := New()
	if d == nil {
		t.Error("New() returned nil")
	}
	if d.Host() != DefaultHost {
		t.Errorf("Expected host %s, got %s", DefaultHost, d.Host())
	}
}

func TestWithHost(t *testing.T) {
	customHost := "api.deepl.com"
	d := New(WithHost(customHost))
	if d.Host() != customHost {
		t.Errorf("Expected host %s, got %s", customHost, d.Host())
	}
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 30 * time.Second}
	d := New(WithHTTPClient(customClient))
	if d.Client() != customClient {
		t.Error("Custom HTTP client not set correctly")
	}
}

func TestWithProxyURL(t *testing.T) {
	proxyURL := "http://proxy:8080"
	d := New(WithProxyURL(proxyURL))
	if d.ProxyURL() != proxyURL {
		t.Errorf("Expected proxyURL %s, got %s", proxyURL, d.ProxyURL())
	}
}

func TestWithDLSession(t *testing.T) {
	session := "test-session-token"
	d := New(WithDLSession(session))
	if d.DLSession() != session {
		t.Errorf("Expected dlSession %s, got %s", session, d.DLSession())
	}
}

func TestSetHost(t *testing.T) {
	d := New()
	newHost := "api.deepl.com"
	d.SetHost(newHost)
	if d.Host() != newHost {
		t.Errorf("Expected host %s, got %s", newHost, d.Host())
	}
}

func TestSetProxyURL(t *testing.T) {
	d := New()
	proxyURL := "http://proxy:8080"
	d.SetProxyURL(proxyURL)
	if d.ProxyURL() != proxyURL {
		t.Errorf("Expected proxyURL %s, got %s", proxyURL, d.ProxyURL())
	}
}

func TestSetDLSession(t *testing.T) {
	d := New()
	session := "new-session-token"
	d.SetDLSession(session)
	if d.DLSession() != session {
		t.Errorf("Expected dlSession %s, got %s", session, d.DLSession())
	}
}

func TestTranslate(t *testing.T) {
	d := New()
	data, err := d.Translate(context.Background(), "Hello World", "en", "id")
	if err != nil {
		t.Skipf("Skipping due to API error: %v", err)
		return
	}
	if data.Text == "" {
		t.Error("Translation returned empty text")
	}
}
