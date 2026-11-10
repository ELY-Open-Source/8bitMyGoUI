package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestDiagramDemoDrawsAndLogsEachSample(t *testing.T) {
	a, tt := open("/c/501-mermaid")
	for _, name := range []string{"TRADE", "BAG", "QUEST"} {
		click(t, tt, name)
		tt.Frame()
		logs := a.logs["501-mermaid"]
		if len(logs) == 0 || logs[len(logs)-1] != "drew "+name {
			t.Errorf("%s: logs %q, want them to end with its draw", name, logs)
		}
		if _, ok := tt.AccessByRole(ui.RoleImage, name); !ok || tt.HasText("DIAGRAM FAILED") {
			t.Errorf("%s: image %v, texts %q", name, ok, tt.Texts())
		}
	}
}
