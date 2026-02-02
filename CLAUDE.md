# CLAUDE.md

## Project

Go library for Google Translate API wrapper. Module: `gopkg.gilang.dev/translator`

## Commands

```bash
go test -v -run TestTranslator   # Run tests
go mod tidy                       # Tidy dependencies
```

## Structure

- `translate.go` - Public API functions
- `api.go` - HTTP request/response handling
- `params/` - Parameter structs and language definitions

## API

```go
gt.Translate(text, toLanguage string)                        // Auto-detect source
gt.TranslateWithParam(params.Translate{Text, From, To})      // Struct-based
gt.ManualTranslate(text, fromLanguage, toLanguage string)    // Explicit languages
```

## Conventions

- Package imported as `gt`
- Language codes: ISO 639-1 (e.g., "en", "fr", "id")
- Use `"auto"` for automatic language detection
