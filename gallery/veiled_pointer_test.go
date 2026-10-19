package gallery

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestVeiledPointersDoNothing(t *testing.T) {
	for _, e := range Entries() {
		for _, at := range [][2]float32{{30, 20}, {120, 40}, {240, 60}, {60, 120}, {200, 160}, {360, 100}, {100, 260}, {300, 320}} {
			now, inert, logs := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC), false, []string{}
			tt := ui.NewTester(func(c *ui.Context) {
				pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return now }, Inert: inert})
				ui.Column(c).Fill().Children(func() { e.Demo(c, func(s string) { logs = append(logs, s) }) })
			}, 640, 480)
			tt.Move(at[0], at[1])
			before := len(logs)
			// The veil falls with the click, then with a notch.
			inert = true
			tt.ClickAt(at[0], at[1])
			tt.Scroll(at[0], at[1], 0, 40)
			if len(logs) > before {
				t.Errorf("%s: a click or wheel at %v under a falling veil logged %q", e.Slug(), at, logs[before:])
				break
			}
		}
	}
}
