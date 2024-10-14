package irc

import (
	"fmt"
	"net"
	"osu_mp_bot/database"
	"osu_mp_bot/log"
	"strings"
	"time"
)

var ConnFlag = false
var RoomMsgFunc func(string, string, []string)

func connect() {
	var err error
	conn, err = net.Dial("tcp", conf.host+":"+conf.port)
	if err != nil {
		fmt.Println("connect err")
		fmt.Println("err:", err)
		return
	}
	ConnFlag = true
	send("CAP LS 302")
	send("PASS " + conf.pass)
	send("NICK " + conf.nick)
	send("USER " + conf.nick + "bot 0 * :" + conf.nick)
	var t time.Time
	go func() {
		for ConnFlag {
			t = <-time.NewTimer(60 * time.Second).C
			send("PING " + t.String())
		}
	}()
	go func() {
		for ConnFlag {
			buf := make([]byte, 1024)
			cnt, err := conn.Read(buf)
			if err != nil {
				fmt.Println("rev err")
				fmt.Println("err:", err)
				ConnFlag = false
				break
			}
			msg := string(buf[:cnt])
			msgs := strings.Split(msg, "\n")
			if msgs[len(msgs)-1] != "" {
				cnt, err = conn.Read(buf)
				if err != nil {
					fmt.Println("rev err")
					fmt.Println("err:", err)
					ConnFlag = false
					err := conn.Close()
					if err != nil {
						return
					}
					return
				}
				msg = msg + string(buf[:cnt])
				msgs = strings.Split(msg, "\n")
			}
			for _, item := range msgs {
				if item == "" {
					break
				}
				args := strings.Split(item[1:], " ")
				if len(args) < 2 {
					fmt.Print(args[0])
					break
				}
				switch args[1] {
				case "PRIVMSG":
					user := strings.Split(args[0], "!")[0]
					fmt.Println("get msg from " + user + " to " + args[2] + " " + strings.Join(args[3:], " "))
					if args[2] != conf.nick {
						if args[2][0:4] == "#mp_" {
							go RoomMsgFunc(args[2][4:], user, args[3:])
							log.ChannelMsgLog(args[2], user, strings.Join(args[3:], " "))
						}
					} else {
						go chat(user, args[3:])
						log.UserMsgLog(user, strings.Join(args[3:], " "))
					}
				case "JOIN":
					fmt.Println("joined channel " + args[2][1:])
				case "403":
					fmt.Println(strings.Join(args[4:], " ")[1:])
					database.DelRoom(args[7][4:])
					database.SetMatchState(args[7][4:], "closed")
				case "PART":
				case "001": //Welcome to the osu!Bancho.
					fmt.Println(strings.Join(args[3:], " ")[1:])
				case "QUIT", "PONG", "MODE", "332", "333", "353", "366", "375", "372", "376":
				default:
					switch args[0] {
					case "PING":
						fmt.Println("get ping form " + strings.Join(args[1:], " "))
						send("PONG " + strings.Join(args[1:], " "))
					default:
						fmt.Println(strings.Join(args, " "))
					}
				}
			}
		}
		err := conn.Close()
		if err != nil {
			return
		}
	}()
}
