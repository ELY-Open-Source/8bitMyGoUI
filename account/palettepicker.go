package account

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PalettePickerOptions struct {
	Disabled bool
}

type PalettePickerResult struct {
	Element *ui.Element
	Changed bool
}

const tileW, tileH = 56, 42

// PalettePicker shows each palette as a live scene; choosing applies.
func PalettePicker(c *ui.Context, text string, chosen *int, o PalettePickerOptions) PalettePickerResult {
	k := pixel.K(c)
	choice("PalettePicker", chosen, len(pixel.Palettes))
	// Each tile draws with a kit in its own palette.
	kits := make([]pixel.Kit, len(pixel.Palettes))
	for i, pal := range pixel.Palettes {
		in := pixel.Current(c)
		in.Palette = pal
		pixel.Scope(c, in, func() { kits[i] = pixel.K(c) })
	}
	old := *chosen
	// Arrows choose at once; the cursor is the choice.
	shell, base, st, _ := pixel.Pad(c, text, 4, len(pixel.Palettes), tileW, tileH, chosen, o.Disabled)
	pal := pixel.Palettes[*chosen]
	name := pal.Name
	if k.Locale() == pixel.Zh {
		name = pal.Zh
	}
	base.Description(name + ", " + k.Msg("in_use"))
	on := *chosen
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		for i, tk := range kits {
			tile(tk.Canvas(p, r), i%4*tileW+2, i/4*tileH+2, pixel.Palettes[i].Night)
		}
		g := k.Canvas(p, r)
		x, y := on%4*tileW, on/4*tileH
		for _, d := range []int{0, 1} {
			g.Rect(x+d, y+d, tileW-2*d, 1, pixel.Accent)
			g.Rect(x+d, y+tileH-1-d, tileW-2*d, 1, pixel.Accent)
			g.Rect(x+d, y+d, 1, tileH-2*d, pixel.Accent)
			g.Rect(x+tileW-1-d, y+d, 1, tileH-2*d, pixel.Accent)
		}
		g.Finish(st)
	})
	res := PalettePickerResult{Element: shell, Changed: *chosen != old}
	if res.Changed {
		slog.Debug("8bit: palette picker", "label", text, "palette", pixel.Palettes[*chosen].Name)
	}
	return res
}

// tile is a tiny scene in the canvas's palette.
func tile(g pixel.Canvas, x, y int, night bool) {
	g.Frame(x, y, 52, 38, pixel.Window, pixel.Back)
	g.Rect(x+3, y+3, 46, 7, pixel.Accent)
	g.Micro(x+5, y+4, "8BIT", pixel.AccentInk, 1)
	if night {
		g.Disc(x+40, y+17, 3, pixel.Lit)
	} else {
		g.Disc(x+40, y+17, 3, pixel.Warn)
	}
	g.Dither(x+3, y+26, 46, 9, pixel.Good, pixel.Shade, 4)
	g.Rect(x+8, y+18, 10, 8, pixel.Face)
	g.Rect(x+10, y+21, 3, 5, pixel.Edge)
}
