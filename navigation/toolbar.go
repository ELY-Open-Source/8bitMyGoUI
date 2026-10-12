package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ToolbarOptions struct{}

// Toolbar holds tool buttons as one Tab stop; arrows move.
func Toolbar(c *ui.Context, text string, o ToolbarOptions, content func()) *ui.Element {
	k := pixel.K(c)
	return ui.Row(c).AlignSelf(ui.Start).Gap(k.D(2)).Padding(k.D(3)).AlignItems(ui.Center).Role(ui.RoleToolbar).Label(text).FocusGroup(ui.Horizontal).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Raised, pixel.Face)
	}).Children(content)
}

type BackForwardOptions struct {
	Disabled bool
}

type BackForwardResult struct {
	Element *ui.Element
	// Moved reports a step through the history.
	Moved bool
}

// BackForward walks a router's history with ◀ and ▶.
func BackForward(c *ui.Context, r *ui.Router, o BackForwardOptions) BackForwardResult {
	k := pixel.K(c)
	res := BackForwardResult{}
	res.Element = ui.Row(c).Gap(k.D(2)).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).Label(k.Msg("history")).Children(func() {
		back := input.IconButtonOptions{Size: input.Small, Disabled: o.Disabled || !r.CanGoBack()}
		if input.IconButton(c, icon.ChevronLeft, k.Msg("back"), back).Clicked {
			r.Back()
			res.Moved = true
		}
		forth := input.IconButtonOptions{Size: input.Small, Disabled: o.Disabled || !r.CanGoForward()}
		if input.IconButton(c, icon.ChevronRight, k.Msg("forward"), forth).Clicked {
			r.Forward()
			res.Moved = true
		}
	})
	if res.Moved {
		slog.Debug("8bit: history", "path", r.Path())
	}
	return res
}
