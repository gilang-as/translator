# translator v2

[![Actions Status](https://github.com/gilang-as/translator/actions/workflows/test.yaml/badge.svg)](https://github.com/gilang-as/translator/actions)

> [!NOTE]
> This repository hosts **v2** (multi-backend). The legacy v1 Google Translate-only module is **deprecated**; use v2 for all new development.

A **free** and **unlimited** translation API supporting multiple backends (Google Translate, DeepL).

## Features

- **Multiple backends** - Google Translate and DeepL support
- **Auto language detection** - Automatically detect source language
- **Spelling & language correction** - Built-in correction suggestions
- **Concurrency-safe** - Thread-safe clients with mutex protection
- **Proxy support** - Configure proxy for both backends
- **Runtime configuration** - Change settings without recreating clients

## Requirements

- Go 1.21 or later

## Install

```bash
go get gopkg.gilang.dev/translator/v2
```

## Migration from v1

If you are upgrading from v1:

1. Update imports to use the v2 module path:
    - `gopkg.gilang.dev/translator/v2`
    - `gopkg.gilang.dev/translator/v2/googletranslate`
    - `gopkg.gilang.dev/translator/v2/deepl`
    - `gopkg.gilang.dev/translator/v2/params`
2. Update your module requirement:

```bash
go get gopkg.gilang.dev/translator/v2@latest
```

3. Replace any legacy Google Translate-only usage with the v2 `googletranslate` subpackage shown below.

## Quick Start

```go
import gt "gopkg.gilang.dev/translator/v2"

ctx := context.Background()

// Translate with Google (default)
result, _ := gt.Translate(ctx, "Hello World", "id")
fmt.Println(result.Text) // "Halo Dunia"

// Switch to DeepL
gt.UseDeepL()
result, _ = gt.Translate(ctx, "Hello World", "id")
fmt.Println(result.Text)
```

## API

### Package-level Functions

```go
// Auto-detect source language
gt.Translate(ctx, "Hello", "fr")

// Explicit source and target languages
gt.ManualTranslate(ctx, "Hello", "en", "fr")

// Using params struct
gt.TranslateWithParam(ctx, params.Translate{
    Text: "Hello",
    From: "en",  // optional
    To:   "fr",
})
```

### Switch Default Translator

```go
gt.UseGoogle()  // Use Google Translate (default)
gt.UseDeepL()   // Use DeepL

// With custom options
gt.UseGoogle(googletranslate.WithHost("google.co.id"))
gt.UseDeepL(deepl.WithDLSession("session-token"))
```

### Use Specific Translator

```go
google := gt.NewGoogleTranslator()
deepl := gt.NewDeepLTranslator()

// Translate with specific translator
gt.TranslateWith(ctx, google, "Hello", "fr")
gt.TranslateWith(ctx, deepl, "Hello", "fr")
```

### Google Translate Client

```go
import "gopkg.gilang.dev/translator/v2/googletranslate"

client := googletranslate.New(
    googletranslate.WithHost("google.co.id"),
    googletranslate.WithHTTPClient(&http.Client{Timeout: 30*time.Second}),
    googletranslate.WithProxyURL("http://proxy:8080"),
)

result, err := client.Translate(ctx, "Hello", "en", "id")

// Runtime configuration
client.SetHost("google.com")
client.SetProxyURL("http://proxy:8080")
```

### DeepL Client

```go
import "gopkg.gilang.dev/translator/v2/deepl"

client := deepl.New(
    deepl.WithProxyURL("http://proxy:8080"),
    deepl.WithDLSession("session-token"),  // For Pro features
)

result, err := client.Translate(ctx, "Hello", "en", "id")
fmt.Println(result.Text)
fmt.Println(result.Alternatives)  // Alternative translations
fmt.Println(result.Method)        // "Free" or "Pro"

// Runtime configuration
client.SetProxyURL("http://proxy:8080")
client.SetDLSession("new-session")
```

## Response

The `Translated` struct contains:

| Field | Type | Description |
|-------|------|-------------|
| `Text` | string | The translated text |
| `Pronunciation` | *string | Pronunciation (Google only) |
| `Alternatives` | []string | Alternative translations (DeepL only) |
| `Method` | string | "Free" or "Pro" (DeepL only) |
| `From.Language.Iso` | string | Detected source language code |
| `From.Language.DidYouMean` | bool | Language correction suggested |
| `From.Text.AutoCorrected` | bool | Text was auto-corrected |
| `From.Text.Value` | *string | Corrected text value |
| `From.Text.DidYouMean` | bool | Text correction suggested |

## Examples

See the [example](./example) directory for complete examples:

- [basic](./example/basic) - Basic usage patterns
- [google](./example/google) - Google Translate configuration
- [deepl](./example/deepl) - DeepL configuration
- [multi](./example/multi) - Using multiple translators

## Language Codes

Use ISO 639-1 language codes (e.g., "en", "fr", "id", "ja"). Common codes are available in the `params` package:

```go
import "gopkg.gilang.dev/translator/v2/params"

params.ENGLISH     // "en"
params.INDONESIAN  // "id"
params.JAPANESE    // "ja"
params.FRENCH      // "fr"
params.GERMAN      // "de"
params.SPANISH     // "es"
params.CHINESE     // "zh-cn"
```

Use `"auto"` for automatic language detection.

## Benchmarks

Run benchmarks with:

```bash
go test -bench=. -benchmem ./...
```

### Google Translate

```
BenchmarkNew-12                     20029788        58.15 ns/op      112 B/op       2 allocs/op
BenchmarkNewWithOptions-12           8231678       144.0 ns/op       200 B/op       5 allocs/op
BenchmarkHost-12                    85322860        13.76 ns/op        0 B/op       0 allocs/op
BenchmarkSetHost-12                 47107504        25.30 ns/op        0 B/op       0 allocs/op
BenchmarkHostConcurrent-12          26058314        40.53 ns/op        0 B/op       0 allocs/op
BenchmarkSetHostConcurrent-12       16455430        76.69 ns/op        0 B/op       0 allocs/op
BenchmarkProxyURL-12                92660499        13.27 ns/op        0 B/op       0 allocs/op
BenchmarkSetProxyURL-12             45487726        25.00 ns/op        0 B/op       0 allocs/op
BenchmarkTranslate-12                      3    371242869 ns/op  20632666 B/op    3565 allocs/op
BenchmarkTranslateConcurrent-12            4    267211196 ns/op  20572356 B/op    3480 allocs/op
```

### DeepL

```
BenchmarkNew-12                     18704745        62.26 ns/op      128 B/op       2 allocs/op
BenchmarkNewWithOptions-12           5821923       216.1 ns/op       264 B/op       7 allocs/op
BenchmarkHost-12                    91805458        13.63 ns/op        0 B/op       0 allocs/op
BenchmarkSetHost-12                 49368652        25.03 ns/op        0 B/op       0 allocs/op
BenchmarkProxyURL-12                89641826        13.24 ns/op        0 B/op       0 allocs/op
BenchmarkSetProxyURL-12             49028032        25.49 ns/op        0 B/op       0 allocs/op
BenchmarkDLSession-12               90087140        13.66 ns/op        0 B/op       0 allocs/op
BenchmarkSetDLSession-12            48464596        24.78 ns/op        0 B/op       0 allocs/op
BenchmarkHostConcurrent-12          28143388        47.32 ns/op        0 B/op       0 allocs/op
BenchmarkSetHostConcurrent-12       16936105        75.13 ns/op        0 B/op       0 allocs/op
BenchmarkTranslate-12                      2    908235471 ns/op    211928 B/op    1111 allocs/op
BenchmarkTranslateConcurrent-12           14    114018568 ns/op    176837 B/op     994 allocs/op
```

## Notes

- **DeepL rate limiting**: DeepL may temporarily block IPs that make too many requests. Handle errors appropriately.
- **Context support**: All translation functions require `context.Context` for cancellation and timeout control.
- **Concurrency**: Both clients are safe for concurrent use.

## Credits

Parts of the code are ported from [gtranslate](https://github.com/bregydoc/gtranslate) and [google-translate-api](https://github.com/matheuss/google-translate-api) (MIT license).

## License

MIT © [Gilang Adi S](https://github.com/gilang-as)
