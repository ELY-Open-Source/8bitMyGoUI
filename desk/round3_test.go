package desk_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/desk"
)

func TestNarrowSizesStayInsideOrPanic(t *testing.T) {
	mustPanic(t, "thin window", func() {
		view(600, 400, newWorld(), func(c *ui.Context) {
			x, y := 0, 0
			desk.AppWindow(c, "W", &x, &y, desk.AppWindowOptions{Width: 1, Height: 10}, func() {})
		})
	})
	mustPanic(t, "thin boot", func() {
		view(600, 400, newWorld(), func(c *ui.Context) { desk.BootScreen(c, "B", nil, desk.BootScreenOptions{Width: 100}) })
	})
	for _, tc := range []struct {
		label string
		build func(c *ui.Context)
	}{
		{"BIOS", func(c *ui.Context) {
			desk.BootScreen(c, "BIOS", []string{"A LONG POST LINE THAT RUNS ON"}, desk.BootScreenOptions{Width: 120})
		}},
		{"Snake", func(c *ui.Context) { desk.Snake(c, "Snake", desk.SnakeOptions{Cols: 6, Rows: 4, Seed: 1}) }},
	} {
		tt := view(800, 500, newWorld(), tc.build)
		x, y, w, h := dots(find(t, tt, tc.label))
		if n := spill(tt.Image(), x, y, w, h); n > 0 {
			t.Errorf("%s spills %d dots", tc.label, n)
		}
	}
}

func TestClocksSteppingBackHoldStill(t *testing.T) {
	wd := newWorld()
	tt := view(600, 300, wd, func(c *ui.Context) {
		desk.DiskDefrag(c, "Disk", true, desk.DiskDefragOptions{Blocks: 64, Seed: 4})
	})
	later(tt, wd, 500*time.Millisecond)
	before := append([]byte(nil), tt.Image().Pix...)
	later(tt, wd, -2*time.Second)
	if !bytes.Equal(before, tt.Image().Pix) {
		t.Fatal("the disk moved as the clock stepped back")
	}
	// Booted behind the clock, the count shows zero.
	fresh := view(600, 300, newWorld(), func(c *ui.Context) { desk.BootScreen(c, "BIOS", nil, desk.BootScreenOptions{}) })
	wb := newWorld()
	back := view(600, 300, wb, func(c *ui.Context) { desk.BootScreen(c, "BIOS", nil, desk.BootScreenOptions{}) })
	later(back, wb, -time.Hour)
	if !bytes.Equal(fresh.Image().Pix, back.Image().Pix) {
		t.Fatal("a boot behind the clock shows a count under zero")
	}
}
