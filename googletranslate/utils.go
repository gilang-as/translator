package googletranslate

import (
	"regexp"
	"strings"
)

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
