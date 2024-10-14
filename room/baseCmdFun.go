package room

import (
	"osu_mp_bot/database"
	"osu_mp_bot/irc"
	"osu_mp_bot/osutool"
	"strings"
)

func baseCmdFun(roomId, user string, args []string) {
	switch strings.ToLower(args[0]) {
	case "help", "h":
		irc.Speak(user, "使用!r或!re来获取上一局游戏成绩")
		irc.Speak(user, "使用!s或!score来获取当前谱面的最好成绩")
		irc.Speak(user, "使用!i或!info来获取当前谱面信息")
		irc.Speak(user, "使用!dl或!download来获取当前谱面下载链接")
		switch database.GetMpMode(roomId) {
		case -1:
			irc.Speak(user, "使用!start [<mode>] 开始使用bot游玩")
		case 0:
			irc.Speak(user, "使用!pick <bid>来指定想玩谱面")
			irc.Speak(user, "使用!skip来跳过当前谱面")
		case 1, 3:
			irc.Speak(user, "使用!skip来跳过自己房主轮次")
		}
	case "start":
		if len(args) < 2 {
			setMpMode(roomId, user, "0")
		} else {
			setMpMode(roomId, user, args[1])
		}
		sendMsg(roomId, "!mp settings")
	case "stop":
		database.SetMpMode(roomId, -1)
		sendMsg(roomId, "停止使用bot游玩")
	case "r", "re":
		osutool.SendUserRecent("#mp_"+roomId, user)
	case "s", "score":
		var bid string
		if len(args) < 2 {
			bid = database.GetNowBid(roomId)
		} else {
			bid = args[1]
		}
		osutool.SendUserScore("#mp_"+roomId, user, bid)
	case "i", "info":
		if len(args) > 1 {
			sendMapInfo(roomId, args[2])
		} else {
			sendMapInfo(roomId, database.GetNowBid(roomId))
		}
	case "dl", "download":
		sendMapLink(roomId, database.GetNowBid(roomId))
	}
}
