# CLAUDE.md

## Project

Go library for translation API wrapper supporting multiple backends (Google Translate, DeepL). Module: `gopkg.gilang.dev/translator`

## Commands

```bash
go test -v ./...           # Run tests
go test -bench=. ./...     # Run benchmarks
go mod tidy                # Tidy dependencies
```

## Structure

- `translate.go` - Public API with multi-translator support
- `googletranslate/` - Google Translate client (concurrency-safe)
- `deepl/` - DeepL client (concurrency-safe)
- `params/` - Parameter structs and language definitions
- `example/` - Usage examples

## API

### Package-level functions (use default translator)

```go
gt.Translate(ctx, text, toLanguage)                           // Auto-detect source
gt.TranslateWithParam(ctx, params.Translate{Text, From, To})  // Struct-based
gt.ManualTranslate(ctx, text, fromLanguage, toLanguage)       // Explicit languages
```

### Switch default translator

```go
gt.UseGoogle()                                    // Use Google (default)
gt.UseDeepL()                                     // Use DeepL
gt.UseGoogle(googletranslate.WithHost("..."))    // With options
gt.UseDeepL(deepl.WithDLSession("..."))          // With options
```

### Use specific translator

```go
google := gt.NewGoogleTranslator(opts...)
deepl := gt.NewDeepLTranslator(opts...)
gt.TranslateWith(ctx, google, text, toLanguage)
```

### Google Translate client

```go
import "gopkg.gilang.dev/translator/googletranslate"

client := googletranslate.New(
    googletranslate.WithHost("google.co.id"),
    googletranslate.WithHTTPClient(&http.Client{Timeout: 30*time.Second}),
    googletranslate.WithProxyURL("http://proxy:8080"),
)
client.Translate(ctx, text, from, to)

// Runtime configuration (concurrency-safe)
client.SetHost("google.com")
client.SetProxyURL("http://proxy:8080")
```

### DeepL client

```go
import "gopkg.gilang.dev/translator/deepl"

client := deepl.New(
    deepl.WithProxyURL("http://proxy:8080"),
    deepl.WithDLSession("session-token"),  // For Pro features
)
client.Translate(ctx, text, from, to)

// Runtime configuration (concurrency-safe)
client.SetProxyURL("http://proxy:8080")
client.SetDLSession("new-session")
```

## Conventions

- Package imported as `gt`
- Language codes: ISO 639-1 (e.g., "en", "fr", "id")
- Use `"auto"` for automatic language detection
- All functions require `context.Context` as first parameter
- Both clients are concurrency-safe with mutex protection
- DeepL may rate limit requests; handle errors appropriately
