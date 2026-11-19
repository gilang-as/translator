package googletranslate

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/tidwall/gjson"
)

func TestTranslationText(t *testing.T) {
	const gendered = `[null,[[["la autora",null,"(feminine)",null,null,null,null,1],["el autor",null,"(masculine)",null,null,null,null,2]],"es",2,"en"]]`
	wantAlternatives := []GenderAlternative{
		{Text: "la autora", Gender: GenderFeminine},
		{Text: "el autor", Gender: GenderMasculine},
	}
	tests := []struct {
		name         string
		response     string
		gender       Gender
		want         string
		alternatives []GenderAlternative
		wantErr      bool
	}{
		{"default", gendered, GenderUnknown, "la autora", wantAlternatives, false},
		{"masculine", gendered, GenderMasculine, "el autor", wantAlternatives, false},
		{"feminine", gendered, GenderFeminine, "la autora", wantAlternatives, false},
		{"reverse order", `[null,[[["el autor",null,null,null,null,null,null,2],["la autora",null,null,null,null,null,null,1]],"es",2,"en"]]`, GenderFeminine, "la autora", []GenderAlternative{{Text: "el autor", Gender: GenderMasculine}, {Text: "la autora", Gender: GenderFeminine}}, false},
		{"missing match", `[null,[[["la autora",null,null,null,null,null,null,1]],"es",2,"en"]]`, GenderMasculine, "la autora", []GenderAlternative{{Text: "la autora", Gender: GenderFeminine}}, false},
		{"unlabelled", `[null,[[["primero"],["segundo"]],"es",2,"en"]]`, GenderMasculine, "primero", []GenderAlternative{{Text: "primero"}, {Text: "segundo"}}, false},
		{"sentences", `[null,[[[null,null,null,null,null,[["Primera frase."],["Segunda frase."]]]],"es",1,"en"]]`, GenderMasculine, "Primera frase. Segunda frase.", nil, false},
		{"missing", `[null,[[[null,null,null,null,null,[]]],"es",1,"en"]]`, GenderUnknown, "", nil, true},
		{"empty sentences", `[null,[[[null,null,null,null,null,[[" "]]]],"es",1,"en"]]`, GenderUnknown, "", nil, true},
		{"incomplete", `[null,[[[null,null,null,null,null,[["Primera frase."],[null]]]],"es",1,"en"]]`, GenderUnknown, "", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, alternatives, err := translationText(gjson.Parse(tt.response), tt.gender)
			if (err != nil) != tt.wantErr || got != tt.want || !reflect.DeepEqual(alternatives, tt.alternatives) {
				t.Fatalf("got %q, alternatives %+v, error %v", got, alternatives, err)
			}
		})
	}
}

func TestTranslateWithGenderRejectsInvalidGender(t *testing.T) {
	if _, err := New().TranslateWithGender(context.Background(), "text", "en", "es", Gender("invalid")); err == nil {
		t.Fatal("expected invalid gender to be rejected")
	}
}

func TestTranslationRequestEscapesText(t *testing.T) {
	text := "He wrote \"It\".\\\n\n\t\x00日本語"
	var outer [][][]json.RawMessage
	if err := json.Unmarshal([]byte(translationRequest(text, "auto", "es")), &outer); err != nil {
		t.Fatal(err)
	}
	var inner string
	if err := json.Unmarshal(outer[0][0][1], &inner); err != nil {
		t.Fatal(err)
	}
	var args [][]json.RawMessage
	if err := json.Unmarshal([]byte(inner), &args); err != nil {
		t.Fatal(err)
	}
	for i, want := range []string{text, "auto", "es"} {
		var got string
		if err := json.Unmarshal(args[0][i], &got); err != nil || got != want {
			t.Fatalf("argument %d = %q, %v; want %q", i, got, err, want)
		}
	}
}

func TestTranslateChunks(t *testing.T) {
	for _, text := range []string{
		"He wrote \"It\".\n\nThe end",
		strings.Repeat("A long sentence. ", 400) + "End",
		strings.Repeat("ñ", 3001),
		"\n \n",
	} {
		t.Run(text[:min(len(text), 20)], func(t *testing.T) {
			calls := 0
			translate := func(_ context.Context, chunk, from, to string, gender Gender) (*Translated, error) {
				if len(chunk) > maxChunkSize || !utf8.ValidString(chunk) {
					t.Fatalf("invalid chunk of %d bytes", len(chunk))
				}
				if (calls == 0 && from != "auto") || (calls > 0 && from != "en") ||
					to != "es" || gender != GenderMasculine {
					t.Fatalf("unexpected request: from=%q to=%q gender=%q", from, to, gender)
				}
				calls++
				pronunciation := "pronunciation"
				return &Translated{
					Text: " " + strings.ToUpper(chunk) + " ", Pronunciation: &pronunciation,
					GenderAlternatives: []GenderAlternative{{Text: "alternative", Gender: gender}},
					From:               TranslateFrom{Language: TranslateFromLanguage{Iso: "en"}},
				}, nil
			}
			result, err := translateChunks(context.Background(), text, "auto", "es", GenderMasculine, translate)
			if err != nil {
				t.Fatal(err)
			}
			if result.Text != strings.ToUpper(text) {
				t.Fatalf("translation changed text layout: %q", result.Text)
			}
			if calls > 0 && result.From.Language.Iso != "en" {
				t.Fatal("detected language must be preserved")
			}
			if calls > 1 && (result.Pronunciation != nil || result.GenderAlternatives != nil) {
				t.Fatal("chunk-specific metadata must not describe a complete multi-request translation")
			}
			if calls == 1 && (result.Pronunciation == nil || len(result.GenderAlternatives) != 1) {
				t.Fatal("single-request metadata must be preserved")
			}
			if len(text) > maxChunkSize && calls < 2 {
				t.Fatal("long text must require multiple requests")
			}
		})
	}
}

func TestTranslateChunksErrors(t *testing.T) {
	sentinel := errors.New("request failed")
	for _, tc := range []struct {
		name   string
		result *Translated
		err    error
	}{
		{name: "client error", err: sentinel},
		{name: "nil result"},
		{name: "empty translation", result: &Translated{}},
		{name: "missing detection", result: &Translated{Text: "Hola"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			translate := func(context.Context, string, string, string, Gender) (*Translated, error) {
				calls++
				if calls == 1 && tc.name != "missing detection" {
					return &Translated{Text: "Hola", From: TranslateFrom{Language: TranslateFromLanguage{Iso: "en"}}}, nil
				}
				return tc.result, tc.err
			}
			result, err := translateChunks(context.Background(), "Hello\nWorld", "auto", "es", GenderUnknown, translate)
			if err == nil || result != nil {
				t.Fatalf("expected no partial result and an error, got %+v, %v", result, err)
			}
			if tc.err != nil && !errors.Is(err, sentinel) {
				t.Fatalf("request error was lost: %v", err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := translateChunks(ctx, "Hello", "en", "es", GenderUnknown, func(context.Context, string, string, string, Gender) (*Translated, error) {
		t.Fatal("cancelled requests must not reach the service")
		return nil, nil
	}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}
