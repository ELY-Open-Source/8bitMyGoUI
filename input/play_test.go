package input_test

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// at clicks dot x, y of the element labeled l.
func at(t *testing.T, tt *ui.Tester, l string, x, y int) {
	t.Helper()
	r := find(t, tt, l)
	tt.ClickAt(r.X+float32(2*x), r.Y+float32(2*y))
}

func TestRatingClicksArrowsAndClears(t *testing.T) {
	stars, hearts, changes := 0, 0, 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		if input.Rating(c, "Stars", &stars, input.RatingOptions{}).Changed {
			changes++
		}
		input.Rating(c, "Love", &hearts, input.RatingOptions{Max: 3, Hearts: true})
	})
	click(t, tt, "3/5")
	click(t, tt, "3/5")
	if stars != 0 {
		t.Fatalf("clicking the rating again left %d", stars)
	}
	click(t, tt, "4/5")
	tt.Key(0, ui.KeyRight)
	if stars != 5 || changes != 4 {
		t.Fatalf("click and Right gave %d after %d changes", stars, changes)
	}
	tt.Key(0, ui.KeyHome)
	if stars != 1 {
		t.Fatalf("Home gave %d", stars)
	}
	// Only a filled mark has an Accent middle dot.
	lit := func(l string) bool {
		r := find(t, tt, l)
		return slots(t, tt.Image(), ui.Rect{X: r.X + r.W/2 - 1, Y: r.Y + r.H/2 - 1, W: 2, H: 2})[0] == pixel.Accent
	}
	if !lit("1/5") || lit("2/5") {
		t.Fatal("marks do not show the rating")
	}
	faint := func(l string) (out []bool) {
		for _, s := range slots(t, tt.Image(), find(t, tt, l)) {
			out = append(out, s == pixel.Faint)
		}
		return out
	}
	if equal(faint("3/5"), faint("3/3")) {
		t.Fatal("hearts look like stars")
	}
}

