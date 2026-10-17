package data

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TreeNode struct {
	Label string
	// Value follows the label in Slot, as a key's value.
	Value    string
	Slot     pixel.Slot
	Children []TreeNode
}

type TreeOptions struct {
	// Width is in dots; zero means 160.
	Width    int
	Disabled bool
}

type TreeResult struct {
	Element          ui.Element
	Changed, Toggled bool
}

type treeKey struct{}

// Tree of ▸▾ items; *selected is a path, "0/2".
func Tree(c *ui.Context, text string, roots []TreeNode, selected *string, o TreeOptions) TreeResult {
	k := pixel.K(c)
	nodes := map[string]TreeNode{}
	var index func(ns []TreeNode, at string)
	index = func(ns []TreeNode, at string) {
		for i, n := range ns {
			p := at + strconv.Itoa(i)
			nodes[p] = n
			index(n.Children, p+"/")
		}
	}
	index(roots, "")
	children := func(p string) []string {
		var out []string
		for i := range nodes[p].Children {
			out = append(out, p+"/"+strconv.Itoa(i))
		}
		return out
	}
	var tops []string
	for i := range roots {
		tops = append(tops, strconv.Itoa(i))
	}
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleTree).Label(text).Disabled(o.Disabled)
	})
	opened := ui.Local(base, treeKey{}, func() map[string]bool { return map[string]bool{} })
	// A choice set inside closed branches opens its way.
	if _, ok := nodes[*selected]; ok {
		for i, r := range *selected {
			if r == '/' {
				(*opened)[(*selected)[:i]] = true
			}
		}
	}
	rows := flatten(tops, children, *opened)
	if len(rows) == 0 {
		return TreeResult{Element: shell}
	}
	// An unknown choice stays, unlit, until a key or click.
	at := -1
	for i, r := range rows {
		if r.key == *selected {
			at = i
		}
	}
	old := *selected
	res := TreeResult{Element: shell}
	if pixel.Shortcut(base, 0, ui.KeyDown) {
		at = min(at+1, len(rows)-1)
	}
	if pixel.Shortcut(base, 0, ui.KeyUp) {
		at = max(at-1, 0)
	}
	right, left := pixel.Shortcut(base, 0, ui.KeyRight), pixel.Shortcut(base, 0, ui.KeyLeft)
	if (right || left) && at >= 0 {
		at, res.Toggled = steer(rows, at, *opened, right)
	}
	enter, space := pixel.Shortcut(base, 0, ui.KeyEnter), pixel.Shortcut(base, 0, ui.KeySpace)
	toggle := func(i int) {
		if key := rows[i].key; rows[i].branch {
			(*opened)[key], res.Toggled = !(*opened)[key], true
		}
	}
	if (enter || space) && at >= 0 {
		toggle(at)
	}
	// Enter and Space are taken, so clicks come from pointers.
	if base.Clicked() {
		x, y, _ := base.PointerPosition()
		if i := int(y/k.Dot()) / pixel.LineHeight; i >= 0 && i < len(rows) {
			at = i
			if int(x/k.Dot()) < rows[i].depth*10+12 {
				toggle(i)
			}
		}
	}
	if res.Toggled {
		rows = flatten(tops, children, *opened)
		at = min(at, len(rows)-1)
	}
	if at >= 0 {
		*selected = rows[at].key
	}
	res.Changed = *selected != old
	base.Description(strings.TrimSpace(nodes[*selected].Label + " " + nodes[*selected].Value))
	shown := rows
	shell.Size(k.D(pixel.Or("Tree Width", o.Width, 160)), k.D(len(shown)*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, rw := range shown {
			y, x, ink := i*pixel.LineHeight, rw.depth*10, pixel.Ink
			if i == at {
				g.Rect(0, y, g.W, pixel.LineHeight, pixel.Accent)
				ink = pixel.AccentInk
			}
			if rw.branch {
				if (*opened)[rw.key] {
					g.Sprite(x+3, y+6, down, 0, ink, 1)
				} else {
					g.Sprite(x+4, y+5, pointer, 0, ink, 1)
				}
			}
			n := nodes[rw.key]
			g.Text(x+12, y, n.Label, ink, 1)
			// On the chosen row the value takes its ink.
			if vs := n.Slot; n.Value != "" {
				if i == at {
					vs = ink
				}
				g.Text(x+18+pixel.Measure(n.Label), y, n.Value, vs, 1)
			}
		}
		g.Finish(st)
	})
	if res.Changed || res.Toggled {
		slog.Debug("8bit: tree", "label", text, "path", *selected, "toggled", res.Toggled)
	}
	return res
}
