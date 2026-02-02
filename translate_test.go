package gt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.gilang.dev/translator/params"
)

func TestTranslateWithParam(t *testing.T) {
	value := params.Translate{
		Text: "Halo Dunia",
		//From: "id",
		To: params.ENGLISH,
	}
	translated, err := TranslateWithParam(context.Background(), value)
	if err != nil {
		t.Error(err)
	}
	if translated.Text != "" {
		t.Log(translated)
	}
}

func TestTranslate(t *testing.T) {
	translated, err := Translate(context.Background(), "Hello World", params.INDONESIAN)
	if err != nil {
		t.Error(err)
	}
	if translated.Text != "" {
		t.Log(translated)
	}
}

func TestManualTranslate(t *testing.T) {
	translated, err := ManualTranslate(context.Background(), "Halo Semuanya", params.INDONESIAN, params.JAVANESE)
	if err != nil {
		t.Error(err)
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
