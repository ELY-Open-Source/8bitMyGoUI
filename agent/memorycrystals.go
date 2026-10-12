package agent

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type MemoryCrystalsOptions struct {
	// Slots is how many crystals; zero means 8.
	Slots    int
	Disabled bool
}

type MemoryCrystalsResult struct {
	Element *ui.Element
	// Picked is the crystal chosen this frame, else -1.
	Picked int
}

var (
	twinkle = pixel.MustSprite(`
		..#..
		..#..
		##.##
		..#..
		..#..`)
)

type crystalKey struct{}

const crystalW, crystalH = 26, 30

// MemoryCrystals keeps memories in crystals; empty ones stay dark.
func MemoryCrystals(c *ui.Context, text string, memories []string, o MemoryCrystalsOptions) MemoryCrystalsResult {
	k := pixel.K(c)
	n := pixel.Or("MemoryCrystals Slots", o.Slots, 8)
	if len(memories) > n || slices.Contains(memories, "") {
		panic(fmt.Sprintf("agent: MemoryCrystals has %d memories for %d slots, or a blank", len(memories), n))
	}
	cols := min(n, 8)
	col := ui.Column(c).Gap(k.D(4))
	cur := ui.Local(col, crystalKey{}, func() int { return 0 })
	*cur = min(*cur, n-1)
	res := MemoryCrystalsResult{Picked: -1}
	res.Element = col.Children(func() {
		shell, base, st, picked := pixel.Pad(c, text, cols, n, crystalW, crystalH, cur, o.Disabled)
		res.Picked = picked
		at, held := *cur, len(memories)
		line, ink := fmt.Sprintf("%d/%d: %s", at+1, n, k.Msg("empty_crystal")), pixel.Faint
		if at < held {
			line, ink = fmt.Sprintf("%d/%d: %s", at+1, n, memories[at]), pixel.Ink
		}
		base.Description(line)
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			if st.Focus {
				g.Cursor(at%cols*crystalW, at/cols*crystalH, crystalW, crystalH)
			}
			// Held crystals twinkle in turn; Reduce Motion stills them.
			step := g.Step(4)
			for i := range n {
				x, y := i%cols*crystalW, i/cols*crystalH
				if i >= held {
					g.Sprite(x+1, y+3, icon.Crystal, 0, pixel.Well, 2)
					continue
				}
				g.Sprite(x+1, y+3, icon.Crystal, 0, pixel.Magic, 2)
				if (step+i*3)%8 == 0 {
					g.Sprite(x+15, y+5, twinkle, 0, pixel.Lit, 1)
				}
			}
			g.Finish(st)
		})
		ui.Box(c).Size(k.D(cols*crystalW), k.D(pixel.LineHeight+6)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
			g.Text(4, 3, pixel.Wrap(line, g.W-8, 1)[0], ink, 1)
		})
	})
	if res.Picked >= 0 {
		slog.Debug("8bit: memory crystal picked", "label", text, "slot", res.Picked)
	}
	return res
}
