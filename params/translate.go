package params

import "errors"

// ErrUnsupportedOption reports a translation preference a provider cannot honor.
var ErrUnsupportedOption = errors.New("unsupported translation option")

// Gender identifies a gendered translation alternative.
type Gender string

const (
	GenderUnknown   Gender = ""
	GenderFeminine  Gender = "feminine"
	GenderMasculine Gender = "masculine"
)

// GenderAlternative is a translation offered for a particular gender.
type GenderAlternative struct {
	Text   string `json:"text"`
	Gender Gender `json:"gender"`
}

type Translate struct {
	Text string `json:"text"`
	From string `json:"from"`
	To   string `json:"to"`
}
