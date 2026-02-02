// Example: Quick start with translator library
//
// For more examples, see:
//   - example/basic/   - Basic usage
//   - example/google/  - Google Translate configuration
//   - example/deepl/   - DeepL Translate configuration
//   - example/multi/   - Using multiple translators
package main

import (
	"context"
	"fmt"
	"log"

	gt "gopkg.gilang.dev/translator"
	"gopkg.gilang.dev/translator/params"
)

func main() {
	ctx := context.Background()

	// Simple translation (uses Google Translate by default)
	result, err := gt.Translate(ctx, "Hello World", params.INDONESIAN)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Translated: %s\n", result.Text)

	// Switch to DeepL
	gt.UseDeepL()
	result, err = gt.Translate(ctx, "Hello World", params.INDONESIAN)
	if err != nil {
		fmt.Printf("DeepL error: %v\n", err)
	} else {
		fmt.Printf("DeepL: %s\n", result.Text)
	}

	// Use specific translator
	google := gt.NewGoogleTranslator()
	result, err = gt.TranslateWith(ctx, google, "Good morning", params.JAPANESE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Google: %s\n", result.Text)
}
