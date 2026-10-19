package feedback_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/feedback"
)

func TestCountdownsHoldWhenTheClockStepsBack(t *testing.T) {
	for name, tc := range map[string]struct {
		build func(c *ui.Context)
		start string
	}{
		"ContinuePrompt": {func(c *ui.Context) { feedback.ContinuePrompt(c, feedback.ContinuePromptOptions{}) }, "10"},
		"ErrorState":     {func(c *ui.Context) { feedback.ErrorState(c, "LOST", feedback.ErrorStateOptions{}) }, "9"},
	} {
		clk := newClock()
		tt := view(500, 400, clk, false, tc.build)
		clk.t = clk.t.Add(-time.Hour)
		tt.Frame()
		// The count, a text's last word, never passes its start.
		start, _ := strconv.Atoi(tc.start)
		for _, s := range tt.Texts() {
			f := strings.Fields(s)
			if len(f) == 0 {
				continue
			}
			if n, err := strconv.Atoi(f[len(f)-1]); err == nil && n > start {
				t.Errorf("%s counts %d after the clock stepped back", name, n)
			}
		}
	}
}

func TestToastsExpireAfterTheClockStepsBack(t *testing.T) {
	clk := newClock()
	queue := []feedback.Toast{{Message: "SAVED"}}
	tt := view(500, 300, clk, false, func(c *ui.Context) { feedback.Toasts(c, &queue, feedback.ToastsOptions{}) })
	clk.t = clk.t.Add(-time.Hour)
	tt.Frame()
	clk.t = clk.t.Add(4 * time.Second)
	tt.Frame()
	if len(queue) != 0 {
		t.Fatalf("a toast outlived its three seconds: %d left", len(queue))
	}
}
