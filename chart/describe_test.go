package chart

import (
	"math"
	"testing"
)

func TestNumbersStayExact(t *testing.T) {
	for v, want := range map[float64]string{
		0: "0", 12: "12", 3.5: "3.5", -0.5: "-0.5", 0.01: "0.01", 1e-7: "0.0000001", 2e-7: "0.0000002",
		1000001: "1000001", 12345: "12345", 0.1 + 0.2: "0.3",
	} {
		if got := num(v); got != want {
			t.Errorf("num(%v) = %q, want %q", v, got, want)
		}
	}
}

func TestTicksKeepTheirStep(t *testing.T) {
	for _, tc := range []struct {
		v, step float64
		want    string
	}{
		{15000, 5000, "15K"}, {12300, 100, "12.3K"}, {1.5e6, 1e5, "1.5M"}, {1000000.2, 0.2, "1000000.2"},
		{0.6000000000000001, 0.2, "0.6"}, {2e-7, 1e-7, "0.0000002"}, {math.Copysign(0, -1), 1, "0"}, {40, 10, "40"},
	} {
		if got := tick(tc.v, tc.step); got != tc.want {
			t.Errorf("tick(%v, %v) = %q, want %q", tc.v, tc.step, got, tc.want)
		}
	}
}

func TestDescriptionsCarryTheData(t *testing.T) {
	for _, tc := range []struct{ got, want string }{
		{read([]string{"A", "B"}, []float64{0.01, 0.02}), "A 0.01, B 0.02"},
		{readSeries([]string{"X"}, []Series{{Name: "HP", Values: []float64{3}}, {Name: "MP", Values: []float64{1}}}), "HP: X 3; MP: X 1"},
		{readSlices([]Slice{{Name: "Gold", Value: 52}}), "Gold 52"},
		{where([]Point{{X: 1, Y: 3, Size: 9}}, true), "(1, 3) 9"},
		{record([]int{1, 1, -1, 0}), "2-1-1: 1, 1, -1, 0"},
		{ohlc([]Candle{{Label: "MON", Open: 10, High: 14, Low: 9, Close: 13}}), "MON 10, 14, 9, 13"},
		{fives([]string{"D6"}, [][]float64{{1, 2, 3, 4, 5}}), "D6 1, 2, 3, 4, 5"},
	} {
		if tc.got != tc.want {
			t.Errorf("got %q, want %q", tc.got, tc.want)
		}
	}
}

func TestFlatRangesGetDistinctTicks(t *testing.T) {
	for _, r := range [][2]float64{{1e6, 1e6}, {40, 40}, {1e-5, 1e-5}, {-3, -3}, {1e6, 1e6 + 1}} {
		from, step, n := ticks(r[0], r[1], 5)
		top := from + float64(n-1)*step
		if !(from < r[1] && r[0] < top) || tick(from, step) == tick(from+step, step) {
			t.Errorf("%v ticks %v to %v by %v", r, from, top, step)
		}
	}
	// All zeros keep zero as the floor, for bars.
	if from, _, _ := ticks(0, 0, 5); from != 0 {
		t.Errorf("flat zero ticks from %v", from)
	}
}
