package gt

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/text/language"
	"gopkg.gilang.dev/translator/v2/deepl"
	"gopkg.gilang.dev/translator/v2/googletranslate"
	"gopkg.gilang.dev/translator/v2/params"
)

// Gender identifies a gendered translation alternative.
type Gender = params.Gender

const (
	GenderUnknown   = params.GenderUnknown
	GenderFeminine  = params.GenderFeminine
	GenderMasculine = params.GenderMasculine
)

// GenderAlternative is a translation offered for a particular gender.
type GenderAlternative = params.GenderAlternative

// ErrUnsupportedOption reports a translation preference a provider cannot honor.
var ErrUnsupportedOption = params.ErrUnsupportedOption

// Translator is the common interface for all translation clients.
type Translator interface {
	Translate(ctx context.Context, text, from, to string) (*Translated, error)
	// TranslateWithGender returns ErrUnsupportedOption when gender selection is unsupported.
	TranslateWithGender(ctx context.Context, text, from, to string, gender Gender) (*Translated, error)
}

// Translated represents a translation result.
type Translated struct {
	Text               string              `json:"text"`
	Pronunciation      *string             `json:"pronunciation"`
	Alternatives       []string            `json:"alternatives,omitempty"`
	GenderAlternatives []GenderAlternative `json:"gender_alternatives,omitempty"`
	From               TranslateFrom       `json:"from"`
	Method             string              `json:"method,omitempty"`
}

// TranslateFrom contains source language and text information.
type TranslateFrom struct {
	Language TranslateFromLanguage `json:"language"`
	Text     TranslateFromText     `json:"text"`
}

// TranslateFromLanguage contains detected language information.
type TranslateFromLanguage struct {
	DidYouMean bool   `json:"did_you_mean"`
	Iso        string `json:"iso"`
}

// TranslateFromText contains text correction information.
type TranslateFromText struct {
	AutoCorrected bool    `json:"auto_corrected"`
	Value         *string `json:"value"`
	DidYouMean    bool    `json:"did_you_mean"`
}

// TranslatorType represents the type of translator to use.
type TranslatorType string

const (
	Google TranslatorType = "google"
	DeepL  TranslatorType = "deepl"
)

// googleTranslateAdapter wraps googletranslate.GoogleTranslate to implement Translator.
type googleTranslateAdapter struct {
	client interface {
		TranslateWithGender(context.Context, string, string, string, Gender) (*googletranslate.Translated, error)
	}
}

func (g *googleTranslateAdapter) Translate(ctx context.Context, text, from, to string) (*Translated, error) {
	return g.TranslateWithGender(ctx, text, from, to, GenderUnknown)
}

func (g *googleTranslateAdapter) TranslateWithGender(ctx context.Context, text, from, to string, gender Gender) (*Translated, error) {
	result, err := g.client.TranslateWithGender(ctx, text, from, to, gender)
	if err != nil {
		return nil, err
	}
	return &Translated{
		Text:               result.Text,
		Pronunciation:      result.Pronunciation,
		GenderAlternatives: result.GenderAlternatives,
		From: TranslateFrom{
			Language: TranslateFromLanguage{
				DidYouMean: result.From.Language.DidYouMean,
				Iso:        result.From.Language.Iso,
			},
			Text: TranslateFromText{
				AutoCorrected: result.From.Text.AutoCorrected,
				Value:         result.From.Text.Value,
				DidYouMean:    result.From.Text.DidYouMean,
			},
		},
	}, nil
}

// deeplAdapter wraps deepl.DeepL to implement Translator.
type deeplAdapter struct {
	client *deepl.DeepL
}

func (d *deeplAdapter) Translate(ctx context.Context, text, from, to string) (*Translated, error) {
	return d.TranslateWithGender(ctx, text, from, to, GenderUnknown)
}

func (d *deeplAdapter) TranslateWithGender(ctx context.Context, text, from, to string, gender Gender) (*Translated, error) {
	result, err := d.client.TranslateWithGender(ctx, text, from, to, gender)
	if err != nil {
		return nil, err
	}
	return &Translated{
		Text:         result.Text,
		Alternatives: result.Alternatives,
		From: TranslateFrom{
			Language: TranslateFromLanguage{
				Iso: result.From.Language.Iso,
			},
		},
		Method: result.Method,
	}, nil
}

