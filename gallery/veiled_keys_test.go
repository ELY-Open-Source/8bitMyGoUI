package gallery

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestVeiledKeysDoNothing(t *testing.T) {
	// KeyUnknown stands for typed text, "a1".
	keys := []ui.Key{ui.KeyEnter, ui.KeySpace, ui.KeyUp, ui.KeyDown, ui.KeyLeft, ui.KeyRight, ui.Key1, ui.KeyUnknown}
	for _, e := range Entries() {
		for i := range 3 * len(keys) {
			stop, key := i/len(keys), keys[i%len(keys)]
			at, inert, logs := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC), false, []string{}
			tt := ui.NewTester(func(c *ui.Context) {
				pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }, Inert: inert})
				ui.Column(c).Fill().Children(func() { e.Demo(c, func(s string) { logs = append(logs, s) }) })
			}, 640, 480)
			// Each key at each of the first three Tab stops.
			for range 1 + stop {
				tt.Key(0, ui.KeyTab)
			}
			before := len(logs)
			// The veil falls in the frame the key comes.
			inert = true
			if key == ui.KeyUnknown {
				tt.Type("a1")
			} else {
				tt.Key(0, key)
			}
			if len(logs) > before {
				t.Errorf("%s: %v under a falling veil logged %q", e.Slug(), key, logs[before:])
				break
			}
		}
	}
}
