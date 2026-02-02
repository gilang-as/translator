package gt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.gilang.dev/translator/v2/params"
)

func TestTranslateWithParam(t *testing.T) {
	value := params.Translate{
		Text: "Halo Dunia",
		To:   params.ENGLISH,
	}
	translated, err := TranslateWithParam(context.Background(), value)
	if err != nil {
		t.Error(err)
		return
	}
	if translated.Text != "" {
		t.Log(translated)
	}
}

func TestTranslate(t *testing.T) {
	translated, err := Translate(context.Background(), "Hello World", params.INDONESIAN)
	if err != nil {
		t.Error(err)
		return
	}
	if translated.Text != "" {
		t.Log(translated)
	}
}

func TestManualTranslate(t *testing.T) {
	translated, err := ManualTranslate(context.Background(), "Halo Semuanya", params.INDONESIAN, params.JAVANESE)
	if err != nil {
		t.Error(err)
		return
	}
	if translated.Text != "" {
		t.Log(translated)
	}
}

func TestTranslateWithParam2(t *testing.T) {
	value := params.Translate{
		Text: "这是第一句话。 这是第二句话。",
		From: "zh-cn",
		To:   "en",
	}

	translated, err := TranslateWithParam(context.Background(), value)
	assert.NoError(t, err, "should not return error")

	expected := "This is the first sentence. This is the second sentence."
	assert.Equal(t, expected, translated.Text)
}

func TestUseGoogle(t *testing.T) {
	UseGoogle()
	translated, err := Translate(context.Background(), "Hello", params.INDONESIAN)
	if err != nil {
		t.Error(err)
		return
	}
	assert.NotEmpty(t, translated.Text)
}

func TestUseDeepL(t *testing.T) {
	UseDeepL()
	translated, err := Translate(context.Background(), "Hello", params.INDONESIAN)
	if err != nil {
		t.Skipf("Skipping due to API error: %v", err)
		return
	}
	assert.NotEmpty(t, translated.Text)
	// Reset to Google for other tests
	UseGoogle()
}

func TestNewGoogleTranslator(t *testing.T) {
	translator := NewGoogleTranslator()
	assert.NotNil(t, translator)
}

func TestNewDeepLTranslator(t *testing.T) {
	translator := NewDeepLTranslator()
	assert.NotNil(t, translator)
}

func TestTranslateWith(t *testing.T) {
	google := NewGoogleTranslator()
	translated, err := TranslateWith(context.Background(), google, "Hello", params.INDONESIAN)
	if err != nil {
		t.Error(err)
		return
	}
	assert.NotEmpty(t, translated.Text)
}

func TestTranslateWithDeepL(t *testing.T) {
	deepl := NewDeepLTranslator()
	translated, err := TranslateWith(context.Background(), deepl, "Hello", params.INDONESIAN)
	if err != nil {
		t.Skipf("Skipping due to API error: %v", err)
		return
	}
	assert.NotEmpty(t, translated.Text)
}

func TestSetDefaultTranslator(t *testing.T) {
	original := GetDefaultTranslator()
	defer SetDefaultTranslator(original)

	newTranslator := NewGoogleTranslator()
	SetDefaultTranslator(newTranslator)

	assert.Equal(t, newTranslator, GetDefaultTranslator())
}
