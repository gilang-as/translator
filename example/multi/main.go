// Example: Using multiple translators together
package main

import (
	"context"
	"fmt"
	"log"

	gt "gopkg.gilang.dev/translator/v2"
	"gopkg.gilang.dev/translator/v2/deepl"
	"gopkg.gilang.dev/translator/v2/googletranslate"
	"gopkg.gilang.dev/translator/v2/params"
)

func main() {
	ctx := context.Background()

	// Create both translators
	google := gt.NewGoogleTranslator(
		googletranslate.WithHost("google.com"),
	)
	deeplClient := gt.NewDeepLTranslator(
		deepl.WithProxyURL(""), // No proxy
	)

	text := "The quick brown fox jumps over the lazy dog"
	targetLang := params.INDONESIAN

	fmt.Printf("Original: %s\n", text)
	fmt.Printf("Target: %s\n\n", targetLang)

	// Translate with Google
	fmt.Println("=== Google Translate ===")
	googleResult, err := gt.TranslateWith(ctx, google, text, targetLang)
	if err != nil {
		log.Printf("Google error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", googleResult.Text)
		fmt.Printf("Detected: %s\n\n", googleResult.From.Language.Iso)
	}

	// Translate with DeepL
	fmt.Println("=== DeepL Translate ===")
	deeplResult, err := gt.TranslateWith(ctx, deeplClient, text, targetLang)
	if err != nil {
		log.Printf("DeepL error: %v\n", err)
		fmt.Println("(DeepL may rate limit requests)")
	} else {
		fmt.Printf("Result: %s\n", deeplResult.Text)
		fmt.Printf("Method: %s\n", deeplResult.Method)
		if len(deeplResult.Alternatives) > 0 {
			fmt.Printf("Alternatives: %v\n", deeplResult.Alternatives)
		}
	}
	fmt.Println()

	// Demonstrate switching default translator
	fmt.Println("=== Switching Default Translator ===")

	// Use Google (default)
	gt.UseGoogle()
	r1, _ := gt.Translate(ctx, "Hello", params.SPANISH)
	fmt.Printf("Google (default): Hello -> %s\n", r1.Text)

	// Switch to DeepL
	gt.UseDeepL()
	r2, err := gt.Translate(ctx, "Hello", params.SPANISH)
	if err != nil {
		fmt.Printf("DeepL (default): error - %v\n", err)
	} else {
		fmt.Printf("DeepL (default): Hello -> %s\n", r2.Text)
	}

	// Switch back to Google
	gt.UseGoogle()
	r3, _ := gt.Translate(ctx, "Goodbye", params.SPANISH)
	fmt.Printf("Google (default): Goodbye -> %s\n", r3.Text)
}
