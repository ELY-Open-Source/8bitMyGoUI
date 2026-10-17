// Package media plays, shows and records sound and pictures.
package media

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Track is a song: title, artist and length.
type Track struct {
	Title, Artist string
	Length        time.Duration
}

// mmss prints d as m:ss, or h:mm:ss past an hour.
func mmss(d time.Duration) string {
	s := int(d / time.Second)
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s/60%60, s%60)
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

// prev is Next turned around.
var prev = icon.Next.Mirror()

// scale is s, or def for zero; above 4 panics.
func scale(name string, s, def int) int {
	s = pixel.Or(name+" Scale", s, def)
	if s > 4 {
		panic(fmt.Sprintf("media: %s Scale %d outside 1–4", name, s))
	}
	return s
}

// art is a named picture w by h dots.
func art(c *ui.Context, k pixel.Kit, label string, w, h int) ui.Element {
	if label == "" {
		panic("media: a picture needs a name")
	}
	return ui.Box(c).Size(k.D(w), k.D(h)).Label(label).Role(ui.RoleImage)
}

// dark is the darkest slot and a line on it.
func dark(k pixel.Kit) (body, line pixel.Slot) {
	if k.Palette().Night {
		return pixel.Well, pixel.Face
	}
	return pixel.Ink, pixel.Faint
}

// sameSize panics unless every sprite matches the first.
func sameSize(name string, sps []*pixel.Sprite) {
	if len(sps) == 0 {
		panic("media: " + name + " needs pictures")
	}
	for _, sp := range sps {
		if sp.W != sps[0].W || sp.H != sps[0].H {
			panic(fmt.Sprintf("media: %s mixes %d×%d and %d×%d pictures", name, sps[0].W, sps[0].H, sp.W, sp.H))
		}
	}
}

type turnKey struct{}

// turn is steps taken, their first frame, and play.
type turn struct {
	held, from int
	playing    bool
}

// turned counts steps while playing; stopped, it holds still.
func turned(c *ui.Context, k pixel.Kit, e ui.Element, playing bool, fps int) int {
	s := ui.Local(e, turnKey{}, func() turn { return turn{} })
	on := playing && !k.Reduce()
	if on {
		n := k.Tick(c, fps)
		// A clock stepping back rebases, so the count never falls.
		if !s.playing || n-s.from < s.held {
			s.from = n - s.held
		}
		s.held = n - s.from
	}
	s.playing = on
	return s.held
}
