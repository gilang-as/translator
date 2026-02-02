// Example: Google Translate with custom configuration
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	gt "gopkg.gilang.dev/translator/v2"
	"gopkg.gilang.dev/translator/v2/googletranslate"
	"gopkg.gilang.dev/translator/v2/params"
)

func main() {
	ctx := context.Background()

	// Create Google Translate client with custom options
	fmt.Println("=== Google Translate with Custom Config ===")
	client := googletranslate.New(
		googletranslate.WithHost("google.co.id"),
		googletranslate.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
	)

	result, err := client.Translate(ctx, "Good morning", "en", "id")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Translated: %s\n", result.Text)
	if result.Pronunciation != nil {
		fmt.Printf("Pronunciation: %s\n", *result.Pronunciation)
	}
	fmt.Println()

	// Using the adapter through package-level function
	fmt.Println("=== Using Google Translator Adapter ===")
	translator := gt.NewGoogleTranslator(
		googletranslate.WithHost("google.com"),
	)

	result2, err := gt.TranslateWith(ctx, translator, "How are you?", params.FRENCH)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("English -> French: %s\n", result2.Text)
	fmt.Println()

	// Set Google as default with custom config
	fmt.Println("=== Set Google as Default ===")
	gt.UseGoogle(
		googletranslate.WithHost("google.co.jp"),
	)

	result3, err := gt.Translate(ctx, "Thank you", params.JAPANESE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("English -> Japanese: %s\n", result3.Text)
}
