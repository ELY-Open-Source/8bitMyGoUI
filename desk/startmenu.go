package desk

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Program is a START entry; with Programs, a submenu.
type Program struct {
	Name     string
	Icon     *pixel.Sprite
	Programs []Program
}

type StartMenuOptions struct {
	Disabled bool
}

type StartMenuResult struct {
	Element *ui.Element
	// Launched is the program run this frame, as "GAMES/SNAKE".
	Launched string
}

type startKey struct{}

// start is the open menu: each open level's lit row.
type start struct {
	open bool
	path []int
	ptr  [2]float32
}

const rowH = pixel.LineHeight + 2

// StartMenu is the START button and its cascading programs.
func StartMenu(c *ui.Context, text string, programs []Program, o StartMenuOptions) StartMenuResult {
	k := pixel.K(c)
	checkPrograms(programs)
	off := o.Disabled || pixel.Current(c).Inert
	b := input.Button(c, k.Msg("start"), input.ButtonOptions{Icon: icon.Gamepad, Label: text, Disabled: o.Disabled})
	m := ui.Local(b.Element, startKey{}, func() start { return start{} })
	if off {
		// Closing on disable must not return focus to it.
		if m.open {
			c.Root().Focus()
		}
		m.open = false
	}
	res := StartMenuResult{Element: b.Element}
	// It rises from a taskbar, so Up opens it too.
	switch {
	case b.Clicked && m.open:
		m.open = false
	case b.Clicked, !m.open && !off && b.Element.Shortcut(0, ui.KeyUp):
		m.open, m.path = true, []int{0}
	}
	if m.open {
		m.fit(programs)
		res.Launched = m.keys(b.Element.Shortcut, programs)
	}
	if run := m.list(c, k, b.Element, programs); run != "" {
		res.Launched = run
	}
	if res.Launched != "" {
		slog.Debug("8bit: start menu launched", "label", text, "program", res.Launched)
	}
	return res
}

// checkPrograms fails on blank, slashed or twin names.
func checkPrograms(ps []Program) {
	if len(ps) == 0 {
		panic("desk: a START menu level is empty")
	}
	seen := map[string]bool{}
	for _, p := range ps {
		sprite("StartMenu "+p.Name, p.Icon)
		if p.Name == "" || strings.Contains(p.Name, "/") || seen[p.Name] {
			panic(fmt.Sprintf("desk: START menu program %q is blank, slashed or twice", p.Name))
		}
		seen[p.Name] = true
		if p.Programs != nil {
			checkPrograms(p.Programs)
		}
	}
}

// fit trims the path to the programs now given.
func (m *start) fit(ps []Program) {
	items := ps
	for d, i := range m.path {
		if i >= len(items) {
			m.path = append(m.path[:d], len(items)-1)
			return
		}
		if d+1 < len(m.path) {
			if items = items[i].Programs; items == nil {
				m.path = m.path[:d+1]
				return
			}
		}
	}
}

// level is the programs at a depth of the path.
func level(ps []Program, path []int, depth int) []Program {
	for _, i := range path[:depth] {
		ps = ps[i].Programs
	}
	return ps
}

// launch closes the menu and names the path's program.
func (m *start) launch(ps []Program) string {
	names := make([]string, len(m.path))
	for d, i := range m.path {
		names[d] = level(ps, m.path, d)[i].Name
	}
	m.open = false
	return strings.Join(names, "/")
}

// keys walks the deepest level; Right or Enter goes in.
func (m *start) keys(short func(ui.Modifiers, ui.Key) bool, ps []Program) string {
	down, up, left, right := short(0, ui.KeyDown), short(0, ui.KeyUp), short(0, ui.KeyLeft), short(0, ui.KeyRight)
	enter, space := short(0, ui.KeyEnter), short(0, ui.KeySpace)
	d := len(m.path) - 1
	items := level(ps, m.path, d)
	switch {
	case down:
		m.path[d] = (m.path[d] + 1) % len(items)
	case up:
		m.path[d] = (m.path[d] - 1 + len(items)) % len(items)
	case left && d > 0:
		m.path = m.path[:d]
	case (right || enter || space) && items[m.path[d]].Programs != nil:
		m.path = append(m.path, 0)
	case enter || space:
		return m.launch(ps)
	}
	return ""
}

