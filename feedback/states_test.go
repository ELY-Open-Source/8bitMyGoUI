package feedback_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestAlertAndBannerReportClicks(t *testing.T) {
	closed, action, dismissed := 0, "", 0
	tt := view(500, 300, newClock(), false, func(c *ui.Context) {
		if feedback.Alert(c, "Low HP", "Drink a potion.", feedback.AlertOptions{Kind: feedback.Danger, Closable: true}).Closed {
			closed++
		}
		r := feedback.Banner(c, "New quest!", feedback.BannerOptions{Kind: feedback.Info, Actions: []string{"VIEW"}, Dismissable: true})
		if r.Action != "" {
			action = r.Action
		}
		if r.Dismissed {
			dismissed++
		}
	})
	if count(slots(t, tt.Image(), find(t, tt, "Low HP")), pixel.Bad) == 0 {
		t.Fatal("a danger alert shows no Bad")
	}
	click(t, tt, "CLOSE")
	click(t, tt, "VIEW")
	if closed != 1 || action != "VIEW" {
		t.Fatalf("closed %d, action %q", closed, action)
	}
	tt = view(500, 100, newClock(), false, func(c *ui.Context) {
		if feedback.Banner(c, "New quest!", feedback.BannerOptions{Dismissable: true}).Dismissed {
			dismissed++
		}
	})
	click(t, tt, "CLOSE")
	if dismissed != 1 {
		t.Fatal("the banner's × did not dismiss it")
	}
}

func TestEmptyStateAction(t *testing.T) {
	acted := 0
	tt := view(400, 300, newClock(), false, func(c *ui.Context) {
		if feedback.EmptyState(c, "Your bag is empty.", feedback.EmptyStateOptions{Action: "SHOP"}).Acted {
			acted++
		}
	})
	find(t, tt, "NO ITEMS")
	click(t, tt, "SHOP")
	if acted != 1 {
		t.Fatal("the action did not report")
	}
}

func TestErrorStateCountsDownAndContinues(t *testing.T) {
	clk := newClock()
	cont := 0
	tt := view(600, 500, clk, true, func(c *ui.Context) {
		if feedback.ErrorState(c, "The slime won.", feedback.ErrorStateOptions{}).Continue {
			cont++
		}
	})
	find(t, tt, "CONTINUE? 9")
	wait(tt, clk, 3*time.Second)
	find(t, tt, "CONTINUE? 6")
	wait(tt, clk, 10*time.Second)
	find(t, tt, "CONTINUE? 0")
	click(t, tt, "CONTINUE")
	if cont != 1 {
		t.Fatal("CONTINUE did not report Continue")
	}
}

func TestSuccessStateStarsPopInTurn(t *testing.T) {
	stars := func(reduce bool, d time.Duration) int {
		clk := newClock()
		tt := view(600, 400, clk, reduce, func(c *ui.Context) {
			feedback.SuccessState(c, "All coins found.", feedback.SuccessStateOptions{Stars: 3})
		})
		wait(tt, clk, d)
		return count(slots(t, tt.Image(), find(t, tt, "STAGE CLEAR!")), pixel.Warn)
	}
	early, late := stars(false, 300*time.Millisecond), stars(false, 2*time.Second)
	if early == 0 || early >= late {
		t.Fatalf("stars did not pop in turn: %d then %d", early, late)
	}
	if stars(true, 0) != late {
		t.Fatal("Reduce Motion did not show every star at once")
	}
}

func TestInlineMessageMarks(t *testing.T) {
	tt := view(400, 100, newClock(), false, func(c *ui.Context) {
		feedback.InlineMessage(c, "Name taken", feedback.InlineMessageOptions{Kind: feedback.Danger})
		feedback.InlineMessage(c, "Looks good", feedback.InlineMessageOptions{Kind: feedback.Success})
	})
	bad, good := slots(t, tt.Image(), find(t, tt, "Name taken")), slots(t, tt.Image(), find(t, tt, "Looks good"))
	if count(bad, pixel.Bad) == 0 || count(good, pixel.Good) == 0 {
		t.Fatal("marks do not take their kind's slot")
	}
	// The mark is the first seven dots across.
	mark := func(l string, ink pixel.Slot) (out []bool) {
		r := find(t, tt, l)
		for _, v := range slots(t, tt.Image(), ui.Rect{X: r.X, Y: r.Y, W: 14, H: r.H}) {
			out = append(out, v == ink)
		}
		return out
	}
	if equal(mark("Name taken", pixel.Bad), mark("Looks good", pixel.Good)) {
		t.Fatal("✓ and ! look alike")
	}
}

