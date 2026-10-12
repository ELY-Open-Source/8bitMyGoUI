package gallery

import (
	"strings"
	"testing"
)

func TestShopDemosLogTheirActions(t *testing.T) {
	// Clicks that only focus a component may log nothing.
	quiet := map[string]bool{"Potions": true, "KEYPAD": true}
	for _, tc := range []struct {
		path  string
		steps []string
	}{
		{"/c/437-wallet", []string{"FIND GOLD"}}, {"/c/438-shopitem", []string{"BUY POTION"}},
		{"/c/439-cart", []string{"REMOVE ETHER"}}, {"/c/441-checkout", []string{"NEXT ▶"}},
		{"/c/442-currencyconvert", []string{"CONVERT"}}, {"/c/444-quantity", []string{"Potions", "key:Right"}},
		{"/c/445-pricingtiers", []string{"SILVER, 5 GEMS"}}, {"/c/446-coupon", []string{"COPY"}},
		{"/c/447-dailyreward", []string{"CLAIM"}}, {"/c/449-tradewindow", []string{"CONFIRM", "BROM CONFIRMS"}},
		{"/c/450-vendingmachine", []string{"INSERT COIN"}},
	} {
		a, tt := open(tc.path)
		slug := tc.path[3:]
		for _, step := range tc.steps {
			before := len(a.logs[slug])
			if k, ok := strings.CutPrefix(step, "key:"); ok {
				tt.Key(0, letters[k])
			} else {
				click(t, tt, step)
			}
			if len(a.logs[slug]) <= before && !quiet[step] {
				t.Errorf("%s: %s logged nothing: %q", slug, step, a.logs[slug])
			}
		}
	}
}
