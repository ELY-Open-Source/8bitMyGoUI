package data

import (
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type TransferOptions struct {
	// Height is each list's in dots; zero means 120.
	Height   int
	Disabled bool
}

type TransferResult struct {
	Element ui.Element
	Changed bool
}

type transferKey struct{}

// Transfer moves chosen items between lists with → ←.
func Transfer(c *ui.Context, text string, left, right *[]string, o TransferOptions) TransferResult {
	k := pixel.K(c)
	row := ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(text)
	sel := ui.Local(row, transferKey{}, func() [2]int { return [2]int{} })
	res := TransferResult{Element: row}
	// side lists one side's items, its choice kept in range.
	side := func(name string, items []string, at *int) {
		*at = max(0, min(*at, len(items)-1))
		List(c, name, len(items), at, ListOptions{Height: o.Height, Disabled: o.Disabled}, func(i int) string { return items[i] }, func(i int) {
			pixel.Text(c, items[i], pixel.TextOptions{})
		}).Element.Width(k.D(100))
	}
	shift := func(from, to *[]string, at int) {
		if at < 0 || at >= len(*from) {
			return
		}
		*to = append(*to, (*from)[at])
		*from = slices.Delete(*from, at, at+1)
		res.Changed = true
	}
	row.Children(func() {
		side(text+" "+k.Msg("left"), *left, &sel[0])
		ui.Column(c).Gap(k.D(4)).Children(func() {
			if input.IconButton(c, icon.ArrowRight, k.Msg("right"), input.IconButtonOptions{Disabled: o.Disabled || len(*left) == 0}).Clicked {
				shift(left, right, sel[0])
			}
			if input.IconButton(c, icon.ArrowLeft, k.Msg("left"), input.IconButtonOptions{Disabled: o.Disabled || len(*right) == 0}).Clicked {
				shift(right, left, sel[1])
			}
		})
		side(text+" "+k.Msg("right"), *right, &sel[1])
	})
	if res.Changed {
		slog.Debug("8bit: transfer", "label", text, "left", len(*left), "right", len(*right))
	}
	return res
}
