package shop

import (
	"fmt"
	"hash/fnv"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// ReceiptLine is a bought ware.
type ReceiptLine struct {
	Name       string
	Qty, Price int
}

type ReceiptOptions struct {
	// Width is in dots; zero means 180.
	Width int
	// Time prints under the shop; zero prints none.
	Time time.Time
}

// Receipt prints lines dot-matrix style, perforated down the sides.
func Receipt(c *ui.Context, shop string, lines []ReceiptLine, o ReceiptOptions) ui.Element {
	k := pixel.K(c)
	if shop == "" {
		panic("shop: Receipt has no shop name")
	}
	w, total := pixel.Or("Receipt Width", o.Width, 180), 0
	rows := make([][2]string, len(lines))
	for i, l := range lines {
		amount("Receipt price", l.Price)
		if l.Name == "" || l.Qty < 1 || l.Price > (math.MaxInt-total)/l.Qty {
			panic(fmt.Sprintf("shop: Receipt line %q, %d at %d", l.Name, l.Qty, l.Price))
		}
		total += l.Qty * l.Price
		rows[i] = [2]string{fmt.Sprintf("%d×%s", l.Qty, l.Name), strconv.Itoa(l.Qty * l.Price)}
	}
	sum, thanks := strconv.Itoa(total), k.Msg("thank_you")
	said := make([]string, len(rows))
	for i, r := range rows {
		said[i] = r[0] + " " + r[1]
	}
	h := pixel.LineHeight*(len(rows)+5) + 24
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(shop).Description(strings.Join(append(said, k.Msg("total")+" "+sum), ", ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		// Paper, with feed holes punched down both sides.
		g.Rect(0, 0, w, h-4, pixel.Lit)
		for y := 4; y < h-8; y += 8 {
			g.Rect(3, y, 2, 2, pixel.Shade)
			g.Rect(w-5, y, 2, 2, pixel.Shade)
		}
		// A torn edge of teeth along the bottom.
		for x := 0; x < w; x += 4 {
			g.Rect(x, h-4, 2, 2, pixel.Lit)
			g.Rect(x+1, h-2, 1, 2, pixel.Lit)
		}
		line := func(y int, l, rt string, ink pixel.Slot) {
			l = pixel.Wrap(l, w-28-pixel.Measure(rt), 1)[0]
			g.Text(10, y, l, ink, 1)
			g.Text(w-10-pixel.Measure(rt), y, rt, ink, 1)
			// Dot leaders run from the name to the price.
			for x := 12 + pixel.Measure(l); x < w-12-pixel.Measure(rt); x += 3 {
				g.Dot(x, y+12, pixel.Faint)
			}
		}
		y := 4
		g.Text((w-pixel.Measure(shop))/2, y, shop, pixel.Ink, 1)
		if !o.Time.IsZero() {
			y += pixel.LineHeight
			d := o.Time.Format("2006-01-02 15:04")
			g.Text((w-pixel.Measure(d))/2, y, d, pixel.Faint, 1)
		}
		y += pixel.LineHeight + 2
		for x := 10; x < w-10; x += 4 {
			g.Rect(x, y, 2, 1, pixel.Shade)
		}
		y += 3
		for _, rw := range rows {
			line(y, rw[0], rw[1], pixel.Ink)
			y += pixel.LineHeight
		}
		line(y+2, k.Msg("total"), sum, pixel.Ink)
		y += pixel.LineHeight + 4
		g.Text((w-pixel.Measure(thanks))/2, y, thanks, pixel.Ink, 1)
		// A barcode drawn from the total, for show.
		hs := fnv.New32a()
		hs.Write([]byte(sum))
		bits := hs.Sum32()
		for i, x := 0, (w-64)/2; i < 32; i, x = i+1, x+2 {
			if bits>>i&1 == 1 || i%8 == 0 {
				g.Rect(x, y+pixel.LineHeight+2, 1, 10, pixel.Ink)
			}
		}
	})
}