var (
	mu                sync.RWMutex
	defaultTranslator Translator = &googleTranslateAdapter{client: googletranslate.New()}
)

// SetDefaultTranslator sets the default translator used by package-level functions.
func SetDefaultTranslator(t Translator) {
	mu.Lock()
	defer mu.Unlock()
	defaultTranslator = t
}

// GetDefaultTranslator returns the current default translator.
func GetDefaultTranslator() Translator {
	mu.RLock()
	defer mu.RUnlock()
	return defaultTranslator
}

// UseGoogle sets the default translator to Google Translate with optional configuration.
func UseGoogle(opts ...googletranslate.Option) {
	client := googletranslate.New(opts...)
	SetDefaultTranslator(&googleTranslateAdapter{client: client})
}

// UseDeepL sets the default translator to DeepL with optional configuration.
func UseDeepL(opts ...deepl.Option) {
	client := deepl.New(opts...)
	SetDefaultTranslator(&deeplAdapter{client: client})
}

// NewGoogleTranslator creates a new Google Translate adapter.
func NewGoogleTranslator(opts ...googletranslate.Option) Translator {
	return &googleTranslateAdapter{client: googletranslate.New(opts...)}
}

// NewDeepLTranslator creates a new DeepL adapter.
func NewDeepLTranslator(opts ...deepl.Option) Translator {
	return &deeplAdapter{client: deepl.New(opts...)}
}

func getTranslator() Translator {
	mu.RLock()
	defer mu.RUnlock()
	return defaultTranslator
}

// TranslateWithParam translates text using parameters struct.
func TranslateWithParam(ctx context.Context, value params.Translate) (*Translated, error) {
	var (
		text string
		from = "auto"
		to   string
	)
	if value.Text == "" {
		return nil, fmt.Errorf("Text Value is required!")
	}
	text = value.Text

	if value.To == "" {
		return nil, fmt.Errorf("To Value is required!")
	}
	if _, err := language.Parse(value.To); err != nil {
		return nil, fmt.Errorf("To Value isn't valid!")
	}
	to = value.To

	if value.From != "" {
		if _, err := language.Parse(value.From); err != nil {
			return nil, fmt.Errorf("From Value isn't valid!")
		}
		from = value.From
	}
	return getTranslator().Translate(ctx, text, from, to)
}

// Translate translates text with auto-detected source language.
func Translate(ctx context.Context, text, toLanguage string) (*Translated, error) {
	if text == "" {
		return nil, fmt.Errorf("Text Value is required!")
	}
	if toLanguage == "" {
		return nil, fmt.Errorf("To Value is required!")
	}
	if _, err := language.Parse(toLanguage); err != nil {
		return nil, fmt.Errorf("To Value isn't valid!")
	}
	return getTranslator().Translate(ctx, text, "auto", toLanguage)
}

// ManualTranslate translates text with explicit source and target languages.
func ManualTranslate(ctx context.Context, text, fromLanguage, toLanguage string) (*Translated, error) {
	if text == "" {
		return nil, fmt.Errorf("Text Value is required!")
	}
	if fromLanguage == "" {
		return nil, fmt.Errorf("From Value is required!")
	}
	if _, err := language.Parse(fromLanguage); err != nil {
		return nil, fmt.Errorf("From Value isn't valid!")
	}
	if toLanguage == "" {
		return nil, fmt.Errorf("To Value is required!")
	}
	if _, err := language.Parse(toLanguage); err != nil {
		return nil, fmt.Errorf("To Value isn't valid!")
	}
	return getTranslator().Translate(ctx, text, fromLanguage, toLanguage)
}

// TranslateWith translates using a specific translator.
func TranslateWith(ctx context.Context, translator Translator, text, toLanguage string) (*Translated, error) {
	if text == "" {
		return nil, fmt.Errorf("Text Value is required!")
	}
	if toLanguage == "" {
		return nil, fmt.Errorf("To Value is required!")
	}
	if _, err := language.Parse(toLanguage); err != nil {
		return nil, fmt.Errorf("To Value isn't valid!")
	}
	return translator.Translate(ctx, text, "auto", toLanguage)
}
