# CLAUDE.md

## Project

Go library for Google Translate API wrapper. Module: `gopkg.gilang.dev/translator`

## Commands

```bash
go test -v ./...    # Run tests
go mod tidy         # Tidy dependencies
```

## Structure

- `translate.go` - Public API functions (package-level)
- `gtranslate.go` - Re-exports for backward compatibility
- `googletranslate/` - GoogleTranslate client (concurrency-safe)
- `params/` - Parameter structs and language definitions

## API

### Package-level functions (use default client)

```go
gt.Translate(ctx, text, toLanguage)                        // Auto-detect source
gt.TranslateWithParam(ctx, params.Translate{Text, From, To})  // Struct-based
gt.ManualTranslate(ctx, text, fromLanguage, toLanguage)    // Explicit languages
```

### Custom client with googletranslate package (recommended)

```go
import "gopkg.gilang.dev/translator/googletranslate"

client := googletranslate.New(
    googletranslate.WithHost("google.co.id"),
    googletranslate.WithHTTPClient(&http.Client{Timeout: 30*time.Second}),
)
client.Translate(ctx, text, from, to)

// Runtime configuration (concurrency-safe)
client.SetHost("google.com")
client.SetClient(customHTTPClient)
```

### Legacy API (backward compatible)

```go
client := gt.NewGoogleTranslate(gt.WithHost("google.co.id"))
client.Translate(ctx, text, from, to)
```

## Conventions

- Package imported as `gt`
- Language codes: ISO 639-1 (e.g., "en", "fr", "id")
- Use `"auto"` for automatic language detection
- All functions require `context.Context` as first parameter
- GoogleTranslate client is concurrency-safe
