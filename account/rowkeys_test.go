package account_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/account"
)

func TestSessionFocusFollowsItsRow(t *testing.T) {
	wd := newWorld()
	ss := []account.Session{{Device: "Laptop", Seen: wd.clk.t}, {Device: "Phone", Seen: wd.clk.t}, {Device: "Tablet", Seen: wd.clk.t}}
	gone := ""
	tt := view(600, 400, wd, func(c *ui.Context) {
		if r := account.SessionList(c, "Sessions", ss, account.SessionListOptions{}); r.SignedOut >= 0 {
			gone = ss[r.SignedOut].Device
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	// The laptop leaves; the focus stays on the phone's button.
	ss = ss[1:]
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if gone != "Phone" {
		t.Fatalf("Enter after the laptop left signed out %q", gone)
	}
	defer func() {
		if r := fmt.Sprint(recover()); !strings.HasPrefix(r, "account: SessionList") {
			t.Fatalf("two phones gave %q", r)
		}
	}()
	twins := []account.Session{{Device: "Phone"}, {Device: "Phone"}}
	view(600, 400, newWorld(), func(c *ui.Context) { account.SessionList(c, "S", twins, account.SessionListOptions{}) })
}
