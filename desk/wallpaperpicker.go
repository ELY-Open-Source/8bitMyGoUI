package desk

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type WallpaperPickerOptions struct {
	Disabled bool
}

type WallpaperPickerResult struct {
	Element ui.Element
	Changed bool
}

type pickerKey struct{}

const swatchW, swatchH = 44, 36

// WallpaperPicker chooses a dither pattern by its swatch.
func WallpaperPicker(c *ui.Context, text string, chosen *Wallpaper, o WallpaperPickerOptions) WallpaperPickerResult {
	k := pixel.K(c)
	chosen.name(k)
	col := ui.Column(c).Gap(k.D(4))
	cur := ui.Local(col, pickerKey{}, func() int { return int(*chosen) })
	res := WallpaperPickerResult{}
	res.Element = col.Children(func() {
		shell, base, st, picked := pixel.Pad(c, text, 3, len(papers), swatchW, swatchH, cur, o.Disabled)
		if picked >= 0 && Wallpaper(picked) != *chosen {
			*chosen, res.Changed = Wallpaper(picked), true
		}
		at, pick := *cur, *chosen
		said := Wallpaper(at).name(k)
		if Wallpaper(at) == pick {
			said += ", " + k.Msg("chosen")
		}
		base.Description(said)
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range papers {
				x, y := i%3*swatchW, i/3*swatchH
				if i == at && st.Focus {
					g.Cursor(x, y, swatchW, swatchH)
				}
				ring := pixel.Edge
				if Wallpaper(i) == pick {
					ring = pixel.Accent
				}
				g.Rect(x+2, y+2, swatchW-4, swatchH-4, ring)
				paper(g, x+4, y+4, swatchW-8, swatchH-8, Wallpaper(i))
			}
			g.Finish(st)
		})
		pixel.Text(c, Wallpaper(at).name(k), pixel.TextOptions{Slot: pixel.Faint})
	})
	if res.Changed {
		slog.Debug("8bit: wallpaper chosen", "label", text, "wallpaper", int(*chosen))
	}
	return res
}
