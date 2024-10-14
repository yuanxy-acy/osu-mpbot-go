package room

import (
	"osu_mp_bot/config"
	"osu_mp_bot/database"
	"osu_mp_bot/util"
	"strings"
)

func banchoMsg(roomId string, args []string) {
	msg := strings.Join(args, " ")
	switch msg {
	case "Aborted the match":
	case "The match has started!":
		sendMsg(roomId, "!mp settings")
	case "All players are ready":
		switch database.GetMpMode(roomId) {
		case -1:
			return
		case 2:
			state := database.GetMatchState(roomId)
			if state != "red pick play" && state != "blue pick play" {
				return
			}
		}
		sendMsg(roomId, "Good Luck Have Fun!!!")
		sendMsg(roomId, "!mp start 7")
	case "The match has finished!":
		switch database.GetMpMode(roomId) {
		case 0:
			nextMap(roomId)
		case 1, 3:
			nextHost(roomId)
		case 2:
			switch database.GetMatchState(roomId) {
			case "free pick":
				nextHost(roomId)
			case "red picked":
				sendMsg(roomId, "蓝方选图")
				database.SetMatchState(roomId, "blue pick")
			case "blue picked":
				sendMsg(roomId, "红方选图")
				database.SetMatchState(roomId, "red pick")
			}
		}
	case "Cleared match host":
		database.SetRoomHost(roomId, "")
	case "Closed the match":
		database.SetMatchState(roomId, "closed")
		database.DelRoom(roomId)
	default:
		if len(args) < 2 {
			if strings.ToLower(args[0]) == strings.ToLower(config.IrcNick) {
				sendMsg(roomId, "!mp password 114514")
				sendMsg(roomId, "!mp mods freemod")
				sendMsg(roomId, "!mp settings")
			}
			return
		}
		if util.StringMatch(msg, "joined in slot") {
			database.AddHostRotate(roomId, args[0])
			switch database.GetMpMode(roomId) {
			case -1:
				return
			case 1, 3:
				if database.GetRoomHost(roomId) == "" {
					nextHost(roomId)
				} else {
					showHostRotate(roomId)
				}
			case 2:
				if database.GetMatchState(roomId) == "free pick" {
					showHostRotate(roomId)
				}
			}
			userName := strings.Split(msg, " joined")[0]
			sendMsg(roomId, "你来啦~~ [https://osu.ppy.sh/users/"+strings.Replace(userName, " ", "_", -1)+" "+userName+"]  快快找个地方坐下来吧~~~ .Ov≤.")
			sendMsg(roomId, "可以输入 !h 来获取帮助哦")
			return
		}
		if util.StringMatch(msg, "left the game") {
			userName := strings.Split(msg, " left")[0]
			database.DelHostRotate(userName)
			switch database.GetMpMode(roomId) {
			case 1, 3:
				if userName == database.GetRoomHost(roomId) {
					database.SetRoomHost(roomId, "")
					nextHost(roomId)
				}
			case 2:
				if database.GetMatchState(roomId) == "free pick" {
					if userName == database.GetRoomHost(roomId) {
						nextHost(roomId)
					}
				}
			}
			return
		}
		if util.StringMatch(msg, "became the host") {
			host := strings.Split(msg, " became the host")[0]
			switch database.GetPlayMode(roomId) {
			case 0:
				database.SetRoomHost(roomId, host)
			case 1:
				if host != database.GetRoomHost(roomId) {
					nextHost(roomId)
				}
			}
			return
		}
		if util.StringMatch(msg, "move to slot") {
			return
		}
		if util.StringMatch(msg, "finished playing ") {
			return
		}
		if args[1] == "beatmap" {
			bid := strings.Split(args[3], "/")[4]
			database.SetNowBid(roomId, bid)
			if database.GetMpMode(roomId) == -1 {
				return
			}
			sendMapInfo(roomId, bid)
			sendMapLink(roomId, bid)
			return
		}
		switch args[0] {
		case "Slot":
			player := strings.Split(strings.Split(msg, "https://osu.ppy.sh/u/")[1], " ")[1]
			switch database.GetMpMode(roomId) {
			case 1, 3:
				if player != database.GetRoomHost(roomId) {
					database.AddHostRotate(roomId, player)
				}
			case 2:
				if database.GetMatchState(roomId) == "free pick" && player != database.GetRoomHost(roomId) {
					database.AddHostRotate(roomId, player)
				}
			}
		case "Beatmap":
			bid := strings.Split(args[len(args)-1], "https://osu.ppy.sh/b/")[1]
			bid = bid[:len(bid)-1]
			title, star, cs, ar, od, hp, ifGot := getMapDiff(bid)
			if !ifGot {
				sendMsg(roomId, "获取谱面信息失败")
			} else {
				switch database.GetMpMode(roomId) {
				case 0:
					sendMsg(roomId, "!mp map "+database.GetNowBid(roomId))
					pick(roomId, args[1], database.GetRoomHost(roomId))
					sendMsg(roomId, "!mp clearhost")
					return
				case 1, 3:
					if !mapInAble(roomId, star, cs, ar, od, hp) {
						sendMsg(roomId, "!mp map "+database.GetNowBid(roomId))
						sendMsg(roomId, "[osu://b/"+bid+" "+title+"] 不符合房间要求")
						return
					}
				case 2:
					if database.GetMatchState(roomId) == "free pick" && !mapInAble(roomId, star, cs, ar, od, hp) {
						sendMsg(roomId, "!mp map "+database.GetNowBid(roomId))
						sendMsg(roomId, "[osu://b/"+bid+" "+title+"] 不符合房间要求")
						return
					}
				}
			}
			database.SetNowBid(roomId, bid)
			if database.GetMpMode(roomId) == -1 {
				return
			}
			sendMapInfo(roomId, bid)
			sendMapLink(roomId, bid)
		case "Beatmap:":
			bid := strings.Split(args[1], "https://osu.ppy.sh/b/")[1]
			database.SetNowBid(roomId, bid)
		case "Room":
			database.SetRoomName(roomId, strings.Split(strings.Join(args[2:], " "), ",")[0])
		}
	}
}
