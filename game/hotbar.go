package game

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type HotbarOptions struct {
	Disabled bool
}

type HotbarResult struct {
	Element *ui.Element
	Changed bool
	// Used is the slot used this frame, or -1.
	Used int
}

type wheelKey struct{}

// wheel holds unspent scroll and a turn not yet reported.
type wheel struct {
	left  float32
	moved bool
}

// Hotbar rows nine slots; 1–9, arrows or the wheel choose.
func Hotbar(c *ui.Context, text string, items []data.Item, selected *int, o HotbarOptions) HotbarResult {
	k := pixel.K(c)
	n := len(items)
	if n == 0 || n > 9 || *selected < 0 || *selected >= n {
		panic(fmt.Sprintf("game: Hotbar slot %d of %d; it holds 1–9", *selected, n))
	}
	old := *selected
	shell, base, st, used := pixel.Pad(c, text, n, n, data.SlotSize, data.SlotSize, selected, o.Disabled)
	for i := range n {
		if base.Shortcut(0, ui.Key1+ui.Key(i)) {
			*selected = i
		}
	}
	wh := ui.Local(base, wheelKey{}, func() wheel { return wheel{} })
	notch := k.D(pixel.LineHeight)
	base.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll {
			return false
		}
		// A notch moves one slot; the bar wraps around.
		for wh.left += ev.DY; wh.left >= notch || wh.left <= -notch; {
			d := 1
			if wh.left < 0 {
				d = -1
			}
			wh.left -= float32(d) * notch
			*selected, wh.moved = ((*selected+d)%n+n)%n, true
		}
		return true
	})
	res := HotbarResult{Element: shell, Changed: *selected != old || wh.moved, Used: used}
	wh.moved = false
	base.Description(items[*selected].String())
	if res.Changed {
		slog.Debug("8bit: hotbar", "label", text, "slot", *selected)
	}
	sel, shown := *selected, append([]data.Item(nil), items...)
	shell.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, it := range shown {
			x := i * data.SlotSize
			data.DrawSlot(g, x, 0, it, data.ItemSlotOptions{Selected: i == sel})
			g.Micro(x+3, 3, strconv.Itoa(i+1), pixel.Faint, 1)
		}
		g.Finish(st)
	})
	return res
}
