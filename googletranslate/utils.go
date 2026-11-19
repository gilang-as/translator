package googletranslate

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// maxChunkSize leaves room below Google's 5000-character request limit.
const maxChunkSize = 4500

type translateChunkFunc func(context.Context, string, string, string, Gender) (*Translated, error)

func translationRequest(text, from, to string) string {
	value, _ := json.Marshal([]any{[]any{text, from, to, true}, []any{nil}})
	quoted, _ := json.Marshal(string(value))
	return fmt.Sprintf(`[[["MkEWBc",%s,null,"generic"]]]`, quoted)
}

func translateChunks(ctx context.Context, text, from, to string, gender Gender, translate translateChunkFunc) (*Translated, error) {
	var result *Translated
	var lines []string
	requests := 0
	for _, line := range strings.Split(text, "\n") {
		var translated strings.Builder
		for _, chunk := range splitLine(line, maxChunkSize) {
			if strings.TrimSpace(chunk) == "" {
				translated.WriteString(chunk)
				continue
			}
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			part, err := translate(ctx, chunk, from, to, gender)
			if err != nil {
				return nil, err
			}
			if part == nil || strings.TrimSpace(part.Text) == "" {
				return nil, fmt.Errorf("empty translation from %s to %s", from, to)
			}
			if from == "auto" {
				if strings.TrimSpace(part.From.Language.Iso) == "" {
					return nil, fmt.Errorf("source language not detected")
				}
				from = part.From.Language.Iso
			}
			if result == nil {
				copy := *part
				result = &copy
			}
			requests++
			translated.WriteString(trimLines(part.Text))
			if strings.HasSuffix(chunk, " ") {
				translated.WriteByte(' ')
			}
		}
		lines = append(lines, translated.String())
	}
	if result == nil {
		return &Translated{Text: text, From: TranslateFrom{Language: TranslateFromLanguage{Iso: from}}}, nil
	}
	result.Text = strings.Join(lines, "\n")
	if requests > 1 {
		result.Pronunciation = nil
		result.GenderAlternatives = nil
	}
	return result, nil
}

func trimLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}

// splitLine retains separators so chunk boundaries do not introduce new line breaks.
func splitLine(line string, size int) []string {
	var parts []string
	for len(line) > size {
		cut := strings.LastIndex(line[:size], ". ")
		if cut >= 0 {
			cut += 2
		} else if cut = strings.LastIndex(line[:size], " "); cut >= 0 {
			cut++
		} else {
			cut = size
			for cut > 0 && !utf8.RuneStart(line[cut]) {
				cut--
			}
		}
		parts = append(parts, line[:cut])
		line = line[cut:]
	}
	return append(parts, line)
}

// extract extracts a value from a string using regex.
func extract(key string, value string) string {
	regex, err := regexp.Compile(`"` + key + `":".*?"`)
	if err != nil {
		return ""
	}
	res := regex.FindString(value)
	if res == "" {
		return ""
	}
	replace := strings.ReplaceAll(res, `"`+key+`":"`, "")
	return replace[:len(replace)-1]
}
