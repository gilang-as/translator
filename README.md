# translator

[![Actions Status](https://github.com/gilang-as/google-translate/actions/workflows/test.yaml/badge.svg)](https://github.com/gilang-as/google-translate/actions)

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
go get gopkg.gilang.dev/translator
```

## Quick Start

```go
import gt "gopkg.gilang.dev/translator"

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
import "gopkg.gilang.dev/translator/googletranslate"

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
import "gopkg.gilang.dev/translator/deepl"

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
import "gopkg.gilang.dev/translator/params"

params.ENGLISH     // "en"
params.INDONESIAN  // "id"
params.JAPANESE    // "ja"
params.FRENCH      // "fr"
params.GERMAN      // "de"
params.SPANISH     // "es"
params.CHINESE     // "zh-cn"
```

Use `"auto"` for automatic language detection.

## Notes

- **DeepL rate limiting**: DeepL may temporarily block IPs that make too many requests. Handle errors appropriately.
- **Context support**: All translation functions require `context.Context` for cancellation and timeout control.
- **Concurrency**: Both clients are safe for concurrent use.

## Credits

Parts of the code are ported from [gtranslate](https://github.com/bregydoc/gtranslate) and [google-translate-api](https://github.com/matheuss/google-translate-api) (MIT license).

## License

MIT © [Gilang Adi S](https://github.com/gilang-as)
