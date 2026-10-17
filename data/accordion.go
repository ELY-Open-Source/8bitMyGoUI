package data

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AccordionOptions struct {
	// Many lets several sections stay open at once.
	Many     bool
	Disabled bool
}

type AccordionResult struct {
	Element ui.Element
	Changed bool
}

// Accordion stacks collapsible sections; one stays open unless Many.
func Accordion(c *ui.Context, text string, titles []string, open []bool, o AccordionOptions, panel func(i int)) AccordionResult {
	k := pixel.K(c)
	if len(open) != len(titles) {
		panic(fmt.Sprintf("data: Accordion has %d titles, %d open flags", len(titles), len(open)))
	}
	col := ui.Column(c).AlignItems(ui.Stretch).Gap(k.D(1)).Role(ui.RoleGroup).Label(text)
	res := AccordionResult{Element: col}
	col.Children(func() {
		for i, t := range titles {
			if !Collapsible(c, t, &open[i], CollapsibleOptions{Disabled: o.Disabled}, func() { panel(i) }).Changed {
				continue
			}
			res.Changed = true
			for j := range open {
				if j != i && open[i] && !o.Many {
					open[j] = false
				}
			}
		}
	})
	if res.Changed {
		slog.Debug("8bit: accordion", "label", text, "open", open)
	}
	return res
}
