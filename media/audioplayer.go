package media

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Player is what an AudioPlayer holds: play, place and volume.
type Player struct {
	Playing bool
	At      time.Duration
	// Level is 0–8 volume bars.
	Level int
	Muted bool
}

type AudioPlayerOptions struct {
	Disabled bool
}

type AudioPlayerResult struct {
	Element *ui.Element
	// Changed is play, seek or volume; others are presses.
	Changed, Stopped, Previous, Next bool
}

// AudioPlayer shows a track over transport, scrubber and volume.
func AudioPlayer(c *ui.Context, text string, t Track, pl *Player, o AudioPlayerOptions) AudioPlayerResult {
	k := pixel.K(c)
	var res AudioPlayerResult
	res.Element = layout.Panel(c, layout.PanelOptions{Style: pixel.Raised}, func() {
		// An LCD shows what plays, ♪ while it does.
		ui.Box(c).FillWidth().Height(k.D(2*pixel.LineHeight + 6)).Role(ui.RoleText).Label(t.Title + ", " + t.Artist).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
			x := 4
			if pl.Playing {
				g.Text(x, 3, "♪", pixel.Accent, 1)
				x += 10
			}
			g.Text(x, 3, pixel.Wrap(t.Title, g.W-x-4, 1)[0], pixel.Ink, 1)
			g.Text(4, 3+pixel.LineHeight, pixel.Wrap(t.Artist, g.W-8, 1)[0], pixel.Faint, 1)
		})
		seek := Scrubber(c, k.Msg("seek"), &pl.At, t.Length, ScrubberOptions{Disabled: o.Disabled}).Changed
		ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
			tr := Transport(c, text, &pl.Playing, TransportOptions{Disabled: o.Disabled})
			// Stop goes back to the start.
			if tr.Stopped {
				pl.At = 0
			}
			vol := Volume(c, k.Msg("volume"), &pl.Level, &pl.Muted, VolumeOptions{Disabled: o.Disabled}).Changed
			res.Changed = seek || tr.Changed || vol || tr.Stopped
			res.Stopped, res.Previous, res.Next = tr.Stopped, tr.Previous, tr.Next
		})
	}).Gap(k.D(4)).Role(ui.RoleGroup).Label(text)
	if res.Changed || res.Previous || res.Next {
		slog.Debug("8bit: audio player", "label", text, "track", t.Title, "playing", pl.Playing, "at", pl.At)
	}
	return res
}
