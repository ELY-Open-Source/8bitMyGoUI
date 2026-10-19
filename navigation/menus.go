package navigation

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type MenuItem struct {
	Label string
	// Icon and Keys, a shortcut's name, sit beside the label.
	Icon     *pixel.Sprite
	Keys     string
	Disabled bool
	// Separator is a rule, not an item.
	Separator bool
}

type menuKey struct{}

// menu is a popup's state: open, lit item, pointer.
type menu struct {
	open bool
	lit  int
	ptr  [2]float32
}

// show opens the menu with its first item lit.
func (m *menu) show(items []MenuItem) {
	m.open, m.lit = true, -1
	m.step(items, 1)
}

// step moves the light d items, past rules and disabled.
func (m *menu) step(items []MenuItem, d int) {
	for i, n := m.lit, 0; n < len(items); n++ {
		i = ((i+d)%len(items) + len(items)) % len(items)
		if !items[i].Separator && !items[i].Disabled {
			m.lit = i
			return
		}
	}
}

// keys moves the light; Enter or Space picks, else -1.
func (m *menu) keys(short func(ui.Modifiers, ui.Key) bool, items []MenuItem) int {
	if !m.open {
		return -1
	}
	down, up := short(0, ui.KeyDown), short(0, ui.KeyUp)
	enter, space := short(0, ui.KeyEnter), short(0, ui.KeySpace)
	switch {
	case down:
		m.step(items, 1)
	case up:
		m.step(items, -1)
	case (enter || space) && m.lit >= 0 && !items[m.lit].Disabled:
		m.open = false
		return m.lit
	}
	return -1
}

// list shows the menu under anchor and returns clicks.
func (m *menu) list(c *ui.Context, anchor ui.Element, items []MenuItem) (panel ui.Element, clicked int) {
	k := pixel.K(c)
	clicked = -1
	panel = ui.PopoverBase(c, anchor, &m.open, func(panel ui.Element) {
		// Only a pointer that moves lights the item under it.
		px, py, over := panel.PointerPosition()
		moved := over && [2]float32{px, py} != m.ptr
		m.ptr = [2]float32{px, py}
		panel.Padding(k.D(3)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
		}).Children(func() {
			ui.Column(c).AlignItems(ui.Stretch).Children(func() {
				for i, it := range items {
					if it.Separator {
						ui.Box(c).Height(k.D(5)).Draw(func(p *ui.Painter, r ui.Rect) {
							g := k.Canvas(p, r)
							g.HLine(2, 2, g.W-4, pixel.Shade)
						})
						continue
					}
					row := ui.Row(c).Role(ui.RoleListItem).Label(it.Label).Height(k.D(pixel.LineHeight)).MinWidth(k.D(menuWidth(it)))
					if it.Disabled {
						row.Description(k.Msg("disabled"))
					}
					if moved && row.Hovered() && !it.Disabled {
						m.lit = i
					}
					// Any click closes; MyGo hands focus back to the opener.
					if row.Clicked() {
						m.open = false
						if !it.Disabled {
							clicked = i
						}
					}
					lit := i == m.lit
					row.Draw(func(p *ui.Painter, r ui.Rect) { menuRow(k.Canvas(p, r), it, lit) })
				}
			})
		})
	})
	// Clicks on a rule or the edge close it too.
	if panel.Valid() && panel.Clicked() {
		m.open = false
	}
	return panel, clicked
}

// menuWidth is a row's room for icon, label and keys.
func menuWidth(it MenuItem) int {
	w := 26 + pixel.Measure(it.Label)
	if it.Keys != "" {
		w += 12 + keysWidth(it.Keys)
	}
	return w
}

// keysWidth measures keys in micro when the font has them.
func keysWidth(s string) int {
	if pixel.MicroOK(s) {
		return pixel.MicroWidth(s)
	}
	return pixel.Measure(s)
}

// menuRow draws an item; a hand marks the lit one.
func menuRow(g pixel.Canvas, it MenuItem, lit bool) {
	ink := pixel.Ink
	if it.Disabled {
		ink = pixel.Faint
	}
	if lit {
		g.Rect(0, 0, g.W, g.H, pixel.Lit)
		g.Sprite(0, (g.H-pixel.Hand.H)/2, pixel.Hand, 0, pixel.Edge, 1)
	}
	if it.Icon != nil {
		g.Sprite(12, (g.H-it.Icon.H)/2, it.Icon, 0, ink, 1)
	}
	g.Text(26, 0, it.Label, ink, 1)
	switch {
	case it.Keys == "":
	case pixel.MicroOK(it.Keys):
		g.Micro(g.W-pixel.MicroWidth(it.Keys)-2, 6, it.Keys, pixel.Faint, 1)
	default:
		g.Text(g.W-pixel.Measure(it.Keys)-2, 0, it.Keys, pixel.Faint, 1)
	}
}
