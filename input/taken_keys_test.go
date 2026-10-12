package input_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
)

func TestDPadKeysDoNotClickWhereThePointerRests(t *testing.T) {
	var pressed []input.Direction
	tt := view(300, 200, nil, func(c *ui.Context) {
		if r := input.DPad(c, "Pad", input.DPadOptions{}); r.Pressed != input.None {
			pressed = append(pressed, r.Pressed)
		}
	})
	r := find(t, tt, "Pad")
	// The pointer rests on the up arm.
	tt.Move(r.X+r.W/2, r.Y+r.H/6)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyEnter)
	if len(pressed) != 0 {
		t.Fatalf("Space and Enter pressed %v", pressed)
	}
}
