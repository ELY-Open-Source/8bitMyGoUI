package chat

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type LinkPreviewOptions struct {
	// Width is in dots; zero means 200.
	Width    int
	Disabled bool
}

type LinkPreviewResult struct {
	Element *ui.Element
	Clicked bool
}

// LinkPreview: thumbnail, title and site; a click opens url.
func LinkPreview(c *ui.Context, url, title, site string, o LinkPreviewOptions) LinkPreviewResult {
	k := pixel.K(c)
	w := pixel.Or("LinkPreview Width", o.Width, 200)
	if w < 100 {
		panic("chat: a LinkPreview needs 100 dots")
	}
	lines := pixel.Wrap(title, w-56, 2)
	// The site sits under the title; the card fits both.
	sy := 4 + len(lines)*pixel.LineHeight
	h := max(40, sy+12)
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleLink).Label(title).Description(url).Disabled(o.Disabled).Cursor(ui.CursorPointer).FocusRing(false)
	})
	shell.Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		fill := pixel.Face
		if st.Hover {
			fill = pixel.Lit
		}
		g.Frame(0, 0, w, h, pixel.Window, fill)
		// The thumbnail is a placeholder: dithered, with a picture icon.
		g.Frame(4, 4, 40, 32, pixel.Sunken, pixel.Well)
		g.Dither(6, 6, 36, 28, pixel.Well, pixel.Shade, 4)
		g.Sprite(18, 14, icon.Image, 0, pixel.Faint, 1)
		for i, l := range lines {
			g.Text(50, 3+i*pixel.LineHeight, l, pixel.Ink, 1)
		}
		tag(g, 50, sy, line(site, w-56), pixel.Faint)
		g.Finish(st)
	})
	res := LinkPreviewResult{Element: shell}
	if base.Clicked() {
		c.OpenURL(url)
		res.Clicked = true
		slog.Debug("8bit: link opened", "url", url)
	}
	return res
}
