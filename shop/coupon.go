package shop

import (
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type CouponOptions struct {
	Disabled bool
}

type CouponResult struct {
	Element *ui.Element
	Copied  bool
}

type couponKey struct{}

// Coupon is a notched ticket: a deal, a code, COPY.
func Coupon(c *ui.Context, title, code string, o CouponOptions) CouponResult {
	k := pixel.K(c)
	if code == "" || title == "" {
		panic("shop: Coupon needs a code and a title")
	}
	var res CouponResult
	row := ui.Row(c).Padding(k.D(6), k.D(14)).Gap(k.D(8)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(title)
	res.Element = row.AlignSelf(ui.Start).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		// Row by row, leaving round notches at both ends.
		for y := range g.H {
			in := 0
			if dy := y - g.H/2; dy*dy < 36 {
				in = int(math.Sqrt(float64(36 - dy*dy)))
			}
			g.Rect(in, y, g.W-2*in, 1, pixel.Face)
			g.Dot(in, y, pixel.Edge)
			g.Dot(g.W-1-in, y, pixel.Edge)
		}
		g.HLine(0, 0, g.W, pixel.Edge)
		g.HLine(0, g.H-1, g.W, pixel.Edge)
	}).Children(func() {
		ui.Column(c).Gap(k.D(2)).Children(func() {
			pixel.Text(c, title, pixel.TextOptions{Slot: pixel.Bad, Scale: 2})
			ui.Box(c).Size(k.D(pixel.Measure(code)+10), k.D(pixel.LineHeight+4)).Role(ui.RoleText).Label(code).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for x := 0; x < g.W; x += 3 {
					g.Rect(x, 0, 2, 1, pixel.Edge)
					g.Rect(x, g.H-1, 2, 1, pixel.Edge)
				}
				for y := 0; y < g.H; y += 3 {
					g.Rect(0, y, 1, 2, pixel.Edge)
					g.Rect(g.W-1, y, 1, 2, pixel.Edge)
				}
				g.Text(5, 2, code, pixel.Ink, 1)
			})
		})
		// A torn-here line between the deal and the button.
		ui.Box(c).Size(k.D(1), k.D(36)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for y := 0; y < g.H; y += 4 {
				g.Rect(0, y, 1, 2, pixel.Shade)
			}
		})
		e := ui.Box(c)
		at := ui.Local(e, couponKey{}, func() time.Time { return time.Time{} })
		// The button says COPIED! for two seconds.
		label := k.Msg("copy")
		if k.Since(c, *at, 4, 8) < 8 {
			label = k.Msg("copied")
		}
		e.Children(func() {
			if input.Button(c, label, input.ButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked {
				c.WriteClipboard(code)
				*at, res.Copied = k.Now(), true
				// The code is a key; logs never see it.
				slog.Debug("8bit: coupon copied", "title", title)
			}
		})
	})
	return res
}
