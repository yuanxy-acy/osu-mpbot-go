package main

import (
	"bufio"
	"os"
	"osu_mp_bot/config"
	"osu_mp_bot/database"
	"osu_mp_bot/irc"
	"osu_mp_bot/room"
	"strings"
)

var command = make(chan string)

func main() {
	run := true
	for run {
		reader := bufio.NewReader(os.Stdin)
		b, _, _ := reader.ReadLine()
		args := strings.Split(string(b), " ")
		switch args[0] {
		case "quit", "q":
			irc.ConnFlag = false
			run = false
			database.Close()
		case "say":
			if len(args) < 3 {
				return
			}
			irc.Speak(args[1], strings.Join(args[2:], " "))
		case "make":
			room.MakeRoom(config.IrcNick + "`s game")
		case "join":
			if len(args) < 2 {
				return
			}
			irc.Join("#mp_" + args[1])
		}
	}
}
