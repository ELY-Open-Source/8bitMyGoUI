package overlay_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestDrawerSlidesInInSteps(t *testing.T) {
	clk, open, closed := newClock(), false, 0
	side := display.Right
	tt := view(640, 400, clk, false, func(c *ui.Context) {
		if input.Button(c, "MENU", input.ButtonOptions{}).Clicked {
			open = true
		}
		if overlay.Drawer(c, "Menu", &open, overlay.DrawerOptions{Side: side, Size: 100}, func() {
			pixel.Text(c, "Items", pixel.TextOptions{})
		}).Closed {
			closed++
		}
	})
	click(t, tt, "MENU")
	var xs []float32
	for range 5 {
		xs = append(xs, find(t, tt, "Items").X)
		clk.t = clk.t.Add(70 * time.Millisecond)
		tt.Frame()
	}
	// Four 50-pixel steps, then rest at 640-200+16.
	want := []float32{456 + 200, 456 + 150, 456 + 100, 456 + 50, 456}
	for i := range want {
		if xs[i] != want[i] {
			t.Fatalf("x by step %v, want %v", xs, want)
		}
	}
	keeps(t, tt.Image())
	tt.Key(0, ui.KeyEscape)
	if open || closed != 1 {
		t.Fatal("Escape left it open")
	}
	side, open = display.Left, true
	tt.Frame()
	clk.t = clk.t.Add(time.Second)
	tt.Frame()
	if x := find(t, tt, "Items").X; x != 16 {
		t.Fatalf("a left drawer opened at %v", x)
	}
}

func TestSheetClosesOnALongPull(t *testing.T) {
	open, closed := true, 0
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		if overlay.Sheet(c, "Options", &open, overlay.SheetOptions{Height: 120}, func() {
			pixel.Text(c, "Sound: ON", pixel.TextOptions{})
		}).Closed {
			closed++
		}
	})
	// The bar tops the sheet, which is 240 pixels.
	top := float32(400 - 240)
	tt.Press(320, top+8)
	tt.Move(320, top+40)
	tt.Release(320, top+40)
	if !open || find(t, tt, "Sound: ON").Y > top+40 {
		t.Fatal("a short pull closed or moved the sheet")
	}
	tt.Press(320, top+8)
	tt.Move(320, top+60)
	tt.Move(320, top+120)
	tt.Release(320, top+120)
	if open || closed != 1 {
		t.Fatalf("a long pull left open %v", open)
	}
}

func TestLightboxTurnsWrapsAndCloses(t *testing.T) {
	open, at, closed := true, 0, 0
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		if overlay.Lightbox(c, "Gallery", 3, &at, &open, overlay.LightboxOptions{}, func(i int) {
			pixel.Text(c, []string{"Castle", "Cave", "Coast"}[i], pixel.TextOptions{})
		}).Closed {
			closed++
		}
	})
	tt.Key(0, ui.KeyLeft)
	if at != 2 || !shows(tt, "Coast") || !shows(tt, "3 / 3") {
		t.Fatalf("Left from the first went to %d", at)
	}
	click(t, tt, "NEXT")
	tt.Key(0, ui.KeyRight)
	if at != 1 || !shows(tt, "Cave") {
		t.Fatalf("▶, Right went to %d", at)
	}
	keeps(t, tt.Image())
	click(t, tt, "CLOSE")
	if open || closed != 1 {
		t.Fatal("× did not close it")
	}
}

func TestDialogueBoxTypesThenTurns(t *testing.T) {
	clk, page, done := newClock(), 0, 0
	pages := []string{"It is dangerous to go alone!", "Take this."}
	tt := view(640, 400, clk, false, func(c *ui.Context) {
		if overlay.DialogueBox(c, pages, &page, overlay.DialogueBoxOptions{Name: "OLD MAN", Portrait: icon.User}).Done {
			done++
		}
	})
	r := find(t, tt, "OLD MAN "+pages[0])
	tt.ClickAt(r.X+4, r.Y+4)
	if page != 0 {
		t.Fatal("the first press turned the page before the text was out")
	}
	tt.Key(0, ui.KeyEnter)
	if page != 1 || !shows(tt, "OLD MAN Take this.") {
		t.Fatalf("Enter after the text left page %d", page)
	}
	clk.t = clk.t.Add(2 * time.Second)
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if done != 1 || !shows(tt, "OLD MAN") {
		t.Fatalf("Enter on the last page: done %d", done)
	}
	keeps(t, tt.Image())
}

