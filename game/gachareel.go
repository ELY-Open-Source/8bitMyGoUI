package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type GachaReelOptions struct{}

type GachaReelResult struct {
	Element ui.Element
	// Done is true in the frame the last reel stops.
	Done bool
}

type spunKey struct{}

// reelStop staggers reel stops: 0.8 s, then 0.4 s more.
func reelStop(i int) time.Duration {
	return 800*time.Millisecond + time.Duration(i)*400*time.Millisecond
}

// GachaReel spins three reels per spin, stopping on stops.
func GachaReel(c *ui.Context, text string, faces []data.Item, stops [3]int, spin int, o GachaReelOptions) GachaReelResult {
	k := pixel.K(c)
	for _, f := range faces {
		if f.Sprite == nil {
			panic(fmt.Sprintf("game: GachaReel face %q has no sprite", f.Name))
		}
	}
	names := make([]string, 3)
	for i, s := range stops {
		if s < 0 || s >= len(faces) || faces[s].Sprite == nil {
			panic(fmt.Sprintf("game: GachaReel stop %d of %d faces", s, len(faces)))
		}
		names[i] = faces[s].Name
	}
	e := ui.Box(c).Role(ui.RoleImage).Label(text)
	roll := start(k, e, spin)
	// First shown, the reels rest; each new spin spins once.
	landed := ui.Local(e, spunKey{}, func() int { return roll.run })
	shown := stops
	spinning := *landed != roll.run && !k.Reduce()
	elapsed := k.Now().Sub(roll.at)
	if spinning {
		n := k.Tick(c, 12)
		for i := range 3 {
			if elapsed < reelStop(i) {
				shown[i] = (n + 5*i) % len(faces)
			}
		}
	}
	res := GachaReelResult{Element: e, Done: *landed != roll.run && (k.Reduce() || elapsed >= reelStop(2))}
	if res.Done {
		*landed = roll.run
	}
	if *landed == roll.run {
		e.Description(strings.Join(names, ", "))
	}
	e.Size(k.D(3*28-4), k.D(28)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, f := range shown {
			x := i * 28
			g.Frame(x, 0, 24, 28, pixel.Sunken, pixel.Well)
			sp := faces[f].Sprite
			g.Sprite(x+(24-sp.W)/2, (28-sp.H)/2, sp, 0, pixel.Ink, 1)
		}
	})
	return res
}
