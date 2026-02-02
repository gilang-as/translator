package gt

import (
	"fmt"

	"golang.org/x/text/language"
	"gopkg.gilang.dev/translator/params"
)

func TranslateWithParam(value params.Translate) (*Translated, error) {
	var (
		text string
		from = "auto"
		to   string
	)
	if value.Text == "" {
		return nil, fmt.Errorf("Text Value is required!")
	} else {
		text = value.Text
	}
	if value.To == "" {
		return nil, fmt.Errorf("To Value is required!")
	} else {
		if _, err := language.Parse(value.To); err != nil {
			return nil, fmt.Errorf("To Value is't valid!")
		}
		to = value.To
	}
	if value.From != "" {
		if _, err := language.Parse(value.From); err != nil {
			return nil, fmt.Errorf("From Value is't valid!")
		}
		from = value.From
	}
	return translateV1(text, from, to)
}

func Translate(text, toLanguage string) (*Translated, error) {
	if text == "" {
		return nil, fmt.Errorf("Text Value is required!")
	}
	if toLanguage == "" {
		return nil, fmt.Errorf("To Value is required!")
	} else {
		if _, err := language.Parse(toLanguage); err != nil {
			return nil, fmt.Errorf("To Value is't valid!")
		}
	}
	return translateV1(text, "auto", toLanguage)
}

func ManualTranslate(text, fromLanguage, toLanguage string) (*Translated, error) {
	if text == "" {
		return nil, fmt.Errorf("Text Value is required!")
	}
	if fromLanguage == "" {
		return nil, fmt.Errorf("From Value is required!")
	} else {
		if _, err := language.Parse(fromLanguage); err != nil {
			return nil, fmt.Errorf("To Value is't valid!")
		}
	}
	if toLanguage == "" {
		return nil, fmt.Errorf("To Value is required!")
	} else {
		if _, err := language.Parse(toLanguage); err != nil {
			return nil, fmt.Errorf("To Value is't valid!")
		}
	}
	return translateV1(text, fromLanguage, toLanguage)
}
