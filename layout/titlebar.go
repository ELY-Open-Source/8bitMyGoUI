package layout

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var (
	closeGlyph = pixel.MustSprite(`
		#...#
		.#.#.
		..#..
		.#.#.
		#...#
	`)
	minimizeGlyph = pixel.MustSprite(`
		.....
		.....
		.....
		.....
		#####
	`)
)

type TitleBarOptions struct {
	// NoMinimize leaves out the minimize box.
	NoMinimize bool
	Disabled   bool
}

type TitleBarResult struct {
	Element  ui.Element
	Close    bool
	Minimize bool
}

// TitleBar is a striped bar: title, close box, minimize box.
func TitleBar(c *ui.Context, title string, o TitleBarOptions) TitleBarResult {
	k := pixel.K(c)
	tw := pixel.Measure(title)
	res := TitleBarResult{}
	res.Element = ui.Row(c).FillWidth().Height(k.D(pixel.LineHeight+2)).AlignItems(ui.Center).Padding(0, k.D(4)).Role(ui.RoleGroup).Label(title).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Face)
		g.HLine(0, g.H-1, g.W, pixel.Edge)
		for y := 4; y < g.H-4; y += 2 {
			g.HLine(2, y, g.W-4, pixel.Shade)
		}
		x := (g.W - tw) / 2
		g.Rect(x-6, 0, tw+12, g.H-1, pixel.Face)
		g.Text(x, 1, title, pixel.Ink, 1)
	}).Children(func() {
		res.Close = windowBox(c, k, closeGlyph, k.Msg("close"), o.Disabled)
		ui.Spacer(c)
		if !o.NoMinimize {
			res.Minimize = windowBox(c, k, minimizeGlyph, k.Msg("minimize"), o.Disabled)
		}
	})
	if res.Close || res.Minimize {
		slog.Debug("8bit: titlebar", "title", title, "close", res.Close, "minimize", res.Minimize)
	}
	return res
}

// windowBox is a framed button showing its glyph on hover.
func windowBox(c *ui.Context, k pixel.Kit, glyph *pixel.Sprite, name string, disabled bool) bool {
	shell, base, st := pixel.Shell(c, func() ui.Element { return ui.ButtonBase(c).Label(name).Disabled(disabled) })
	shell.Size(k.D(11), k.D(11)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		fill := pixel.Face
		if st.Press {
			fill = pixel.Shade
		}
		box(g, pixel.Flat, fill)
		if st.Hover || st.Press || st.Focus {
			g.Sprite(3, 3, glyph, 0, pixel.Ink, 1)
		}
		if st.Disabled {
			g.Checker(0, 0, g.W, g.H, pixel.Face)
		}
		if st.Focus {
			g.FocusMark(0, 0, g.W, g.H)
		}
	})
	return base.Clicked()
}
