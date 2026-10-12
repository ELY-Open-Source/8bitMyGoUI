package code_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/code"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestSpriteSheetHoldsAFrameWhenTheClockStepsBack(t *testing.T) {
	wd, frame, playing := newWorld(), 0, true
	sp := pixel.MustSprite("#\n\na")
	tt := view(300, 200, wd, func(c *ui.Context) { code.SpriteSheet(c, "Frames", sp, &frame, &playing, code.SpriteSheetOptions{}) })
	wd.clk.t = wd.clk.t.Add(-200 * time.Millisecond)
	tt.Frame()
	tt.Frame()
	if frame < 0 {
		t.Fatalf("frame %d", frame)
	}
}
