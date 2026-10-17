package account

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TwoFactorOptions struct {
	// Period is a code's life; zero means 30 s.
	Period   time.Duration
	Disabled bool
}

type TwoFactorResult struct {
	Element ui.Element
	// Submitted is the sixth box filled.
	Changed, Submitted bool
}

// TwoFactor takes a six-digit code; a bar drains its period.
func TwoFactor(c *ui.Context, text string, code *string, o TwoFactorOptions) TwoFactorResult {
	k := pixel.K(c)
	period := o.Period
	if period == 0 {
		period = 30 * time.Second
	}
	if period < time.Second || len(*code) > 6 || strings.Trim(*code, "0123456789") != "" {
		panic(fmt.Sprintf("account: TwoFactor code of %d runes, period %v", len(*code), period))
	}
	var res TwoFactorResult
	// The bar spans the boxes, not the parent.
	res.Element = ui.Column(c).AlignSelf(ui.Start).Gap(k.D(4)).Children(func() {
		// PinInput takes typing, paste and Backspace.
		pin := input.PinInput(c, text, code, input.PinInputOptions{Length: 6, Disabled: o.Disabled})
		res.Changed, res.Submitted = pin.Changed, pin.Complete
		now := k.Now()
		k.At(c, now.Truncate(time.Second).Add(time.Second))
		secs := secondsLeft(now, period)
		ui.Box(c).FillWidth().Height(k.D(6)).Role(ui.RoleProgress).Label(k.Msg("time_left")).Description(fmt.Sprintf("%d s", secs)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			// The bar drains on 4 fps steps of the clock.
			rest := period - time.Duration(g.Busy(4))*(time.Second/4)%period
			bar := pixel.Good
			switch {
			case rest <= 5*time.Second:
				bar = pixel.Bad
			case rest <= 10*time.Second:
				bar = pixel.Warn
			}
			g.Frame(0, 0, g.W, 6, pixel.Sunken, pixel.Well)
			// Floats, so long periods cannot overflow.
			g.Rect(2, 2, int(float64(g.W-4)*float64(rest)/float64(period)), 2, bar)
		})
	})
	// The code never reaches the log, only its length.
	if res.Changed {
		slog.Debug("8bit: two factor", "label", text, "length", len(*code), "submitted", res.Submitted)
	}
	return res
}

// secondsLeft rounds up the time left in now's period, epoch-aligned.
func secondsLeft(now time.Time, period time.Duration) int64 {
	left := period - time.Duration(now.UnixNano()%int64(period))
	// left is at least 1 ns; this form cannot overflow.
	return int64(1 + (left-1)/time.Second)
}
