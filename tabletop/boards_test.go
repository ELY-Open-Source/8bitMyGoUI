package tabletop_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/tabletop"
)

// cell clicks square i of a labeled grid.
func cell(t *testing.T, tt *ui.Tester, label string, cols, side, i int) {
	t.Helper()
	r := find(t, tt, label)
	tt.ClickAt(r.X+float32(2*(i%cols*side+side/2)), r.Y+float32(2*(i/cols*side+side/2)))
}

func TestChessboardMovesCapturesAndSwaps(t *testing.T) {
	var board [64]tabletop.Piece
	board[57] = tabletop.Piece{Kind: tabletop.Knight, White: true}
	board[62] = tabletop.Piece{Kind: tabletop.Knight, White: true}
	board[42] = tabletop.Piece{Kind: tabletop.Pawn}
	var last tabletop.ChessboardResult
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		if r := tabletop.Chessboard(c, "Board", &board, tabletop.ChessboardOptions{}); r.Moved {
			last = r
		}
	})
	// Its own side swaps the held piece; reclicking drops it.
	cell(t, tt, "CHESSBOARD", 8, 16, 62)
	cell(t, tt, "CHESSBOARD", 8, 16, 57)
	cell(t, tt, "CHESSBOARD", 8, 16, 57)
	x, y, _, _ := dots(find(t, tt, "CHESSBOARD"))
	if last.Moved || board[42].Kind != tabletop.Pawn || count(tt.Image(), x, y, 128, 128, pixel.AccentLit) > 0 {
		t.Fatal("a dropped piece is still held")
	}
	cell(t, tt, "CHESSBOARD", 8, 16, 57)
	cell(t, tt, "CHESSBOARD", 8, 16, 42)
	if last.From != 57 || last.To != 42 || board[42] != (tabletop.Piece{Kind: tabletop.Knight, White: true}) || board[57].Kind != tabletop.None {
		t.Fatalf("the knight's capture went %+v", last)
	}
}

func TestGoBoardTakesTurns(t *testing.T) {
	board := make([]tabletop.Stone, 81)
	placed := []int{}
	tt := view(400, 400, newWorld(), func(c *ui.Context) {
		if i := tabletop.GoBoard(c, "Goban", board, tabletop.GoBoardOptions{}).Placed; i >= 0 {
			placed = append(placed, i)
		}
	})
	find(t, tt, "BLACK TO PLAY")
	cell(t, tt, "Goban", 9, 10, 40)
	find(t, tt, "WHITE TO PLAY")
	cell(t, tt, "Goban", 9, 10, 40)
	cell(t, tt, "Goban", 9, 10, 41)
	if len(placed) != 2 || board[40] != tabletop.Black || board[41] != tabletop.White {
		t.Fatalf("placed %v: %v %v", placed, board[40], board[41])
	}
	find(t, tt, "BLACK TO PLAY")
}

func TestRouletteSpinsToItsResult(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		segs := []string{"ZERO", "ONE", "TWO", "THREE", "FOUR"}
		result, cue := 0, 0
		tt := view(400, 300, wd, func(c *ui.Context) {
			tabletop.Roulette(c, "Wheel", segs, result, tabletop.RouletteOptions{Cue: cue})
		})
		find(t, tt, "ZERO")
		result, cue = 3, 1
		tt.Frame()
		_, landed := tt.Find("THREE")
		if landed != reduce {
			t.Errorf("reduce %v: landed at once is %v", reduce, landed)
		}
		later(tt, wd, 3*time.Second)
		find(t, tt, "THREE")
	}
}

func TestCoinFlipLandsOnItsFace(t *testing.T) {
	wd := newWorld()
	heads, cue := true, 0
	tt := view(300, 300, wd, func(c *ui.Context) { tabletop.CoinFlip(c, "Coin", heads, tabletop.CoinFlipOptions{Cue: cue}) })
	find(t, tt, "HEADS")
	heads, cue = false, 1
	tt.Frame()
	if _, ok := tt.Find("TAILS"); ok {
		t.Fatal("the coin landed before it flew")
	}
	later(tt, wd, time.Second)
	find(t, tt, "TAILS")
}

