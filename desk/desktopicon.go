package desk

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DesktopIconOptions struct {
	// Chosen inverts it, as a chosen icon does.
	Chosen   bool
	Disabled bool
}

type DesktopIconResult struct {
	Element *ui.Element
	// Clicked chooses it; Opened is a double click or Enter.
	Clicked, Opened bool
}

const iconW, iconH = 72, 60

// DesktopIcon is an icon over its label; chosen, it inverts.
func DesktopIcon(c *ui.Context, label string, icon *pixel.Sprite, o DesktopIconOptions) DesktopIconResult {
	k := pixel.K(c)
	sp := sprite("DesktopIcon", icon)
	if label == "" {
		panic("desk: DesktopIcon has no label")
	}
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleButton).Label(label).Disabled(o.Disabled)
	})
	res := DesktopIconResult{Element: shell}
	if !base.IsDisabled() {
		res.Opened = pixel.Shortcut(base, 0, ui.KeyEnter) || base.DoubleClicked()
		res.Clicked = base.Clicked()
	}
	if res.Clicked || res.Opened {
		slog.Debug("8bit: desktop icon chosen", "label", label, "opened", res.Opened)
	}
	lines := pixel.Wrap(label, iconW-4, 2)
	shell.Size(k.D(iconW), k.D(iconH)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		tile(g, 0, 0, iconW, sp, lines, o.Chosen)
		g.Finish(st)
	})
	return res
}
