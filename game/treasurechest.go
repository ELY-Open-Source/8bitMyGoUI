package game

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TreasureChestOptions struct {
	Disabled bool
}

type TreasureChestResult struct {
	Element ui.Element
	// Opened is true in the frame a press opens it.
	Opened bool
}

type lidKey struct{}

// lid is whether the chest was open, and since when.
type lid struct {
	open bool
	at   time.Time
}

// TreasureChest opens in frames when pressed; the item rises out.
func TreasureChest(c *ui.Context, text string, it data.Item, open *bool, o TreasureChestOptions) TreasureChestResult {
	k := pixel.K(c)
	if it.Sprite == nil {
		panic(fmt.Sprintf("game: TreasureChest item %q has no sprite", it.Name))
	}
	shell, base, st := pixel.Shell(c, func() ui.Element { return ui.ButtonBase(c).Label(text).Disabled(o.Disabled) })
	res := TreasureChestResult{Element: shell}
	if base.Clicked() && !*open {
		*open, res.Opened = true, true
		slog.Debug("8bit: chest opened", "label", text, "item", it.Name)
	}
	// First shown, the chest sits as open says.
	s := ui.Local(base, lidKey{}, func() lid { return lid{open: *open} })
	if *open != s.open {
		s.open, s.at = *open, k.Now()
	}
	frame, rise := 0, 0
	if *open {
		frame, rise = 2, 10
		if !k.Reduce() {
			frame, rise = k.Since(c, s.at, 8, 2), k.Since(c, s.at, 15, 10)
		}
	}
	if *open {
		base.Description(it.String())
	}
	// Room above the chest holds the fully risen item.
	sp := it.Sprite
	top := sp.H + 13
	shell.Size(k.D(24), k.D(top+14)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		if frame == 2 {
			g.Sprite((24-sp.W)/2, top-3-sp.H-rise, sp, 0, pixel.Ink, 1)
		}
		chest(g, 2, top, frame)
		g.Finish(st)
	})
	return res
}

// chest draws a 20×14 chest: closed, ajar or open.
func chest(g pixel.Canvas, x, y, frame int) {
	g.Frame(x, y+5, 20, 9, pixel.Flat, pixel.Shade)
	switch frame {
	case 0:
		g.Frame(x, y, 20, 6, pixel.Flat, pixel.Shade)
	case 1:
		g.Frame(x, y-3, 20, 6, pixel.Flat, pixel.Shade)
		g.Rect(x+1, y+3, 18, 2, pixel.Ink)
		g.Veil(x+1, y+3, 18, 2, pixel.Warn, 6)
	default:
		g.Frame(x, y-4, 20, 4, pixel.Flat, pixel.Edge)
		g.Rect(x+1, y+5, 18, 3, pixel.Ink)
		g.Veil(x+1, y+5, 18, 3, pixel.Warn, 8)
	}
	g.Rect(x+9, y+5, 2, 9, pixel.Warn)
	if frame == 0 {
		g.Rect(x+9, y, 2, 6, pixel.Warn)
		g.Rect(x+8, y+4, 4, 4, pixel.Warn)
		g.Dot(x+10, y+6, pixel.Edge)
	}
}
