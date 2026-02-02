# CLAUDE.md

## Project

Go library for Google Translate API wrapper. Module: `gopkg.gilang.dev/translator`

## Commands

```bash
go test -v ./...    # Run tests
go mod tidy         # Tidy dependencies
```

## Structure

- `translate.go` - Public API functions
- `gtranslate.go` - GoogleTranslate client and HTTP handling
- `params/` - Parameter structs and language definitions

## API

### Package-level functions (use default client)

```go
gt.Translate(ctx, text, toLanguage)                        // Auto-detect source
gt.TranslateWithParam(ctx, params.Translate{Text, From, To})  // Struct-based
gt.ManualTranslate(ctx, text, fromLanguage, toLanguage)    // Explicit languages
```

### Custom client with configuration

```go
client := gt.NewGoogleTranslate(
    gt.WithHost("google.co.id"),                           // Custom host
    gt.WithHTTPClient(&http.Client{Timeout: 30*time.Second}), // Custom HTTP client
)
client.Translate(ctx, text, from, to)
```

## Conventions

- Package imported as `gt`
- Language codes: ISO 639-1 (e.g., "en", "fr", "id")
- Use `"auto"` for automatic language detection
- All functions require `context.Context` as first parameter
