package data

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

// Hit is one search result.
type Hit struct {
	Title, Text string
}

type SearchResultsOptions struct {
	// Height is the list's in dots; zero means 120.
	Height   int
	Disabled bool
}

type SearchResultsResult = ListResult

// SearchResults counts hits and lists them with query matches marked.
func SearchResults(c *ui.Context, label, query string, hits []Hit, selected *int, o SearchResultsOptions) SearchResultsResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(4))
	var res ListResult
	col.Children(func() {
		pixel.Text(c, fmt.Sprintf(k.Msg("found"), len(hits)), pixel.TextOptions{Slot: pixel.Faint})
		res = List(c, label, len(hits), selected, ListOptions{Height: o.Height, Disabled: o.Disabled}, func(i int) string { return hits[i].Title + ", " + hits[i].Text }, func(i int) {
			ui.Column(c).Children(func() {
				text.Highlight(c, hits[i].Title, query, text.HighlightOptions{})
				text.Highlight(c, hits[i].Text, query, text.HighlightOptions{Slot: pixel.Faint})
			})
		})
	})
	if res.Submitted {
		slog.Debug("8bit: search result opened", "label", label, "hit", *selected)
	}
	res.Element = col
	return res
}
