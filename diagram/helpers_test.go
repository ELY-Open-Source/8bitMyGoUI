package diagram

import (
	"encoding/json"
	"os"
	"testing"
)

func corpus(t *testing.T) map[string][]string {
	t.Helper()
	raw, err := os.ReadFile("testdata/corpus.json")
	if err != nil {
		t.Fatal(err)
	}
	var docs map[string][]string
	if err := json.Unmarshal(raw, &docs); err != nil {
		t.Fatal(err)
	}
	return docs
}
