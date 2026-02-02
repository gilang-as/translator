package googletranslate

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

// TranslateFrom contains source language and text information.
type TranslateFrom struct {
	Language TranslateFromLanguage `json:"language"`
	Text     TranslateFromText     `json:"text"`
}

// Translated represents a translation result.
type Translated struct {
	Text          string        `json:"text"`
	Pronunciation *string       `json:"pronunciation"`
	From          TranslateFrom `json:"from"`
}

// reqData holds the session data extracted from Google Translate page.
type reqData struct {
	FsId string `json:"f.sid"`
	Bl   string `json:"bl"`
	At   string `json:"at"`
}
