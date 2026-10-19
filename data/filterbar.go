package data

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type FilterBarOptions struct {
	Disabled bool
}

type FilterBarResult struct {
	Element ui.Element
	Changed bool
}

// FilterBar toggles filter chips; CLEAR turns them all off.
func FilterBar(c *ui.Context, text string, labels []string, on []bool, o FilterBarOptions) FilterBarResult {
	k := pixel.K(c)
	if len(on) != len(labels) {
		panic(fmt.Sprintf("data: FilterBar has %d labels, %d flags", len(labels), len(on)))
	}
	row := ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(text)
	res := FilterBarResult{Element: row}
	row.Children(func() {
		for i, l := range labels {
			if input.Toggle(c, l, &on[i], input.ToggleOptions{Size: input.Small, Disabled: o.Disabled}).Changed {
				res.Changed = true
			}
		}
		if input.Button(c, k.Msg("clear"), input.ButtonOptions{Size: input.Small, Variant: input.Ghost, Disabled: o.Disabled || !slices.Contains(on, true)}).Clicked {
			clear(on)
			res.Changed = true
		}
	})
	if res.Changed {
		slog.Debug("8bit: filters", "label", text, "on", on)
	}
	return res
}
