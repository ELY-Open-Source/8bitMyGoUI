// Package chat draws messages, channels, mail and calls.
package chat

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Message is a chat line; Mine sits on the right.
type Message struct {
	From, Text string
	At         time.Time
	Mine       bool
}

// listed is a chat list's choice; Submitted is Enter.
type listed struct {
	Element            *ui.Element
	Changed, Submitted bool
}

// list is a data.List of n named rows.
func list(c *ui.Context, k pixel.Kit, text string, n int, selected *int, w, h int, disabled bool, name func(i int) string, row func(i int)) listed {
	if *selected < -1 || *selected >= n {
		panic(fmt.Sprintf("chat: %s selects %d of %d rows", text, *selected, n))
	}
	var r data.ListResult
	e := ui.Box(c).Width(k.D(w)).Children(func() {
		r = data.List(c, text, n, selected, data.ListOptions{Height: h, Disabled: disabled}, name, row)
	})
	return listed{Element: e, Changed: r.Changed, Submitted: r.Submitted}
}

// para wraps s to w dots, hugging its lines.
func para(c *ui.Context, k pixel.Kit, s string, w, maxLines int, ink pixel.Slot) *ui.Element {
	lines := pixel.Wrap(s, w, maxLines)
	lw := 0
	for _, l := range lines {
		lw = max(lw, pixel.Measure(l))
	}
	return ui.Box(c).Role(ui.RoleText).Label(s).Size(k.D(lw), k.D(len(lines)*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, l := range lines {
			g.Text(0, i*pixel.LineHeight, l, ink, 1)
		}
	})
}

// line is s cut to w dots with an ellipsis.
func line(s string, w int) string { return pixel.Wrap(s, w, 1)[0] }

// clock prints t's time of day.
func clock(t time.Time) string { return t.Format("15:04") }

// tag draws s in micro letters, or the font.
func tag(g pixel.Canvas, x, y int, s string, ink pixel.Slot) {
	if pixel.MicroOK(s) {
		g.Micro(x, y+3, s, ink, 1)
		return
	}
	g.Text(x, y, s, ink, 1)
}

// tagW is tag's width for s.
func tagW(s string) int {
	if pixel.MicroOK(s) {
		return pixel.MicroWidth(s)
	}
	return pixel.Measure(s)
}

// rule centers label on a line across w dots.
func rule(c *ui.Context, k pixel.Kit, label string, w int, ink, edge pixel.Slot) *ui.Element {
	tw := tagW(label)
	if w < tw+16 {
		panic(fmt.Sprintf("chat: a %d-dot rule cannot hold %q", w, label))
	}
	return ui.Box(c).Role(ui.RoleText).Label(label).Size(k.D(w), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		x, y := (w-tw)/2, pixel.LineHeight/2
		g.HLine(0, y, x-4, edge)
		g.HLine(x+tw+4, y, w-x-tw-4, edge)
		tag(g, x, 0, label, ink)
	})
}
