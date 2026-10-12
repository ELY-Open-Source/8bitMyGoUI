package chat_test

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/chat"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestReactionsToggleYours(t *testing.T) {
	wd := &world{clk: newClock()}
	rs := []chat.Reaction{{Emoji: "love", Count: 3}, {Emoji: "fire", Count: 1, Mine: true}}
	toggled := -1
	tt := view(300, 100, wd, func(c *ui.Context) {
		if r := chat.Reactions(c, "Reactions", rs, chat.ReactionsOptions{}); r.Toggled >= 0 {
			toggled = r.Toggled
		}
	})
	click(t, tt, "love 3")
	if !rs[0].Mine || rs[0].Count != 4 || toggled != 0 {
		t.Fatalf("after a click: %+v, toggled %d", rs[0], toggled)
	}
	click(t, tt, "fire 1")
	if rs[1].Mine || rs[1].Count != 0 || toggled != 1 {
		t.Fatalf("taking yours back: %+v", rs[1])
	}
	for name, bad := range map[string]chat.Reaction{"negative": {Emoji: "love", Count: -1}, "mine at zero": {Emoji: "love", Mine: true}, "no such emoji": {Emoji: "nope", Count: 1}} {
		mustPanic(t, name, func() {
			view(300, 100, wd, func(c *ui.Context) { chat.Reactions(c, "R", []chat.Reaction{bad}, chat.ReactionsOptions{}) })
		})
	}
}

