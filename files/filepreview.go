package files

import (
	"image"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/display"
	"github.com/ELY-Open-Source/8bitMyGoUI/layout"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Preview is what a file shows: a picture, else text.
type Preview struct {
	Picture image.Image
	Text    string
}

type FilePreviewOptions struct {
	// W and H are dots; zeros mean 160 and 120.
	W, H int
}

type previewKey struct{}

// FilePreview shows a picture, else text, in a sunken pane.
func FilePreview(c *ui.Context, name string, pv Preview, o FilePreviewOptions) ui.Element {
	k := pixel.K(c)
	w, h := pixel.Or("FilePreview W", o.W, 160), pixel.Or("FilePreview H", o.H, 120)
	return ui.Column(c).Gap(k.D(2)).AlignItems(ui.Start).Role(ui.RoleGroup).Label(name).Children(func() {
		// The group says the name; the title only shows it.
		ui.Box(c).Size(k.D(pixel.Measure(name)), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
			k.Canvas(p, r).Text(0, 0, name, pixel.Faint, 1)
		})
		layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well, Padding: 2}, func() {
			pane := ui.Column(c).Size(k.D(w), k.D(h)).AlignItems(ui.Center).Justify(ui.Center)
			pane.Children(func() {
				switch {
				case pv.Picture != nil:
					pw, ph := fit(pv.Picture.Bounds(), w, h)
					display.Picture(c, pv.Picture, display.PictureOptions{W: pw, H: ph, Label: k.Msg("preview")})
				case pv.Text != "":
					s := ui.Local(pane, previewKey{}, func() ui.ScrollState { return ui.ScrollState{} })
					layout.ScrollArea(c, s, layout.ScrollAreaOptions{Label: k.Msg("preview")}, func() {
						pixel.Text(c, pv.Text, pixel.TextOptions{Wrap: true})
					}).Element.Size(k.D(w), k.D(h))
				default:
					pixel.Text(c, k.Msg("no_preview"), pixel.TextOptions{Slot: pixel.Faint})
				}
			})
		})
	})
}

// fit scales b into w×h: whole steps up, any down.
func fit(b image.Rectangle, w, h int) (int, int) {
	bw, bh := max(1, b.Dx()), max(1, b.Dy())
	s := min(float64(w)/float64(bw), float64(h)/float64(bh))
	if s >= 1 {
		s = float64(int(s))
	}
	return max(1, int(float64(bw)*s)), max(1, int(float64(bh)*s))
}
