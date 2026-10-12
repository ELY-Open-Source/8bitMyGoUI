package navigation

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PaginationOptions struct {
	Disabled bool
}

type PaginationResult struct {
	Element *ui.Element
	Changed bool
}

// Pagination steps through pages; the current one is sunken.
func Pagination(c *ui.Context, text string, page *int, total int, o PaginationOptions) PaginationResult {
	k := pixel.K(c)
	check(*page, total)
	res := PaginationResult{}
	set := func(p int) {
		if p != *page {
			*page, res.Changed = p, true
		}
	}
	small := input.IconButtonOptions{Size: input.Small, Disabled: o.Disabled || *page == 1}
	res.Element = ui.Row(c).Gap(k.D(2)).AlignItems(ui.Center).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(text).Children(func() {
		// Keys keep each button's focus as the pages shift.
		ui.Box(c).Key("prev").Children(func() {
			if input.IconButton(c, icon.ChevronLeft, k.Msg("prev"), small).Clicked {
				set(*page - 1)
			}
		})
		for i, p := range pages(*page, total) {
			if p == 0 {
				ui.Box(c).Key(-i).Children(func() { pixel.Text(c, "…", pixel.TextOptions{Slot: pixel.Faint}) })
				continue
			}
			ui.Box(c).Key(p).Children(func() {
				on := p == *page
				if input.Toggle(c, strconv.Itoa(p), &on, input.ToggleOptions{Size: input.Small, Disabled: o.Disabled}).Changed {
					set(p)
				}
			})
		}
		ui.Box(c).Key("next").Children(func() {
			small.Disabled = o.Disabled || *page == total
			if input.IconButton(c, icon.ChevronRight, k.Msg("next"), small).Clicked {
				set(*page + 1)
			}
		})
	})
	if res.Changed {
		slog.Debug("8bit: page", "label", text, "page", *page)
	}
	return res
}

// pages lists 1, the pages round p, and n.
func pages(p, n int) []int {
	var out []int
	for i := 1; i <= n; i++ {
		switch {
		case i == 1 || i == n || i >= p-1 && i <= p+1:
			out = append(out, i)
		case len(out) > 0 && out[len(out)-1] != 0:
			out = append(out, 0)
		}
	}
	return out
}

// check panics on a page outside 1 to total.
func check(page, total int) {
	if total < 1 || page < 1 || page > total {
		panic(fmt.Sprintf("navigation: page %d of %d", page, total))
	}
}

type PageDotsOptions struct {
	Disabled bool
}

type PageDotsResult struct {
	Element *ui.Element
	Changed bool
}

// PageDots shows a dot per page; clicks jump, arrows move.
func PageDots(c *ui.Context, text string, page *int, total int, o PageDotsOptions) PageDotsResult {
	k := pixel.K(c)
	check(*page, total)
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleSlider).Label(text).Description(fmt.Sprintf("%d/%d", *page, total)).Disabled(o.Disabled)
	})
	old := *page
	for _, m := range []struct {
		key ui.Key
		to  func() int
	}{
		{ui.KeyLeft, func() int { return *page - 1 }}, {ui.KeyRight, func() int { return *page + 1 }},
		{ui.KeyHome, func() int { return 1 }}, {ui.KeyEnd, func() int { return total }},
	} {
		if base.Shortcut(0, m.key) {
			*page = max(1, min(total, m.to()))
		}
	}
	// Taking Enter and Space keeps keys from clicking.
	base.Shortcut(0, ui.KeyEnter)
	base.Shortcut(0, ui.KeySpace)
	if base.Clicked() {
		x, _, _ := base.PointerPosition()
		*page = max(1, min(total, int(x/k.Dot())/8+1))
	}
	shell.Size(k.D(total*8), k.D(8)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range total {
			if i+1 == *page {
				g.Rect(i*8+1, 1, 6, 6, pixel.Accent)
			} else {
				g.Rect(i*8+2, 2, 4, 4, pixel.Shade)
			}
		}
		g.Finish(st)
	})
	res := PageDotsResult{Element: shell, Changed: *page != old}
	if res.Changed {
		slog.Debug("8bit: page dots", "label", text, "page", *page)
	}
	return res
}
