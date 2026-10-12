package media

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PlaylistOptions struct {
	// Width and Height are dots; zeros mean 200 by 120.
	Width, Height int
	Disabled      bool
}

type PlaylistResult struct {
	Element *ui.Element
	// Submitted is Enter or a double click on a row.
	Changed, Submitted bool
}

// Playlist lists tracks; ♪ marks the one playing.
func Playlist(c *ui.Context, text string, tracks []Track, playing int, selected *int, o PlaylistOptions) PlaylistResult {
	k := pixel.K(c)
	if playing < -1 || playing >= len(tracks) || *selected < -1 || *selected >= len(tracks) {
		panic(fmt.Sprintf("media: Playlist plays %d and selects %d of %d", playing, *selected, len(tracks)))
	}
	for _, t := range tracks {
		if t.Length < 0 {
			panic(fmt.Sprintf("media: Playlist track %q lasts %v", t.Title, t.Length))
		}
	}
	name := func(i int) string {
		s := tracks[i].Title + ", " + tracks[i].Artist
		if i == playing {
			s += ", " + k.Msg("playing")
		}
		return s
	}
	var r data.ListResult
	e := ui.Box(c).Width(k.D(pixel.Or("Playlist Width", o.Width, 200))).Children(func() {
		r = data.List(c, text, len(tracks), selected, data.ListOptions{Height: o.Height, Disabled: o.Disabled}, name, func(i int) {
			t := tracks[i]
			ui.Box(c).FillWidth().Height(k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, rc ui.Rect) {
				g := k.Canvas(p, rc)
				ink, length, x := pixel.Ink, mmss(t.Length), pixel.Measure("♪")+2
				if i == playing {
					g.Text(0, 0, "♪", pixel.Accent, 1)
					ink = pixel.Accent
				}
				lw := pixel.MicroWidth(length)
				g.Text(x, 0, pixel.Wrap(t.Title+" · "+t.Artist, g.W-x-6-lw, 1)[0], ink, 1)
				g.Micro(g.W-lw-2, 5, length, pixel.Faint, 1)
			})
		})
	})
	if r.Changed || r.Submitted {
		slog.Debug("8bit: playlist", "label", text, "selected", *selected, "submitted", r.Submitted)
	}
	return PlaylistResult{Element: e, Changed: r.Changed, Submitted: r.Submitted}
}
