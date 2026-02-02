package googletranslate

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
			WithHost("google.co.id"),
			WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
		)
	}
}

func BenchmarkHost(b *testing.B) {
	gt := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gt.Host()
	}
}

func BenchmarkSetHost(b *testing.B) {
	gt := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gt.SetHost("google.com")
	}
}

func BenchmarkHostConcurrent(b *testing.B) {
	gt := New()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = gt.Host()
		}
	})
}

func BenchmarkSetHostConcurrent(b *testing.B) {
	gt := New()
	hosts := []string{"google.com", "google.co.id", "google.co.uk", "google.co.jp"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			gt.SetHost(hosts[i%len(hosts)])
			i++
		}
	})
}

func BenchmarkProxyURL(b *testing.B) {
	gt := New(WithProxyURL("http://proxy:8080"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gt.ProxyURL()
	}
}

func BenchmarkSetProxyURL(b *testing.B) {
	gt := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gt.SetProxyURL("http://proxy:8080")
	}
}

func BenchmarkTranslate(b *testing.B) {
	gt := New()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = gt.Translate(ctx, "Hello", "en", "id")
	}
}

func BenchmarkTranslateConcurrent(b *testing.B) {
	gt := New()
	ctx := context.Background()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = gt.Translate(ctx, "Hello", "en", "id")
		}
	})
}
