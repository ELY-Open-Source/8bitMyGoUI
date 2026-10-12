package navigation

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

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
		if pixel.Shortcut(base, 0, m.key) {
			*page = max(1, min(total, m.to()))
		}
	}
	// Taking Enter and Space keeps keys from clicking.
	pixel.Shortcut(base, 0, ui.KeyEnter)
	pixel.Shortcut(base, 0, ui.KeySpace)
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
