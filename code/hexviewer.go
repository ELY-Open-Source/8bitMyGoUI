package code

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type HexViewerOptions struct {
	// Columns are bytes a row, 4–16; zero means 8.
	Columns  int
	Disabled bool
}

type HexViewerResult struct {
	Element *ui.Element
	Changed bool
}

// HexViewer shows offsets, hex and ASCII; arrows move a cursor.
func HexViewer(c *ui.Context, text string, data []byte, cursor *int, o HexViewerOptions) HexViewerResult {
	k := pixel.K(c)
	cols := pixel.Or("HexViewer Columns", o.Columns, 8)
	if cols < 4 || cols > 16 || len(data) == 0 || *cursor < 0 || *cursor >= len(data) {
		panic(fmt.Sprintf("code: HexViewer of %d bytes, %d a row, cursor %d", len(data), cols, *cursor))
	}
	rows := (len(data) + cols - 1) / cols
	var res HexViewerResult
	res.Element = ui.Row(c).Gap(k.D(6)).Children(func() {
		ui.Box(c).Size(k.D(15), k.D(rows*12)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range rows {
				g.Micro(0, i*12+4, fmt.Sprintf("%04X", i*cols), pixel.Faint, 1)
			}
		})
		old := *cursor
		shell, base, st, _ := pixel.Pad(c, text, cols, len(data), 12, 12, cursor, o.Disabled)
		b := data[*cursor]
		base.Description(fmt.Sprintf("%04X: %02X %q", *cursor, b, glyph(b)))
		at := *cursor
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, v := range data {
				x, y, ink := i%cols*12, i/cols*12, pixel.Ink
				switch {
				case i == at:
					g.Rect(x, y, 11, 12, pixel.Accent)
					ink = pixel.AccentInk
				case v == 0:
					ink = pixel.Faint
				}
				g.Micro(x+2, y+4, fmt.Sprintf("%02X", v), ink, 1)
			}
			g.Finish(st)
		})
		ui.Box(c).Size(k.D(cols*7), k.D(rows*12)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, v := range data {
				x, y, ch := i%cols*7, i/cols*12, string(glyph(v))
				ink := pixel.Info
				if i == at {
					g.Rect(x, y, 7, 12, pixel.Accent)
					ink = pixel.AccentInk
				}
				g.Text(x+(7-pixel.Measure(ch))/2, y, ch, ink, 1)
			}
		})
		if res.Changed = *cursor != old; res.Changed {
			slog.Debug("8bit: hex viewer", "label", text, "offset", *cursor)
		}
	})
	return res
}

// glyph is b as a character; unprintable is a dot.
func glyph(b byte) rune {
	if b < 0x20 || b > 0x7E {
		return '.'
	}
	return rune(b)
}
