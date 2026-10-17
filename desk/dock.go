package desk

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// DockApp is an app on the dock.
type DockApp struct {
	Name string
	Icon *pixel.Sprite
}

type DockOptions struct {
	Disabled bool
}

type DockResult struct {
	Element ui.Element
	// Launched is the app pressed this frame, else -1.
	Launched int
}

type dockKey struct{}

const dockCell = 30

// Dock: a shelf of apps; the pointed or focused rises.
func Dock(c *ui.Context, text string, apps []DockApp, o DockOptions) DockResult {
	k := pixel.K(c)
	if len(apps) == 0 {
		panic("desk: Dock has no apps")
	}
	seen := map[string]bool{}
	for i, a := range apps {
		sprite("Dock", a.Icon)
		if a.Name == "" || seen[a.Name] {
			panic(fmt.Sprintf("desk: Dock app %d %q is nameless or twice", i, a.Name))
		}
		seen[a.Name] = true
	}
	w := len(apps) * dockCell
	col := ui.Column(c).Width(k.D(w + 8)).AlignItems(ui.Center)
	// lit is the risen app, read when the names paint.
	lit := ui.Local(col, dockKey{}, func() int { return -1 })
	*lit = -1
	res := DockResult{Launched: -1}
	res.Element = col.Children(func() {
		ui.Box(c).Size(k.D(w+8), k.D(pixel.LineHeight)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			if *lit < 0 {
				return
			}
			g := k.Canvas(p, r)
			name := pixel.Wrap(apps[*lit].Name, w+8, 1)[0]
			x := max(0, min(g.W-pixel.Measure(name), 4+*lit*dockCell+dockCell/2-pixel.Measure(name)/2))
			g.Text(x, 0, name, pixel.Ink, 1)
		})
		ui.Row(c).Padding(0, k.D(4)).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).Label(text).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Rect(0, g.H-6, g.W, 6, pixel.Face)
			g.HLine(0, g.H-6, g.W, pixel.Lit)
			g.HLine(0, g.H-1, g.W, pixel.Edge)
		}).Children(func() {
			for i, a := range apps {
				// Keyed by name: apps come and go.
				ui.Box(c).Key(a.Name).Children(func() {
					shell, base, st := pixel.Shell(c, func() ui.Element {
						return ui.Box(c).Focusable().Role(ui.RoleButton).Label(a.Name).Disabled(o.Disabled)
					})
					up := st.Hover || st.Focus
					if up {
						*lit = i
					}
					// It rises eight dots in four steps, and sinks back.
					rise := 2 * k.Toward(c, shell, up, 15, 4)
					if !base.IsDisabled() && base.Clicked() {
						res.Launched = i
					}
					shell.Size(k.D(dockCell), k.D(40)).Draw(func(p *ui.Painter, r ui.Rect) {
						g := k.Canvas(p, r)
						g.Sprite((dockCell-2*a.Icon.W)/2, 10-rise+2*(12-a.Icon.H), a.Icon, 0, pixel.Ink, 2)
						g.Finish(st)
					})
				})
			}
		})
	})
	if res.Launched >= 0 {
		slog.Debug("8bit: dock launched", "label", text, "app", apps[res.Launched].Name)
	}
	return res
}
