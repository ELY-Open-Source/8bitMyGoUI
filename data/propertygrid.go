package data

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Property is one row; set exactly one value.
type Property struct {
	Name   string
	Text   *string
	Number *float64
	Switch *bool
	// Min and Max bound Number as NumberInput's do.
	Min, Max float64
}

type PropertyGridOptions struct {
	// NameWidth is the names' column in dots; zero means 80.
	NameWidth int
	Disabled  bool
}

type PropertyGridResult struct {
	Element ui.Element
	Changed bool
}

// PropertyGrid edits named text, number and switch values.
func PropertyGrid(c *ui.Context, text string, props []Property, o PropertyGridOptions) PropertyGridResult {
	k := pixel.K(c)
	nw := pixel.Or("PropertyGrid NameWidth", o.NameWidth, 80)
	col := ui.Column(c).Gap(k.D(2)).Role(ui.RoleGroup).Label(text)
	res := PropertyGridResult{Element: col}
	col.Children(func() {
		for _, p := range props {
			ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Width(k.D(nw)).Shrink(0).Children(func() {
					pixel.Text(c, p.Name, pixel.TextOptions{Slot: pixel.Faint})
				})
				switch {
				case p.Text != nil && p.Number == nil && p.Switch == nil:
					res.Changed = input.TextInput(c, p.Name, p.Text, input.TextInputOptions{Disabled: o.Disabled}).Changed || res.Changed
				case p.Number != nil && p.Text == nil && p.Switch == nil:
					res.Changed = input.NumberInput(c, p.Name, p.Number, input.NumberInputOptions{Min: p.Min, Max: p.Max, Disabled: o.Disabled}).Changed || res.Changed
				case p.Switch != nil && p.Text == nil && p.Number == nil:
					res.Changed = input.Switch(c, "", p.Switch, input.SwitchOptions{Label: p.Name, Disabled: o.Disabled}).Changed || res.Changed
				default:
					panic(fmt.Sprintf("data: Property %q needs exactly one value", p.Name))
				}
			})
		}
	})
	if res.Changed {
		slog.Debug("8bit: property grid", "label", text)
	}
	return res
}
