package data

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CarouselOptions struct {
	// Every turns the page each that many seconds; zero never.
	Every    int
	Disabled bool
}

type CarouselResult struct {
	Element ui.Element
	Changed bool
}

type carouselKey struct{}

// Carousel turns pages; auto-turning waits on hover, focus, reduce.
func Carousel(c *ui.Context, text string, pages int, page *int, o CarouselOptions, content func(i int)) CarouselResult {
	k := pixel.K(c)
	if pages < 1 || *page < 0 || *page >= pages {
		panic(fmt.Sprintf("data: Carousel page %d of %d", *page, pages))
	}
	every := pixel.Or("Carousel Every", o.Every, 0)
	col := ui.Column(c).AlignItems(ui.Center).Gap(k.D(4)).Role(ui.RoleGroup).Label(text)
	since := ui.Local(col, carouselKey{}, func() time.Time { return k.Now() })
	old := *page
	turn := func(d int) { *page = (*page + d + pages) % pages }
	wait := every == 0 || o.Disabled || pixel.Current(c).Inert || k.Reduce() || col.Hovered() || col.FocusWithin()
	if wait {
		*since = k.Now()
	} else if k.Since(c, *since, 4, 4*every) == 4*every {
		turn(1)
	}
	col.Children(func() {
		ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
			if input.IconButton(c, icon.ChevronLeft, k.Msg("prev"), input.IconButtonOptions{Disabled: o.Disabled}).Clicked {
				turn(-1)
			}
			ui.Box(c).Key(*page).Children(func() { content(*page) })
			if input.IconButton(c, icon.ChevronRight, k.Msg("next"), input.IconButtonOptions{Disabled: o.Disabled}).Clicked {
				turn(1)
			}
		})
		at := *page
		ui.Box(c).Size(k.D(pages*6-2), k.D(4)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range pages {
				if i == at {
					g.Rect(i*6, 0, 4, 4, pixel.Accent)
				} else {
					g.Rect(i*6, 0, 4, 4, pixel.Edge)
					g.Rect(i*6+1, 1, 2, 2, pixel.Well)
				}
			}
		})
	})
	col.Description(fmt.Sprintf("%d/%d", *page+1, pages))
	res := CarouselResult{Element: col, Changed: *page != old}
	if res.Changed {
		*since = k.Now()
		slog.Debug("8bit: carousel", "label", text, "page", *page)
	}
	return res
}
