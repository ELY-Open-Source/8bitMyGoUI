package media

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type CartridgeOptions struct {
	// Title is up to two micro lines on the sticker.
	Title string
	// Art is the sticker's picture, up to 38 by 22.
	Art *pixel.Sprite
}

// Cartridge is a game cart with a sticker: art, title.
func Cartridge(c *ui.Context, text string, o CartridgeOptions) ui.Element {
	k := pixel.K(c)
	if o.Art != nil && (o.Art.W > 38 || o.Art.H > 22) {
		panic("media: Cartridge Art is larger than 38 by 22")
	}
	title := words(strings.ToUpper(o.Title), 38)
	if len(title) > 2 {
		panic(fmt.Sprintf("media: Cartridge Title %q needs %d lines", o.Title, len(title)))
	}
	return art(c, k, text, 56, 64).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		// A handheld cart's top right corner is cut.
		for y := range 64 {
			w := 56 - max(0, 6-y)
			g.HLine(0, y, w, pixel.Face)
			g.Dot(w-1, y, pixel.Edge)
		}
		g.HLine(0, 0, 50, pixel.Edge)
		g.HLine(0, 63, 56, pixel.Edge)
		g.VLine(0, 0, 64, pixel.Edge)
		g.HLine(1, 1, 49, pixel.Lit)
		g.VLine(1, 1, 62, pixel.Lit)
		g.HLine(1, 62, 54, pixel.Shade)
		for _, y := range []int{4, 6, 8} {
			g.HLine(5, y, 36, pixel.Shade)
		}
		g.Frame(6, 11, 44, 40, pixel.Sunken, pixel.Lit)
		g.Rect(8, 13, 40, 23, pixel.Accent)
		if o.Art != nil {
			g.Sprite(28-o.Art.W/2, 24-o.Art.H/2, o.Art, 0, pixel.AccentInk, 1)
		}
		for i, l := range title {
			g.Micro(28-pixel.MicroWidth(l)/2, 37+i*6+(2-len(title))*3, l, pixel.Ink, 1)
		}
		// An arrow pressed into the plastic: this end goes in.
		for i := range 4 {
			g.HLine(25+i, 54+i, 7-2*i, pixel.Shade)
		}
	})
}

// words fills micro lines w dots wide; long words panic.
func words(s string, w int) []string {
	var lines []string
	for _, word := range strings.Fields(s) {
		if !pixel.MicroOK(word) || pixel.MicroWidth(word) > w {
			panic(fmt.Sprintf("media: %q does not fit %d dots", word, w))
		}
		if n := len(lines); n > 0 && pixel.MicroWidth(lines[n-1]+" "+word) <= w {
			lines[n-1] += " " + word
			continue
		}
		lines = append(lines, word)
	}
	return lines
}
