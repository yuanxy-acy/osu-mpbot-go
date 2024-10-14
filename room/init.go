package room

import (
	"osu_mp_bot/config"
	"osu_mp_bot/database"
	"osu_mp_bot/irc"
)

func init() {
	irc.RoomMsgFunc = msgFun
	var rev map[int]string
	database.GetAllActiveRoom(&rev)
	for i := 0; i < len(rev); i++ {
		irc.Join("#mp_" + rev[i])
	}
}

func MakeRoom(name string) {
	irc.MakeRoom(config.IrcNick, name, "114514")
}

func sendMsg(roomId string, msg string) {
	irc.Speak("#mp_"+roomId, msg)
}
