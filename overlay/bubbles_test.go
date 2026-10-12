package overlay_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestPopoverOpensBelowWithANotch(t *testing.T) {
	open := false
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() {
			b := input.Button(c, "MAP", input.ButtonOptions{})
			if b.Clicked {
				open = !open
			}
			overlay.Popover(c, "Map", b.Element, &open, overlay.PopoverOptions{}, func() {
				pixel.Text(c, "You are here.", pixel.TextOptions{})
			})
		})
	})
	click(t, tt, "MAP")
	a, p := find(t, tt, "MAP"), find(t, tt, "You are here.")
	if p.Y <= a.Y+a.H {
		t.Fatalf("panel text at %v, anchor %v", p, a)
	}
	// The notch tip sits under the anchor's middle.
	if tip := slot(t, tt.Image(), int(a.X+a.W/2), int(a.Y+a.H)+1); tip != pixel.Edge && tip != pixel.Face {
		t.Fatalf("no notch under the anchor: %v", tip)
	}
	keeps(t, tt.Image())
	tt.Key(0, ui.KeyEscape)
	if open || shows(tt, "You are here.") {
		t.Fatal("Escape did not close it")
	}
}

func TestTooltipWaitsForRestOrKeyboard(t *testing.T) {
	clk := newClock()
	tt := view(640, 400, clk, false, func(c *ui.Context) {
		ui.Column(c).Padding(80).Gap(20).Children(func() {
			input.Button(c, "FIRST", input.ButtonOptions{})
			overlay.Tooltip(c, "Casts fire", overlay.TooltipOptions{}, func() {
				input.Button(c, "FIRE", input.ButtonOptions{})
			})
		})
	})
	r := find(t, tt, "FIRE")
	tt.Move(r.X+4, r.Y+4)
	clk.t = clk.t.Add(300 * time.Millisecond)
	tt.Frame()
	if shows(tt, "Casts fire") {
		t.Fatal("the tip showed before 500 ms")
	}
	clk.t = clk.t.Add(300 * time.Millisecond)
	tt.Frame()
	if !shows(tt, "Casts fire") {
		t.Fatal("the tip did not show after 600 ms")
	}
	tt.ClickAt(r.X+4, r.Y+4)
	tt.Move(600, 380)
	if shows(tt, "Casts fire") {
		t.Fatal("focus from a click kept the tip")
	}
	click(t, tt, "FIRST")
	tt.Move(600, 380)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("FIRE") || !shows(tt, "Casts fire") {
		t.Fatal("keyboard focus did not show the tip")
	}
	keeps(t, tt.Image())
}

func TestHoverCardStaysWhileThePointerMovesOntoIt(t *testing.T) {
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() {
			overlay.HoverCard(c, "Aria", overlay.HoverCardOptions{}, func() {
				pixel.Text(c, "@aria", pixel.TextOptions{})
			}, func() {
				pixel.Text(c, "Knight, LV 7", pixel.TextOptions{})
			})
		})
	})
	if shows(tt, "Knight, LV 7") {
		t.Fatal("the card showed unhovered")
	}
	r := find(t, tt, "@aria")
	tt.Move(r.X+4, r.Y+4)
	tt.Frame()
	card := find(t, tt, "Knight, LV 7")
	tt.Move(card.X+4, card.Y+4)
	tt.Frame()
	if !shows(tt, "Knight, LV 7") {
		t.Fatal("moving onto the card hid it")
	}
	tt.Move(600, 380)
	tt.Frame()
	if shows(tt, "Knight, LV 7") {
		t.Fatal("the card stayed after the pointer left")
	}
}

func TestPopconfirmAsksFirst(t *testing.T) {
	yes, no := 0, 0
	disabled := false
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() {
			r := overlay.Popconfirm(c, "SELL", "Sell the crown?", overlay.PopconfirmOptions{Disabled: disabled})
			if r.Confirmed {
				yes++
			}
			if r.Canceled {
				no++
			}
		})
	})
	click(t, tt, "SELL")
	click(t, tt, "YES")
	click(t, tt, "SELL")
	click(t, tt, "NO")
	click(t, tt, "SELL")
	tt.Key(0, ui.KeyEscape)
	if yes != 1 || no != 2 || shows(tt, "Sell the crown?") {
		t.Fatalf("yes %d no %d", yes, no)
	}
	click(t, tt, "SELL")
	disabled = true
	tt.Frame()
	if shows(tt, "Sell the crown?") {
		t.Fatal("the bubble stayed open as its button turned off")
	}
	click(t, tt, "SELL")
	if shows(tt, "Sell the crown?") {
		t.Fatal("a disabled popconfirm opened")
	}
}

func TestCoachmarkCountsAndGoesOn(t *testing.T) {
	next := 0
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() {
			b := input.Button(c, "INVENTORY", input.ButtonOptions{})
			if overlay.Coachmark(c, b.Element, "Your items live here.", 2, 3, overlay.CoachmarkOptions{}).Next {
				next++
			}
		})
	})
	if !shows(tt, "2/3") || find(t, tt, "Your items live here.").Y <= find(t, tt, "INVENTORY").Y {
		t.Fatalf("texts %q", tt.Texts())
	}
	click(t, tt, "NEXT ▶")
	if next != 1 {
		t.Fatal("NEXT did not report")
	}
	keeps(t, tt.Image())
}

func TestSpotlightDimsAllButTheTarget(t *testing.T) {
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		ui.Column(c).Padding(60).Children(func() {
			b := input.Button(c, "SHOP", input.ButtonOptions{})
			overlay.Spotlight(c, b.Element, "Buy potions here.", overlay.SpotlightOptions{})
		})
	})
	tt.Frame()
	img, r := tt.Image(), find(t, tt, "SHOP")
	if in := count(t, img, r, pixel.Ink); in > int(r.W*r.H)/16 {
		t.Fatalf("%d Ink dots veil the target", in)
	}
	far := ui.Rect{X: 400, Y: 300, W: 100, H: 60}
	if n := count(t, img, far, pixel.Ink); n < int(far.W*far.H)/8 {
		t.Fatalf("only %d Ink dots far from the target", n)
	}
	if !shows(tt, "Buy potions here.") {
		t.Fatal("no caption")
	}
	keeps(t, img)
}

func TestPopoverRefusesNegativeWidth(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("a negative width did not panic")
		}
	}()
	open := false
	view(200, 100, newClock(), false, func(c *ui.Context) {
		b := input.Button(c, "X", input.ButtonOptions{})
		overlay.Popover(c, "X", b.Element, &open, overlay.PopoverOptions{Width: -1}, func() {})
	})
}

func TestPopconfirmParksFocusAsItTurnsOff(t *testing.T) {
	disabled := false
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() {
			overlay.Popconfirm(c, "SELL", "Sell?", overlay.PopconfirmOptions{Disabled: disabled})
			pixel.Text(c, "elsewhere", pixel.TextOptions{})
		})
	})
	click(t, tt, "SELL")
	tt.Key(0, ui.KeyTab)
	disabled = true
	tt.Frame()
	if tt.Focused("SELL") && !tt.Focused("elsewhere") {
		t.Fatal("the focus came back to the button as it turned off")
	}
}
