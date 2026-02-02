// Example: Basic usage with default translator (Google Translate)
package main

import (
	"context"
	"fmt"
	"log"

	gt "gopkg.gilang.dev/translator/v2"
	"gopkg.gilang.dev/translator/v2/params"
)

func main() {
	ctx := context.Background()

	// Simple translation with auto-detect source language
	fmt.Println("=== Simple Translation ===")
	result, err := gt.Translate(ctx, "Hello World", params.INDONESIAN)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Original: Hello World\n")
	fmt.Printf("Translated: %s\n", result.Text)
	fmt.Printf("Detected language: %s\n\n", result.From.Language.Iso)

	// Translation with explicit source and target languages
	fmt.Println("=== Manual Translation ===")
	result, err = gt.ManualTranslate(ctx, "Selamat pagi", params.INDONESIAN, params.JAPANESE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Original: Selamat pagi\n")
	fmt.Printf("Translated: %s\n\n", result.Text)

	// Translation with params struct
	fmt.Println("=== Translation with Params ===")
	p := params.Translate{
		Text: "这是第一句话。这是第二句话。",
		From: params.CHINESE,
		To:   params.ENGLISH,
	}
	result, err = gt.TranslateWithParam(ctx, p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Original: %s\n", p.Text)
	fmt.Printf("Translated: %s\n", result.Text)
}