func TestDialogueBoxUnderReduceMotionTurnsAtOnce(t *testing.T) {
	page := 0
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		overlay.DialogueBox(c, []string{"One.", "Two.", "Three."}, &page, overlay.DialogueBoxOptions{})
	})
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	click(t, tt, "One.")
	tt.Key(0, ui.KeyEnter)
	if page != 2 {
		t.Fatalf("two presses under Reduce Motion reached page %d", page)
	}
}

func TestNotificationsExpireAndDismiss(t *testing.T) {
	clk := newClock()
	queue := []overlay.Notification{{Title: "Saved", Text: "Slot 1."}, {Icon: icon.Mail, Title: "Letter", Text: "From mom.", For: 10 * time.Second}}
	var dismissed string
	tt := view(640, 400, clk, false, func(c *ui.Context) {
		if r := overlay.Notifications(c, &queue, overlay.NotificationsOptions{}); r.Dismissed != nil {
			dismissed = r.Dismissed.Title
		}
	})
	a, b := find(t, tt, "Saved"), find(t, tt, "Letter")
	if a.X+a.W != 640-16 || b.Y <= a.Y {
		t.Fatalf("cards at %v and %v", a, b)
	}
	clk.t = clk.t.Add(6 * time.Second)
	tt.Frame()
	if shows(tt, "Saved") || len(queue) != 1 {
		t.Fatalf("the first card outlived five seconds: %d left", len(queue))
	}
	click(t, tt, "CLOSE")
	if dismissed != "Letter" || len(queue) != 0 {
		t.Fatalf("× dismissed %q", dismissed)
	}
}

func TestKeyHintsDrawPadButtonsAndKeycaps(t *testing.T) {
	tt := view(640, 100, newClock(), false, func(c *ui.Context) {
		overlay.KeyHints(c, []overlay.KeyHint{{Key: "A", Action: "CONFIRM"}, {Key: "START", Action: "MENU"}}, overlay.KeyHintsOptions{})
	})
	r := find(t, tt, "A CONFIRM")
	if s := slot(t, tt.Image(), int(r.X)+18, int(r.Y)+r2i(r.H/2)); s != pixel.Ink {
		t.Fatalf("pad button is %v", s)
	}
	if !shows(tt, "START") || !shows(tt, "MENU") {
		t.Fatalf("texts %q", tt.Texts())
	}
	keeps(t, tt.Image())
}

func r2i(f float32) int { return int(f) }

func TestTitleCardDithersInHoldsAndGoes(t *testing.T) {
	clk, cue := newClock(), 0
	tt := view(320, 200, clk, false, func(c *ui.Context) {
		ui.Box(c).Size(300, 160).Children(func() {
			overlay.TitleCard(c, "WORLD 1-1", overlay.TitleCardOptions{Cue: cue, Hold: time.Second})
		})
	})
	area := ui.Rect{X: 0, Y: 0, W: 300, H: 160}
	if n := count(t, tt.Image(), area, pixel.Ink); n != 0 || shows(tt, "WORLD 1-1") {
		t.Fatal("the card showed before its cue")
	}
	cue = 1
	tt.Frame()
	var inks []int
	for range 4 {
		inks = append(inks, count(t, tt.Image(), area, pixel.Ink))
		clk.t = clk.t.Add(time.Second / 8)
		tt.Frame()
	}
	if !(inks[0] < inks[1] && inks[1] < inks[2] && inks[2] < inks[3]) || !shows(tt, "WORLD 1-1") {
		t.Fatalf("ink by step %v", inks)
	}
	clk.t = clk.t.Add(2 * time.Second)
	tt.Frame()
	if shows(tt, "WORLD 1-1") || count(t, tt.Image(), area, pixel.Ink) != 0 {
		t.Fatal("the card stayed past its hold")
	}
}

func TestOffDialogueBoxStaysPut(t *testing.T) {
	for _, inert := range []bool{false, true} {
		page := 0
		tt := view(640, 400, newClock(), inert, func(c *ui.Context) {
			overlay.DialogueBox(c, []string{"One.", "Two."}, &page, overlay.DialogueBoxOptions{Disabled: !inert})
		})
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		click(t, tt, "One.")
		click(t, tt, "One.")
		tt.Key(0, ui.KeyEnter)
		if page != 0 {
			t.Fatalf("inert %v: an off box turned to page %d", inert, page)
		}
	}
}

