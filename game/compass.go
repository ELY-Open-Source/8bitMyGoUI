package game

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Bearing marks a direction on a compass, in degrees.
type Bearing struct {
	Degrees float64
	Slot    pixel.Slot
}

type CompassOptions struct {
	// Width is in dots; zero means 120. It shows 180°.
	Width int
}

// Compass is a heading strip with letters and markers.
func Compass(c *ui.Context, text string, heading float64, marks []Bearing, o CompassOptions) ui.Element {
	k := pixel.K(c)
	for _, v := range append([]float64{heading}, bearings(marks)...) {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			panic(fmt.Sprintf("game: Compass bearing %v", v))
		}
	}
	w := pixel.Or("Compass Width", o.Width, 120)
	letters := strings.Fields(k.Msg("compass"))
	heading = math.Mod(math.Mod(heading, 360)+360, 360)
	// x places bearing b; false means off the strip.
	x := func(b float64) (int, bool) {
		d := math.Mod(math.Mod(b-heading, 360)+540, 360) - 180
		return w/2 + int(math.Round(d*float64(w)/180)), d >= -90 && d <= 90
	}
	top := tagH(letters[0]) + 3
	return ui.Box(c).Role(ui.RoleImage).Label(text).Description(strconv.Itoa(int(math.Round(heading)))+"°").Size(k.D(w), k.D(top+10)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, top, w, 10, pixel.Sunken, pixel.Well)
		for b := 0; b < 360; b += 15 {
			at, ok := x(float64(b))
			if !ok || at < 1 || at > w-2 {
				continue
			}
			tall := 3
			if b%45 == 0 {
				tall = 6
			}
			g.VLine(at, top+10-1-tall, tall, pixel.Edge)
			if b%90 == 0 {
				l := letters[b/90]
				tag(g, max(0, min(w-tagW(l), at-tagW(l)/2)), 0, l, pixel.Ink)
			}
		}
		for _, m := range marks {
			if at, ok := x(m.Degrees); ok {
				g.Rect(max(1, min(w-4, at-1)), top+2, 3, 3, m.Slot)
			}
		}
		// The heading's own mark points down at the middle.
		g.Rect(w/2-2, top-1, 5, 1, pixel.Accent)
		g.Rect(w/2-1, top, 3, 1, pixel.Accent)
		g.Dot(w/2, top+1, pixel.Accent)
	})
}

func bearings(marks []Bearing) []float64 {
	out := make([]float64, len(marks))
	for i, m := range marks {
		out[i] = m.Degrees
	}
	return out
}
