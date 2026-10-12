package display_test

import (
	"image"
	"image/color"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func newClock() *clock { return &clock{time.Unix(2_000_000, 0)} }

type env struct {
	pal    *pixel.Palette
	locale pixel.Locale
	clk    *clock
}

func view(w, h int, e *env, build func(c *ui.Context)) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: e.pal, Locale: e.locale, Now: e.clk.now})
		build(c)
	}, w, h)
}

func spring() *env { return &env{pixel.SpringDay, pixel.En, newClock()} }

func find(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q in %q", label, tt.Texts())
	}
	return r
}

// slots maps a box's dots to palette slots, or fails.
func slots(t *testing.T, img *image.RGBA, pal *pixel.Palette, r ui.Rect) []pixel.Slot {
	t.Helper()
	var out []pixel.Slot
	for y := int(r.Y); y < int(r.Y+r.H); y += 2 {
		for x := int(r.X); x < int(r.X+r.W); x += 2 {
			c, found := img.RGBAAt(x, y), false
			for s, pc := range pal.Colors {
				if pc.R == c.R && pc.G == c.G && pc.B == c.B {
					out, found = append(out, pixel.Slot(s)), true
					break
				}
			}
			if !found {
				t.Fatalf("pixel %d,%d off %s: %v", x, y, pal.Name, c)
			}
		}
	}
	return out
}

func has(s []pixel.Slot, want pixel.Slot) bool {
	for _, v := range s {
		if v == want {
			return true
		}
	}
	return false
}

func TestIconLabels(t *testing.T) {
	tt := view(200, 100, spring(), func(c *ui.Context) {
		display.Icon(c, icon.Heart, display.IconOptions{Scale: 3, Label: "heart"})
		display.Icon(c, icon.Star, display.IconOptions{})
	})
	if r := find(t, tt, "heart"); r.W != float32(icon.Heart.W*3*2) {
		t.Fatalf("icon box %v", r)
	}
	if len(tt.Texts()) != 1 {
		t.Fatalf("an unlabeled icon reached assistive technology: %q", tt.Texts())
	}
}

