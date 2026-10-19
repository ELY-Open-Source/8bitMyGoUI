package input

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type IconButtonOptions struct {
	Variant  Variant
	Size     Size
	Disabled bool
}

type IconButtonResult = ButtonResult

// IconButton is a square button named for assistive technology.
func IconButton(c *ui.Context, icon *pixel.Sprite, name string, o IconButtonOptions) IconButtonResult {
	if name == "" {
		panic("input: IconButton needs a name")
	}
	return Button(c, "", ButtonOptions{Variant: o.Variant, Size: o.Size, Icon: icon, Disabled: o.Disabled, Label: name})
}
