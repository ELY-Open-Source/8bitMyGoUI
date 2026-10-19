package desk

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type BootScreenOptions struct {
	// Memory counts up in kilobytes; zero means 640.
	Memory int
	// Width and Height are dots; zeros mean 320 by 200.
	Width, Height int
}

type BootScreenResult struct {
	Element ui.Element
	// Done is the frame the last line shows.
	Done bool
}

type bootKey struct{}

// boot is when the screen showed, and if it finished.
type boot struct {
	at   time.Time
	done bool
}

// BootScreen counts memory, then prints POST lines by the logo.
func BootScreen(c *ui.Context, title string, lines []string, o BootScreenOptions) BootScreenResult {
	k := pixel.K(c)
	mem := pixel.Or("BootScreen Memory", o.Memory, 640)
	w, h := pixel.Or("BootScreen Width", o.Width, 320), pixel.Or("BootScreen Height", o.Height, 200)
	if title == "" || mem < 1 || w < 120 || 6+(len(lines)+3)*pixel.LineHeight > h {
		panic(fmt.Sprintf("desk: BootScreen %q, %dK, %d lines in %d dots", title, mem, len(lines), h))
	}
	e := ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleStatus).Label(title)
	st := ui.Local(e, bootKey{}, func() boot { return boot{at: k.Now()} })
	// Counting takes 15 steps, a line 4; reduced, at once.
	total := 15 + 4*max(0, len(lines)-1)
	n := total
	if !k.Reduce() {
		n = max(0, k.Since(c, st.at, 15, total))
	}
	res := BootScreenResult{Element: e, Done: n == total && !st.done}
	st.done = n == total
	kb := mem / 15 * min(n, 15)
	if n >= 15 {
		kb = mem
	}
	count := fmt.Sprintf(k.Msg("memory_test"), kb)
	if n >= 15 {
		count += " " + k.Msg("ok")
	}
	shown := lines[:max(0, min(len(lines), (n-15)/4+1))]
	if n < 15 {
		shown = nil
	}
	e.Description(count)
	if res.Done {
		slog.Debug("8bit: boot screen done", "label", title)
	}
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, w, h, pixel.Well)
		g.Sprite(w-44, 6, icon.Star, 0, pixel.Warn, 3)
		g.Text(6, 6, pixel.Wrap(title, w-56, 1)[0], pixel.Ink, 1)
		g.Text(6, 6+2*pixel.LineHeight, pixel.Wrap(count, w-56, 1)[0], pixel.Ink, 1)
		for i, l := range shown {
			g.Text(6, 6+(3+i)*pixel.LineHeight, pixel.Wrap(l, w-12, 1)[0], pixel.Faint, 1)
		}
	})
	return res
}
