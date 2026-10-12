package display

import (
	"image"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PictureOptions struct {
	// W×H is the size in dots; zero takes the image's.
	W, H  int
	Label string
}

type pictureKey struct{}

type pictureCache struct {
	sum  uint64
	pal  *pixel.Palette
	w, h int
	dots []int8
}

// Picture reduces img to the palette with Bayer dither.
func Picture(c *ui.Context, img image.Image, o PictureOptions) *ui.Element {
	k := pixel.K(c)
	b := img.Bounds()
	w, h := pixel.Or("Picture W", o.W, b.Dx()), pixel.Or("Picture H", o.H, b.Dy())
	e := art(c, k, o.Label, w, h, nil)
	pc := ui.Local(e, pictureKey{}, func() pictureCache { return pictureCache{} })
	if sum := fingerprint(img); pc.dots == nil || pc.sum != sum || pc.pal != k.Palette() || pc.w != w || pc.h != h {
		*pc = pictureCache{sum, k.Palette(), w, h, reduce(img, w, h, k.Palette())}
	}
	dots := pc.dots
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for y := range h {
			for x := 0; x < w; {
				s, end := dots[y*w+x], x+1
				for end < w && dots[y*w+end] == s {
					end++
				}
				if s >= 0 {
					g.Rect(x, y, end-x, 1, pixel.Slot(s))
				}
				x = end
			}
		}
	})
}

// reduce maps img to w×h slots; -1 marks clear dots.
func reduce(img image.Image, w, h int, pal *pixel.Palette) []int8 {
	b := img.Bounds()
	out := make([]int8, w*h)
	for y := range h {
		y0, y1 := span(b.Min.Y, b.Dy(), y, h)
		for x := range w {
			x0, x1 := span(b.Min.X, b.Dx(), x, w)
			var sum [4]uint32
			n := uint32(0)
			for sy := y0; sy < y1; sy++ {
				for sx := x0; sx < x1; sx++ {
					r, g, bl, a := img.At(sx, sy).RGBA()
					sum[0], sum[1], sum[2], sum[3] = sum[0]+r>>8, sum[1]+g>>8, sum[2]+bl>>8, sum[3]+a>>8
					n++
				}
			}
			if sum[3]/n < 128 {
				out[y*w+x] = -1
				continue
			}
			nudge := (float64(pixel.Bayer(x, y))+0.5)/16 - 0.5
			out[y*w+x] = int8(nearest(pal, [3]float64{
				float64(sum[0]/n) + nudge*48, float64(sum[1]/n) + nudge*48, float64(sum[2]/n) + nudge*48,
			}))
		}
	}
	return out
}

// fingerprint mixes img's bounds and pixels, to notice edits.
func fingerprint(img image.Image) uint64 {
	b := img.Bounds()
	s := uint64(b.Min.X)<<48 ^ uint64(b.Min.Y)<<32 ^ uint64(b.Dx())<<16 ^ uint64(b.Dy())
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, a := img.At(x, y).RGBA()
			s = (s ^ uint64(r>>8)<<24 ^ uint64(g>>8)<<16 ^ uint64(bl>>8)<<8 ^ uint64(a>>8)) * 1099511628211
		}
	}
	return s
}

// span is the source range under cell i of n.
func span(min, size, i, n int) (int, int) {
	a, b := min+i*size/n, min+(i+1)*size/n
	return a, max(b, a+1)
}

// nearest is the slot closest to rgb.
func nearest(pal *pixel.Palette, rgb [3]float64) pixel.Slot {
	best, dist := pixel.Ink, -1.0
	for s, c := range pal.Colors {
		dr, dg, db := rgb[0]-float64(c.R), rgb[1]-float64(c.G), rgb[2]-float64(c.B)
		if d := dr*dr + dg*dg + db*db; dist < 0 || d < dist {
			best, dist = pixel.Slot(s), d
		}
	}
	return best
}
