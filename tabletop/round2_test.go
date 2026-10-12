package tabletop_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/tabletop"
)

func TestChessHoldsAreLoggedAndReported(t *testing.T) {
	var board [64]tabletop.Piece
	board[0] = tabletop.Piece{Kind: tabletop.Rook}
	held := []int{}
	log := logged(func() {
		tt := view(600, 400, newWorld(), func(c *ui.Context) {
			if r := tabletop.Chessboard(c, "B", &board, tabletop.ChessboardOptions{}); r.Changed {
				held = append(held, r.Held)
			}
		})
		cell(t, tt, "CHESSBOARD", 8, 16, 0)
		cell(t, tt, "CHESSBOARD", 8, 16, 0)
		// A move is a move, not a hold put down.
		cell(t, tt, "CHESSBOARD", 8, 16, 0)
		cell(t, tt, "CHESSBOARD", 8, 16, 8)
	})
	if len(held) != 3 || held[0] != 0 || held[1] != -1 || held[2] != 0 || strings.Count(log, "chess hold") != 3 {
		t.Fatalf("holds %v; log:\n%s", held, log)
	}
}

func TestMeepleSlotIsTheItemSlot(t *testing.T) {
	got := view(200, 200, newWorld(), func(c *ui.Context) {
		tabletop.Meeple(c, "M", tabletop.MeepleOptions{Empty: true})
	}).Image().Pix
	want := view(200, 200, newWorld(), func(c *ui.Context) {
		k := pixel.K(c)
		ui.Box(c).Size(k.D(data.SlotSize), k.D(data.SlotSize)).Draw(func(p *ui.Painter, r ui.Rect) {
			data.DrawSlot(k.Canvas(p, r), 0, 0, data.Item{}, data.ItemSlotOptions{})
		})
	}).Image().Pix
	if !bytes.Equal(got, want) {
		t.Fatal("an empty meeple slot is not data's item slot")
	}
}
