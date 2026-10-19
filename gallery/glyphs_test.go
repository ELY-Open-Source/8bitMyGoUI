package gallery

import (
	"testing"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestGalleryTextHasGlyphs(t *testing.T) {
	// Every visible word comes from the font, never tofu.
	check := func(where, s string) {
		for _, r := range s {
			if r != ' ' && !pixel.Has(r) {
				t.Errorf("%s: the font has no %q in %q", where, r, s)
			}
		}
	}
	for _, e := range Entries() {
		check(e.Slug(), e.Name+e.Zh+e.Summary)
	}
	for _, path := range []string{"/", "/cat/foundations", "/nowhere"} {
		_, tt := open(path)
		for _, s := range tt.Texts() {
			check(path, s)
		}
	}
}
