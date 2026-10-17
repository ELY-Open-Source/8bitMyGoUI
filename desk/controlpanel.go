package desk

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Applet is a setting in the control panel.
type Applet struct {
	Name string
	Icon *pixel.Sprite
}

type ControlPanelOptions struct {
	Disabled bool
}

type ControlPanelResult struct {
	Element ui.Element
	// Opened is the applet opened this frame, else -1.
	Opened int
}

type panelKey struct{}

const appletW, appletH = 66, 46

// ControlPanel is a settings grid; Enter or a click opens.
func ControlPanel(c *ui.Context, text string, applets []Applet, o ControlPanelOptions) ControlPanelResult {
	k := pixel.K(c)
	if len(applets) == 0 {
		panic("desk: ControlPanel has no applets")
	}
	lines := make([][]string, len(applets))
	for i, a := range applets {
		sprite("ControlPanel", a.Icon)
		if a.Name == "" {
			panic(fmt.Sprintf("desk: ControlPanel applet %d has no name", i))
		}
		lines[i] = pixel.Wrap(a.Name, appletW-4, 1)
	}
	box := ui.Box(c)
	cur := ui.Local(box, panelKey{}, func() int { return 0 })
	*cur = min(*cur, len(applets)-1)
	res := ControlPanelResult{Opened: -1}
	res.Element = box.Children(func() {
		shell, base, st, picked := pixel.Pad(c, text, min(len(applets), 4), len(applets), appletW, appletH, cur, o.Disabled)
		res.Opened = picked
		at := *cur
		base.Description(applets[at].Name)
		cols := min(len(applets), 4)
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, a := range applets {
				x, y := i%cols*appletW, i/cols*appletH
				lit := i == at && st.Focus
				if lit {
					g.Cursor(x, y, appletW, appletH)
				}
				tile(g, x, y, appletW, a.Icon, lines[i], false)
			}
			g.Finish(st)
		})
	})
	if res.Opened >= 0 {
		slog.Debug("8bit: control panel opened", "label", text, "applet", applets[res.Opened].Name)
	}
	return res
}
