package data

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type OutlineOptions struct {
	// Height is in dots; zero means 140.
	Height   int
	Disabled bool
}

type OutlineResult struct {
	Element            ui.Element
	Changed, Submitted bool
}

type outlineKey[K comparable] struct{}

// outline is a tree's open branches and its chosen row.
type outline[K comparable] struct {
	opened map[K]bool
	list   ui.ListState
	at     int
}

// Outline is a virtual tree for huge hierarchies.
func Outline[K comparable](c *ui.Context, text string, roots []K, children func(K) []K, label func(K) string, selected *K, o OutlineOptions) OutlineResult {
	k := pixel.K(c)
	off := o.Disabled || pixel.Current(c).Inert
	box := ui.Box(c)
	st := ui.Local(box, outlineKey[K]{}, func() outline[K] { return outline[K]{opened: map[K]bool{}} })
	rows := flatten(roots, children, st.opened)
	st.at = -1
	for i, r := range rows {
		if r.key == *selected {
			st.at = i
		}
	}
	st.list.Selected = &st.at
	st.list.Key = func(i int) any { return rows[i].key }
	st.list.Label = func(i int) string { return label(rows[i].key) }
	if off {
		st.list.Selected = nil
	}
	var e ui.Element
	box.Children(func() {
		e = ui.List(c, &st.list, len(rows), pixel.Lazy(c, func(i int) {
			r, on, ink := rows[i], i == st.at, pixel.Ink
			if on {
				ink = pixel.AccentInk
			}
			ui.Row(c).Height(k.D(pixel.LineHeight)).AlignItems(ui.Center).Draw(func(p *ui.Painter, rc ui.Rect) {
				if on {
					g := k.Canvas(p, rc)
					g.Rect(0, 0, g.W, g.H, pixel.Accent)
				}
			}).Children(func() {
				ui.Box(c).Width(k.D(r.depth * 10)).Shrink(0)
				arrow := ui.Box(c).Size(k.D(12), k.D(pixel.LineHeight)).Shrink(0).Role(ui.RoleNone)
				if r.branch && !off && arrow.Clicked() {
					st.opened[r.key] = !st.opened[r.key]
					slog.Debug("8bit: outline toggled", "label", text, "row", i, "open", st.opened[r.key])
				}
				shut := !st.opened[r.key]
				arrow.Draw(func(p *ui.Painter, rc ui.Rect) {
					g := k.Canvas(p, rc)
					switch {
					case !r.branch:
					case shut:
						g.Sprite(4, 3, pointer, 0, ink, 1)
					default:
						g.Sprite(3, 4, down, 0, ink, 1)
					}
				})
				pixel.Text(c, label(r.key), pixel.TextOptions{Slot: ink})
			})
		})).Height(k.D(pixel.Or("Outline Height", o.Height, 140))).Label(text)
		right, left := pixel.Shortcut(e, 0, ui.KeyRight), pixel.Shortcut(e, 0, ui.KeyLeft)
		if (right || left) && st.at >= 0 && !off {
			at, toggled := steer(rows, st.at, st.opened, right)
			if toggled {
				slog.Debug("8bit: outline toggled", "label", text, "row", at, "open", right)
			}
			st.at, rows = at, flatten(roots, children, st.opened)
		}
		veil(c, off)
	})
	pixel.Park(c, e, off)
	if off {
		e.Description(k.Msg("disabled"))
	}
	well(c, box, pixel.State{Focus: e.FocusRing(false).FocusVisible(), Disabled: o.Disabled})
	res := OutlineResult{Element: box, Submitted: e.Submitted()}
	if st.at >= 0 && rows[st.at].key != *selected {
		*selected, res.Changed = rows[st.at].key, true
		slog.Debug("8bit: outline", "label", text, "row", st.at)
	}
	return res
}
