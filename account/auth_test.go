package account_test

import (
	"bytes"
	"image"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/account"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// world is the settings a test view runs under.
type world struct {
	clk           *clock
	reduce, inert bool
}

func newWorld() *world {
	return &world{clk: &clock{time.Date(2026, 9, 2, 10, 30, 0, 0, time.UTC)}}
}

// view builds on a fake clock, two pixels a dot.
func view(w, h int, wd *world, build func(c *ui.Context)) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: wd.clk.now, Inert: wd.inert})
		build(c)
	}, w, h)
	if wd.reduce {
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tt.Frame()
	}
	return tt
}

// later moves the clock on and builds a frame.
func later(tt *ui.Tester, wd *world, d time.Duration) {
	wd.clk.t = wd.clk.t.Add(d)
	tt.Frame()
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

// none is a color outside the palette.
const none pixel.Slot = 255

// at is the slot at dot x, y, or none.
func at(img *image.RGBA, x, y int) pixel.Slot {
	c := img.RGBAAt(2*x, 2*y)
	for s, want := range pixel.SpringDay.Colors {
		if c.R == want.R && c.G == want.G && c.B == want.B {
			return pixel.Slot(s)
		}
	}
	return none
}

// count is how many dots in a box hold s.
func count(img *image.RGBA, x, y, w, h int, s pixel.Slot) int {
	n := 0
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			if at(img, i, j) == s {
				n++
			}
		}
	}
	return n
}

// dots is a found box in dots.
func dots(r ui.Rect) (x, y, w, h int) {
	return int(r.X / 2), int(r.Y / 2), int(r.W / 2), int(r.H / 2)
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

// logged runs fn with the debug log captured.
func logged(fn func()) string {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(old)
	fn()
	return buf.String()
}

func TestLoginFormSubmitsAndShakesOnFailure(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		name, pass, fails, sent := "", "", 0, 0
		tt := view(500, 300, wd, func(c *ui.Context) {
			if account.LoginForm(c, "WELCOME", &name, &pass, account.LoginFormOptions{Failures: fails}).Submitted {
				sent++
			}
		})
		click(t, tt, "LOGIN")
		if sent != 0 {
			t.Fatal("LOGIN sent an empty form")
		}
		log := logged(func() {
			click(t, tt, "NAME")
			tt.Type("player1")
			click(t, tt, "PASSWORD")
			tt.Type("hunter2")
			tt.Key(0, ui.KeyEnter)
		})
		if sent != 1 || strings.Contains(log, "hunter2") || strings.Contains(log, "player1") {
			t.Fatalf("Enter sent %d; the log must not hold credentials:\n%s", sent, log)
		}
		x0, _, _, _ := dots(find(t, tt, "LOGIN"))
		// The caller says no: the password clears, focus comes back.
		click(t, tt, "NAME")
		fails++
		tt.Frame()
		if pass != "" || name != "player1" || !tt.HasText("WRONG NAME OR PASSWORD") || !tt.Focused("PASSWORD") {
			t.Fatalf("after a failure: %q %q, %q", name, pass, tt.Texts())
		}
		x1, _, _, _ := dots(find(t, tt, "LOGIN"))
		if moved := x1 != x0; moved == reduce {
			t.Fatalf("reduce %v: the form moved %v", reduce, moved)
		}
		later(tt, wd, time.Second)
		if x2, _, _, _ := dots(find(t, tt, "LOGIN")); x2 != x0 {
			t.Fatalf("the shake did not settle: %d, was %d", x2, x0)
		}
		tt.Type("x")
		if tt.HasText("WRONG NAME OR PASSWORD") {
			t.Fatal("typing kept the warning up")
		}
	}
}

