package code

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type FPSMeterOptions struct {
	// Width and Height are dots; zeros mean 120 by 40.
	Width, Height int
}

// FPSMeter draws frame times as bars, newest right, with percentiles.
func FPSMeter(c *ui.Context, text string, frames []time.Duration, o FPSMeterOptions) ui.Element {
	k := pixel.K(c)
	w, h := pixel.Or("FPSMeter Width", o.Width, 120), pixel.Or("FPSMeter Height", o.Height, 40)
	for _, f := range frames {
		if f < 0 {
			panic(fmt.Sprintf("code: FPSMeter frame of %v", f))
		}
	}
	// Two dots a bar; the oldest scroll off the left.
	shown := frames[max(0, len(frames)-(w-4)/2):]
	p50, p95 := percentile(shown, 50), percentile(shown, 95)
	head := "P50 -  P95 -"
	if len(shown) > 0 {
		head = fmt.Sprintf("P50 %.1fMS  P95 %.1fMS", ms(p50), ms(p95))
	}
	top := 2 * time.Second / 60
	if len(shown) > 0 {
		top = max(top, slices.Max(shown))
	}
	return ui.Box(c).Size(k.D(w), k.D(h+8)).Role(ui.RoleImage).Label(text).Description(head).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		body, ink, _ := dark(k)
		g.Micro(0, 0, head, pixel.Faint, 1)
		g.Rect(0, 8, w, h, body)
		y := func(d time.Duration) int { return 8 + h - 2 - int(int64(h-4)*int64(d)/int64(top)) }
		for i, f := range shown {
			slot := pixel.Good
			switch {
			case f > time.Second/30:
				slot = pixel.Bad
			case f > time.Second/60:
				slot = pixel.Warn
			}
			g.Rect(2+2*i, y(f), 1, 8+h-2-y(f), slot)
		}
		if len(shown) > 0 {
			for x := 2; x < w-2; x += 2 {
				g.Dot(x, y(p50), ink)
				g.Dot(x+1, y(p95), pixel.Warn)
			}
		}
	})
}

// percentile is the nearest-rank q of ds; empty is 0.
func percentile(ds []time.Duration, q int) time.Duration {
	if len(ds) == 0 {
		return 0
	}
	s := slices.Sorted(slices.Values(ds))
	return s[max(0, (q*len(s)+99)/100-1)]
}

func ms(d time.Duration) float64 { return float64(d) / float64(time.Millisecond) }
