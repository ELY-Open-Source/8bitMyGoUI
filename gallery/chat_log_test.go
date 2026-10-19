package gallery

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestChatDemosLogNoTypedText(t *testing.T) {
	for _, tc := range []struct{ path, field string }{{"/c/323-composer", "Message"}, {"/c/329-mentioninput", "Mention"}} {
		a, tt := open(tc.path)
		click(t, tt, tc.field)
		tt.Type("xyzzy")
		tt.Key(0, ui.KeyEnter)
		logs := a.logs[tc.path[3:]]
		if len(logs) == 0 || strings.Contains(strings.Join(logs, "\n"), "xyzzy") {
			t.Errorf("%s logged %q", tc.path, logs)
		}
	}
}

func TestPaletteInUseSpeaksTheLocale(t *testing.T) {
	a, tt := open("/")
	click(t, tt, "language")
	if a.locale != pixel.Zh {
		t.Fatal("no Chinese")
	}
	name := pixel.Palettes[a.palette].Name + ", 使用中"
	if _, ok := tt.Find(name); !ok {
		t.Fatalf("no %q in %q", name, tt.Texts())
	}
}
