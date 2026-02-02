# ⚠️ Deprecated

> This module has been renamed to:
> 
> 👉 https://gopkg.gilang.dev/translator
> 
> Last maintained version:
> - v1.x.x (branch: v1)
> 
> Please migrate to the new module.

# google-translate (Legacy)
[![Actions Status](https://github.com/gilang-as/google-translate/actions/workflows/test.yaml/badge.svg)](https://github.com/gilang-as/google-translate/actions)

A **free** and **unlimited** API for Google Translate

Parts of the code are ported from [gtranslate](https://github.com/bregydoc/gtranslate) and [google-translate-api](https://github.com/matheuss/google-translate-api) (also MIT license).

> **Note:** This is the legacy `v1` branch. This module has been renamed to `gopkg.gilang.dev/translator`.
> For the latest version, see the `main` branch.

## Features
- Auto language detection
- Spelling correction
- Language correction
- Fast and reliable – it uses the same servers that [translate.google.com](https://translate.google.com) uses

## Requirements

- Go 1.25 or later

## Install

```bash
go get gopkg.gilang.dev/google-translate
```

## Usage

### Quick Start

```go
import gt "gopkg.gilang.dev/google-translate"

// Translate with auto-detection of source language
translated, err := gt.Translate("Hello World", "fr")
fmt.Println(translated.Text) // "Bonjour le monde"
```

### API Functions

#### Translate
Auto-detect source language and translate to target language.
```go
translated, err := gt.Translate("Hello", "es")
```

#### ManualTranslate
Specify both source and target languages explicitly.
```go
translated, err := gt.ManualTranslate("Hello", "en", "fr")
```

#### TranslateWithParam
Use a struct for full control over parameters.
```go
import "gopkg.gilang.dev/google-translate/params"

value := params.Translate{
    Text: "Halo Dunia",
    From: "id",  // optional, defaults to auto-detect
    To:   "en",
}
translated, err := gt.TranslateWithParam(value)
```

### Full Example

```go
package main

import (
	"encoding/json"
	"fmt"

	gt "gopkg.gilang.dev/google-translate"
	"gopkg.gilang.dev/google-translate/params"
)

func main() {
	value := params.Translate{
		Text: "Halo Dunia",
		To:   "en",
	}
	translated, err := gt.TranslateWithParam(value)
	if err != nil {
		panic(err)
	}
	prettyJSON, _ := json.MarshalIndent(translated, "", "\t")
	fmt.Println(string(prettyJSON))
}
```

### Returns an `object`:
- `text` *(string)* – The translated text.
- `pronunciation` *(string)* – The Pronunciation text.
- `from` *(object)*
    - `language` *(object)*
        - `did_you_mean` *(boolean)* - `true` if the API suggest a correction in the source language
        - `iso` *(string)* - The code of the language that the API has recognized in the `text`
    - `text` *(object)*
        - `auto_corrected` *(boolean)* – `true` if the API has auto corrected the `text`
        - `value` *(string)* – The auto corrected `text` or the `text` with suggested corrections
        - `did_you_mean` *(boolean)* – `true` if the API has suggested corrections to the `text`

## Migrating to the new module

To upgrade to the new module, update your imports:

```diff
- import gt "gopkg.gilang.dev/google-translate"
- import "gopkg.gilang.dev/google-translate/params"
+ import gt "gopkg.gilang.dev/translator"
+ import "gopkg.gilang.dev/translator/params"
```

## License

MIT © [Gilang Adi S](https://github.com/gilang-as)
