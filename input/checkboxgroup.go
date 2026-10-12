package input

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CheckboxGroupOptions struct {
	Disabled bool
}

type CheckboxGroupResult struct {
	Element *ui.Element
	Changed bool
}

// count is how many boxes are checked.
func count(checked []bool) int {
	n := 0
	for _, v := range checked {
		if v {
			n++
		}
	}
	return n
}

// CheckboxGroup puts a parent box over children; partial shows mixed.
func CheckboxGroup(c *ui.Context, text string, items []string, checked []bool, o CheckboxGroupOptions) CheckboxGroupResult {
	if len(items) != len(checked) {
		panic(fmt.Sprintf("input: CheckboxGroup has %d items and %d states", len(items), len(checked)))
	}
	k := pixel.K(c)
	all := count(checked) == len(checked) && len(checked) > 0
	res := CheckboxGroupResult{}
	res.Element = ui.Column(c).Gap(k.D(2)).Role(ui.RoleGroup).Children(func() {
		shell, base, st := pixel.Shell(c, func() *ui.Element {
			return ui.CheckboxBase(c, &all).Label(text).Disabled(o.Disabled)
		})
		shell.Size(k.D(12+pixel.Measure(text)), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
			var mark *pixel.Sprite
			switch n := count(checked); {
			case n == len(checked) && n > 0:
				mark = check
			case n > 0:
				mark = dash
			}
			box(k.Canvas(p, r), text, mark, st)
		})
		if base.Changed() && !base.IsDisabled() {
			for i := range checked {
				checked[i] = all
			}
			res.Changed = true
		}
		ui.Column(c).Gap(k.D(2)).Padding(0, 0, 0, k.D(12)).Children(func() {
			for i, it := range items {
				if Checkbox(c, it, &checked[i], CheckboxOptions{Disabled: o.Disabled}).Changed {
					res.Changed = true
				}
			}
		})
	})
	if res.Changed {
		slog.Debug("8bit: checkbox group", "label", text, "checked", checked)
	}
	return res
}