func TestEmojiPickerPicksByKeysAcrossCategories(t *testing.T) {
	wd := &world{clk: newClock()}
	picked := ""
	tt := view(400, 200, wd, func(c *ui.Context) {
		if r := chat.EmojiPicker(c, "Emoji", chat.EmojiPickerOptions{}); r.Picked != "" {
			picked = r.Picked
		}
	})
	// A click picks the face under it: the grid's middle.
	click(t, tt, "Emoji")
	if picked != "dead" {
		t.Fatalf("the click picked %q", picked)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	// A row of eight below grin is angry.
	if picked != "angry" {
		t.Fatalf("picked %q", picked)
	}
	// From the last face, a shorter category starts over.
	tt.Key(0, ui.KeyEnd)
	click(t, tt, "ITEMS")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if picked != "coin" {
		t.Fatalf("the items tab picked %q", picked)
	}
}

func TestStickerPickerPicksBigFaces(t *testing.T) {
	wd := &world{clk: newClock()}
	picked := ""
	tt := view(400, 400, wd, func(c *ui.Context) {
		if r := chat.StickerPicker(c, "Stickers", chat.StickerPickerOptions{}); r.Picked != "" {
			picked = r.Picked
		}
	})
	if r := find(t, tt, "Stickers"); r.W != 2*4*40 {
		t.Fatalf("four 40-dot stickers a row, got %v wide", r.W)
	}
	click(t, tt, "Stickers")
	if picked != "sleep" {
		t.Fatalf("the click picked %q", picked)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if picked != "love" {
		t.Fatalf("picked %q", picked)
	}
}

func TestMentionInputCompletesMembers(t *testing.T) {
	wd := &world{clk: newClock()}
	v, mentioned := "", ""
	tt := view(400, 200, wd, func(c *ui.Context) {
		if r := chat.MentionInput(c, "Mention", &v, []string{"Aria", "Arthur", "Brom"}, chat.MentionInputOptions{}); r.Mentioned != "" {
			mentioned = r.Mentioned
		}
	})
	click(t, tt, "Mention")
	tt.Type("hi @ar")
	find(t, tt, "@Aria  @Arthur")
	tt.Key(0, ui.KeyTab)
	if v != "hi @Aria " || mentioned != "Aria" {
		t.Fatalf("value %q, mentioned %q", v, mentioned)
	}
	// Past a space the word is done: nothing is suggested.
	tt.Type("ok")
	if _, ok := tt.Find("@Aria  @Arthur"); ok {
		t.Fatal("suggestions outlived the mention")
	}
}

func TestNotificationBellRingsOnlyWhenNewsArrives(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		n, clicks := 0, 0
		tt := view(200, 100, wd, func(c *ui.Context) {
			if chat.NotificationBell(c, "Bell", n, chat.NotificationBellOptions{}).Clicked {
				clicks++
			}
		})
		bell := find(t, tt, "Bell")
		rest := count(tt.Image(), int(bell.X/2), int(bell.Y/2)+5, 3, 12, pixel.Ink)
		swung := func() bool { return count(tt.Image(), int(bell.X/2), int(bell.Y/2)+5, 3, 12, pixel.Ink) != rest }
		n = 3
		tt.Frame()
		if swung() == reduce {
			t.Fatalf("reduce %v: the bell swung %v on news", reduce, !reduce)
		}
		wd.clk.t = wd.clk.t.Add(time.Second)
		tt.Frame()
		n = 1
		tt.Frame()
		if swung() {
			t.Fatalf("reduce %v: fewer notes rang the bell", reduce)
		}
		click(t, tt, "Bell")
		if clicks != 1 {
			t.Fatal("the bell did not report its click")
		}
	}
	mustPanic(t, "negative", func() {
		view(100, 100, &world{clk: newClock()}, func(c *ui.Context) { chat.NotificationBell(c, "Bell", -1, chat.NotificationBellOptions{}) })
	})
}

func TestCallControlsToggleAndHangUp(t *testing.T) {
	wd := &world{clk: newClock()}
	muted, camera, hung := false, true, 0
	tt := view(300, 100, wd, func(c *ui.Context) {
		if chat.CallControls(c, "Call", &muted, &camera, chat.CallControlsOptions{}).HungUp {
			hung++
		}
	})
	click(t, tt, "MUTE")
	click(t, tt, "CAMERA")
	click(t, tt, "HANG UP")
	if !muted || camera || hung != 1 {
		t.Fatalf("muted %v, camera %v, hung up %d", muted, camera, hung)
	}
}

func TestVoiceMeterLightsBars(t *testing.T) {
	wd := &world{clk: newClock()}
	for _, tc := range []struct {
		level           float64
		good, warn, hot int
	}{{0, 0, 0, 0}, {0.5, 4, 0, 0}, {1, 6, 1, 1}} {
		tt := view(100, 50, wd, func(c *ui.Context) { chat.VoiceMeter(c, "Voice", tc.level, chat.VoiceMeterOptions{}) })
		img := tt.Image()
		// Each lit bar shows one dot on the bottom row.
		if g, w, b := count(img, 0, 11, 23, 1, pixel.Good)/2, count(img, 0, 11, 23, 1, pixel.Warn)/2, count(img, 0, 11, 23, 1, pixel.Bad)/2; g != tc.good || w != tc.warn || b != tc.hot {
			t.Fatalf("level %v lit %d, %d, %d bars", tc.level, g, w, b)
		}
	}
	for _, bad := range []float64{-0.1, 1.1, math.NaN()} {
		mustPanic(t, "level", func() {
			view(100, 50, wd, func(c *ui.Context) { chat.VoiceMeter(c, "Voice", bad, chat.VoiceMeterOptions{}) })
		})
	}
}

func TestAttachmentChipSizesAndRemoves(t *testing.T) {
	wd := &world{clk: newClock()}
	removed := 0
	tt := view(500, 200, wd, func(c *ui.Context) {
		for _, f := range []struct {
			name string
			size int64
		}{{"map.png", 2_516_582}, {"song.mp3", 12 << 20}, {"a.txt", 812}, {"b.bin", 1024}} {
			if chat.AttachmentChip(c, f.name, f.size, chat.AttachmentChipOptions{}).Removed {
				removed++
			}
		}
	})
	for _, s := range []string{"2.4 MB", "12 MB", "812 B", "1.0 KB"} {
		find(t, tt, s)
	}
	click(t, tt, "REMOVE a.txt")
	if removed != 1 {
		t.Fatal("× did not remove")
	}
}

func TestLinkPreviewOpensItsURL(t *testing.T) {
	wd := &world{clk: newClock()}
	tt := view(500, 200, wd, func(c *ui.Context) {
		chat.LinkPreview(c, "https://example.com/a", "A long title that wraps over two lines here", "example.com", chat.LinkPreviewOptions{})
	})
	r := find(t, tt, "A long title that wraps over two lines here")
	if r.H < 2*44 {
		t.Fatalf("a two-line title needs room for the site: %v high", r.H)
	}
	click(t, tt, "A long title that wraps over two lines here")
	tt.Key(0, ui.KeyEnter)
	if got := tt.OpenedURLs(); len(got) != 2 || got[0] != "https://example.com/a" {
		t.Fatalf("opened %q", got)
	}
}

func TestPollMovesYourVote(t *testing.T) {
	wd := &world{clk: newClock()}
	opts := []chat.PollOption{{Text: "Mage", Votes: 2}, {Text: "Knight", Votes: 1}}
	mine, voted := -1, -1
	tt := view(400, 200, wd, func(c *ui.Context) {
		if r := chat.Poll(c, "Class?", opts, &mine, chat.PollOptions{}); r.Voted >= 0 {
			voted = r.Voted
		}
	})
	click(t, tt, "Knight, 1")
	if mine != 1 || opts[1].Votes != 2 || voted != 1 {
		t.Fatalf("first vote: mine %d, %+v", mine, opts)
	}
	click(t, tt, "Mage, 2")
	if mine != 0 || opts[0].Votes != 3 || opts[1].Votes != 1 {
		t.Fatalf("moved vote: mine %d, %+v", mine, opts)
	}
	voted = -1
	click(t, tt, "Mage, 3")
	if voted != -1 || opts[0].Votes != 3 {
		t.Fatal("voting again for your choice counted twice")
	}
	find(t, tt, "4 VOTES")
	mustPanic(t, "lone option", func() {
		one := -1
		view(300, 100, wd, func(c *ui.Context) { chat.Poll(c, "Q", opts[:1], &one, chat.PollOptions{}) })
	})
}
