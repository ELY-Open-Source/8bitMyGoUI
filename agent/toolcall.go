package agent

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Outcome is how a cast went.
type Outcome int

const (
	Casting Outcome = iota
	OK
	Failed
)

type ToolCallOptions struct {
	// Args show under the spell, cut to a line.
	Args string
	// Width is in dots; zero means 200.
	Width int
}

var glow = []pixel.Slot{pixel.Magic, pixel.Info, pixel.Accent}

// ToolCall is a spell card: CAST, then OK or FAILED.
func ToolCall(c *ui.Context, tool string, out Outcome, o ToolCallOptions) ui.Element {
	k := pixel.K(c)
	if tool == "" || out < Casting || out > Failed {
		panic(fmt.Sprintf("agent: ToolCall %q with outcome %d", tool, out))
	}
	w := pixel.Or("ToolCall Width", o.Width, 200)
	spell := fmt.Sprintf(k.Msg("cast_tool"), tool)
	word, ink, mark := k.Msg("casting"), pixel.Info, (*pixel.Sprite)(nil)
	switch out {
	case OK:
		word, ink, mark = k.Msg("ok"), pixel.Good, icon.Check
	case Failed:
		word, ink, mark = k.Msg("failed"), pixel.Bad, icon.Close
	}
	// The outcome takes the right; the spell keeps some room.
	side := 13
	if mark != nil {
		side = 16 + pixel.Measure(word)
	}
	if w < side+60 {
		panic(fmt.Sprintf("agent: ToolCall Width %d is under %d", w, side+60))
	}
	h, said := pixel.LineHeight+6, word
	if o.Args != "" {
		h, said = h+pixel.LineHeight, o.Args+", "+word
	}
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(spell).Description(said).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, w, h, pixel.Raised, pixel.Face)
		g.Rect(2, 2, 2, h-4, ink)
		x, rune := w-4-side, pixel.Magic
		if mark == nil {
			// The rune and dots glow in turn while it casts.
			step := g.Busy(8)
			rune = pixel.Cycle(glow, step)[0]
			for i := range 3 {
				dot := pixel.Faint
				if i == step%3 {
					dot = pixel.Info
				}
				g.Rect(x+i*5, 10, 3, 3, dot)
			}
		} else {
			g.Sprite(x, 5, mark, 0, ink, 1)
			g.Text(x+16, 3, word, ink, 1)
		}
		g.Sprite(7, 5, icon.Sparkle, 0, rune, 1)
		g.Text(22, 3, pixel.Wrap(spell, x-26, 1)[0], pixel.Ink, 1)
		if o.Args != "" {
			g.Text(22, 3+pixel.LineHeight, pixel.Wrap(o.Args, w-26, 1)[0], pixel.Faint, 1)
		}
	})
}
