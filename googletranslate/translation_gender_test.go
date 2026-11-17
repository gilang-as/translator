package googletranslate

import (
	"context"
	"reflect"
	"testing"

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