func TestToastsQueueExpireAndAct(t *testing.T) {
	clk := newClock()
	q := []feedback.Toast{{Message: "one"}, {Message: "two"}, {Message: "three", For: 5 * time.Second}, {Message: "four", Action: "UNDO"}}
	var acted *feedback.Toast
	tt := view(500, 400, clk, false, func(c *ui.Context) {
		if r := feedback.Toasts(c, &q, feedback.ToastsOptions{}); r.Acted != nil {
			acted = r.Acted
		}
	})
	if _, ok := tt.Find("four"); ok || len(q) != 4 {
		t.Fatal("more than three toasts showed")
	}
	wait(tt, clk, 3*time.Second)
	find(t, tt, "three")
	find(t, tt, "four")
	if _, ok := tt.Find("one"); ok || len(q) != 2 {
		t.Fatalf("toasts did not leave on time: %d left", len(q))
	}
	click(t, tt, "UNDO")
	if acted == nil || acted.Message != "four" || len(q) != 1 {
		t.Fatal("the action did not report or clear its toast")
	}
	wait(tt, clk, 3*time.Second)
	if len(q) != 0 {
		t.Fatal("a five-second toast outlived its time")
	}
}

func TestContinuePromptAnswersOrTimesOut(t *testing.T) {
	clk := newClock()
	var yes, timedOut int
	tt := view(400, 400, clk, false, func(c *ui.Context) {
		r := feedback.ContinuePrompt(c, feedback.ContinuePromptOptions{Seconds: 5})
		if r.Yes {
			yes++
		}
		if r.TimedOut {
			timedOut++
		}
	})
	find(t, tt, "5")
	wait(tt, clk, 2*time.Second)
	find(t, tt, "3")
	wait(tt, clk, 4*time.Second)
	wait(tt, clk, time.Second)
	click(t, tt, "YES")
	if timedOut != 1 || yes != 0 {
		t.Fatalf("timed out %d times; YES after time counted %d", timedOut, yes)
	}
	clk2 := newClock()
	tt = view(400, 400, clk2, false, func(c *ui.Context) {
		if feedback.ContinuePrompt(c, feedback.ContinuePromptOptions{}).Yes {
			yes++
		}
	})
	click(t, tt, "YES")
	wait(tt, clk2, 20*time.Second)
	if yes != 1 {
		t.Fatal("YES did not report")
	}
}

func TestTwinToastsBothShow(t *testing.T) {
	q := []feedback.Toast{{Message: "same"}, {Message: "same"}}
	tt := view(500, 300, newClock(), false, func(c *ui.Context) {
		feedback.Toasts(c, &q, feedback.ToastsOptions{})
	})
	n := 0
	for _, s := range tt.Texts() {
		if s == "same" {
			n++
		}
	}
	// Each toast and its text both carry the message.
	if n != 4 {
		t.Fatalf("%d labels for two twin toasts", n)
	}
}

func TestDisabledFeedbackButtonsDoNothing(t *testing.T) {
	events := 0
	q := []feedback.Toast{{Message: "dropped", Action: "UNDO"}}
	tt := view(800, 2000, newClock(), false, func(c *ui.Context) {
		on := func(b bool) {
			if b {
				events++
			}
		}
		on(feedback.Alert(c, "A", "a", feedback.AlertOptions{Closable: true, Disabled: true}).Closed)
		r := feedback.Banner(c, "B", feedback.BannerOptions{Actions: []string{"GO"}, Disabled: true})
		on(r.Action != "")
		on(feedback.EmptyState(c, "E", feedback.EmptyStateOptions{Action: "SHOP", Disabled: true}).Acted)
		on(feedback.ErrorState(c, "X", feedback.ErrorStateOptions{Disabled: true}).Continue)
		on(feedback.ContinuePrompt(c, feedback.ContinuePromptOptions{Disabled: true}).Yes)
		on(feedback.Toasts(c, &q, feedback.ToastsOptions{Disabled: true}).Acted != nil)
	})
	for _, l := range []string{"CLOSE", "GO", "SHOP", "CONTINUE", "YES", "UNDO"} {
		click(t, tt, l)
	}
	if events != 0 {
		t.Fatalf("disabled buttons reported %d events", events)
	}
}