func TestPasswordGridSpellsAndEnds(t *testing.T) {
	wd := newWorld()
	pass, ended := "", false
	tt := view(500, 300, wd, func(c *ui.Context) {
		ended = account.PasswordGrid(c, "Grid", &pass, account.PasswordGridOptions{Length: 4}).Submitted || ended
	})
	r := find(t, tt, "SYMBOLS")
	cell := func(i int) (float32, float32) { return r.X + float32(22*(i%8)+11), r.Y + float32(22*(i/8)+11) }
	log := logged(func() {
		for _, i := range []int{9, 20, 18, 19} {
			tt.ClickAt(cell(i))
		}
	})
	if pass != "JUST" || strings.Contains(log, "JUST") {
		t.Fatalf("spelled %q; log:\n%s", pass, log)
	}
	tt.ClickAt(cell(0))
	if pass != "JUST" {
		t.Fatal("a full password took another symbol")
	}
	tt.Key(0, ui.KeyBackspace)
	click(t, tt, "BACK")
	if pass != "JU" {
		t.Fatalf("Backspace and BACK left %q", pass)
	}
	click(t, tt, "END")
	if ended {
		t.Fatal("END worked on a short password")
	}
	tt.ClickAt(cell(18))
	tt.ClickAt(cell(19))
	click(t, tt, "END")
	if !ended {
		t.Fatal("END on a full password did not end")
	}
	for name, p := range map[string]string{"lowercase": "ab", "too long": "ABCDE"} {
		mustPanic(t, name, func() {
			view(500, 300, newWorld(), func(c *ui.Context) {
				account.PasswordGrid(c, "G", &p, account.PasswordGridOptions{Length: 4})
			})
		})
	}
}

func TestTwoFactorTakesDigitsAndPastes(t *testing.T) {
	wd := newWorld()
	code, sent := "", 0
	tt := view(500, 300, wd, func(c *ui.Context) {
		if account.TwoFactor(c, "Code", &code, account.TwoFactorOptions{}).Submitted {
			sent++
		}
	})
	click(t, tt, "Code")
	log := logged(func() { tt.Type("42a7") })
	if code != "427" || strings.Contains(log, "427") {
		t.Fatalf("typed %q; log:\n%s", code, log)
	}
	tt.Key(0, ui.KeyBackspace)
	tt.SetClipboard("your code: 1-3-5-2-4-6, enjoy")
	tt.Key(ui.Cmd, ui.KeyV)
	// Paste fills the four boxes left, from 1352.
	if code != "421352" || sent != 1 {
		t.Fatalf("paste gave %q, sent %d", code, sent)
	}
	tt.Type("9")
	if code != "421352" || sent != 1 {
		t.Fatalf("a seventh digit went in: %q, %d", code, sent)
	}
	// Three seconds left in the period: the bar is Bad.
	wd.clk.t = time.Unix(60*1000+27, 0)
	tt.Frame()
	x, y, w, _ := dots(find(t, tt, "TIME LEFT"))
	if count(tt.Image(), x, y+2, w, 2, pixel.Bad) == 0 {
		t.Fatal("the bar is not Bad with three seconds left")
	}
	mustPanic(t, "letters in the code", func() {
		view(500, 300, newWorld(), func(c *ui.Context) {
			s := "12a"
			account.TwoFactor(c, "C", &s, account.TwoFactorOptions{})
		})
	})
}

func TestSessionListSignsOut(t *testing.T) {
	wd := newWorld()
	out := -1
	tt := view(600, 300, wd, func(c *ui.Context) {
		if r := account.SessionList(c, "Devices", []account.Session{
			{Device: "Laptop", Seen: wd.clk.t, Current: true}, {Device: "Phone", Seen: wd.clk.t.Add(-2 * time.Hour)},
		}, account.SessionListOptions{}); r.SignedOut >= 0 {
			out = r.SignedOut
		}
	})
	find(t, tt, "THIS DEVICE")
	if !tt.HasText("2 H AGO") {
		t.Fatalf("no last-seen time in %q", tt.Texts())
	}
	click(t, tt, "SIGN OUT Phone")
	if out != 1 {
		t.Fatalf("SIGN OUT Phone signed out %d", out)
	}
}

func TestFeedbackFormSendsOnceRated(t *testing.T) {
	wd := newWorld()
	rating, words, sent := 0, "", false
	tt := view(500, 300, wd, func(c *ui.Context) {
		sent = account.FeedbackForm(c, "Feedback", &rating, &words, account.FeedbackFormOptions{}).Sent || sent
	})
	click(t, tt, "MESSAGE")
	tt.Type("more slimes")
	click(t, tt, "SEND ▶")
	if sent {
		t.Fatal("SEND worked without a rating")
	}
	click(t, tt, "4/5")
	log := logged(func() { click(t, tt, "SEND ▶") })
	if !sent || rating == 0 || strings.Contains(log, "slimes") {
		t.Fatalf("sent %v at %d stars; log:\n%s", sent, rating, log)
	}
}
