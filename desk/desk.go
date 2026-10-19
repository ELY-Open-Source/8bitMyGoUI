// Package desk is a retro desktop: windows, a taskbar, games.
package desk

import (
	"fmt"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Wallpaper is a dithered desktop pattern.
type Wallpaper int

const (
	Checks Wallpaper = iota
	Bricks
	Weave
	Dots
	Waves
	Fade
)

var papers = [...]string{"paper_checks", "paper_bricks", "paper_weave", "paper_dots", "paper_waves", "paper_fade"}

// name is the pattern's name; unknown ones panic.
func (wp Wallpaper) name(k pixel.Kit) string {
	if wp < Checks || wp > Fade {
		panic(fmt.Sprintf("desk: unknown Wallpaper %d", wp))
	}
	return k.Msg(papers[wp])
}

// paper fills a box with the pattern from its corner.
func paper(g pixel.Canvas, x, y, w, h int, wp Wallpaper) {
	switch wp {
	case Checks:
		g.Rect(x, y, w, h, pixel.Well)
		for j := 0; j < h; j += 4 {
			for i := j / 4 % 2 * 4; i < w; i += 8 {
				g.Rect(x+i, y+j, min(4, w-i), min(4, h-j), pixel.Info)
			}
		}
	case Bricks:
		g.Rect(x, y, w, h, pixel.Shade)
		for j := 0; j < h; j += 5 {
			for i := j/5%2*4 - 4; i < w; i += 8 {
				l, r := max(i, 0), min(i+7, w)
				g.Rect(x+l, y+j, r-l, min(4, h-j), pixel.Face)
				g.Rect(x+l, y+j, r-l, 1, pixel.Lit)
			}
		}
	case Weave:
		g.Rect(x, y, w, h, pixel.Well)
		for j := range h {
			for i := range w {
				if (i+j)%6 == 0 || (i-j+6*h)%6 == 0 {
					g.Dot(x+i, y+j, pixel.Accent)
				}
			}
		}
	case Dots:
		g.Rect(x, y, w, h, pixel.Well)
		for j := 1; j < h; j += 4 {
			for i := 1 + j/4%2*3; i < w; i += 6 {
				g.Dot(x+i, y+j, pixel.Faint)
			}
		}
	case Waves:
		g.Rect(x, y, w, h, pixel.Info)
		for j := 0; j < h; j += 6 {
			for i := range w {
				g.Dot(x+i, y+j+[]int{0, 1, 2, 1}[i/2%4], pixel.Magic)
			}
		}
	case Fade:
		// Seventeen Bayer bands run Magic down into Info.
		for b := range 17 {
			top, foot := b*h/17, (b+1)*h/17
			g.Dither(x, y+top, w, foot-top, pixel.Magic, pixel.Info, b)
		}
	default:
		panic(fmt.Sprintf("desk: unknown Wallpaper %d", wp))
	}
}

// tile draws an icon over a label of two lines.
func tile(g pixel.Canvas, x, y, w int, sp *pixel.Sprite, lines []string, chosen bool) {
	ink, back, text := pixel.Ink, pixel.Face, pixel.Ink
	// Chosen, the icon takes the Accent and the label inverts.
	if chosen {
		ink, back, text = pixel.Accent, pixel.Accent, pixel.AccentInk
	}
	g.Sprite(x+(w-2*sp.W)/2, y+2+(12-sp.H), sp, 0, ink, 2)
	for i, l := range lines {
		lw := pixel.Measure(l)
		lx, ly := x+(w-lw)/2, y+28+i*pixel.LineHeight
		g.Rect(lx-1, ly+1, lw+2, pixel.LineHeight-1, back)
		g.Text(lx, ly, l, text, 1)
	}
}

// sprite checks a 12-dot icon, named for the panic.
func sprite(name string, sp *pixel.Sprite) *pixel.Sprite {
	if sp == nil || sp.W > 12 || sp.H > 12 {
		panic(fmt.Sprintf("desk: %s needs an icon of 12 dots at most", name))
	}
	return sp
}

// barLook paints a title bar: Accent while active, else Shade.
func barLook(g pixel.Canvas, title string, active bool) {
	back, ink := pixel.Shade, pixel.Lit
	if active {
		back, ink = pixel.Accent, pixel.AccentInk
	}
	g.Rect(0, 0, g.W, g.H, back)
	g.Text(3, (g.H-pixel.LineHeight)/2, pixel.Wrap(title, g.W-6, 1)[0], ink, 1)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
