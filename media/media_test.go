package media_test

import (
	"image"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/media"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// world is the settings a test view runs under.
type world struct {
	clk    *clock
	reduce bool
	pal    *pixel.Palette
}

func newWorld() *world {
	return &world{clk: &clock{time.Date(2026, 9, 2, 10, 30, 0, 0, time.UTC)}, pal: pixel.SpringDay}
}

// view builds on a fake clock, two pixels a dot.
func view(w, h int, wd *world, build func(c *ui.Context)) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: wd.pal, Now: wd.clk.now})
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
func at(img *image.RGBA, pal *pixel.Palette, x, y int) pixel.Slot {
	c := img.RGBAAt(2*x, 2*y)
	for s, want := range pal.Colors {
		if c.R == want.R && c.G == want.G && c.B == want.B {
			return pixel.Slot(s)
		}
	}
	return none
}

// count is how many dots in a box hold s.
func count(img *image.RGBA, pal *pixel.Palette, x, y, w, h int, s pixel.Slot) int {
	n := 0
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			if at(img, pal, i, j) == s {
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

func TestTransportPlaysStopsAndSkips(t *testing.T) {
	wd := newWorld()
	playing := false
	var got media.TransportResult
	tt := view(400, 100, wd, func(c *ui.Context) {
		if r := media.Transport(c, "Deck", &playing, media.TransportOptions{}); r.Changed || r.Stopped || r.Previous || r.Next {
			got = r
		}
	})
	// ▶ has 50 ink dots; ❚❚ has 60.
	ink := func() int {
		x, y, w, h := dots(find(t, tt, "PLAY"))
		return count(tt.Image(), wd.pal, x, y, w, h, pixel.Ink)
	}
	if n := ink(); n != 50 {
		t.Fatalf("stopped, PLAY shows %d ink dots, want ▶'s 50", n)
	}
	click(t, tt, "PLAY")
	if !playing || !got.Changed || ink() != 60 {
		t.Fatalf("PLAY: playing %v, %+v, %d ink dots", playing, got, ink())
	}
	got = media.TransportResult{}
	click(t, tt, "STOP")
	if playing || !got.Stopped || !got.Changed {
		t.Fatalf("STOP: playing %v, %+v", playing, got)
	}
	got = media.TransportResult{}
	click(t, tt, "STOP")
	if !got.Stopped || got.Changed {
		t.Fatalf("STOP while stopped: %+v; nothing changed", got)
	}
	for _, b := range []string{"PREVIOUS", "NEXT"} {
		got = media.TransportResult{}
		click(t, tt, b)
		if got.Previous != (b == "PREVIOUS") || got.Next != (b == "NEXT") || got.Changed {
			t.Fatalf("%s: %+v", b, got)
		}
	}
}

func TestScrubberSeeksByPressAndKeys(t *testing.T) {
	wd := newWorld()
	pos, total := 72*time.Second, 192*time.Second
	changed := false
	tt := view(400, 100, wd, func(c *ui.Context) {
		changed = media.Scrubber(c, "Seek", &pos, total, media.ScrubberOptions{Width: 96}).Changed || changed
	})
	find(t, tt, "1:12")
	find(t, tt, "3:12")
	r := find(t, tt, "Seek")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	if !changed || pos != total/2 || !tt.HasText("1:36") {
		t.Fatalf("pressing the middle: changed %v, at %v; want %v", changed, pos, total/2)
	}
	tt.Key(0, ui.KeyHome)
	if pos != 0 || !tt.HasText("0:00") {
		t.Fatalf("Home: at %v, %q", pos, tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	if pos != 9600*time.Millisecond {
		t.Fatalf("Right steps a twentieth of 3:12, 9.6 s: %v", pos)
	}
	mustPanic(t, "a place past the end", func() {
		view(400, 100, newWorld(), func(c *ui.Context) {
			p := total + time.Second
			media.Scrubber(c, "Seek", &p, total, media.ScrubberOptions{})
		})
	})
}

func TestVolumeStepsMutesAndLights(t *testing.T) {
	wd := newWorld()
	level, muted := 5, false
	changed := false
	tt := view(400, 100, wd, func(c *ui.Context) {
		changed = media.Volume(c, "Volume", &level, &muted, media.VolumeOptions{}).Changed || changed
	})
	bars := func() []pixel.Slot {
		img := tt.Image()
		x, y, _, h := dots(find(t, tt, "Volume"))
		var out []pixel.Slot
		for i := range 8 {
			out = append(out, at(img, wd.pal, x+4*i+1, y+h-2))
		}
		return out
	}
	want := []pixel.Slot{pixel.Good, pixel.Good, pixel.Good, pixel.Good, pixel.Good, pixel.Well, pixel.Well, pixel.Well}
	if got := bars(); !slices.Equal(got, want) {
		t.Fatalf("level 5 lights %v, want %v", got, want)
	}
	r := find(t, tt, "Volume")
	tt.Press(r.X+r.W-1, r.Y+r.H/2)
	tt.Release(r.X+r.W-1, r.Y+r.H/2)
	if level != 8 || !changed {
		t.Fatalf("pressing the top bar: level %d, changed %v", level, changed)
	}
	if got := bars(); got[5] != pixel.Good || got[6] != pixel.Warn || got[7] != pixel.Warn {
		t.Fatalf("the loud quarter is Warn: %v", got)
	}
	tt.Key(0, ui.KeyLeft)
	if level != 7 {
		t.Fatalf("Left: level %d, want 7", level)
	}
	changed = false
	click(t, tt, "MUTE")
	if !muted || !changed || bars()[0] != pixel.Faint {
		t.Fatalf("MUTE: muted %v, changed %v, bars %v", muted, changed, bars())
	}
	mustPanic(t, "level above the bars", func() {
		view(400, 100, newWorld(), func(c *ui.Context) {
			l, m := 9, false
			media.Volume(c, "Volume", &l, &m, media.VolumeOptions{})
		})
	})
}

func TestEqualizerReportsTheBand(t *testing.T) {
	wd := newWorld()
	gains := []float64{0, 0, 0}
	band := -1
	tt := view(400, 200, wd, func(c *ui.Context) {
		if b := media.Equalizer(c, "EQ", []string{"BASS", "MID", "HIGH"}, gains, media.EqualizerOptions{}).Band; b >= 0 {
			band = b
		}
	})
	click(t, tt, "MID")
	tt.Key(0, ui.KeyUp)
	if band != 1 || gains[1] <= 1 || gains[0] != 0 || gains[2] != 0 {
		t.Fatalf("Up on MID: band %d, gains %v", band, gains)
	}
	mustPanic(t, "gain past 12 dB", func() {
		view(400, 200, newWorld(), func(c *ui.Context) {
			media.Equalizer(c, "EQ", []string{"A"}, []float64{13}, media.EqualizerOptions{})
		})
	})
	mustPanic(t, "bands without gains", func() {
		view(400, 200, newWorld(), func(c *ui.Context) {
			media.Equalizer(c, "EQ", []string{"A", "B"}, []float64{0}, media.EqualizerOptions{})
		})
	})
}

var songs = []media.Track{
	{Title: "One", Artist: "Ann", Length: 61 * time.Second},
	{Title: "Two", Artist: "Bob", Length: 3723 * time.Second},
	{Title: "Six", Artist: "Cy", Length: 9 * time.Second},
}

func TestPlaylistMarksPlayingAndSubmits(t *testing.T) {
	wd := newWorld()
	sel := 0
	var got media.PlaylistResult
	tt := view(500, 300, wd, func(c *ui.Context) {
		if r := media.Playlist(c, "List", songs, 2, &sel, media.PlaylistOptions{}); r.Changed || r.Submitted {
			got = r
		}
	})
	find(t, tt, "Six, Cy, PLAYING")
	click(t, tt, "Two, Bob")
	if sel != 1 || !got.Changed {
		t.Fatalf("clicking Two: selected %d, %+v", sel, got)
	}
	got = media.PlaylistResult{}
	tt.Key(0, ui.KeyEnter)
	if !got.Submitted {
		t.Fatalf("Enter did not submit: %+v", got)
	}
	// Past the ♪, the playing row's title is Accent too.
	x, y, _, _ := dots(find(t, tt, "Six, Cy, PLAYING"))
	if count(tt.Image(), wd.pal, x+27, y, 30, 12, pixel.Accent) == 0 {
		t.Fatal("the playing title is not Accent")
	}
	mustPanic(t, "playing past the end", func() {
		view(500, 300, newWorld(), func(c *ui.Context) {
			s := 0
			media.Playlist(c, "List", songs, 3, &s, media.PlaylistOptions{})
		})
	})
}

func TestAudioPlayerStopRewinds(t *testing.T) {
	wd := newWorld()
	pl := media.Player{At: 30 * time.Second, Level: 4}
	var got media.AudioPlayerResult
	tt := view(600, 300, wd, func(c *ui.Context) {
		if r := media.AudioPlayer(c, "Deck", songs[1], &pl, media.AudioPlayerOptions{}); r.Changed || r.Previous || r.Next {
			got = r
		}
	})
	find(t, tt, "Two, Bob")
	find(t, tt, "1:02:03")
	click(t, tt, "PLAY")
	if !pl.Playing || !got.Changed {
		t.Fatalf("PLAY: %+v, %+v", pl, got)
	}
	got = media.AudioPlayerResult{}
	click(t, tt, "STOP")
	if pl.Playing || pl.At != 0 || !got.Stopped || !got.Changed {
		t.Fatalf("STOP did not end play and rewind: %+v, %+v", pl, got)
	}
	got = media.AudioPlayerResult{}
	click(t, tt, "NEXT")
	if !got.Next || got.Changed {
		t.Fatalf("NEXT: %+v", got)
	}
	got = media.AudioPlayerResult{}
	click(t, tt, "MUTE")
	if !pl.Muted || !got.Changed {
		t.Fatalf("MUTE: %+v, %+v", pl, got)
	}
}

func TestVideoFrameShowsTheScreenAndSeeks(t *testing.T) {
	wd := newWorld()
	playing, pos := false, time.Duration(0)
	changed := false
	tt := view(600, 400, wd, func(c *ui.Context) {
		screen := func() { ui.Box(c).Size(400, 400).Label("film").Role(ui.RoleImage) }
		changed = media.VideoFrame(c, "Clip", &playing, &pos, 20*time.Second, media.VideoFrameOptions{}, screen).Changed || changed
	})
	// The screen, 160 by 90 dots, clips the film.
	if f := find(t, tt, "film"); f.H != 180 {
		t.Fatalf("the screen shows %v of the film, want 180 pixels tall", f)
	}
	click(t, tt, "PLAY")
	if !playing || !changed {
		t.Fatalf("PLAY: playing %v, changed %v", playing, changed)
	}
	changed = false
	r := find(t, tt, "SEEK")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	if pos != 10*time.Second || !changed {
		t.Fatalf("seeking the middle: at %v, changed %v", pos, changed)
	}
}
