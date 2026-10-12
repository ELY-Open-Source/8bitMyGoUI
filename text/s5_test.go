package text_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

func TestTypewriterKeepsItsSpeedOnS5Steps(t *testing.T) {
	clk, done := newClock(), 0
	tt := view(400, 100, clk, func(c *ui.Context) {
		// Twenty letters at the default twenty a second.
		if text.Typewriter(c, "ABCDEFGHIJKLMNOPQRST", text.TypewriterOptions{}).Done {
			done++
		}
	})
	clk.add(900 * time.Millisecond)
	tt.Frame()
	if done != 0 {
		t.Fatal("the text was out early")
	}
	clk.add(150 * time.Millisecond)
	tt.Frame()
	if done != 1 {
		t.Fatal("the text was not out after a second")
	}
}

func TestTextAnimationsSurviveTheClockSteppingBack(t *testing.T) {
	clk := newClock()
	tt := view(400, 200, clk, func(c *ui.Context) {
		text.Typewriter(c, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", text.TypewriterOptions{})
		text.ScrambleText(c, "SCRAMBLED EGGS", text.ScrambleTextOptions{})
		text.Marquee(c, "HELLO WORLD", text.MarqueeOptions{Width: 40})
	})
	for _, d := range []time.Duration{200 * time.Millisecond, -time.Hour, 200 * time.Millisecond} {
		clk.add(d)
		tt.Frame()
	}
}

func TestMarqueeHoldsWhenTheClockStepsBack(t *testing.T) {
	clk := newClock()
	tt := view(400, 100, clk, func(c *ui.Context) { text.Marquee(c, "HELLO WORLD", text.MarqueeOptions{Width: 60}) })
	clk.add(time.Second)
	tt.Frame()
	before := bytes.Clone(tt.Image().Pix)
	clk.add(-time.Hour)
	tt.Frame()
	if !bytes.Equal(before, tt.Image().Pix) {
		t.Fatal("the marquee ran backward as the clock stepped back")
	}
}

func TestTypingSpeedsHaveACeiling(t *testing.T) {
	defer func() {
		if r := fmt.Sprint(recover()); !strings.HasPrefix(r, "text: Typewriter CPS") {
			t.Fatalf("CPS 1001 gave %q", r)
		}
	}()
	view(400, 100, newClock(), func(c *ui.Context) { text.Typewriter(c, "HI", text.TypewriterOptions{CPS: 1001}) })
}
