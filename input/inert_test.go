package input_test

import (
	"bytes"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestTextAreaHoldsStillWhenInert(t *testing.T) {
	v, inert := "1\n2\n3\n4\n5\n6\n7\n8", false
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Inert: inert})
		input.TextArea(c, "Notes", &v, input.TextAreaOptions{Lines: 2})
	}, 400, 200)
	r, ok := tt.Find("Notes")
	if !ok {
		t.Fatal("no Notes")
	}
	inert = true
	tt.Frame()
	before := tt.Image()
	// The caret keeps the view at the end; wheel upward.
	tt.Move(r.X+10, r.Y+10)
	tt.Scroll(r.X+10, r.Y+10, 0, -80)
	tt.Frame()
	if !bytes.Equal(before.Pix, tt.Image().Pix) {
		t.Fatal("the wheel scrolled an inert text area")
	}
}
