// Example: DeepL Translate with custom configuration
package main

import (
	"context"
	"fmt"
	"log"

	gt "gopkg.gilang.dev/translator"
	"gopkg.gilang.dev/translator/deepl"
	"gopkg.gilang.dev/translator/params"
)

func main() {
	ctx := context.Background()

	// Create DeepL client with custom options
	fmt.Println("=== DeepL Translate ===")
	client := deepl.New(
		// Optional: Set proxy URL
		// deepl.WithProxyURL("http://proxy:8080"),
		// Optional: Set DeepL session for Pro features
		// deepl.WithDLSession("your-session-token"),
	)

	result, err := client.Translate(ctx, "Hello, how are you?", "en", "id")
	if err != nil {
		log.Printf("DeepL error: %v\n", err)
		fmt.Println("Note: DeepL may rate limit requests. Try again later.")
		return
	}
	fmt.Printf("Translated: %s\n", result.Text)
	fmt.Printf("Method: %s\n", result.Method)
	if len(result.Alternatives) > 0 {
		fmt.Println("Alternatives:")
		for i, alt := range result.Alternatives {
			fmt.Printf("  %d. %s\n", i+1, alt)
		}
	}
	fmt.Println()

	// Using the adapter through package-level function
	fmt.Println("=== Using DeepL Translator Adapter ===")
	translator := gt.NewDeepLTranslator()

	result2, err := gt.TranslateWith(ctx, translator, "Good evening", params.GERMAN)
	if err != nil {
		log.Printf("DeepL error: %v\n", err)
		return
	}
	fmt.Printf("English -> German: %s\n", result2.Text)
	fmt.Println()

	// Set DeepL as default
	fmt.Println("=== Set DeepL as Default ===")
	gt.UseDeepL()

	result3, err := gt.Translate(ctx, "Welcome", params.FRENCH)
	if err != nil {
		log.Printf("DeepL error: %v\n", err)
		return
	}
	fmt.Printf("English -> French: %s\n", result3.Text)
}
