package layout_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/layout"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestScrollAreaFollowsItsEndWhileThere(t *testing.T) {
	s, n, follow := ui.ScrollState{}, 10, true
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay})
		layout.ScrollArea(c, &s, layout.ScrollAreaOptions{Label: "Log", Follow: follow}, func() {
			for range n {
				ui.Box(c).Height(32)
			}
		}).Element.Size(200, 100)
	}, 300, 200)
	tt.Frame()
	if s.MaxY == 0 || s.Y != s.MaxY {
		t.Fatalf("a following view opened at %v of %v", s.Y, s.MaxY)
	}
	n++
	tt.Frame()
	if s.Y != s.MaxY {
		t.Fatalf("at the end, a new row left the view at %v of %v", s.Y, s.MaxY)
	}
	// Scrolled up it stays; back at the end, it follows.
	s.Y = 0
	n++
	tt.Frame()
	if s.Y != 0 {
		t.Fatalf("scrolled up, the view jumped to %v", s.Y)
	}
	s.Y = s.MaxY
	tt.Frame()
	n++
	tt.Frame()
	if s.Y != s.MaxY {
		t.Fatalf("back at the end, the view stopped following at %v of %v", s.Y, s.MaxY)
	}
	// Without Follow, growth leaves the view alone.
	follow, s = false, ui.ScrollState{}
	n++
	tt.Frame()
	if s.Y != 0 {
		t.Fatalf("a plain view followed to %v", s.Y)
	}
}
