package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestEveryFocusHasAName(t *testing.T) {
	for _, e := range Entries() {
		tt := ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay})
			// A named start, so Tab always lands and laps.
			input.Button(c, "SENTINEL", input.ButtonOptions{})
			ui.Column(c).Fill().Children(func() { e.Demo(c, func(string) {}) })
		}, 640, 480)
		tt.Key(0, ui.KeyTab)
		for i := range 12 {
			tt.Key(0, ui.KeyTab)
			if tt.Focused("SENTINEL") {
				break
			}
			// Focused counts ancestors, so a named focus has a hit.
			named := false
			for _, s := range tt.Texts() {
				named = named || tt.Focused(s)
			}
			if !named {
				t.Errorf("%s: Tab %d focuses something with no name", e.Slug(), i+1)
				break
			}
		}
	}
}
