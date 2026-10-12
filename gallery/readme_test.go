package gallery

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestReadmeCountsMatchTheRegistry(t *testing.T) {
	b, err := os.ReadFile("../README.md")
	if err != nil {
		t.Fatal(err)
	}
	readme, total := string(b), 0
	for _, cat := range Categories[1:] {
		n := count(cat.Key)
		total += n
		if row := fmt.Sprintf("| %s | %s | %d |", cat.Name, cat.Zh, n); !strings.Contains(readme, row) {
			t.Errorf("README lacks %q", row)
		}
	}
	if all := fmt.Sprintf("| **All** | **全部** | **%d** |", total); !strings.Contains(readme, all) {
		t.Errorf("README lacks %q", all)
	}
	if n := count("foundations"); n != 8 || !strings.Contains(readme, "Eight foundation pages") {
		t.Errorf("README says eight foundation pages; the registry has %d", n)
	}
}
