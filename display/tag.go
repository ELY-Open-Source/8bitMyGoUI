package display

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var cross = pixel.MustSprite(`
	#...#
	.#.#.
	..#..
	.#.#.
	#...#
`)

type TagOptions struct {
	Icon      *pixel.Sprite
	Removable bool
	// Slot is the chip fill; zero means Face.
	Slot     pixel.Slot
	Disabled bool
}

type TagResult struct {
	Element *ui.Element
	Removed bool
}

// Tag is a chip with text, icon and remove box.
func Tag(c *ui.Context, label string, o TagOptions) TagResult {
	k := pixel.K(c)
	fill := slotOr(o.Slot, pixel.Face)
	res := TagResult{}
	res.Element = ui.Row(c).Height(k.D(pixel.LineHeight+2)).Padding(0, k.D(4)).Gap(k.D(3)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(label).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		if g.W >= 3 && g.H >= 3 {
			g.Frame(0, 0, g.W, g.H, pixel.Notched, fill)
		}
	}).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if o.Disabled {
			g := k.Canvas(p, r)
			g.Checker(0, 0, g.W, g.H, fill)
		}
	}).Children(func() {
		if o.Icon != nil {
			Icon(c, o.Icon, IconOptions{})
		}
		pixel.Text(c, label, pixel.TextOptions{})
		if o.Removable {
			res.Removed = removeBox(c, k, k.Msg("remove")+" "+label, o.Disabled)
		}
	})
	if res.Removed {
		slog.Debug("8bit: tag removed", "label", label)
	}
	return res
}

// removeBox is the × that removes a tag.
func removeBox(c *ui.Context, k pixel.Kit, name string, disabled bool) bool {
	shell, base, st := pixel.Shell(c, func() *ui.Element { return ui.ButtonBase(c).Label(name).Disabled(disabled) })
	shell.Size(k.D(9), k.D(9)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		if st.Hover || st.Press {
			g.Rect(0, 0, g.W, g.H, pixel.Shade)
		}
		g.Sprite(2, 2, cross, 0, pixel.Ink, 1)
		if st.Focus {
			g.FocusMark(0, 0, g.W, g.H)
		}
	})
	return base.Clicked()
}
