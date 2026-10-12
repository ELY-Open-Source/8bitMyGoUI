package data

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PagerOptions struct {
	// Sizes offered; nil means 10, 20 and 50.
	Sizes    []int
	Disabled bool
}

type PagerResult struct {
	Element *ui.Element
	Changed bool
}

// Pager pages through total items: "1–20 OF 312".
func Pager(c *ui.Context, text string, total int, page, size *int, o PagerOptions) PagerResult {
	k := pixel.K(c)
	sizes := o.Sizes
	if sizes == nil {
		sizes = []int{10, 20, 50}
	}
	at := slices.Index(sizes, *size)
	if total < 0 || *page < 0 || at < 0 || slices.Min(sizes) < 1 {
		panic(fmt.Sprintf("data: Pager page %d, size %d of %v, total %d", *page, *size, sizes, total))
	}
	old, oldSize := *page, *size
	last := func() int { return max(0, (total-1) / *size) }
	*page = min(*page, last())
	row := ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(text)
	row.Children(func() {
		pixel.Text(c, k.Msg("per_page"), pixel.TextOptions{Slot: pixel.Faint})
		labels := make([]string, len(sizes))
		for i, s := range sizes {
			labels[i] = strconv.Itoa(s)
		}
		if input.Segmented(c, k.Msg("per_page"), labels, &at, input.SegmentedOptions{Disabled: o.Disabled}).Changed {
			*size, *page = sizes[at], 0
		}
		from, to := min(*page**size+1, total), min((*page+1)**size, total)
		pixel.Text(c, fmt.Sprintf(k.Msg("range"), from, to, total), pixel.TextOptions{})
		if input.IconButton(c, icon.ChevronLeft, k.Msg("prev"), input.IconButtonOptions{Disabled: o.Disabled || *page == 0}).Clicked {
			*page--
		}
		if input.IconButton(c, icon.ChevronRight, k.Msg("next"), input.IconButtonOptions{Disabled: o.Disabled || *page >= last()}).Clicked {
			*page++
		}
	})
	res := PagerResult{Element: row, Changed: *page != old || *size != oldSize}
	if res.Changed {
		slog.Debug("8bit: pager", "label", text, "page", *page, "size", *size)
	}
	return res
}
