package icon

import "git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"

type Named struct {
	Name   string
	Sprite *pixel.Sprite
}

// All lists every icon by name, for galleries and tests.
var All = []Named{
	{"check", Check}, {"close", Close}, {"plus", Plus}, {"minus", Minus},
	{"arrow-left", ArrowLeft}, {"arrow-right", ArrowRight}, {"arrow-up", ArrowUp}, {"arrow-down", ArrowDown},
	{"chevron-left", ChevronLeft}, {"chevron-right", ChevronRight}, {"chevron-up", ChevronUp}, {"chevron-down", ChevronDown},
	{"triangle-right", TriangleRight}, {"triangle-down", TriangleDown}, {"menu", Menu}, {"more", More},
	{"search", Search}, {"settings", Settings}, {"edit", Edit}, {"trash", Trash},
	{"copy", Copy}, {"clipboard", Clipboard}, {"save", Save}, {"download", Download},
	{"upload", Upload}, {"refresh", Refresh}, {"link", Link}, {"filter", Filter},
	{"sort", Sort}, {"send", Send}, {"power", Power}, {"grid", Grid}, {"list", List},
	{"home", Home}, {"user", User}, {"users", Users}, {"bell", Bell},
	{"mail", Mail}, {"mail-open", MailOpen}, {"chat", Chat}, {"folder", Folder},
	{"folder-open", FolderOpen}, {"file", File}, {"file-text", FileText}, {"image", Image},
	{"music", Music}, {"video", Video}, {"code", Code}, {"terminal", Terminal},
	{"bug", Bug}, {"lock", Lock}, {"unlock", Unlock}, {"key", Key},
	{"eye", Eye}, {"eye-off", EyeOff}, {"calendar", Calendar}, {"clock", Clock},
	{"hourglass", Hourglass}, {"camera", Camera}, {"mic", Mic}, {"phone", Phone},
	{"map", Map}, {"flag", Flag}, {"pin", Pin}, {"globe", Globe},
	{"book", Book}, {"cart", Cart}, {"tag", Tag}, {"gift", Gift},
	{"robot", Robot}, {"branch", Branch}, {"wifi", Wifi}, {"battery", Battery}, {"gamepad", Gamepad},
	{"info", Info}, {"warning", Warning}, {"error", Error}, {"question", Question},
	{"heart", Heart}, {"heart-empty", HeartEmpty}, {"star", Star}, {"star-empty", StarEmpty},
	{"sparkle", Sparkle}, {"volume", Volume}, {"mute", Mute}, {"play", Play},
	{"pause", Pause}, {"stop", Stop}, {"next", Next}, {"prev", Prev},
	{"coin", Coin}, {"gem", Gem}, {"potion", Potion}, {"sword", Sword},
	{"shield", Shield}, {"chest", Chest}, {"crown", Crown}, {"trophy", Trophy},
	{"skull", Skull}, {"scroll", Scroll}, {"flame", Flame}, {"drop", Drop},
	{"leaf", Leaf}, {"flower", Flower}, {"sun", Sun}, {"moon", Moon},
	{"cloud", Cloud}, {"snowflake", Snowflake}, {"bolt", Bolt}, {"palette", Palette},
}
