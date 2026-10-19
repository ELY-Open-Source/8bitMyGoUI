package data_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestDescriptionListSetsTermsBesideText(t *testing.T) {
	tt := view(400, 200, newClock(), false, func(c *ui.Context) {
		data.DescriptionList(c, [][2]string{{"Class", "Knight"}, {"Home", "Castle Town"}}, data.DescriptionListOptions{})
	})
	term, desc := find(t, tt, "Home"), find(t, tt, "Castle Town")
	if desc.X != term.X+176 || desc.Y != term.Y {
		t.Fatalf("description at %v beside term at %v", desc, term)
	}
	keeps(t, tt.Image())
}

func TestLeaderboardMedalsAndPlayer(t *testing.T) {
	scores := []data.Score{{Name: "ACE", Score: 9000}, {Name: "BOB", Score: 8000}, {Name: "CAT", Score: 7000}, {Name: "YOU", Score: 6000}, {Name: "DOG", Score: 5000}}
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		data.Leaderboard(c, scores, data.LeaderboardOptions{Player: 3})
	})
	for _, l := range []string{"GOLD", "SILVER", "BRONZE", "0006000", "4 YOU 6000"} {
		if !shows(tt, l) {
			t.Fatalf("no %q in %q", l, tt.Texts())
		}
	}
	r, img := find(t, tt, "4 YOU 6000"), tt.Image()
	if s := slot(t, img, int(r.X)+r2i(r.W)-4, int(r.Y)+4); s != pixel.Accent {
		t.Fatalf("your row is %v", s)
	}
	keeps(t, img)
}

func r2i(f float32) int { return int(f) }

func TestTimelineMarksEachEvent(t *testing.T) {
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		data.Timeline(c, []data.Event{{Date: "1986", Text: "Zelda"}, {Date: "1991", Text: "A Link to the Past"}}, data.TimelineOptions{})
	})
	for _, l := range []string{"1986", "Zelda", "1991", "A Link to the Past"} {
		if !shows(tt, l) {
			t.Fatalf("no %q", l)
		}
	}
	r, img := find(t, tt, "1986"), tt.Image()
	// The marker's middle sits 13 dots left of the date.
	if s := slot(t, img, int(r.X)-26+8, int(r.Y)+10); s != pixel.Accent {
		t.Fatalf("marker is %v", s)
	}
	keeps(t, img)
}

func TestStatCompareInksDeltas(t *testing.T) {
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		data.StatCompare(c, []data.Change{{Label: "ATK", Old: 12, New: 15}, {Label: "DEF", Old: 9, New: 7}, {Label: "SPD", Old: 5, New: 5}}, data.StatCompareOptions{})
	})
	for _, l := range []string{"ATK 12 → 15 ▲3", "DEF 9 → 7 ▼2", "SPD 5 → 5 "} {
		if !shows(tt, l) {
			t.Fatalf("no %q in %q", l, tt.Texts())
		}
	}
	img := tt.Image()
	ink := func(label string) map[pixel.Slot]bool {
		r, out := find(t, tt, label), map[pixel.Slot]bool{}
		for y := int(r.Y); y < int(r.Y+r.H); y += 2 {
			for x := int(r.X); x < int(r.X+r.W); x += 2 {
				out[slot(t, img, x, y)] = true
			}
		}
		return out
	}
	if up, dn := ink("ATK 12 → 15 ▲3"), ink("DEF 9 → 7 ▼2"); !up[pixel.Good] || up[pixel.Bad] || !dn[pixel.Bad] || dn[pixel.Good] {
		t.Fatalf("delta inks: up %v, down %v", up, dn)
	}
	keeps(t, img)
}

func TestBestiaryHidesUnseen(t *testing.T) {
	m := data.Monster{Sprite: icon.Skull, Name: "Skeleton", Stats: [][2]string{{"HP", "30"}}, Notes: "Rattles."}
	for _, seen := range []bool{true, false} {
		m.Seen = seen
		tt := view(300, 240, newClock(), false, func(c *ui.Context) { data.Bestiary(c, m, data.BestiaryOptions{}) })
		if shows(tt, "Skeleton") != seen || shows(tt, "HP 30") != seen || shows(tt, "Rattles.") != seen || shows(tt, "???") == seen {
			t.Fatalf("seen %v: texts %q", seen, tt.Texts())
		}
		keeps(t, tt.Image())
	}
}

func TestFeedListsEntries(t *testing.T) {
	tt := view(300, 200, newClock(), false, func(c *ui.Context) {
		data.Feed(c, []data.Activity{{Icon: icon.Star, Text: "Aria found a star", Time: "2m"}, {Icon: icon.Coin, Text: "Brom sold a coin", Time: "5m"}}, data.FeedOptions{})
	})
	a, b := find(t, tt, "Aria found a star"), find(t, tt, "Brom sold a coin")
	if !shows(tt, "2m") || !shows(tt, "5m") || b.Y <= a.Y || b.X != a.X {
		t.Fatalf("entries at %v and %v", a, b)
	}
	keeps(t, tt.Image())
}

func TestCardGridFitsColumns(t *testing.T) {
	// Widths are pixels, two per dot; cards need 96 dots.
	for _, tc := range []struct{ w, cols int }{{800, 4}, {500, 2}, {180, 1}} {
		tt := view(tc.w+32, 300, newClock(), false, func(c *ui.Context) {
			ui.Column(c).Padding(16).Children(func() {
				data.CardGrid(c, 6, data.CardGridOptions{}, func(i int) {
					pixel.Text(c, fmt.Sprintf("card %d", i), pixel.TextOptions{})
				})
			})
		})
		tt.Frame()
		first, row := find(t, tt, "card 0"), 0
		for i := range 6 {
			if find(t, tt, fmt.Sprintf("card %d", i)).Y == first.Y {
				row++
			}
		}
		if row != tc.cols {
			t.Fatalf("%d px wide: %d cards in the first row, want %d", tc.w, row, tc.cols)
		}
		keeps(t, tt.Image())
	}
}

func TestItemSlotRefusesNaNCooldown(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("a NaN cooldown did not panic")
		}
	}()
	view(100, 100, newClock(), false, func(c *ui.Context) {
		data.ItemSlot(c, data.Item{Name: "Ether"}, data.ItemSlotOptions{Cooldown: math.NaN()})
	})
}
