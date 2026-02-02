package gt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestTranslate_v1(t *testing.T) {
	gt := NewGoogleTranslate()
	data, err := gt.Translate(context.Background(), "Halo Dunia", "id", "ja")
	if err != nil {
		t.Error(err)
		t.Fail()
	} else {
		prettyJSON, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			log.Fatal("Failed to generate json", err)
		}
		fmt.Println(string(prettyJSON))
	}
}
