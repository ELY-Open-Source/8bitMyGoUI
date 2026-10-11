package input

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type PasswordInputOptions struct {
	Placeholder string
	MaxLength   int
	Width       int
	Disabled    bool
	// RevealLabel names the eye button; zero means "Show".
	RevealLabel string
}

type revealKey struct{}

type PasswordInputResult = TextInputResult

// PasswordInput masks with ★; the eye button shows the text.
func PasswordInput(c *ui.Context, text string, value *string, o PasswordInputOptions) PasswordInputResult {
	wrap := ui.Box(c)
	shown := ui.Local(wrap, revealKey{}, func() bool { return false })
	var res TextInputResult
	wrap.Children(func() {
		mask := '★'
		if *shown {
			mask = 0
		}
		ed, e := field(c, value, pixel.EditorOptions{Label: text, Placeholder: o.Placeholder, MaxLength: o.MaxLength, Mask: mask, Disabled: o.Disabled}, o.Width, nil, nil, func() {
			eye := icon.Eye
			if *shown {
				eye = icon.EyeOff
			}
			if IconButton(c, eye, label(o.RevealLabel, "Show"), IconButtonOptions{Variant: Ghost, Size: Small, Disabled: o.Disabled}).Clicked {
				*shown = !*shown
			}
		})
		res = result(e, ed, text)
	})
	res.Element = wrap
	return res
}
