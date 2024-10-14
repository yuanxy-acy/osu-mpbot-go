package irc

import (
	"osu_mp_bot/database"
	"osu_mp_bot/osutool"
	"strings"
)

func chat(user string, args []string) {
	if len(args[0]) < 3 {
		return
	}
	args[0] = strings.Replace(args[0], "！", "!", 1)
	if args[0][1] == '!' {
		args[0] = args[0][2:]
		switch strings.ToLower(args[0]) {
		case "h", "help":
			Speak(user, "!make <roomName> 创建房间")
			Speak(user, "!r或!re 来获取上一局游戏成绩")
			Speak(user, "!s <bid> 获取bid对应谱面成绩")
		case "make":
			var info RoomInfo
			if len(args) < 2 {
				info = MakeRoom(user, conf.nick+"`s game", "114514")
			} else if len(args) < 3 {
				info = MakeRoom(user, args[1], "114514")
			} else {
				info = MakeRoom(user, args[1], args[2])
			}
			Speak(user, "已创建：[osu://mp/"+info.roomid+"/"+info.password+" "+info.name+"] 密码："+info.password)
		case "r", "re":
			osutool.SendUserRecent(user, user)
		case "s", "score":
			if len(args) < 2 {
				Speak(user, "缺少参数<bid>")
				return
			}
			osutool.SendUserScore(user, user, args[1])
		}
	} else if user == "BanchoBot" {
		args[0] = args[0][1:]
		switch args[0] {
		case "Created":
			info := <-makeQueue
			info.roomid = args[4][22:]
			createdQueue <- info
			Speak("#mp_"+args[4][22:], "!mp password "+info.password)
			Speak("#mp_"+args[4][22:], "!mp mods freemod")
			database.AddRoom(args[4][22:], strings.Join(args[5:], " "), 0)
		}
	}
}
