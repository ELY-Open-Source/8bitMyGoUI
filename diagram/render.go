// Package diagram draws Mermaid diagrams on the dot grid.
package diagram

import (
	"context"

	pele "github.com/Broderick-Westrope/pele-go"
	_ "github.com/Broderick-Westrope/pele-go/all"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/egoist/mygo/ui"
)

// glyph is the pixel font's size; Pele's unit, a dot.
const glyph = 12

// measurer sizes text as the pixel font draws it.
type measurer struct{}

func (measurer) Width(text string, size float64, _ pele.Style) float64 {
	return float64(pixel.Measure(text) * scale(float32(size)))
}

// render draws Mermaid source onto dots in the palette's slots.
func render(source string, pal palette, slots *[pixel.Slots]ui.Color) (*scene, error) {
	off := false
	r, err := pele.Render(context.Background(), source, pele.Options{Measurer: measurer{}, FontSize: pele.Ptr(float64(glyph)),
		Links: &off, Images: &off})
	if err != nil {
		return nil, err
	}
	return parse(r.SVG, pal, slots)
}
