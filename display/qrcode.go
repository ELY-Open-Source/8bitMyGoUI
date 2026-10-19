package display

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"rsc.io/qr"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type QRCodeOptions struct {
	// Scale is dots per module; zero means 2.
	Scale int
}

type qrKey struct{}

type qrCache struct {
	text string
	code *qr.Code
}

// QRCode encodes s at level M, with a quiet zone.
func QRCode(c *ui.Context, s string, o QRCodeOptions) ui.Element {
	k := pixel.K(c)
	sc := pixel.Or("QRCode Scale", o.Scale, 2)
	e := ui.Box(c)
	qc := ui.Local(e, qrKey{}, func() qrCache { return qrCache{} })
	if qc.code == nil || qc.text != s {
		code, err := qr.Encode(s, qr.M)
		if err != nil {
			panic(fmt.Sprintf("display: QRCode %q: %v", s, err))
		}
		*qc = qrCache{s, code}
	}
	code := qc.code
	n := (code.Size + 8) * sc
	dark, light := extremes(k.Palette())
	return e.Size(k.D(n), k.D(n)).Role(ui.RoleImage).Label(s).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, n, n, light)
		for y := range code.Size {
			for x := 0; x < code.Size; {
				if !code.Black(x, y) {
					x++
					continue
				}
				end := x + 1
				for end < code.Size && code.Black(end, y) {
					end++
				}
				g.Rect((4+x)*sc, (4+y)*sc, (end-x)*sc, sc, dark)
				x = end
			}
		}
	})
}

// extremes are the darkest and lightest slots, for scanners.
func extremes(pal *pixel.Palette) (dark, light pixel.Slot) {
	lum := func(s pixel.Slot) int {
		c := pal.Colors[s]
		return 299*int(c.R) + 587*int(c.G) + 114*int(c.B)
	}
	for s := range pixel.Slot(pixel.Slots) {
		if lum(s) < lum(dark) {
			dark = s
		}
		if lum(s) > lum(light) {
			light = s
		}
	}
	return dark, light
}
