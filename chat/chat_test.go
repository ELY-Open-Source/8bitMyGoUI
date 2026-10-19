package chat_test

import (
	"bytes"
	"image"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/chat"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// newClock starts on Wednesday 2 September 2026, mid-morning.
func newClock() *clock { return &clock{time.Date(2026, 9, 2, 10, 30, 0, 0, time.UTC)} }

// world is the settings a test view runs under.
type world struct {
	clk               *clock
	reduce, zh, inert bool
}

// view builds on a fake clock, two pixels a dot.
func view(w, h int, wd *world, build func(c *ui.Context)) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		s := pixel.Settings{Palette: pixel.SpringDay, Now: wd.clk.now, Inert: wd.inert}
		if wd.zh {
			s.Locale = pixel.Zh
		}
		pixel.Use(c, s)
		build(c)
	}, w, h)
	if wd.reduce {
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tt.Frame()
	}
	return tt
}

func click(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatalf("%v in %q", err, tt.Texts())
	}
}

func find(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q in %q", label, tt.Texts())
	}
	return r
}

// count is how many dots of a box hold s.
func count(img *image.RGBA, x, y, w, h int, s pixel.Slot) int {
	n, want := 0, pixel.SpringDay.Colors[s]
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			if c := img.RGBAAt(2*i, 2*j); c.R == want.R && c.G == want.G && c.B == want.B {
				n++
			}
		}
	}
	return n
}

func mustPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s did not panic", name)
		}
	}()
	fn()
}

func TestSpeechBubbleSidesNameAndTime(t *testing.T) {
	wd := &world{clk: newClock()}
	at := wd.clk.t
	tt := view(300, 200, wd, func(c *ui.Context) {
		ui.Column(c).Width(pixel.K(c).D(140)).Children(func() {
			chat.SpeechBubble(c, chat.Message{From: "Aria", Text: "Hi there", At: at}, chat.SpeechBubbleOptions{Width: 80})
			chat.SpeechBubble(c, chat.Message{From: "You", Text: "Hello", At: at, Mine: true}, chat.SpeechBubbleOptions{Width: 80})
		})
	})
	theirs, mine := find(t, tt, "Hi there"), find(t, tt, "Hello")
	if theirs.X > 40 || mine.X < 140 {
		t.Fatalf("theirs at %v, mine at %v: theirs sit left, yours right", theirs.X, mine.X)
	}
	find(t, tt, "Aria")
	if !tt.HasText("10:30") {
		t.Fatalf("no time in %q", tt.Texts())
	}
	// Yours fill Accent; theirs Lit.
	img := tt.Image()
	if count(img, int(mine.X/2), int(mine.Y/2), 20, 6, pixel.Accent) == 0 || count(img, int(theirs.X/2), int(theirs.Y/2), 20, 6, pixel.Accent) != 0 {
		t.Fatal("only your bubble is Accent")
	}
}

func TestChatLogFollowsItsEndOnlyWhileThere(t *testing.T) {
	wd := &world{clk: newClock()}
	n := 8
	tt := view(400, 300, wd, func(c *ui.Context) {
		chat.ChatLog(c, "Log", chat.ChatLogOptions{Width: 160, Height: 80}, func() {
			for i := range n {
				pixel.Text(c, "line "+strings.Repeat("I", i+1), pixel.TextOptions{})
			}
		})
	})
	// Find clips to the view; whole lines are 32 pixels.
	shown := func(label string) bool {
		r, l := find(t, tt, label), find(t, tt, "Log")
		return r.H == 32 && r.Y >= l.Y && r.Y+r.H <= l.Y+l.H
	}
	if !shown("line IIIIIIII") {
		t.Fatal("the log did not open at its end")
	}
	n++
	tt.Frame()
	if !shown("line IIIIIIIII") {
		t.Fatal("at the end, the log did not follow a new line")
	}
	// Scrolled up, a new line leaves the view alone.
	for range 4 {
		click(t, tt, "SCROLL UP")
	}
	n++
	tt.Frame()
	if shown("line IIIIIIIIII") {
		t.Fatal("scrolled up, the log jumped to a new line")
	}
}

func TestChatLogKeepsFollowingThroughInert(t *testing.T) {
	wd := &world{clk: newClock()}
	n := 8
	tt := view(400, 300, wd, func(c *ui.Context) {
		chat.ChatLog(c, "Log", chat.ChatLogOptions{Width: 160, Height: 80}, func() {
			for i := range n {
				pixel.Text(c, "line "+strings.Repeat("I", i+1), pixel.TextOptions{})
			}
		})
	})
	shown := func(label string) bool {
		r, l := find(t, tt, label), find(t, tt, "Log")
		return r.H == 32 && r.Y >= l.Y && r.Y+r.H <= l.Y+l.H
	}
	// A line landing while inert shows once the log wakes.
	wd.inert = true
	n++
	tt.Frame()
	tt.Frame()
	wd.inert = false
	tt.Frame()
	if !shown("line IIIIIIIII") {
		t.Fatal("after Inert the log did not catch up with its end")
	}
	n++
	tt.Frame()
	if !shown("line IIIIIIIIII") {
		t.Fatal("after Inert the log stopped following")
	}
}

