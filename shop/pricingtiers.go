package shop

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Tier is a plan: its price and its perks.
type Tier struct {
	Price int
	Perks []string
}

type PricingTiersOptions struct {
	Currency Currency
	Disabled bool
}

type PricingTiersResult struct {
	Element *ui.Element
	Changed bool
}

const tierW = 112

var metals = [3]string{"bronze", "silver", "gold"}

// PricingTiers lays bronze, silver and gold plans to choose.
func PricingTiers(c *ui.Context, text string, tiers [3]Tier, chosen *int, o PricingTiersOptions) PricingTiersResult {
	k := pixel.K(c)
	if *chosen < -1 || *chosen > 2 {
		panic(fmt.Sprintf("shop: PricingTiers choice %d", *chosen))
	}
	sp, word := coin(k, o.Currency)
	most := 0
	for _, t := range tiers {
		amount("PricingTiers price", t.Price)
		for _, p := range t.Perks {
			if p == "" {
				panic("shop: PricingTiers has an empty perk")
			}
		}
		most = max(most, len(t.Perks))
	}
	// Band, name, price, perks, then the foot line under them.
	h := 8 + pixel.LineHeight + 4 + 2*pixel.LineHeight + 2 + most*pixel.LineHeight + pixel.LineHeight + 7
	old := *chosen
	var res PricingTiersResult
	res.Element = ui.Row(c).Gap(k.D(6)).FocusGroup(ui.Horizontal).Role(ui.RoleRadioGroup).Label(text).Children(func() {
		for i, t := range tiers {
			name, price := k.Msg(metals[i]), strconv.Itoa(t.Price)
			sel := *chosen
			// The shell draws the card; the radio stays bare.
			shell, base, st := pixel.Shell(c, func() *ui.Element {
				return ui.RadioBase(c, &sel, i).Label(name + ", " + price + " " + word).Disabled(o.Disabled)
			})
			// The radio picks before it is told it is off.
			if sel != *chosen && !base.IsDisabled() {
				*chosen = sel
			}
			base.Description(strings.Join(t.Perks, ", "))
			on := i == *chosen
			shell.Size(k.D(tierW), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Frame(0, 0, tierW, h, pixel.Raised, pixel.Face)
				// A band in medal metal; the name sits on Face.
				band := []pixel.Slot{pixel.Warn, pixel.Lit, pixel.Warn}[i]
				g.Rect(2, 2, tierW-4, 6, band)
				if i == 0 {
					g.Veil(2, 2, tierW-4, 6, pixel.Bad, 8)
				}
				g.Text((tierW-pixel.Measure(name))/2, 10, name, pixel.Ink, 1)
				y := 10 + pixel.LineHeight + 4
				g.Sprite((tierW-pixel.Measure(price)*2-14)/2, y+6, sp, 0, pixel.Ink, 1)
				g.Text((tierW-pixel.Measure(price)*2-14)/2+14, y, price, pixel.Ink, 2)
				y += 2*pixel.LineHeight + 2
				for _, perk := range t.Perks {
					g.Text(6, y, "•", pixel.Accent, 1)
					g.Text(14, y, pixel.Wrap(perk, tierW-18, 1)[0], pixel.Faint, 1)
					y += pixel.LineHeight
				}
				foot, ink := k.Msg("choose"), pixel.Faint
				if on {
					foot, ink = k.Msg("chosen"), pixel.Accent
					g.Rect(0, 0, tierW, 2, pixel.Accent)
					g.Rect(0, h-2, tierW, 2, pixel.Accent)
				}
				g.Text((tierW-pixel.Measure(foot))/2, h-pixel.LineHeight-3, foot, ink, 1)
				g.Finish(st)
			})
		}
	})
	res.Changed = *chosen != old
	if res.Changed {
		slog.Debug("8bit: pricing tier", "label", text, "tier", *chosen)
	}
	return res
}
