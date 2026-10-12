package pixel

import "testing"

func TestArtRoundTrips(t *testing.T) {
	art := "#a.\n.kr\n\n...\nmA#\n"
	sp := MustSprite(art)
	if got := sp.Art(); got != art {
		t.Fatalf("Art() = %q, want %q", got, art)
	}
	one := SpriteOf(2, 1, func(x, y int) (Slot, bool) { return Good, x == 1 })
	if got := one.Art(); got != ".g\n" {
		t.Fatalf("SpriteOf Art() = %q", got)
	}
}

func TestHasKnowsTofu(t *testing.T) {
	if !Has('A') || !Has('中') || !Has('♪') || Has('拖') || Has('\u0378') {
		t.Fatal("Has does not match the font")
	}
}
