package account

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Language is a tongue by its own name, flagged.
type Language struct {
	Name string
	Flag *pixel.Sprite
}

type LanguagePickerOptions struct {
	// Width and Height are dots; zeros mean 160 by 120.
	Width, Height int
	Disabled      bool
}

type LanguagePickerResult = data.ListResult

// LanguagePicker lists languages with flags; ▶ marks the chosen one.
func LanguagePicker(c *ui.Context, text string, langs []Language, chosen *int, o LanguagePickerOptions) LanguagePickerResult {
	k := pixel.K(c)
	choice("LanguagePicker", chosen, len(langs))
	var res LanguagePickerResult
	ui.Box(c).Width(k.D(pixel.Or("LanguagePicker Width", o.Width, 160))).Children(func() {
		res = data.List(c, text, len(langs), chosen, data.ListOptions{Height: o.Height, Disabled: o.Disabled}, func(i int) string { return langs[i].Name }, func(i int) {
			l := langs[i]
			ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
				if l.Flag != nil {
					ui.Box(c).Size(k.D(l.Flag.W+2), k.D(l.Flag.H+2)).Draw(func(p *ui.Painter, r ui.Rect) {
						g := k.Canvas(p, r)
						g.Rect(0, 0, g.W, g.H, pixel.Edge)
						g.Sprite(1, 1, l.Flag, 0, pixel.Ink, 1)
					})
				}
				// The row is named already; its text only shows.
				mark(c, k, l.Name, pixel.Ink, 0)
			})
		})
	})
	return res
}