func TestSpritePlaysAndMirrors(t *testing.T) {
	e := spring()
	sp := pixel.MustSprite("r.\n..\n\n.r\n..")
	tt := view(100, 100, e, func(c *ui.Context) {
		ui.Column(c).Children(func() {
			display.Sprite(c, sp, display.SpriteOptions{FPS: 4, Label: "s"})
			display.Sprite(c, sp, display.SpriteOptions{Mirror: true, Label: "m"})
		})
	})
	at := func(label string) []pixel.Slot { return slots(t, tt.Image(), e.pal, find(t, tt, label)) }
	f0 := at("s")
	if f0[0] != pixel.Bad || at("m")[1] != pixel.Bad {
		t.Fatalf("frame 0 %v, mirrored %v", f0, at("m"))
	}
	e.clk.t = e.clk.t.Add(time.Second / 4)
	tt.Frame()
	if at("s")[1] != pixel.Bad {
		t.Fatal("the sprite did not play its second frame")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if at("s")[0] != pixel.Bad {
		t.Fatal("Reduce Motion did not hold the first frame")
	}
}

func TestPictureRecolors(t *testing.T) {
	e := spring()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for i := range img.Pix {
		img.Pix[i] = 200
	}
	tt := view(100, 100, e, func(c *ui.Context) { display.Picture(c, img, display.PictureOptions{Label: "pic"}) })
	day := slots(t, tt.Image(), e.pal, find(t, tt, "pic"))
	e.pal = pixel.WinterNight
	tt.Frame()
	night := slots(t, tt.Image(), e.pal, find(t, tt, "pic"))
	// Quantized afresh, a gray lands on other slots at night.
	if slices.Equal(day, night) {
		t.Fatalf("the picture kept its day slots at night: %v", night[:4])
	}
}

// listImage holds a slice, so its values do not compare.
type listImage struct{ px []color.RGBA }

func (l listImage) ColorModel() color.Model { return color.RGBAModel }
func (l listImage) Bounds() image.Rectangle { return image.Rect(0, 0, 4, 4) }
func (l listImage) At(x, y int) color.Color { return l.px[y*4+x] }

func TestPictureFollowsEditsAndAnyImage(t *testing.T) {
	e := spring()
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	list := listImage{make([]color.RGBA, 16)}
	for i := range list.px {
		list.px[i] = color.RGBA{0, 0, 0, 255}
	}
	tt := view(200, 100, e, func(c *ui.Context) {
		display.Picture(c, img, display.PictureOptions{Label: "rgba"})
		display.Picture(c, list, display.PictureOptions{Label: "list"})
	})
	before := slots(t, tt.Image(), e.pal, find(t, tt, "rgba"))
	for i := range img.Pix {
		img.Pix[i] = 255
	}
	tt.Frame()
	if slices.Equal(before, slots(t, tt.Image(), e.pal, find(t, tt, "rgba"))) {
		t.Fatal("editing the image in place did not redraw it")
	}
	if got := slots(t, tt.Image(), e.pal, find(t, tt, "list")); got[0] != pixel.Ink {
		t.Fatalf("a black list image drew %v", got[0])
	}
}

func TestBadgeAndTrophyTakeAnyText(t *testing.T) {
	tt := view(400, 200, spring(), func(c *ui.Context) {
		display.Badge(c, 0, display.BadgeOptions{Text: "新"})
		display.Badge(c, 7, display.BadgeOptions{})
		display.Trophy(c, "冠军 É", display.TrophyOptions{})
	})
	e := spring()
	if cjk, digit := find(t, tt, "新"), find(t, tt, "7"); cjk.H <= digit.H {
		t.Fatalf("the CJK badge %v is no taller than the micro one %v", cjk, digit)
	}
	if !has(slots(t, tt.Image(), e.pal, find(t, tt, "新")), pixel.AccentInk) {
		t.Fatal("the CJK badge drew no text")
	}
	// The cup has no Ink; only the plaque's words do.
	if !has(slots(t, tt.Image(), e.pal, find(t, tt, "冠军 É")), pixel.Ink) {
		t.Fatal("the trophy plaque drew no text")
	}
}

func TestSignpostIsOneTabStop(t *testing.T) {
	tt := view(300, 300, spring(), func(c *ui.Context) {
		display.Signpost(c, []string{"Town", "Cave", "Shop"}, display.SignpostOptions{})
		ui.ButtonBase(c).Label("after").Size(20, 20)
	})
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Town") {
		t.Fatal("Tab did not enter at the first board")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Cave") {
		t.Fatal("Down did not move to the next board")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("after") {
		t.Fatal("Tab did not leave the signpost")
	}
}

func TestBadgeTagStatus(t *testing.T) {
	e := spring()
	removed, disabled := 0, false
	tt := view(400, 200, e, func(c *ui.Context) {
		display.Badge(c, 128, display.BadgeOptions{})
		display.Badge(c, 5, display.BadgeOptions{Text: "NEW"})
		if display.Tag(c, "fire", display.TagOptions{Removable: true, Disabled: disabled}).Removed {
			removed++
		}
		display.StatusDot(c, display.Online, display.StatusDotOptions{Live: true})
	})
	find(t, tt, "99+")
	find(t, tt, "NEW")
	if err := tt.Click("REMOVE fire"); err != nil || removed != 1 {
		t.Fatalf("remove: %v, %d", err, removed)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("REMOVE fire")
	if removed != 1 {
		t.Fatal("a disabled tag was removed")
	}
	dot := func() []pixel.Slot { return slots(t, tt.Image(), e.pal, find(t, tt, "ONLINE")) }
	if !has(dot(), pixel.Good) {
		t.Fatal("step 0 hid the online light")
	}
	e.clk.t = e.clk.t.Add(500 * time.Millisecond)
	tt.Frame()
	if has(dot(), pixel.Good) {
		t.Fatal("a live online dot did not blink")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if !has(dot(), pixel.Good) {
		t.Fatal("Reduce Motion left the light off")
	}
	e.locale = pixel.Zh
	tt.Frame()
	find(t, tt, "在线")
}

func TestSignpostBoards(t *testing.T) {
	clicked, disabled := -1, false
	tt := view(300, 300, spring(), func(c *ui.Context) {
		if r := display.Signpost(c, []string{"Town", "Cave", "Shop"}, display.SignpostOptions{Disabled: disabled}); r.Clicked >= 0 {
			clicked = r.Clicked
		}
	})
	if err := tt.Click("Cave"); err != nil || clicked != 1 {
		t.Fatalf("clicked %d, %v", clicked, err)
	}
	disabled, clicked = true, -1
	tt.Frame()
	_ = tt.Click("Shop")
	if clicked != -1 {
		t.Fatal("a disabled board clicked")
	}
	for _, n := range []int{0, 4} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%d boards did not panic", n)
				}
			}()
			view(100, 100, spring(), func(c *ui.Context) { display.Signpost(c, make([]string, n), display.SignpostOptions{}) })
		}()
	}
}

