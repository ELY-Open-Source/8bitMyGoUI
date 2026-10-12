package layout

import "github.com/egoist/mygo/ui"

type LayersOptions struct{}

// Layers stacks layers in a box; the first sizes it.
func Layers(c *ui.Context, o LayersOptions, layers ...func()) *ui.Element {
	if len(layers) == 0 {
		panic("layout: Layers needs a layer")
	}
	return ui.Box(c).Children(func() {
		for i, l := range layers {
			b := ui.Box(c).Fill()
			if i > 0 {
				b.Absolute().Top(0).Left(0)
			}
			b.Children(l)
		}
	})
}
