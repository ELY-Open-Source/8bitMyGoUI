package diagram

import (
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/egoist/mygo/ui"
)

// tokens maps Pele's tokens to palette slots' colors.
func tokens(p *pixel.Palette) palette {
	c := p.C
	return palette{bg: c(pixel.Back), fg: c(pixel.Ink), muted: c(pixel.Faint), line: c(pixel.Ink), surface: c(pixel.Face),
		alt: c(pixel.Well), border: c(pixel.Edge), accent: c(pixel.Accent), critical: c(pixel.Bad),
		series: [8]ui.Color{c(pixel.Accent), c(pixel.Good), c(pixel.Warn), c(pixel.Bad), c(pixel.Info), c(pixel.Magic),
			c(pixel.AccentLit), c(pixel.Shade)}}
}
