package layout

import (
	"reflect"
	"testing"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func panics(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("%s did not panic", name)
		}
	}()
	fn()
}

func TestOrDef(t *testing.T) {
	if orDef("x", 0, 7) != 7 || orDef("x", 3, 7) != 3 {
		t.Fatal("orDef did not default zero or keep a value")
	}
	panics(t, "negative option", func() { orDef("x", -1, 7) })
	panics(t, "unknown frame style", func() { border(pixel.FrameStyle(99)) })
}

func TestEdges(t *testing.T) {
	type e struct{ right, below bool }
	want := []e{{true, true}, {true, true}, {false, true}, {true, false}, {false, false}}
	for i, w := range want {
		if r, b := edges(i, 5, 3); r != w.right || b != w.below {
			t.Errorf("cell %d of 5 in 3 columns: right %v below %v", i, r, b)
		}
	}
}

func TestLanes(t *testing.T) {
	got := lanes([]int{40, 24, 56, 32, 24, 48, 32, 40}, 3, 4)
	want := [][]int{{0, 4, 7}, {1, 3, 6}, {2, 5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("lanes %v, want %v", got, want)
	}
	panics(t, "zero height", func() { lanes([]int{8, 0}, 2, 4) })
}

func TestShares(t *testing.T) {
	for _, tc := range []struct {
		total, n int
		want     []int
	}{{10, 3, []int{4, 3, 3}}, {9, 3, []int{3, 3, 3}}, {-5, 2, []int{0, 0}}} {
		if got := shares(tc.total, tc.n); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("shares(%d, %d) = %v", tc.total, tc.n, got)
		}
	}
}

func TestThumb(t *testing.T) {
	for _, tc := range []struct {
		h, view, y, maxY, size, at float32
	}{
		{60, 100, 0, 1180, 16, 0},
		{60, 100, 1180, 1180, 16, 44},
		{10, 50, 30, 1000, 10, 0},
		{60, 100, 0, 0, 60, 0},
		{100, 100, 50, 100, 50, 24},
	} {
		if size, at := thumb(tc.h, tc.view, tc.y, tc.maxY, 2); size != tc.size || at != tc.at {
			t.Errorf("thumb(h %v view %v y %v max %v) = %v at %v, want %v at %v", tc.h, tc.view, tc.y, tc.maxY, size, at, tc.size, tc.at)
		}
	}
}

func TestFit(t *testing.T) {
	for _, tc := range []struct {
		w, h         int
		roomW, roomH float32
		base, want   int
	}{
		{64, 40, 256, 200, 2, 4},
		{64, 40, 200, 200, 2, 2},
		{10, 10, 0, 0, 1, 1},
		{160, 144, 1000, 1000, 1, 4},
		{160, 144, 400, 300, 1, 2},
	} {
		if got := fit(tc.w, tc.h, tc.roomW, tc.roomH, tc.base); got != tc.want {
			t.Errorf("fit(%d×%d in %v×%v from %d) = %d, want %d", tc.w, tc.h, tc.roomW, tc.roomH, tc.base, got, tc.want)
		}
	}
}
