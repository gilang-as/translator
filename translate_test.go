package gt

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.gilang.dev/translator/v2/googletranslate"
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

type googleTranslateFunc func(context.Context, string, string, string, Gender) (*googletranslate.Translated, error)

func (f googleTranslateFunc) TranslateWithGender(ctx context.Context, text, from, to string, gender Gender) (*googletranslate.Translated, error) {
	return f(ctx, text, from, to, gender)
}

func TestGoogleTranslatorGenderPreferences(t *testing.T) {
	for _, tt := range []struct {
		name   string
		gender Gender
		plain  bool
	}{
		{"plain translation", GenderUnknown, true},
		{"unknown", GenderUnknown, false},
		{"feminine", GenderFeminine, false},
		{"masculine", GenderMasculine, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			alternatives := []GenderAlternative{{Text: "la autora", Gender: GenderFeminine}}
			var translator Translator = &googleTranslateAdapter{client: googleTranslateFunc(func(ctx context.Context, text, from, to string, gender Gender) (*googletranslate.Translated, error) {
				if text != "The author" || from != "auto" || to != "es" || gender != tt.gender {
					t.Fatalf("unexpected request: %q, %q, %q, %q", text, from, to, gender)
				}
				return &googletranslate.Translated{
					Text: "la autora", GenderAlternatives: alternatives,
					From: googletranslate.TranslateFrom{Language: googletranslate.TranslateFromLanguage{Iso: "en"}},
				}, nil
			})}
			var result *Translated
			var err error
			if tt.plain {
				result, err = translator.Translate(context.Background(), "The author", "auto", "es")
			} else {
				result, err = translator.TranslateWithGender(context.Background(), "The author", "auto", "es", tt.gender)
			}
			if err != nil || result == nil || result.Text != "la autora" || result.From.Language.Iso != "en" ||
				!reflect.DeepEqual(result.GenderAlternatives, alternatives) {
				t.Fatalf("unexpected translation: %+v, %v", result, err)
			}
		})
	}
	sentinel := errors.New("translation failed")
	var translator Translator = &googleTranslateAdapter{client: googleTranslateFunc(func(context.Context, string, string, string, Gender) (*googletranslate.Translated, error) {
		return nil, sentinel
	})}
	if result, err := translator.TranslateWithGender(context.Background(), "text", "en", "es", GenderFeminine); result != nil || !errors.Is(err, sentinel) {
		t.Fatalf("expected the provider error without a result: %+v, %v", result, err)
	}
}

func TestDeepLTranslatorGenderPreferences(t *testing.T) {
	translator := NewDeepLTranslator()
	for _, tt := range []struct {
		name        string
		gender      Gender
		unsupported bool
	}{
		{"masculine", GenderMasculine, true},
		{"feminine", GenderFeminine, true},
		{"invalid", Gender("invalid"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := translator.TranslateWithGender(context.Background(), "The author", "en", "es", tt.gender)
			if result != nil || err == nil {
				t.Fatalf("expected an explicit error without a result: %+v, %v", result, err)
			}
			if tt.unsupported {
				if !errors.Is(err, ErrUnsupportedOption) || !strings.Contains(err.Error(), "DeepL") || !strings.Contains(err.Error(), string(tt.gender)) {
					t.Fatalf("expected an unsupported gender preference error, got %v", err)
				}
			} else if errors.Is(err, ErrUnsupportedOption) || !strings.Contains(err.Error(), "invalid translation gender") {
				t.Fatalf("invalid values must be distinguished from unsupported preferences: %v", err)
			}
		})
	}
}
