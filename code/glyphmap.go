package code

import (
	"fmt"
	"log/slog"
	"unicode"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type GlyphMapOptions struct {
	Disabled bool
}

type GlyphMapResult struct {
	Element *ui.Element
	Changed bool
}

// GlyphMap pages code points by 128; the chosen shows large.
func GlyphMap(c *ui.Context, text string, at *rune, o GlyphMapOptions) GlyphMapResult {
	k := pixel.K(c)
	if *at < 0 || *at > unicode.MaxRune {
		panic(fmt.Sprintf("code: GlyphMap at %U", *at))
	}
	old := *at
	var res GlyphMapResult
	res.Element = ui.Row(c).Gap(k.D(8)).AlignItems(ui.Start).Children(func() {
		page := *at &^ 127
		cur := int(*at - page)
		shell, base, st, _ := pixel.Pad(c, text, 16, 128, 14, 14, &cur, o.Disabled)
		*at = page + rune(cur)
		// PageUp and PageDown turn a page, keeping the column.
		if !base.IsDisabled() {
			if base.Shortcut(0, ui.KeyPageUp) && *at >= 128 {
				*at -= 128
			}
			if base.Shortcut(0, ui.KeyPageDown) && *at+128 <= unicode.MaxRune {
				*at += 128
			}
		}
		base.Description(name(*at))
		shown := *at
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			first := shown &^ 127
			for i := range 128 {
				x, y, r := i%16*14, i/16*14, first+rune(i)
				ink := pixel.Ink
				if r == shown {
					g.Rect(x, y, 14, 14, pixel.Accent)
					ink = pixel.AccentInk
				}
				if drawn(r) {
					s := string(r)
					g.Text(x+(14-pixel.Measure(s))/2, y+1, s, ink, 1)
				} else {
					g.Dot(x+7, y+7, pixel.Faint)
				}
			}
			g.Finish(st)
		})
		ui.Column(c).Gap(k.D(3)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(k.D(48), k.D(48)).Role(ui.RoleImage).Label(name(shown)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Frame(0, 0, 48, 48, pixel.Sunken, pixel.Well)
				if drawn(shown) {
					s := string(shown)
					g.Text(24-pixel.Measure(s)*3/2, 6, s, pixel.Ink, 3)
				}
			})
			code := fmt.Sprintf("U+%04X", shown)
			slot := pixel.Ink
			if !pixel.Has(shown) {
				slot = pixel.Bad
			}
			ui.Box(c).Size(k.D(pixel.MicroWidth(code)), k.D(5)).Draw(func(p *ui.Painter, r ui.Rect) {
				k.Canvas(p, r).Micro(0, 0, code, slot, 1)
			})
		})
	})
	if res.Changed = *at != old; res.Changed {
		slog.Debug("8bit: glyph map", "label", text, "at", fmt.Sprintf("%U", *at))
	}
	return res
}

// drawn is a graphic rune the font holds.
func drawn(r rune) bool {
	return utf8.ValidRune(r) && (unicode.IsGraphic(r) || r == ' ') && pixel.Has(r)
}

// name reads a code point aloud, with its glyph.
func name(r rune) string {
	if drawn(r) {
		return fmt.Sprintf("%U %c", r, r)
	}
	return fmt.Sprintf("%U", r)
}