func TestFieldsLogChangesWithoutText(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(old)
	draft, mention := "", ""
	tt := view(400, 300, &world{clk: newClock()}, func(c *ui.Context) {
		chat.Composer(c, "Message", &draft, chat.ComposerOptions{})
		chat.MentionInput(c, "Mention", &mention, []string{"Aria"}, chat.MentionInputOptions{})
	})
	click(t, tt, "Message")
	tt.Type("hi")
	click(t, tt, "Mention")
	tt.Type("yo")
	out := buf.String()
	for _, want := range []string{"composer changed", "mention input changed", "runes=2"} {
		if !strings.Contains(out, want) {
			t.Errorf("no %q in the logs:\n%s", want, out)
		}
	}
	if strings.Contains(out, "hi") || strings.Contains(out, "yo") {
		t.Error("the logs carry what was typed")
	}
}

func TestComposerSendsOnEnterOrSend(t *testing.T) {
	wd := &world{clk: newClock()}
	draft, sent := "", []string{}
	tt := view(400, 200, wd, func(c *ui.Context) {
		if r := chat.Composer(c, "Message", &draft, chat.ComposerOptions{}); r.Sent != "" {
			sent = append(sent, r.Sent)
		}
	})
	// Empty, SEND ▶ does nothing.
	click(t, tt, "SEND ▶")
	click(t, tt, "Message")
	tt.Type("hi")
	tt.Key(ui.Shift, ui.KeyEnter)
	tt.Type("all")
	tt.Key(0, ui.KeyEnter)
	if len(sent) != 1 || sent[0] != "hi\nall" || draft != "" {
		t.Fatalf("sent %q, draft %q", sent, draft)
	}
	tt.Type("bye")
	click(t, tt, "SEND ▶")
	if len(sent) != 2 || sent[1] != "bye" || draft != "" {
		t.Fatalf("sent %q, draft %q", sent, draft)
	}
	// The field keeps the focus; next words land there.
	tt.Type("again")
	if draft != "again" {
		t.Fatalf("after sending, typing gave %q", draft)
	}
}

func TestTypingIndicatorHopsUnlessReduced(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		tt := view(200, 100, wd, func(c *ui.Context) { chat.TypingIndicator(c, "Aria", chat.TypingIndicatorOptions{}) })
		r := find(t, tt, "Aria IS TYPING")
		before := count(tt.Image(), int(r.X/2), int(r.Y/2), 17, 3, pixel.Faint)
		moved := false
		for range 4 {
			wd.clk.t = wd.clk.t.Add(time.Second / 6)
			tt.Frame()
			moved = moved || count(tt.Image(), int(r.X/2), int(r.Y/2), 17, 3, pixel.Faint) != before
		}
		if moved == reduce {
			t.Fatalf("reduce %v: dots hopped %v", reduce, moved)
		}
	}
}

func TestDateSeparatorNamesTheDay(t *testing.T) {
	for _, tc := range []struct {
		zh   bool
		days int
		want string
	}{{false, 0, "TODAY"}, {false, -1, "YESTERDAY"}, {false, -7, "WED 26 AUG"}, {true, -7, "8月26日 周三"}, {true, 0, "今天"}} {
		wd := &world{clk: newClock(), zh: tc.zh}
		day := wd.clk.t.AddDate(0, 0, tc.days)
		tt := view(400, 100, wd, func(c *ui.Context) { chat.DateSeparator(c, day, chat.DateSeparatorOptions{}) })
		find(t, tt, tc.want)
	}
	// At midnight today turns into yesterday by itself.
	wd := &world{clk: newClock()}
	wd.clk.t = time.Date(2026, 9, 2, 23, 59, 59, 0, time.UTC)
	day := wd.clk.t
	pixel.ResetDue()
	tt := view(400, 100, wd, func(c *ui.Context) { chat.DateSeparator(c, day, chat.DateSeparatorOptions{}) })
	if until := pixel.Due().Sub(tt.Now()); until <= 0 || until > 2*time.Second {
		t.Fatalf("the next build is due in %v, not at midnight", until)
	}
	wd.clk.t = wd.clk.t.Add(2 * time.Second)
	tt.Frame()
	find(t, tt, "YESTERDAY")
	mustPanic(t, "narrow", func() {
		view(100, 100, wd, func(c *ui.Context) { chat.UnreadDivider(c, chat.UnreadDividerOptions{Width: 30}) })
	})
}

func TestDividerAndSystemLines(t *testing.T) {
	wd := &world{clk: newClock()}
	tt := view(500, 100, wd, func(c *ui.Context) {
		chat.UnreadDivider(c, chat.UnreadDividerOptions{})
		chat.SystemMessage(c, "PLAYER 2 HAS JOINED", chat.SystemMessageOptions{})
	})
	r := find(t, tt, "NEW MESSAGES")
	if count(tt.Image(), int(r.X/2), int(r.Y/2), 200, 16, pixel.Bad) < 100 {
		t.Fatal("the unread rule is not Bad")
	}
	find(t, tt, "PLAYER 2 HAS JOINED")
}