// list shows each open level beside its parent row.
func (m *start) list(c *ui.Context, k pixel.Kit, anchor *ui.Element, ps []Program) string {
	run := ""
	ui.PopoverBase(c, anchor, &m.open, func(panel *ui.Element) {
		// Only a pointer that moves lights the row under it.
		px, py, over := panel.PointerPosition()
		moved := over && [2]float32{px, py} != m.ptr
		m.ptr = [2]float32{px, py}
		path := slices.Clone(m.path)
		// Submenus grow away from START, so open rows hold still.
		below, align := panel.Bounds().Y > anchor.Bounds().Y, ui.End
		if below {
			align = ui.Start
		}
		panel.Children(func() {
			ui.Row(c).AlignItems(align).Children(func() {
				shift := 0
				for d := range path {
					items := level(ps, path, d)
					col := ui.Column(c).Padding(k.D(3))
					// A submenu lines up with its parent row.
					if d > 0 && below {
						shift += path[d-1] * rowH
						col.Margin(k.D(shift), 0, 0, 0)
					}
					if d > 0 && !below {
						shift += (len(level(ps, path, d-1)) - 1 - path[d-1]) * rowH
						col.Margin(0, 0, k.D(shift), 0)
					}
					col.Draw(func(p *ui.Painter, r ui.Rect) {
						g := k.Canvas(p, r)
						g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
					}).Children(func() {
						ui.Row(c).Children(func() {
							if d == 0 {
								banner(c, k, len(items))
							}
							ui.Column(c).Children(func() {
								for i, p := range items {
									if r := m.row(c, k, p, d, i, moved); r == launched {
										run = m.launch(ps)
									}
								}
							})
						})
					})
				}
			})
		})
	})
	return run
}

type rowEvent int

const (
	quiet rowEvent = iota
	launched
)

// row is one program; hover or click opens a submenu.
func (m *start) row(c *ui.Context, k pixel.Kit, p Program, d, i int, moved bool) rowEvent {
	w := 20 + pixel.Measure(p.Name) + 16
	row := ui.Row(c).Height(k.D(rowH)).MinWidth(k.D(w)).Role(ui.RoleListItem).Label(p.Name)
	// open lights the row and shows its submenu, if any.
	open := func() {
		m.path = append(m.path[:d], i)
		if p.Programs != nil {
			m.path = append(m.path, 0)
		}
	}
	// A level a row before closed takes no input.
	ev, alive := quiet, len(m.path) > d
	shut := alive && m.path[d] == i && len(m.path) == d+1 && p.Programs != nil
	if alive && moved && row.Hovered() && (m.path[d] != i || shut) {
		open()
		c.Invalidate()
	}
	if alive && row.Clicked() {
		open()
		if p.Programs == nil {
			ev = launched
		}
	}
	row.Draw(func(pt *ui.Painter, r ui.Rect) {
		g := k.Canvas(pt, r)
		ink := pixel.Ink
		if len(m.path) > d && m.path[d] == i {
			g.Rect(0, 0, g.W, g.H, pixel.Accent)
			ink = pixel.AccentInk
		}
		g.Sprite(2, (g.H-p.Icon.H)/2, p.Icon, 0, ink, 1)
		g.Text(18, 1, p.Name, ink, 1)
		if p.Programs != nil {
			for j := range 4 {
				g.VLine(g.W-8+j, g.H/2-3+j, 7-2*j, ink)
			}
		}
	})
	return ev
}

// banner is START's dithered side strip.
func banner(c *ui.Context, k pixel.Kit, rows int) {
	ui.Box(c).Size(k.D(10), k.D(rows*rowH)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for b := range 17 {
			top, foot := b*g.H/17, (b+1)*g.H/17
			g.Dither(0, top, g.W, foot-top, pixel.Accent, pixel.Magic, 16-b)
		}
	})
}
