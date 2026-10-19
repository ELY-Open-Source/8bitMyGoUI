package code

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Bank is a stretch of memory and its use.
type Bank struct {
	Name       string
	Used, Size int
}

type MemoryMapOptions struct {
	// Columns of blocks; zero means 4.
	Columns int
}

// MemoryMap draws each bank as a block filled by use.
func MemoryMap(c *ui.Context, text string, banks []Bank, o MemoryMapOptions) ui.Element {
	k := pixel.K(c)
	cols := pixel.Or("MemoryMap Columns", o.Columns, 4)
	for _, b := range banks {
		name := strings.ToUpper(b.Name)
		if b.Size <= 0 || b.Used < 0 || b.Used > b.Size || !pixel.MicroOK(name) || pixel.MicroWidth(name) > 36 {
			panic(fmt.Sprintf("code: MemoryMap bank %+v", b))
		}
	}
	return ui.Grid(c).Columns(cols).Width(k.D(cols * 44)).Gap(k.D(4)).Role(ui.RoleGroup).Label(text).Children(func() {
		for _, b := range banks {
			pct := b.Used * 100 / b.Size
			label := fmt.Sprintf("%s, %d of %d KB, %d%%", b.Name, b.Used/1024, b.Size/1024, pct)
			ui.Box(c).Size(k.D(40), k.D(45)).Role(ui.RoleImage).Label(label).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				g.Frame(0, 0, 40, 32, pixel.Raised, pixel.Face)
				g.Frame(3, 3, 34, 26, pixel.Sunken, pixel.Well)
				fill := pixel.Good
				switch {
				case pct >= 80:
					fill = pixel.Bad
				case pct >= 50:
					fill = pixel.Warn
				}
				// The fill rises from the floor, dithered on top.
				if n := b.Used * 22 / b.Size; n > 0 {
					g.Rect(5, 27-n, 30, n, fill)
					g.Dither(5, 27-n, 30, 1, fill, pixel.Well, 8)
				}
				g.Micro(20-pixel.MicroWidth(strings.ToUpper(b.Name))/2, 34, strings.ToUpper(b.Name), pixel.Ink, 1)
				s := fmt.Sprintf("%d%%", pct)
				g.Micro(20-pixel.MicroWidth(s)/2, 40, s, pixel.Faint, 1)
			})
		}
	})
}
