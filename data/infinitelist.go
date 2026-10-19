package data

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/feedback"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type InfiniteListOptions struct {
	// Height is in dots; zero means 120.
	Height int
	// Done hides the LOADING row: there is no more.
	Done bool
}

type InfiniteListResult struct {
	Element ui.Element
	// More asks for rows past the n there are.
	More bool
}

type moreKey struct{}

// asking remembers the length last asked past.
type asking struct {
	at   int
	more bool
	list ui.ListState
}

// InfiniteList asks for more when its LOADING row shows.
func InfiniteList(c *ui.Context, text string, n int, o InfiniteListOptions, row func(i int)) InfiniteListResult {
	k := pixel.K(c)
	box := ui.Box(c)
	s := ui.Local(box, moreKey{}, func() asking { return asking{at: -1} })
	res := InfiniteListResult{Element: box, More: s.more}
	s.more = false
	rows := n + 1
	if o.Done {
		rows = n
	}
	var e ui.Element
	box.Children(func() {
		e = ui.List(c, &s.list, rows, pixel.Lazy(c, func(i int) {
			if i < n {
				row(i)
				return
			}
			// Built as it nears the view: ask once per length.
			if s.at != n {
				s.at, s.more = n, true
				c.Invalidate()
			}
			ui.Row(c).Gap(k.D(4)).Padding(k.D(2), k.D(6)).AlignItems(ui.Center).Children(func() {
				feedback.DotsLoader(c, k.Msg("loading"), feedback.DotsLoaderOptions{})
				pixel.Text(c, k.Msg("loading"), pixel.TextOptions{Slot: pixel.Faint})
			})
		})).Padding(0, k.D(3)).Height(k.D(pixel.Or("InfiniteList Height", o.Height, 120))).Label(text)
		veil(c, pixel.Current(c).Inert)
	})
	well(c, box, pixel.State{Focus: e.FocusRing(false).FocusVisible()})
	if res.More {
		slog.Debug("8bit: list wants more", "label", text, "rows", n)
	}
	return res
}
