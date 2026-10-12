package gallery

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestCalendarDemoReadsLocalDays(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	// 21:00 in New York is already the 14th in UTC.
	now := time.Date(2026, 10, 13, 21, 0, 0, 0, ny)
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return now }})
		demoCalendarHeatmap(c, func(string) {})
	}, 800, 200)
	r, ok := tt.Find("Play days 2026")
	if !ok {
		t.Fatalf("no heatmap in %q", tt.Texts())
	}
	// Oct 14 sits at week 41, Wednesday, after four blanks.
	x, y := int(r.X)/2+41*5, int(r.Y)/2+8+3*5
	img, shade := tt.Image(), pixel.SpringDay.Colors[pixel.Shade]
	for j := y; j < y+4; j++ {
		for i := x; i < x+4; i++ {
			if c := img.RGBAAt(2*i, 2*j); c.R != shade.R || c.G != shade.G || c.B != shade.B {
				t.Fatalf("tomorrow's square holds %v at %d,%d", c, i, j)
			}
		}
	}
}
