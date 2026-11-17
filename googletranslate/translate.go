package googletranslate

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
)

// check retrieves session data from Google Translate page.
func (gt *GoogleTranslate) check(ctx context.Context) (*reqData, error) {
	gt.mu.RLock()
	host := gt.host
	proxyURL := gt.proxyURL
	gt.mu.RUnlock()

	baseURL := "https://translate." + host

	client := req.C()
	if proxyURL != "" {
		client.SetProxyURL(proxyURL)
	}

	headers := http.Header{
		"Accept":             []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"Accept-Language":    []string{"en-US,en;q=0.9,id;q=0.8"},
		"Sec-Ch-Ua":          []string{`.Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		"Sec-Ch-Ua-Mobile":   []string{"?0"},
		"Sec-Ch-Ua-Platform": []string{`"Windows"`},
		"Sec-Fetch-Dest":     []string{"document"},
		"Sec-Fetch-User":     []string{"?1"},
		"User-Agent":         []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"},
	}

	r := client.R()
	r.SetContext(ctx)
	r.Headers = headers

	resp, err := r.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error: bad network")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: request failed with status code %d", resp.StatusCode)
	}

	body := resp.String()

	return &reqData{
		FsId: extract("FdrFJe", body),
		Bl:   extract("cfb2h", body),
		At:   extract("SNlM0e", body),
	}, nil
}

// Translate translates text from one language to another.
func (gt *GoogleTranslate) Translate(ctx context.Context, text string, from string, to string) (*Translated, error) {
	return gt.TranslateWithGender(ctx, text, from, to, GenderUnknown)
}

// TranslateWithGender prefers the requested gender when Google offers alternatives.
// If no matching alternative is available, the first translation is used.
func (gt *GoogleTranslate) TranslateWithGender(ctx context.Context, text string, from string, to string, gender Gender) (*Translated, error) {
	if gender != GenderUnknown && gender != GenderFeminine && gender != GenderMasculine {
		return nil, fmt.Errorf("invalid translation gender %q", gender)
	}
	gt.mu.RLock()
	host := gt.host
	proxyURL := gt.proxyURL
	gt.mu.RUnlock()

	rpcId := "MkEWBc"
	baseURL := "https://translate." + host

	checkData, err := gt.check(ctx)
	if err != nil {
		return nil, err
	}

	// Build query parameters
	params := url.Values{}
	params.Set("rpcids", rpcId)
	params.Set("f.sid", checkData.FsId)
	params.Set("bl", checkData.Bl)
	params.Set("hl", "en-US")
	params.Set("soc-app", "1")
	params.Set("soc-platform", "1")
	params.Set("soc-device", "1")
	params.Set("_reqid", strconv.Itoa(100000+int(float64(9000)*float64(int64(1)%100)/100)))
	params.Set("rt", "c")

	fullURL := baseURL + "/_/TranslateWebserverUi/data/batchexecute?" + params.Encode()

	// Build request body
	value := fmt.Sprintf(`[["%s","%s","%s",true],[null]]`, text, from, to)
	fReq := fmt.Sprintf(`[[["MkEWBc",%q,null,"generic"]]]`, value)

	body := url.Values{}
	body.Set("f.req", fReq)

	client := req.C()
	if proxyURL != "" {
		client.SetProxyURL(proxyURL)
	}

	headers := http.Header{
		"Sec-Ch-Ua":          []string{`"Google Chrome";v="95", "Chromium";v="95", ";Not A Brand";v="99"`},
		"Content-Type":       []string{"application/x-www-form-urlencoded;charset=UTF-8"},
		"X-Same-Domain":      []string{"1"},
		"Sec-Ch-Ua-Mobile":   []string{"?1"},
		"User-Agent":         []string{"Mozilla/5.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.004) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Mobile Safari/537.36"},
		"Sec-Ch-Ua-Platform": []string{`"Android"`},
		"Accept":             []string{"*/*"},
		"Origin":             []string{"https://translate.google.com"},
		"Sec-Fetch-Site":     []string{"same-origin"},
		"Sec-Fetch-Mode":     []string{"cors"},
		"Sec-Fetch-Dest":     []string{"empty"},
		"Accept-Language":    []string{"en-US,en;q=0.9"},
	}

	r := client.R()
	r.SetContext(ctx)
	r.Headers = headers

	resp, err := r.SetBody(bytes.NewBufferString(body.Encode())).Post(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error: bad network")
	}

	raw := resp.String()
	if len(raw) < 6 {
		return nil, fmt.Errorf("error: invalid response")
	}

	// Parse response
	lines := strings.Split(raw[6:], "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("error: parsing response")
	}

	// Parse first level JSON
	result := gjson.Parse(lines[1])
	innerJSON := result.Get("0.2").String()
	if innerJSON == "" {
		return nil, fmt.Errorf("error: request on google translate api isn't working, please check your parameter")
	}

	// Parse inner JSON
	data := gjson.Parse(innerJSON)

	translatedText, alternatives, err := translationText(data, gender)
	if err != nil {
		return nil, err
	}

	// Extract pronunciation
	var pronunciation *string
	pronValue := data.Get("1.0.0.1").String()
	if pronValue != "" {
		pronunciation = &pronValue
	}

	// Extract source language ISO
	textIso := data.Get("1.3").String()

	// Extract did you mean and autocorrect
	didYouMean := false
	didYouMeanLanguage := false
	autoCorrected := false
	var autoCorrectedValue *string

	if data.Get("0.0").Exists() && data.Get("0.0").String() != "" {
		autoCorrected = true
		didYouMeanLanguage = true
		txt := data.Get("0.0").String()
		autoCorrectedValue = &txt
	} else if data.Get("0.1.0.0.1").Exists() {
		txt := data.Get("0.1.0.0.1").String()
		if txt != "" {
			r := regexp.MustCompile(`<.*?>`)
			cleanTxt := r.ReplaceAllString(txt, "")
			autoCorrectedValue = &cleanTxt
			didYouMean = true
		}
	}

	return &Translated{
		Text:               translatedText,
		Pronunciation:      pronunciation,
		GenderAlternatives: alternatives,
		From: TranslateFrom{
			Language: TranslateFromLanguage{
				DidYouMean: didYouMeanLanguage,
				Iso:        textIso,
			},
			Text: TranslateFromText{
				AutoCorrected: autoCorrected,
				Value:         autoCorrectedValue,
				DidYouMean:    didYouMean,
			},
		},
	}, nil
}

func translationText(data gjson.Result, gender Gender) (string, []GenderAlternative, error) {
	entries := data.Get("1.0").Array()
	if len(entries) == 0 {
		return "", nil, fmt.Errorf("error: translation response contains no text")
	}

	if entries[0].Get("0").Type == gjson.String {
		var alternatives []GenderAlternative
		for _, entry := range entries {
			text := entry.Get("0")
			if text.Type != gjson.String || strings.TrimSpace(text.String()) == "" {
				continue
			}
			var alternativeGender Gender
			// Google uses 1 for feminine and 2 for masculine.
			switch entry.Get("7").Int() {
			case 1:
				alternativeGender = GenderFeminine
			case 2:
				alternativeGender = GenderMasculine
			}
			alternatives = append(alternatives, GenderAlternative{Text: text.String(), Gender: alternativeGender})
		}
		if len(alternatives) == 0 {
			return "", nil, fmt.Errorf("error: translation response contains no text")
		}
		selected := alternatives[0].Text
		if gender != GenderUnknown {
			for _, alternative := range alternatives {
				if alternative.Gender == gender {
					selected = alternative.Text
					break
				}
			}
		}
		return selected, alternatives, nil
	}

	sentences := entries[0].Get("5").Array()
	if len(sentences) == 0 {
		return "", nil, fmt.Errorf("error: translation response contains no text")
	}
	var parts []string
	for _, sentence := range sentences {
		text := sentence.Get("0")
		if text.Type != gjson.String || text.String() == "" {
			return "", nil, fmt.Errorf("error: translation response contains an incomplete sentence")
		}
		parts = append(parts, text.String())
	}
	text := strings.TrimSpace(strings.Join(parts, " "))
	if text == "" {
		return "", nil, fmt.Errorf("error: translation response contains no text")
	}
	return text, nil, nil
}
