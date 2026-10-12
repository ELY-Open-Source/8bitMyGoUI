package datetime

import (
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type YearViewOptions struct{}

// A small month's day cell and box, in dots.
const miniW, miniH, monthW, monthH = 10, 7, 7*miniW + 4, 6*miniH + pixel.LineHeight + 2

// YearView sets a year's twelve months small; today is marked.
func YearView(c *ui.Context, year int, o YearViewOptions) *ui.Element {
	k := pixel.K(c)
	now := today(c, k)
	return ui.Box(c).Size(k.D(4*monthW), k.D(3*monthH)).Role(ui.RoleImage).Label(strconv.Itoa(year)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for m := range 12 {
			ox, oy := m%4*monthW, m/4*monthH
			first := time.Date(year, time.Month(m+1), 1, 0, 0, 0, 0, now.Location())
			g.Text(ox, oy, monthName(k, first.Month()), pixel.Faint, 1)
			skip := int(first.Weekday())
			for d := first; d.Month() == first.Month(); d = d.AddDate(0, 0, 1) {
				i := skip + d.Day() - 1
				x, y := ox+i%7*miniW, oy+pixel.LineHeight+i/7*miniH
				n := strconv.Itoa(d.Day())
				ink := pixel.Ink
				if sameDay(d, now) {
					g.Rect(x-1, y-1, miniW, miniH, pixel.Accent)
					ink = pixel.AccentInk
				}
				g.Micro(x+(miniW-1-pixel.MicroWidth(n)), y, n, ink, 1)
			}
		}
	})
}
