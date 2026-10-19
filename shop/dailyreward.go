package shop

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type DailyRewardOptions struct {
	Disabled bool
}

type DailyRewardResult struct {
	Element ui.Element
	Claimed bool
}

// An open chest, its lid up and its gold showing.
var openChest = pixel.MustSprite(`
		..eeeeeeee..
		.esslsslsse.
		.esssssssse.
		.eeeeeeeeee.
		.eyyyyylyye.
		.eylyyyyyye.
		.eeeeyyeeee.
		.essseyesse.
		.esssssssse.
		.esssssssse.
		.eeeeeeeeee.
		............`)

const dayW, dayH = 44, 2*pixel.LineHeight + 28

// DailyReward lines up seven chests; today's can be claimed.
func DailyReward(c *ui.Context, text string, today int, claimed []bool, o DailyRewardOptions) DailyRewardResult {
	k := pixel.K(c)
	if today < 1 || today > 7 || len(claimed) != 7 {
		panic(fmt.Sprintf("shop: DailyReward day %d of %d marks", today, len(claimed)))
	}
	said := make([]string, 7)
	for i, got := range claimed {
		state := k.Msg("locked")
		switch {
		case i >= today && got:
			panic(fmt.Sprintf("shop: DailyReward day %d is claimed before it came", i+1))
		case got:
			state = k.Msg("claimed")
		case i < today-1:
			state = k.Msg("missed")
		case i == today-1:
			state = k.Msg("today")
		}
		said[i] = fmt.Sprintf(k.Msg("day_n"), i+1) + " " + state
	}
	marks := append([]bool(nil), claimed...)
	var res DailyRewardResult
	res.Element = ui.Column(c).Gap(k.D(6)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(text).Description(strings.Join(said, ", ")).Children(func() {
		ui.Box(c).Size(k.D(7*dayW+6*4), k.D(dayH)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i := range 7 {
				x, day := i*(dayW+4), fmt.Sprintf(k.Msg("day_n"), i+1)
				g.Frame(x, 0, dayW, dayH, pixel.Raised, pixel.Face)
				g.Text(x+(dayW-pixel.Measure(day))/2, 3, day, pixel.Faint, 1)
				cx, cy := x+(dayW-24)/2, pixel.LineHeight+6
				switch {
				case marks[i]:
					g.Sprite(cx, cy, openChest, 0, pixel.Ink, 2)
					g.Sprite(x+dayW-14, dayH-14, icon.Check, 0, pixel.Good, 1)
				case i == today-1:
					// Today's chest glows, ringing on 4 fps steps.
					ring := pixel.Accent
					if g.Step(4)%2 == 1 {
						ring = pixel.AccentLit
					}
					g.Rect(x, 0, dayW, 2, ring)
					g.Rect(x, dayH-2, dayW, 2, ring)
					g.Sprite(cx, cy, icon.Chest, 0, pixel.Ink, 2)
				default:
					// Missed and later days stand in shadow.
					g.Silhouette(cx, cy, icon.Chest, 0, pixel.Shade, 2)
					if i >= today {
						g.Sprite(x+dayW-14, dayH-14, icon.Lock, 0, pixel.Faint, 1)
					}
				}
			}
		})
		done := claimed[today-1]
		if input.Button(c, k.Msg("claim"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled || done}).Clicked {
			claimed[today-1], res.Claimed = true, true
			slog.Debug("8bit: daily reward", "label", text, "day", today)
		}
		if done {
			pixel.Text(c, k.Msg("see_you"), pixel.TextOptions{Slot: pixel.Faint})
		}
	})
	return res
}
