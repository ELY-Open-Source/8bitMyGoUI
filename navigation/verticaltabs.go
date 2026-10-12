package navigation

import (
	"github.com/egoist/mygo/ui"
)

type VerticalTabsOptions struct {
	Disabled bool
}

type VerticalTabsResult struct {
	Element *ui.Element
	Changed bool
}

// VerticalTabs are tabs down the panel's left edge.
func VerticalTabs(c *ui.Context, text string, labels []string, selected *int, o VerticalTabsOptions, panel func()) VerticalTabsResult {
	return VerticalTabsResult(tabs(c, text, labels, selected, true, TabsOptions(o), panel))
}
