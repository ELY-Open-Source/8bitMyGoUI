package pixel

import "fmt"

var messages = map[Locale]map[string]string{
	En: {
		"ok": "OK", "cancel": "CANCEL", "yes": "YES", "no": "NO", "close": "CLOSE",
		"continue": "CONTINUE?", "game_over": "GAME OVER", "loading": "NOW LOADING",
		"search": "SEARCH", "empty": "NO ITEMS", "minimize": "MINIMIZE", "resize": "RESIZE",
		"scroll_up": "SCROLL UP", "scroll_down": "SCROLL DOWN",
		"up": "UP", "down": "DOWN", "left": "LEFT", "right": "RIGHT", "disabled": "DISABLED",
		"remove": "REMOVE", "online": "ONLINE", "away": "AWAY", "busy": "BUSY", "offline": "OFFLINE",
		"gold": "GOLD", "silver": "SILVER", "bronze": "BRONZE",
		"common": "COMMON", "uncommon": "UNCOMMON", "rare": "RARE", "epic": "EPIC", "legendary": "LEGENDARY",
		"hold": "HOLD TO CONFIRM", "del": "DEL", "end": "END", "lives": "+30 LIVES!",
		"retry": "RETRY", "stage_clear": "STAGE CLEAR!", "level_up": "LEVEL UP!", "achievement": "ACHIEVEMENT UNLOCKED",
		"combo": "COMBO!", "saving": "SAVING", "saved": "SAVED", "go": "GO!", "wait": "PLEASE WAIT",
		"prev": "PREVIOUS", "next": "NEXT", "back": "BACK", "forward": "FORWARD", "history": "HISTORY",
		"on": "ON", "off": "OFF", "clear": "CLEAR", "per_page": "PER PAGE", "range": "%d–%d OF %d", "found": "FOUND: %d",
		"unknown": "???", "sure": "ARE YOU SURE?", "paused": "PAUSED", "resume": "RESUME", "options": "OPTIONS", "quit": "QUIT",
		"next_step": "NEXT ▶",
		"months":    "JAN FEB MAR APR MAY JUN JUL AUG SEP OCT NOV DEC", "weekdays": "SU MO TU WE TH FR SA",
		"start": "START", "stop": "STOP", "lap": "LAP", "reset": "RESET", "pause": "PAUSE",
		"just_now": "JUST NOW", "min_ago": "%d MIN AGO", "hours_ago": "%d H AGO", "yesterday": "YESTERDAY", "days_ago": "%d DAYS AGO",
		"in_min": "IN %d MIN", "in_hours": "IN %d H", "tomorrow": "TOMORROW", "in_days": "IN %d DAYS",
		"spring": "SPRING", "summer": "SUMMER", "autumn": "AUTUMN", "winter": "WINTER", "minutes": "MIN", "seconds": "SEC",
	},
	Zh: {
		"ok": "确定", "cancel": "取消", "yes": "是", "no": "否", "close": "关闭",
		"continue": "继续？", "game_over": "游戏结束", "loading": "读取中",
		"search": "搜索", "empty": "空空如也", "minimize": "最小化", "resize": "调整大小",
		"scroll_up": "向上滚动", "scroll_down": "向下滚动",
		"up": "上", "down": "下", "left": "左", "right": "右", "disabled": "已禁用",
		"remove": "移除", "online": "在线", "away": "离开", "busy": "忙碌", "offline": "离线",
		"gold": "金", "silver": "银", "bronze": "铜",
		"common": "普通", "uncommon": "优秀", "rare": "稀有", "epic": "史诗", "legendary": "传说",
		"hold": "长按确认", "del": "删", "end": "完", "lives": "命数 +30！",
		"retry": "重试", "stage_clear": "过关！", "level_up": "升级！", "achievement": "成就解锁",
		"combo": "连击！", "saving": "存档中", "saved": "已存档", "go": "开始！", "wait": "请稍候",
		"prev": "上一页", "next": "下一页", "back": "后退", "forward": "前进", "history": "历史",
		"on": "开", "off": "关", "clear": "清除", "per_page": "每页", "range": "%d–%d / 共 %d", "found": "找到：%d",
		"unknown": "？？？", "sure": "确定吗？", "paused": "暂停", "resume": "继续游戏", "options": "设置", "quit": "退出",
		"next_step": "下一步 ▶",
		"months":    "一月 二月 三月 四月 五月 六月 七月 八月 九月 十月 十一月 十二月", "weekdays": "日 一 二 三 四 五 六",
		"start": "开始", "stop": "停止", "lap": "计次", "reset": "复位", "pause": "暂停",
		"just_now": "刚刚", "min_ago": "%d 分钟前", "hours_ago": "%d 小时前", "yesterday": "昨天", "days_ago": "%d 天前",
		"in_min": "%d 分钟后", "in_hours": "%d 小时后", "tomorrow": "明天", "in_days": "%d 天后",
		"spring": "春", "summer": "夏", "autumn": "秋", "winter": "冬", "minutes": "分", "seconds": "秒",
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
