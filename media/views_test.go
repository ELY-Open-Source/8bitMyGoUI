package media_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/media"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// solid is a w by h picture in one slot.
func solid(w, h int, s pixel.Slot) *pixel.Sprite {
	return pixel.SpriteOf(w, h, func(x, y int) (pixel.Slot, bool) { return s, true })
}

func TestImageCompareSplitsWhereDragged(t *testing.T) {
	wd := newWorld()
	split, changed := 0.5, false
	tt := view(300, 200, wd, func(c *ui.Context) {
		changed = media.ImageCompare(c, "Compare", solid(10, 6, pixel.Good), solid(10, 6, pixel.Bad), &split, media.ImageCompareOptions{Scale: 2}).Changed || changed
	})
	r := find(t, tt, "Compare")
	x, y, w, h := dots(r)
	img := tt.Image()
	// Row 0 is above the grip.
	if w != 20 || h != 12 || at(img, wd.pal, x+2, y) != pixel.Good || at(img, wd.pal, x+10, y) != pixel.Edge || at(img, wd.pal, x+15, y) != pixel.Bad {
		t.Fatal("before is not left of the split, after right")
	}
	tt.Press(r.X+r.W/4, r.Y+2)
	tt.Release(r.X+r.W/4, r.Y+2)
	if split != 0.25 || !changed || at(tt.Image(), wd.pal, x+7, y) != pixel.Bad {
		t.Fatalf("pressing a quarter in split at %v", split)
	}
	tt.Key(0, ui.KeyHome)
	if split != 0 || at(tt.Image(), wd.pal, x+2, y) != pixel.Bad || at(tt.Image(), wd.pal, x+6, y+h/2) != pixel.Edge {
		t.Fatal("at 0, all is after and the grip stays whole")
	}
	tt.Key(0, ui.KeyEnd)
	if split != 1 || at(tt.Image(), wd.pal, x+17, y) != pixel.Good {
		t.Fatal("at 1, all is before")
	}
	mustPanic(t, "pictures of two sizes", func() {
		view(300, 200, newWorld(), func(c *ui.Context) {
			s := 0.5
			media.ImageCompare(c, "C", solid(2, 2, pixel.Ink), solid(3, 2, pixel.Ink), &s, media.ImageCompareOptions{})
		})
	})
}

func TestGalleryChoosesByArrowsAndClicks(t *testing.T) {
	wd := newWorld()
	pics := []*pixel.Sprite{solid(4, 3, pixel.Good), solid(4, 3, pixel.Bad), solid(4, 3, pixel.Info)}
	chosen, changed := 0, false
	tt := view(300, 200, wd, func(c *ui.Context) {
		changed = media.Gallery(c, "Views", pics, []string{"Green", "Red", "Blue"}, &chosen, media.GalleryOptions{Scale: 2}).Changed || changed
	})
	big := func(name string) pixel.Slot {
		x, y, _, _ := dots(find(t, tt, name))
		return at(tt.Image(), wd.pal, x+4, y+4)
	}
	if big("Green") != pixel.Good {
		t.Fatal("the big picture is not the chosen one")
	}
	// The middle of the strip is the second thumbnail.
	click(t, tt, "Views")
	if chosen != 1 || !changed || big("Red") != pixel.Bad {
		t.Fatalf("clicking the middle chose %d", chosen)
	}
	tt.Key(0, ui.KeyRight)
	if chosen != 2 || big("Blue") != pixel.Info {
		t.Fatalf("Right chose %d", chosen)
	}
	changed = false
	tt.Key(0, ui.KeyRight)
	if chosen != 2 || changed {
		t.Fatalf("Right at the end: %d, changed %v", chosen, changed)
	}
	// The chosen thumbnail wears the cursor ring.
	x, y, _, _ := dots(find(t, tt, "Views"))
	if at(tt.Image(), wd.pal, x+2*10, y) != pixel.Accent || at(tt.Image(), wd.pal, x, y) == pixel.Accent {
		t.Fatal("the ring is not on the chosen thumbnail")
	}
	for name, fn := range map[string]func(c *ui.Context){
		"names short of pictures": func(c *ui.Context) {
			n := 0
			media.Gallery(c, "V", pics, []string{"Green"}, &n, media.GalleryOptions{})
		},
		"a choice past the end": func(c *ui.Context) {
			n := 3
			media.Gallery(c, "V", pics, []string{"Green", "Red", "Blue"}, &n, media.GalleryOptions{})
		},
	} {
		mustPanic(t, name, func() { view(300, 200, newWorld(), fn) })
	}
}
