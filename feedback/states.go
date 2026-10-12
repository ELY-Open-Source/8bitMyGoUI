package feedback

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

var openChest = pixel.MustSprite(`
	..eeeeeeeeeeee..
	.esssssssssssse.
	.elllllllllllle.
	..eeeeeeeeeeee..
	................
	eeeeeeeeeeeeeeee
	ewwwwwwwwwwwwwwe
	ewwwwwwwwwwwwwwe
	eeeeeeeyyeeeeeee
	essssssyysssssse
	esssssssssssssse
	eeeeeeeeeeeeeeee
`)

type EmptyStateOptions struct {
	// Action labels a button; Acted reports it.
	Action   string
	Disabled bool
}

type EmptyStateResult struct {
	Element *ui.Element
	Acted   bool
}

// EmptyState shows NO ITEMS over an empty chest.
func EmptyState(c *ui.Context, message string, o EmptyStateOptions) EmptyStateResult {
	k := pixel.K(c)
	res := EmptyStateResult{}
	res.Element = ui.Column(c).Gap(k.D(6)).Padding(k.D(12)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(k.Msg("empty")).Children(func() {
		display.Sprite(c, openChest, display.SpriteOptions{Scale: 2})
		text.ShadowText(c, k.Msg("empty"), text.ShadowTextOptions{Slot: pixel.Faint, Scale: 2})
		pixel.Text(c, message, pixel.TextOptions{Slot: pixel.Faint})
		if o.Action != "" {
			res.Acted = input.Button(c, o.Action, input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled}).Clicked
		}
	})
	return res
}

type ErrorStateOptions struct {
	Disabled bool
}

type ErrorStateResult struct {
	Element  *ui.Element
	Continue bool
}

type startKey struct{}

// ErrorState is GAME OVER with a CONTINUE? countdown and RETRY.
func ErrorState(c *ui.Context, message string, o ErrorStateOptions) ErrorStateResult {
	k := pixel.K(c)
	col := ui.Column(c).Gap(k.D(8)).Padding(k.D(16)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(k.Msg("game_over"))
	start := ui.Local(col, startKey{}, func() time.Time { return k.Now() })
	left := max(0, 9-max(0, int(k.Now().Sub(*start)/time.Second)))
	if left > 0 {
		k.Tick(c, 10)
	}
	res := ErrorStateResult{Element: col}
	col.Children(func() {
		text.ShadowText(c, k.Msg("game_over"), text.ShadowTextOptions{Slot: pixel.Bad, Scale: 3})
		pixel.Text(c, message, pixel.TextOptions{Slot: pixel.Faint, Wrap: true, MaxLines: 6}).Width(k.D(min(pixel.Measure(message), 280)))
		display.Icon(c, icon.Skull, display.IconOptions{Scale: 3})
		text.ShadowText(c, k.Msg("continue")+" "+strconv.Itoa(left), text.ShadowTextOptions{Slot: pixel.Accent, Scale: 2})
		res.Continue = input.Button(c, k.Msg("retry"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled}).Clicked
	})
	if res.Continue {
		slog.Debug("8bit: continue", "left", left)
	}
	return res
}

type SuccessStateOptions struct {
	// Stars pop in one by one; zero means 3.
	Stars int
}

// SuccessState is STAGE CLEAR! with stars popping in steps.
func SuccessState(c *ui.Context, message string, o SuccessStateOptions) *ui.Element {
	k := pixel.K(c)
	n := pixel.Or("SuccessState Stars", o.Stars, 3)
	col := ui.Column(c).Gap(k.D(8)).Padding(k.D(16)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(k.Msg("stage_clear"))
	start := ui.Local(col, startKey{}, func() time.Time { return k.Now() })
	// Stars come four a second; Reduce Motion shows them all.
	shown := n + 1
	if !k.Reduce() {
		shown = min(n+1, max(0, int(k.Now().Sub(*start)*4/time.Second)))
		if shown <= n {
			k.Tick(c, 8)
		}
	}
	return col.Children(func() {
		text.ShadowText(c, k.Msg("stage_clear"), text.ShadowTextOptions{Slot: pixel.Good, Scale: 3})
		w := icon.Star.W*2 + 4
		ui.Box(c).Size(k.D(n*w), k.D(icon.Star.H*2+4)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range min(shown, n) {
				if i == shown-1 {
					for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
						g.Sprite(i*w+2+d[0], 2+d[1], icon.Star, 0, pixel.Lit, 2)
					}
				}
				g.Sprite(i*w+2, 2, icon.Star, 0, pixel.Warn, 2)
			}
		})
		pixel.Text(c, message, pixel.TextOptions{})
	})
}