func TestQRCodeQuietZoneAndFinder(t *testing.T) {
	e := spring()
	tt := view(400, 400, e, func(c *ui.Context) { display.QRCode(c, "HELLO", display.QRCodeOptions{Scale: 1}) })
	r := find(t, tt, "HELLO")
	s := slots(t, tt.Image(), e.pal, r)
	n := int(r.W / 2)
	if n != 21+8 {
		t.Fatalf("version 1 code with quiet zone is %d modules, want 29", n)
	}
	at := func(x, y int) pixel.Slot { return s[y*n+x] }
	for i := range 4 {
		if at(i, i) != pixel.AccentInk || at(n-1-i, i) != pixel.AccentInk {
			t.Fatal("the quiet zone is not light")
		}
	}
	for i := range 7 {
		if at(4+i, 4) != pixel.Ink || at(4, 4+i) != pixel.Ink || at(4+i, 10) != pixel.Ink {
			t.Fatal("the top-left finder is not a dark ring")
		}
	}
	if at(5, 5) != pixel.AccentInk || at(7, 7) != pixel.Ink {
		t.Fatal("the finder's gap or center is wrong")
	}
}

func TestBubbleTailMovesContent(t *testing.T) {
	side := display.Bottom
	tt := view(300, 200, spring(), func(c *ui.Context) {
		display.Bubble(c, display.BubbleOptions{Tail: side}, func() { pixel.Text(c, "HI", pixel.TextOptions{}) })
	})
	below := find(t, tt, "HI")
	side = display.Top
	tt.Frame()
	if above := find(t, tt, "HI"); above.Y != below.Y+5*2 {
		t.Fatalf("a top tail put the text at %v, bottom at %v", above.Y, below.Y)
	}
	side = display.Left
	tt.Frame()
	if left := find(t, tt, "HI"); left.X != below.X+5*2 {
		t.Fatalf("a left tail put the text at %v", left.X)
	}
}

func TestPanicsOnBadInput(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"emoji":   func(c *ui.Context) { display.Emoji(c, "nope", display.EmojiOptions{}) },
		"status":  func(c *ui.Context) { display.StatusDot(c, 9, display.StatusDotOptions{}) },
		"medal":   func(c *ui.Context) { display.Medal(c, 9, display.MedalOptions{}) },
		"side":    func(c *ui.Context) { display.Bubble(c, display.BubbleOptions{Tail: 9}, func() {}) },
		"scale":   func(c *ui.Context) { display.Icon(c, icon.Star, display.IconOptions{Scale: 5}) },
		"qr long": func(c *ui.Context) { display.QRCode(c, string(make([]byte, 4000)), display.QRCodeOptions{}) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s did not panic", name)
				}
			}()
			view(200, 200, spring(), build)
		}()
	}
}
