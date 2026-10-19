package code

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type SpriteSheetOptions struct {
	// Scale is dots per sprite dot, 1–4; zero means 2.
	Scale int
	// Columns of frames; zero means 4.
	Columns int
	// FPS is the playback speed, 4–15; zero means 6.
	FPS      int
	Disabled bool
}

type SpriteSheetResult struct {
	Element ui.Element
	// Changed is a frame chosen or play pressed.
	Changed bool
}

type sheetKey struct{}

// sheet is the frame play began on, and when.
type sheet struct {
	from  int
	since time.Time
}

// SpriteSheet lays frames in a grid; playing steps through them.
func SpriteSheet(c *ui.Context, text string, sp *pixel.Sprite, frame *int, playing *bool, o SpriteSheetOptions) SpriteSheetResult {
	k := pixel.K(c)
	n, s := sp.Frames(), pixel.Or("SpriteSheet Scale", o.Scale, 2)
	fps, cols := pixel.Or("SpriteSheet FPS", o.FPS, 6), min(n, pixel.Or("SpriteSheet Columns", o.Columns, 4))
	if s > 4 || *frame < 0 || *frame >= n {
		panic(fmt.Sprintf("code: SpriteSheet frame %d of %d at scale %d", *frame, n, s))
	}
	if fps < 4 || fps > 15 {
		panic(fmt.Sprintf("code: SpriteSheet FPS %d outside 4–15", fps))
	}
	var res SpriteSheetResult
	res.Element = ui.Row(c).Gap(k.D(8)).AlignItems(ui.Start)
	sh := ui.Local(res.Element, sheetKey{}, func() sheet { return sheet{} })
	// Playback is the content asked for, so it keeps time.
	switch {
	case !*playing:
		sh.since = time.Time{}
	case sh.since.IsZero():
		sh.from, sh.since = *frame, k.Now()
	}
	if *playing {
		per := time.Second / time.Duration(fps)
		m := int(k.Elapsed(&sh.since) / per)
		*frame = (sh.from + m) % n
		k.At(c, sh.since.Add(time.Duration(m+1)*per))
	}
	cw, ch := sp.W*s+4, sp.H*s+4
	res.Element.Children(func() {
		old := *frame
		shell, base, st, _ := pixel.Pad(c, text, cols, n, cw, ch, frame, o.Disabled)
		base.Description(fmt.Sprintf("%d / %d", *frame+1, n))
		if *frame != old {
			res.Changed, sh.from, sh.since = true, *frame, k.Now()
		}
		at := *frame
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range n {
				x, y := i%cols*cw, i/cols*ch
				g.Frame(x, y, cw, ch, pixel.Flat, pixel.Well)
				if i == at {
					g.Cursor(x, y, cw, ch)
				}
				g.Sprite(x+2, y+2, sp, i, pixel.Ink, s)
			}
			g.Finish(st)
		})
		ui.Column(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
			// The preview doubles the grid's scale, up to 4.
			big := min(4, 2*s)
			ui.Box(c).Size(k.D(sp.W*big+8), k.D(sp.H*big+8)).Role(ui.RoleImage).Label(fmt.Sprintf("%s %d", text, at+1)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
				g.Sprite(4, 4, sp, at, pixel.Ink, big)
			})
			sym := icon.Play
			if *playing {
				sym = icon.Pause
			}
			if input.Toggle(c, "", playing, input.ToggleOptions{Icon: sym, Label: k.Msg("play"), Disabled: o.Disabled}).Changed {
				res.Changed = true
			}
		})
	})
	if res.Changed {
		slog.Debug("8bit: sprite sheet", "label", text, "frame", *frame, "playing", *playing)
	}
	return res
}
