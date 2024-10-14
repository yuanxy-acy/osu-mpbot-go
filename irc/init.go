package irc

import (
	"fmt"
	"net"
	"osu_mp_bot/config"
	"osu_mp_bot/log"
	"osu_mp_bot/osutool"
)

type Config struct {
	host string
	port string
	pass string
	nick string
}

type RoomInfo struct {
	user     string
	roomid   string
	password string
	name     string
}

var conf Config
var conn net.Conn
var makeQueue chan RoomInfo
var createdQueue chan RoomInfo

func init() {
	makeQueue = make(chan RoomInfo, 1)
	createdQueue = make(chan RoomInfo, 1)
	osutool.IrcSengFunc = Speak
	conf = Config{
		host: "irc.ppy.sh",
		port: "6667",
		pass: config.IrcPass,
		nick: config.IrcNick,
	}
	connect()
}

func Speak(target string, text string) {
	say("PRIVMSG", target, text)
	log.BotMsgLog(target, text)
}

func Join(channel string) {
	send("JOIN " + channel)
}

func MakeRoom(user, name, password string) RoomInfo {
	info := RoomInfo{
		user:     user,
		roomid:   "",
		password: password,
		name:     name,
	}
	makeQueue <- info
	Speak("BanchoBot", "!mp make "+name)
	return <-createdQueue
}

func send(msg string) {
	if !ConnFlag {
		connect()
	}
	buf := []byte(msg + "\r\n")
	_, err := conn.Write(buf)
	if err != nil {
		fmt.Println("send err")
		fmt.Println("err:", err)
		return
	}
}

func say(kind string, target string, text string) {
	send(kind + " " + target + " :" + text)
}
