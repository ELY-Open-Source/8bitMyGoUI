package pixel

import "testing"

func TestMessagesCoverBothLocales(t *testing.T) {
	for key := range messages[En] {
		if _, ok := messages[Zh][key]; !ok {
			t.Errorf("message %q missing in Chinese", key)
		}
	}
	for key := range messages[Zh] {
		if _, ok := messages[En][key]; !ok {
			t.Errorf("message %q missing in English", key)
		}
	}
	if got := (Kit{locale: Zh}).Msg("game_over"); got != "游戏结束" {
		t.Fatalf("Zh game_over = %q", got)
	}
	if got := (Kit{locale: En}).Msg("continue"); got != "CONTINUE?" {
		t.Fatalf("En continue = %q", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown message did not panic")
		}
	}()
	Kit{locale: En}.Msg("no_such_key")
}

func TestMessagesHaveGlyphs(t *testing.T) {
	for loc, m := range messages {
		for key, s := range m {
			for _, r := range s {
				if _, ok := font12.find(r); !ok && r != ' ' {
					t.Errorf("locale %d message %q: the font has no %q", loc, key, r)
				}
			}
		}
	}
}
