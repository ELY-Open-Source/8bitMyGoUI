package input

import (
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MultiSelectOptions struct {
	// Width is in dots; zero means 160.
	Width    int
	Disabled bool
}

type MultiSelectResult struct {
	Element *ui.Element
	Changed bool
}

type multiKey struct{}

// MultiSelect checks many options in a popup; tags show them.
func MultiSelect(c *ui.Context, text string, chosen *[]string, options []string, o MultiSelectOptions) MultiSelectResult {
	k := pixel.K(c)
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.ButtonBase(c).Role(ui.RolePopUpButton).Label(text).Description(strings.Join(*chosen, ", ")).Disabled(o.Disabled)
	})
	open := ui.Local(base, multiKey{}, func() bool { return false })
	if o.Disabled || pixel.Current(c).Inert {
		// Closing on disable must not return focus to it.
		if *open {
			c.Root().Focus()
		}
		*open = false
	}
	if base.Clicked() || !*open && base.Shortcut(0, ui.KeyDown) {
		*open = !*open
	}
	shell.Size(k.D(pixel.Or("MultiSelect Width", o.Width, 160)), k.D(pixel.LineHeight+4)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		ink, dy := face(g, Normal, st, *open)
		if len(*chosen) == 0 {
			g.Text(6, 2+dy, text, pixel.Faint, 1)
		}
		tags(g, *chosen, 3, 2+dy, g.W-14)
		g.Sprite(g.W-10, (g.H-tinyDown.H)/2+dy, tinyDown, 0, ink, 1)
		finish(g, st)
	})
	res := MultiSelectResult{Element: shell}
	ui.PopoverBase(c, base, open, func(panel *ui.Element) {
		menu(c, panel, func() {
			col := ui.Column(c).Gap(k.D(1)).FocusGroup(ui.Vertical)
			col.Children(func() {
				for i, opt := range options {
					on := slices.Contains(*chosen, opt)
					item, cb, cst := pixel.Shell(c, func() *ui.Element { return ui.CheckboxBase(c, &on).Label(opt) })
					if i == 0 {
						cb.AutoFocus()
					}
					if cb.Changed() {
						res.Changed = true
						if on {
							*chosen = append(*chosen, opt)
						} else {
							*chosen = slices.DeleteFunc(*chosen, func(s string) bool { return s == opt })
						}
					}
					item.Size(k.D(14+pixel.Measure(opt)), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
						var mark *pixel.Sprite
						if on {
							mark = check
						}
						box(k.Canvas(p, r), opt, mark, cst)
					})
				}
			})
		})
	})
	if res.Changed {
		slog.Debug("8bit: multiselect", "label", text, "chosen", *chosen)
	}
	return res
}

// tags draws chips left to right; +N counts the overflow.
func tags(g pixel.Canvas, all []string, x, y, right int) {
	for i, t := range all {
		w := pixel.Measure(t) + 6
		room := right
		if i < len(all)-1 {
			room -= pixel.MicroWidth("+"+strconv.Itoa(len(all)-i-1)) + 2
		}
		if x+w > room {
			g.Micro(x, y+5, "+"+strconv.Itoa(len(all)-i), pixel.Ink, 1)
			return
		}
		g.Frame(x, y, w, pixel.LineHeight, pixel.Notched, pixel.Lit)
		g.Text(x+3, y, t, pixel.Ink, 1)
		x += w + 2
	}
}
