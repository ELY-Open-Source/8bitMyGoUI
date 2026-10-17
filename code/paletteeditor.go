package code

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PaletteEditorOptions struct {
	Disabled bool
}

type PaletteEditorResult struct {
	Element ui.Element
	Changed bool
}

type tintKey struct{}

// PaletteEditor sets *pal in channel levels; a preview follows.
func PaletteEditor(c *ui.Context, text string, pal *pixel.Palette, o PaletteEditorOptions) PaletteEditorResult {
	k := pixel.K(c)
	if pal == nil {
		panic("code: PaletteEditor needs a palette")
	}
	// Inside the scope everything draws in the edited colors.
	in := pixel.Current(c)
	in.Palette = pal
	var res PaletteEditorResult
	res.Element = ui.Row(c).Gap(k.D(8)).AlignItems(ui.Start)
	sel := ui.Local(res.Element, tintKey{}, func() int { return int(pixel.Accent) })
	res.Element.Children(func() {
		ui.Column(c).Gap(k.D(4)).Children(func() {
			pixel.Scope(c, in, func() {
				// Paint runs after the scope: keep its kit now.
				sk := pixel.K(c)
				shell, base, st, _ := pixel.Pad(c, text, 8, int(pixel.Slots), 14, 14, sel, o.Disabled)
				base.Description(fmt.Sprintf("%v, %s", pixel.Slot(*sel), hex(pal.Colors[*sel])))
				at := *sel
				shell.Draw(func(p *ui.Painter, r ui.Rect) {
					g := sk.Canvas(p, r)
					for s := range int(pixel.Slots) {
						x, y := s%8*14, s/8*14
						if s == at {
							g.Rect(x, y, 14, 14, pixel.Accent)
						}
						g.Frame(x+1, y+1, 12, 12, pixel.Flat, pixel.Slot(s))
					}
					g.Finish(st)
				})
			})
			col := &pal.Colors[*sel]
			pixel.Text(c, fmt.Sprintf("%v %s", pixel.Slot(*sel), hex(*col)), pixel.TextOptions{Slot: pixel.Faint})
			for _, ch := range []struct {
				name string
				v    *uint8
			}{{"R", &col.R}, {"G", &col.G}, {"B", &col.B}} {
				// Levels 0–15, each worth 17.
				level := math.Round(float64(*ch.v) / 17)
				ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
					// Drawn, not named: the slider carries the name.
					ui.Box(c).Size(k.D(pixel.Measure(ch.name)), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
						k.Canvas(p, r).Text(0, 0, ch.name, pixel.Faint, 1)
					})
					if input.StepSlider(c, ch.name, &level, input.StepSliderOptions{Max: 15, Step: 1, Width: 96, Disabled: o.Disabled}).Changed {
						*ch.v, res.Changed = uint8(level)*17, true
					}
				})
			}
		})
		pixel.Scope(c, in, func() { preview(c, k.Msg("preview")) })
	})
	if res.Changed {
		slog.Debug("8bit: palette editor", "label", text, "slot", pixel.Slot(*sel), "color", hex(pal.Colors[*sel]))
	}
	return res
}

// preview is a small window drawn in the edited colors.
func preview(c *ui.Context, text string) {
	k := pixel.K(c)
	ui.Box(c).Size(k.D(112), k.D(84)).Role(ui.RoleImage).Label(text).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Back)
		g.Frame(0, 0, 112, 84, pixel.Window, pixel.Face)
		g.Rect(3, 3, 106, 12, pixel.Accent)
		g.Text(6, 3, "HELLO!", pixel.AccentInk, 1)
		g.Text(6, 18, "PRESS START", pixel.Ink, 1)
		g.Text(6, 30, "HI 9999", pixel.Faint, 1)
		g.Frame(6, 44, 46, 14, pixel.Raised, pixel.Face)
		g.Text(12, 45, "OK", pixel.Ink, 1)
		g.Frame(58, 44, 46, 14, pixel.Sunken, pixel.Well)
		for i, s := range []pixel.Slot{pixel.Good, pixel.Warn, pixel.Bad, pixel.Info, pixel.Magic, pixel.AccentLit} {
			g.Rect(6+i*17, 66, 15, 10, s)
		}
	})
}

// hex prints a color as #RRGGBB.
func hex(col ui.Color) string { return fmt.Sprintf("#%02X%02X%02X", col.R, col.G, col.B) }
