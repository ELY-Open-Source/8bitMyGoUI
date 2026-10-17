package data

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type Column struct {
	Title string
	// Width is in dots; zero means 60.
	Width int
}

// Sort is the sorting column and its direction.
type Sort struct {
	Column int
	Desc   bool
}

type TableOptions struct {
	// Height is the rows' in dots; zero means 120.
	Height   int
	Disabled bool
}

type TableResult struct {
	Element ui.Element
	// Sorted reports a header click; sort the rows by *sort.
	Changed, Submitted, Sorted bool
}

type tableKey struct{}

// Table sorts by its headers over virtual, dithered zebra rows.
func Table(c *ui.Context, text string, cols []Column, rows int, sort *Sort, selected *int, o TableOptions, cell func(row, col int) string) TableResult {
	k := pixel.K(c)
	off := o.Disabled || pixel.Current(c).Inert
	ws, total := make([]int, len(cols)), 0
	for i, col := range cols {
		ws[i] = pixel.Or("Column Width", col.Width, 60)
		total += ws[i]
	}
	box := ui.Box(c).Width(k.D(total + 4))
	st := ui.Local(box, tableKey{}, func() ui.ListState { return ui.ListState{} })
	st.Selected, st.Label = selected, func(r int) string { return cell(r, 0) }
	if off {
		st.Selected = nil
	}
	res := TableResult{Element: box}
	var list ui.Element
	box.Children(func() {
		ui.Row(c).FocusGroup(ui.Horizontal).Role(ui.RoleRow).Children(func() {
			for i, col := range cols {
				shell, base, hs := pixel.Shell(c, func() ui.Element {
					return ui.ButtonBase(c).Role(ui.RoleColumnHeader).Label(col.Title).Disabled(o.Disabled)
				})
				if base.Clicked() {
					*sort = Sort{Column: i, Desc: sort.Column == i && !sort.Desc}
					res.Sorted = true
				}
				on, desc := sort.Column == i, sort.Desc
				if dir := "up"; on {
					if desc {
						dir = "down"
					}
					base.Description(k.Msg(dir))
				}
				shell.Size(k.D(ws[i]), k.D(pixel.LineHeight+4)).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					g.Frame(0, 0, g.W, g.H, pixel.Raised, pixel.Face)
					g.Clip(0, 0, g.W-10, g.H, func() { g.Text(4, 2, col.Title, pixel.Ink, 1) })
					if on {
						arrow := up
						if desc {
							arrow = down
						}
						g.Sprite(g.W-9, 7, arrow, 0, pixel.Accent, 1)
					}
					g.Finish(hs)
				})
			}
		})
		list = ui.List(c, st, rows, pixel.Lazy(c, func(r int) {
			texts := make([]string, len(cols))
			for j := range cols {
				texts[j] = cell(r, j)
			}
			chosen := r == *selected
			ui.Row(c).Height(k.D(pixel.LineHeight)).Label(strings.Join(texts, ", ")).Draw(func(p *ui.Painter, rc ui.Rect) {
				g := k.Canvas(p, rc)
				switch {
				case chosen:
					g.Rect(0, 0, g.W, g.H, pixel.Accent)
				case r%2 == 1:
					g.Veil(0, 0, g.W, g.H, pixel.Shade, 4)
				}
				ink, x := pixel.Ink, 0
				if chosen {
					ink = pixel.AccentInk
				}
				for j, s := range texts {
					g.Clip(x, 0, ws[j], g.H, func() { g.Text(x+4, 0, s, ink, 1) })
					x += ws[j]
				}
			})
		})).Width(k.D(total)).Height(k.D(pixel.Or("Table Height", o.Height, 120))).Label(text)
		veil(c, off)
	})
	pixel.Park(c, list, off)
	if off {
		list.Description(k.Msg("disabled"))
	}
	well(c, box, pixel.State{Focus: list.FocusRing(false).FocusVisible(), Disabled: o.Disabled})
	res.Changed, res.Submitted = list.Changed(), list.Submitted()
	if res.Changed {
		slog.Debug("8bit: table chose", "label", text, "row", *selected)
	}
	if res.Sorted {
		slog.Debug("8bit: table sorted", "label", text, "column", sort.Column, "desc", sort.Desc)
	}
	return res
}
