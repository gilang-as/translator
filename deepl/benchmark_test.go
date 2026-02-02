package deepl

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

func BenchmarkNewWithOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(
			WithHost("api.deepl.com"),
			WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
			WithProxyURL("http://proxy:8080"),
			WithDLSession("test-session"),
		)
	}
}

func BenchmarkHost(b *testing.B) {
	d := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Host()
	}
}

func BenchmarkSetHost(b *testing.B) {
	d := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.SetHost("api.deepl.com")
	}
}

func BenchmarkProxyURL(b *testing.B) {
	d := New(WithProxyURL("http://proxy:8080"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.ProxyURL()
	}
}

func BenchmarkSetProxyURL(b *testing.B) {
	d := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.SetProxyURL("http://proxy:8080")
	}
}

func BenchmarkDLSession(b *testing.B) {
	d := New(WithDLSession("test-session"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.DLSession()
	}
}

func BenchmarkSetDLSession(b *testing.B) {
	d := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.SetDLSession("test-session")
	}
}

func BenchmarkHostConcurrent(b *testing.B) {
	d := New()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = d.Host()
		}
	})
}

func BenchmarkSetHostConcurrent(b *testing.B) {
	d := New()
	hosts := []string{"www2.deepl.com", "api.deepl.com", "deepl.com"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			d.SetHost(hosts[i%len(hosts)])
			i++
		}
	})
}

func BenchmarkTranslate(b *testing.B) {
	d := New()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.Translate(ctx, "Hello", "en", "id")
	}
}

func BenchmarkTranslateConcurrent(b *testing.B) {
	d := New()
	ctx := context.Background()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = d.Translate(ctx, "Hello", "en", "id")
		}
	})
}