func TestBingoCardCallsALineOnce(t *testing.T) {
	var numbers [25]int
	for i := range numbers {
		numbers[i] = i + 1
	}
	var marked [25]bool
	bingos, changes := 0, 0
	tt := view(400, 400, newWorld(), func(c *ui.Context) {
		r := tabletop.BingoCard(c, "Card", numbers, &marked, tabletop.BingoCardOptions{})
		if r.Bingo {
			bingos++
		}
		if r.Changed {
			changes++
		}
	})
	if !marked[12] {
		t.Fatal("FREE is not marked")
	}
	// The middle row needs four marks beside FREE.
	for _, i := range []int{10, 11, 13} {
		cell(t, tt, "Card", 5, 18, i)
	}
	if bingos != 0 {
		t.Fatal("BINGO before the line was full")
	}
	cell(t, tt, "Card", 5, 18, 14)
	tt.Frame()
	if bingos != 1 || !tt.HasText("BINGO!") {
		t.Fatalf("a full row called %d bingos", bingos)
	}
	cell(t, tt, "Card", 5, 18, 12)
	if !marked[12] || changes != 4 {
		t.Fatalf("FREE came unmarked, or changed: %d changes", changes)
	}
}

func TestTicTacToeAnswersAndContinues(t *testing.T) {
	over, winner := 0, -1
	tt := view(400, 400, newWorld(), func(c *ui.Context) {
		if r := tabletop.TicTacToe(c, "Grid", tabletop.TicTacToeOptions{}); r.Over {
			over, winner = over+1, r.Winner
		}
	})
	// X plays 0, 8, 1; O answers 4, 2, 6.
	for _, i := range []int{0, 8, 1} {
		cell(t, tt, "Grid", 3, 22, i)
	}
	if over != 1 || winner != 2 {
		t.Fatalf("over %d, winner %d: %q", over, winner, tt.Texts())
	}
	// X is a Bad cross; O is an Info ring.
	x, y, _, _ := dots(find(t, tt, "Grid"))
	if count(tt.Image(), x, y, 22, 22, pixel.Bad) == 0 || count(tt.Image(), x+22, y+22, 22, 22, pixel.Info) == 0 {
		t.Fatal("the marks are not drawn")
	}
	find(t, tt, "GAME OVER")
	find(t, tt, "CONTINUE?")
	cell(t, tt, "Grid", 3, 22, 0)
	if tt.HasText("CONTINUE?") {
		t.Fatal("a click after the end did not start over")
	}
}

func TestScoreTrackStandsMeeples(t *testing.T) {
	tt := view(500, 200, newWorld(), func(c *ui.Context) {
		tabletop.ScoreTrack(c, "Track", []int{3, 3, 42}, tabletop.ScoreTrackOptions{})
	})
	x, y, _, _ := dots(find(t, tt, "Track"))
	img := tt.Image()
	// Space 3 holds two pawns; 42 laps to space 12.
	if count(img, x+60, y, 20, 20, pixel.Bad) == 0 || count(img, x+60, y, 20, 20, pixel.Info) == 0 {
		t.Error("space 3 lacks its two pawns")
	}
	if count(img, x+7*20, y+20, 20, 20, pixel.Good) < 20 {
		t.Error("space 12 lacks the third meeple")
	}
}

func TestMeepleStandsOrWaits(t *testing.T) {
	tt := view(300, 200, newWorld(), func(c *ui.Context) {
		tabletop.Meeple(c, "Two", tabletop.MeepleOptions{Player: 2})
		tabletop.Meeple(c, "Seat", tabletop.MeepleOptions{Player: 2, Empty: true})
	})
	x, y, w, h := dots(find(t, tt, "Two"))
	if count(tt.Image(), x, y, w, h, pixel.Good) == 0 {
		t.Error("player 3's meeple is not Good")
	}
	x, y, w, h = dots(find(t, tt, "Seat"))
	if count(tt.Image(), x, y, w, h, pixel.Good) > 0 {
		t.Error("an empty seat shows a meeple")
	}
}