func TestSheetVanishesOnTheReleaseThatClosesIt(t *testing.T) {
	open := true
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		overlay.Sheet(c, "Options", &open, overlay.SheetOptions{Height: 120}, func() {
			pixel.Text(c, "Sound: ON", pixel.TextOptions{})
		})
	})
	top := float32(400 - 240)
	tt.Press(320, top+8)
	tt.Move(320, top+120)
	tt.Release(320, top+120)
	if open || shows(tt, "Sound: ON") {
		t.Fatal("the sheet stayed on screen after the pull closed it")
	}
}

func TestPanelsWithoutControlsTakeTheFocus(t *testing.T) {
	for _, kind := range []string{"drawer", "sheet", "dialog", "alert", "confirm", "choice", "pause", "prompt", "progress"} {
		open, opened := false, 0
		tt := view(640, 400, newClock(), false, func(c *ui.Context) {
			if input.Button(c, "OPEN", input.ButtonOptions{}).Clicked {
				open = true
				opened++
			}
			body := func() { pixel.Text(c, "Just words.", pixel.TextOptions{}) }
			name := "Menu"
			switch kind {
			case "drawer":
				overlay.Drawer(c, name, &open, overlay.DrawerOptions{Side: display.Left}, body)
			case "sheet":
				overlay.Sheet(c, name, &open, overlay.SheetOptions{}, body)
			case "dialog":
				overlay.Dialog(c, name, &open, overlay.DialogOptions{}, body)
			case "alert":
				overlay.AlertDialog(c, name, "Hi.", &open, overlay.AlertDialogOptions{Disabled: true})
			case "confirm":
				overlay.ConfirmDialog(c, "Hi.", &open, overlay.ConfirmDialogOptions{Disabled: true})
			case "choice":
				overlay.ChoicePrompt(c, name, []string{"A"}, &open, overlay.ChoicePromptOptions{Disabled: true})
			case "pause":
				overlay.PauseMenu(c, &open, overlay.PauseMenuOptions{Disabled: true})
			case "prompt":
				s := ""
				overlay.PromptDialog(c, name, &s, &open, overlay.PromptDialogOptions{Disabled: true})
			case "progress":
				overlay.ProgressDialog(c, name, 0.5, &open, overlay.ProgressDialogOptions{Disabled: true})
			}
		})
		click(t, tt, "OPEN")
		tt.Key(0, ui.KeyEnter)
		if opened != 1 || tt.Focused("OPEN") {
			t.Fatalf("%s: Enter reached the opener behind it (%d opens)", kind, opened)
		}
	}
}

func TestOffSheetTakesNoPull(t *testing.T) {
	for _, inert := range []bool{false, true} {
		open := true
		tt := view(640, 400, newClock(), inert, func(c *ui.Context) {
			overlay.Sheet(c, "Options", &open, overlay.SheetOptions{Height: 120, Disabled: !inert}, func() {})
		})
		top := float32(400 - 240)
		tt.Press(320, top+8)
		tt.Move(320, top+120)
		tt.Release(320, top+120)
		if !open {
			t.Fatalf("inert %v: a pull closed an off sheet", inert)
		}
	}
}

func TestDisabledPanelsTakeNoClicks(t *testing.T) {
	open, at, saving, next := true, 0, true, 0
	queue := []overlay.Notification{{Title: "Saved", Text: "Slot 1."}}
	tt := view(640, 400, newClock(), false, func(c *ui.Context) {
		overlay.Lightbox(c, "Gallery", 3, &at, &open, overlay.LightboxOptions{Disabled: true}, func(int) {})
	})
	tt.Key(0, ui.KeyRight)
	click(t, tt, "NEXT")
	click(t, tt, "CLOSE")
	tt = view(640, 400, newClock(), false, func(c *ui.Context) {
		overlay.ProgressDialog(c, "Saving", 0.2, &saving, overlay.ProgressDialogOptions{Disabled: true})
	})
	click(t, tt, "CANCEL")
	tt = view(640, 400, newClock(), false, func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() {
			b := input.Button(c, "BAG", input.ButtonOptions{})
			if overlay.Coachmark(c, b.Element, "Bag.", 1, 2, overlay.CoachmarkOptions{Disabled: true}).Next {
				next++
			}
		})
		overlay.Notifications(c, &queue, overlay.NotificationsOptions{Disabled: true})
	})
	click(t, tt, "NEXT ▶")
	click(t, tt, "CLOSE")
	if at != 0 || !open || !saving || next != 0 || len(queue) != 1 {
		t.Fatalf("disabled: at %d open %v saving %v next %d queue %d", at, open, saving, next, len(queue))
	}
}
