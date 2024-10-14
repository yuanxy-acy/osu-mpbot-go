package room

import (
	"osu_mp_bot/database"
	"osu_mp_bot/webapi"
	"strings"
)

func setMpMode(roomId, user, pm string) {
	switch strings.ToLower(pm) {
	case "0", "p&r", "r&p", "pick&random", "random&pick":
		database.SetMpMode(roomId, 0)
		sendMsg(roomId, "当前游玩模式为：pick&random map")
		sendMsg(roomId, "!mp set 0 0 16")
		sendMsg(roomId, "!mp mods freemod")
		database.SetRoomHost(roomId, "")
		sendMsg(roomId, "!mp clearhost")
	case "1", "ahr", "auto", "host", "rotate":
		database.SetMpMode(roomId, 1)
		sendMsg(roomId, "当前游玩模式为：auto host rotate")
		sendMsg(roomId, "!mp set 0 0 16")
		sendMsg(roomId, "!mp mods freemod")
		setHost(roomId, user)
	case "2", "team", "rvb", "bvr", "teamvs":
		database.SetMpMode(roomId, 2)
		sendMsg(roomId, "当前游玩模式为：team vs")
		sendMsg(roomId, "!mp set 2 3 16")
		sendMsg(roomId, "!mp mods freemod")
		poolId := database.GetMapPool(roomId)
		sendMapPoolSetting(roomId, poolId)
		if poolId != -1 {
			sendMsg(roomId, "!mp clearhost")
			sendMsg(roomId, "请双方指定队长，使用!set captain <red/blue> <user>指定或转让，当不存在队长时可省略<user>选项指定自己为队长")
			database.SetMatchState(roomId, "setting captain")
		} else {
			setHost(roomId, user)
		}
	case "3", "relay", "co-op":
		database.SetMpMode(roomId, 3)
		sendMsg(roomId, "当前游戏模式为：relay mode - auto host rotate")
		sendMsg(roomId, "!mp set 1 0 16")
		sendMsg(roomId, "!mp mods none")
		setHost(roomId, user)
	}
}

func setHost(roomId, user string) {
	host := database.GetRoomHost(roomId)
	if host == "" {
		info := webapi.GetUserInfo(user)
		database.SetRoomHost(roomId, info.Name)
		database.DelHostRotate(info.Name)
		sendMsg(roomId, "!mp host "+info.Name)
	} else {
		sendMsg(roomId, "!mp host "+host)
	}
	showHostRotate(roomId)
}
