package media

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type VideoFrameOptions struct {
	// Width is the screen in dots, 16:9; zero means 160.
	Width    int
	Disabled bool
}

type VideoFrameResult struct {
	Element ui.Element
	// Changed is play, pause or a seek.
	Changed bool
}

// VideoFrame puts screen in 16:9 over play and a scrubber.
func VideoFrame(c *ui.Context, text string, playing *bool, at *time.Duration, total time.Duration, o VideoFrameOptions, screen func()) VideoFrameResult {
	k := pixel.K(c)
	w := pixel.Or("VideoFrame Width", o.Width, 160)
	// Play and the scrubber's two times need 128 dots.
	if w < 128 {
		panic(fmt.Sprintf("media: VideoFrame Width %d is under 128", w))
	}
	h := w * 9 / 16
	res := VideoFrameResult{}
	res.Element = ui.Column(c).Gap(k.D(3)).Role(ui.RoleGroup).Label(text).Children(func() {
		ui.Box(c).Size(k.D(w+4), k.D(h+4)).Padding(k.D(2)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			body, _ := dark(k)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, body)
		}).Children(func() {
			ui.Box(c).Size(k.D(w), k.D(h)).Clip().Children(screen)
		})
		ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Children(func() {
			sp := icon.Play
			if *playing {
				sp = icon.Pause
			}
			play := input.Toggle(c, "", playing, input.ToggleOptions{Icon: sp, Label: k.Msg("play"), Disabled: o.Disabled}).Changed
			seek := Scrubber(c, k.Msg("seek"), at, total, ScrubberOptions{Width: w - 64, Disabled: o.Disabled}).Changed
			res.Changed = play || seek
		})
	})
	if res.Changed {
		slog.Debug("8bit: video frame", "label", text, "playing", *playing, "at", *at)
	}
	return res
}
