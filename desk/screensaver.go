package desk

import (
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type ScreensaverOptions struct{}

type ScreensaverResult struct {
	Element ui.Element
	// Woke is the input that ended it this frame.
	Woke bool
}

type saverKey struct{}

// bounces ink the logo, a new one at each wall.
var bounces = []pixel.Slot{pixel.Accent, pixel.Good, pixel.Warn, pixel.Info, pixel.Magic, pixel.Bad}

// Screensaver bounces a logo around; any input stops it.
func Screensaver(c *ui.Context, logo string, open *bool, o ScreensaverOptions) ScreensaverResult {
	k := pixel.K(c)
	if logo == "" {
		panic("desk: Screensaver has no logo")
	}
	n := 0
	anchor, _, woke := cover(c, k, logo, open, true, func(g pixel.Canvas) {
		g.Rect(0, 0, g.W, g.H, pixel.Well)
		bw, bh := 2*pixel.Measure(logo)+12, 2*pixel.LineHeight+8
		rx, ry := max(1, g.W-bw), max(1, g.H-bh)
		x, y := fold(n+rx/3, rx), fold(n+ry/2, ry)
		hits := (n+rx/3)/rx + (n+ry/2)/ry
		ink := bounces[hits%len(bounces)]
		g.Frame(x, y, bw, bh, pixel.Raised, pixel.Face)
		g.Text(x+6, y+4, logo, ink, 2)
		// A corner hit, the rarest bounce, earns sparkles.
		if (x == 0 || x == rx) && (y == 0 || y == ry) {
			for _, d := range [][2]int{{-4, -4}, {bw + 1, -4}, {-4, bh + 1}, {bw + 1, bh + 1}} {
				g.Rect(x+d[0]+1, y+d[1], 1, 3, pixel.Lit)
				g.Rect(x+d[0], y+d[1]+1, 3, 1, pixel.Lit)
			}
		}
	})
	st := ui.Local(anchor, saverKey{}, func() saver { return saver{} })
	switch {
	case !*open:
		st.on = false
	case !st.on:
		*st = saver{on: true, at: k.Now()}
	}
	// It moves a dot a step; reduced, it holds still.
	if *open && !k.Reduce() {
		n = max(0, k.Since(c, st.at, 15, math.MaxInt32))
	}
	if woke {
		slog.Debug("8bit: screensaver woke", "logo", logo)
	}
	return ScreensaverResult{Element: anchor, Woke: woke}
}

// saver is whether it shows, and since when.
type saver struct {
	on bool
	at time.Time
}

// fold reflects t back and forth across 0 to r.
func fold(t, r int) int {
	m := t % (2 * r)
	if m > r {
		return 2*r - m
	}
	return m
}
