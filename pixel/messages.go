package pixel

import "fmt"

var messages = map[Locale]map[string]string{
	En: {
		"ok": "OK", "cancel": "CANCEL", "yes": "YES", "no": "NO", "close": "CLOSE",
		"continue": "CONTINUE?", "game_over": "GAME OVER", "loading": "NOW LOADING",
		"search": "SEARCH", "empty": "NO ITEMS",
	},
	Zh: {
		"ok": "确定", "cancel": "取消", "yes": "是", "no": "否", "close": "关闭",
		"continue": "继续？", "game_over": "游戏结束", "loading": "读取中",
		"search": "搜索", "empty": "空空如也",
	},
}

// Msg is built-in copy in the frame's locale.
func (k Kit) Msg(key string) string {
	s, ok := messages[k.locale][key]
	if !ok {
		panic(fmt.Sprintf("pixel: no message %q for locale %d", key, k.locale))
	}
	return s
}
