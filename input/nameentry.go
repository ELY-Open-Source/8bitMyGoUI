package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type NameEntryOptions struct {
	Disabled bool
}

type NameEntryResult struct {
	Element            *ui.Element
	Changed, Confirmed bool
}

type nameKey struct{}

const (
	nameLen = 3
	// The grid ends with DEL and END after A–Z.
	nameDel, nameEnd = 26, 27
)

// NameEntry spells arcade initials from a letter grid; END confirms.
func NameEntry(c *ui.Context, text string, name *string, o NameEntryOptions) NameEntryResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(6)).AlignItems(ui.Center)
	cur := ui.Local(col, nameKey{}, func() int { return 0 })
	res := NameEntryResult{Element: col}
	col.Children(func() {
		ui.Box(c).Size(k.D(nameLen*16-4), k.D(16)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range nameLen {
				x := i * 16
				if i < len(*name) {
					l := (*name)[i : i+1]
					g.Text(x+(12-pixel.Measure(l))/2, 0, l, pixel.Accent, 1)
				}
				ink := pixel.Edge
				if i == len(*name) {
					ink = pixel.Accent
				}
				g.Rect(x, 14, 12, 2, ink)
			}
		})
		shell, base, st, picked := pad(c, text, 7, nameEnd+1, 14, pixel.LineHeight, cur, o.Disabled)
		switch {
		case picked == nameDel && len(*name) > 0:
			*name, res.Changed = (*name)[:len(*name)-1], true
		case picked == nameEnd && len(*name) > 0:
			res.Confirmed = true
		case picked >= 0 && picked < nameDel && len(*name) < nameLen:
			*name, res.Changed = *name+string(rune('A'+picked)), true
			if len(*name) == nameLen {
				*cur = nameEnd
			}
		}
		base.Description(*name)
		focused := base.Focused()
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range nameEnd + 1 {
				x, y := i%7*14, i/7*pixel.LineHeight
				if focused && i == *cur {
					cursor(g, x, y, 14, pixel.LineHeight)
				}
				switch i {
				case nameDel:
					caption(g, x, y, k.Msg("del"), pixel.Faint)
				case nameEnd:
					caption(g, x, y, k.Msg("end"), pixel.Accent)
				default:
					l := string(rune('A' + i))
					g.Text(x+(14-pixel.Measure(l))/2, y, l, pixel.Ink, 1)
				}
			}
			finish(g, st)
		})
	})
	if res.Confirmed {
		slog.Debug("8bit: name entered", "label", text, "name", *name)
	}
	return res
}

// caption centers a cell's word, in micro when possible.
func caption(g pixel.Canvas, x, y int, s string, ink pixel.Slot) {
	if pixel.MicroOK(s) {
		g.Micro(x+(14-pixel.MicroWidth(s))/2, y+5, s, ink, 1)
		return
	}
	g.Text(x+(14-pixel.Measure(s))/2, y, s, ink, 1)
}
