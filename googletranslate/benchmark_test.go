package googletranslate

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func BenchmarkNew(b *testing.B) {
	for b.Loop() {
		_ = New()
	}
}

func BenchmarkNewWithOptions(b *testing.B) {
	for b.Loop() {
		_ = New(
			WithHost("google.co.id"),
			WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
		)
	}
}

func BenchmarkHost(b *testing.B) {
	gt := New()
	b.ResetTimer()
	for b.Loop() {
		_ = gt.Host()
	}
}

func BenchmarkSetHost(b *testing.B) {
	gt := New()
	b.ResetTimer()
	for b.Loop() {
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

func BenchmarkTranslate(b *testing.B) {
	gt := New()
	ctx := context.Background()
	b.ResetTimer()
	for b.Loop() {
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