func TestRatingNavigationNeverClears(t *testing.T) {
	one, five := 1, 1
	tt := view(400, 200, nil, func(c *ui.Context) {
		input.Rating(c, "One", &one, input.RatingOptions{Max: 1})
		input.Rating(c, "Five", &five, input.RatingOptions{})
	})
	click(t, tt, "1/1")
	click(t, tt, "1/1")
	tt.Key(0, ui.KeyRight)
	if one != 1 {
		t.Fatalf("Right on a lone mark left %d", one)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	if five != 1 {
		t.Fatalf("Home onto the chosen mark left %d", five)
	}
}

func TestKnobDragsKeysAndSnapsToSixteenAngles(t *testing.T) {
	v, w := 0.5, 0.52
	tt := view(400, 200, nil, func(c *ui.Context) {
		input.Knob(c, "Gain", &v, input.KnobOptions{Step: 0.05})
		input.Knob(c, "Twin", &w, input.KnobOptions{})
	})
	if !equal(ink(t, tt.Image(), find(t, tt, "Gain")), ink(t, tt.Image(), find(t, tt, "Twin"))) {
		t.Fatal("0.5 and 0.52 share an angle but look apart")
	}
	r := find(t, tt, "Gain")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	// Each one-dot move is under half a step.
	for i := 1; i <= 16; i++ {
		tt.Move(x, y-float32(2*i))
	}
	tt.Release(x, y-32)
	if math.Abs(v-0.75) > 1e-9 {
		t.Fatalf("dragging up 16 dots gave %v, want 0.75", v)
	}
	tt.Key(0, ui.KeyDown)
	if math.Abs(v-0.7) > 1e-9 {
		t.Fatalf("Down gave %v", v)
	}
	tt.Key(0, ui.KeyEnd)
	if v != 1 || equal(ink(t, tt.Image(), find(t, tt, "Gain")), ink(t, tt.Image(), find(t, tt, "Twin"))) {
		t.Fatalf("End gave %v, or the pointer did not turn", v)
	}
}

func TestPinInputTypesBacksAndPastes(t *testing.T) {
	pin, complete := "", 0
	tt := view(400, 200, nil, func(c *ui.Context) {
		if input.PinInput(c, "PIN", &pin, input.PinInputOptions{}).Complete {
			complete++
		}
	})
	click(t, tt, "PIN")
	tt.Type("12a3")
	tt.Key(0, ui.KeyBackspace)
	if pin != "12" {
		t.Fatalf("typing and Backspace left %q", pin)
	}
	tt.SetClipboard("98765")
	tt.Command("paste")
	if pin != "1298" || complete != 1 {
		t.Fatalf("paste left %q, complete %d times", pin, complete)
	}
	if _, ok := tt.TextCaret(); !ok {
		t.Fatal("the boxes take no text from input methods")
	}
}

func TestNameEntrySpellsAndConfirms(t *testing.T) {
	name, confirmed := "", 0
	tt := view(400, 300, nil, func(c *ui.Context) {
		if input.NameEntry(c, "Initials", &name, input.NameEntryOptions{}).Confirmed {
			confirmed++
		}
	})
	at(t, tt, "Initials", 3*14+7, 7)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	if name != "DEL" || confirmed != 0 {
		t.Fatalf("spelled %q", name)
	}
	tt.Key(0, ui.KeyEnter)
	if confirmed != 1 {
		t.Fatal("the cursor did not jump to END after three letters")
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyEnter)
	if name != "DE" {
		t.Fatalf("DEL left %q", name)
	}
}

func TestVirtualKeyboardTypes(t *testing.T) {
	text, sent := "", 0
	tt := view(400, 300, nil, func(c *ui.Context) {
		if input.VirtualKeyboard(c, "Keys", &text, input.VirtualKeyboardOptions{}).Submitted {
			sent++
		}
	})
	at(t, tt, "Keys", 7, 14+7)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if text != "QW" {
		t.Fatalf("typed %q", text)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if sent != 1 {
		t.Fatal("enter did not submit")
	}
}

func TestDPadReportsArms(t *testing.T) {
	clk := &clock{time.Unix(3_000_000, 0)}
	var got []input.Direction
	tt := view(400, 200, clk, func(c *ui.Context) {
		if d := input.DPad(c, "Pad", input.DPadOptions{}).Pressed; d != input.None {
			got = append(got, d)
		}
	})
	at(t, tt, "Pad", 18, 18)
	at(t, tt, "Pad", 18, 4)
	tt.Key(0, ui.KeyRight)
	// The pointer rests on Up; Enter must not click it.
	tt.Key(0, ui.KeyEnter)
	if !equal(got, []input.Direction{input.Up, input.Right}) {
		t.Fatalf("got %v, want up then right", got)
	}
	r := find(t, tt, "Pad")
	arm := ui.Rect{X: r.X + 2*26, Y: r.Y + 2*14, W: 2, H: 2}
	if slots(t, tt.Image(), arm)[0] != pixel.Shade {
		t.Fatal("a key press did not light its arm")
	}
	clk.t = clk.t.Add(200 * time.Millisecond)
	tt.Frame()
	if slots(t, tt.Image(), arm)[0] != pixel.Face {
		t.Fatal("the arm stayed down")
	}
}

func TestHoldButtonNeedsTheHold(t *testing.T) {
	clk := &clock{time.Unix(3_000_000, 0)}
	confirmed := 0
	tt := view(400, 200, clk, func(c *ui.Context) {
		if input.HoldButton(c, "ERASE", input.HoldButtonOptions{}).Confirmed {
			confirmed++
		}
	})
	r := find(t, tt, "ERASE")
	step := func(d time.Duration) {
		clk.t = clk.t.Add(d)
		tt.Frame()
	}
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	step(500 * time.Millisecond)
	if !has(slots(t, tt.Image(), r), pixel.Accent) {
		t.Fatal("half a hold shows no fill")
	}
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	step(900 * time.Millisecond)
	if confirmed != 0 {
		t.Fatal("letting go did not restart the hold")
	}
	step(200 * time.Millisecond)
	step(200 * time.Millisecond)
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	if confirmed != 1 {
		t.Fatalf("a full hold confirmed %d times", confirmed)
	}
	// Keys: repeats every 50 ms after a 500 ms delay.
	tt.Key(0, ui.KeySpace)
	step(500 * time.Millisecond)
	for range 6 {
		tt.Key(0, ui.KeySpace)
		step(50 * time.Millisecond)
	}
	// Released at 750 ms; frames continue past one second.
	for range 12 {
		step(100 * time.Millisecond)
	}
	if confirmed != 1 {
		t.Fatal("keys let go early still confirmed")
	}
	tt.Key(0, ui.KeySpace)
	for range 25 {
		step(50 * time.Millisecond)
		tt.Key(0, ui.KeySpace)
	}
	if confirmed != 2 {
		t.Fatalf("a held key confirmed %d times in all", confirmed)
	}
}
