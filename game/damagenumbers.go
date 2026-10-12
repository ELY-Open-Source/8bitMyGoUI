package game

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Hit is a number rising from dot X, Y.
type Hit struct {
	Value int
	Crit  bool
	X, Y  int
	at    time.Time
}

type DamageNumbersOptions struct {
	// Width and Height are dots; zeros mean 160 by 80.
	Width, Height int
}

// life is a number's time; it dissolves the last half.
const life = 800 * time.Millisecond

// DamageNumbers floats hits up, dissolves, then drops them.
func DamageNumbers(c *ui.Context, text string, hits *[]Hit, o DamageNumbersOptions) *ui.Element {
	k := pixel.K(c)
	w, h := pixel.Or("DamageNumbers Width", o.Width, 160), pixel.Or("DamageNumbers Height", o.Height, 80)
	now := k.Now()
	live := (*hits)[:0]
	for _, hit := range *hits {
		if hit.X < 0 || hit.X >= w || hit.Y < 0 || hit.Y >= h {
			panic(fmt.Sprintf("game: hit at %d,%d is off the %d×%d box", hit.X, hit.Y, w, h))
		}
		if hit.at.IsZero() {
			hit.at = now
		}
		if now.Sub(hit.at) < life {
			live = append(live, hit)
		}
	}
	*hits = live
	shown := append([]Hit(nil), live...)
	parts := make([]string, len(shown))
	for i, hit := range shown {
		parts[i] = strconv.Itoa(hit.Value)
		if hit.Crit {
			parts[i] += "!"
		}
	}
	if len(shown) > 0 {
		k.Tick(c, 20)
	}
	return ui.Box(c).Role(ui.RoleImage).Label(text).Description(strings.Join(parts, ", ")).Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		// Numbers wider than the box are cut at it.
		g.Clip(0, 0, w, h, func() { numbers(g, k, shown, now, w, h) })
	})
}

// numbers draws live hits, risen and dissolved by age.
func numbers(g pixel.Canvas, k pixel.Kit, shown []Hit, now time.Time, w, h int) {
	for _, hit := range shown {
		t := float64(now.Sub(hit.at)) / float64(life)
		ink, sc, s := pixel.Ink, 1, strconv.Itoa(hit.Value)
		if hit.Crit {
			ink, sc, s = pixel.Warn, 2, s+"!"
		}
		// Wide or high numbers are pushed back inside the box.
		tw, th := pixel.Measure(s)*sc, pixel.LineHeight*sc
		x, y := max(0, min(w-tw, hit.X-tw/2)), max(0, min(h-th, hit.Y-th/2))
		// It rises 12 dots and dissolves; reduced, it holds still.
		if k.Reduce() {
			g.Text(x, y, s, ink, sc)
			continue
		}
		level := 16
		if t > 0.5 {
			level = int(16 * (1 - t) * 2)
		}
		g.TextVeil(x, max(0, y-int(12*t)), s, ink, level, sc)
	}
}
