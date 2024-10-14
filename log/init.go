package log

import (
	"fmt"
	"os"
	"osu_mp_bot/config"
	"time"
)

type log struct {
	typeCode int
	channel  string
	name     string
	msg      string
}

var logPath string
var logQueue chan log

func init() {
	logQueue = make(chan log, 9)
	go func() {
		for {
			logPath = "log/"
			l := <-logQueue
			switch l.typeCode {
			case 0:
				logPath += "channel/"
			case 1:
				logPath += "user/"
			}
			openFile, err := os.OpenFile(logPath+l.channel+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				fmt.Println("日志文件打开失败")
				fmt.Println(err)
				return
			}
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			_, err = openFile.WriteString(timeStr + " " + l.name + " " + l.msg + "\n")
			if err != nil {
				fmt.Println("日志文件写入失败")
				fmt.Println(err)
				return
			}
			err = openFile.Close()
			if err != nil {
				fmt.Println("日志文件关闭失败")
				return
			}
		}
	}()
}

func ChannelMsgLog(channel, user, msg string) {
	l := log{
		typeCode: 0,
		channel:  channel,
		name:     user,
		msg:      msg,
	}
	logQueue <- l
}

func UserMsgLog(user, msg string) {
	l := log{
		typeCode: 1,
		channel:  user,
		name:     user,
		msg:      msg,
	}
	logQueue <- l
}

func BotMsgLog(target, msg string) {
	var tc int
	if target[0] == '#' {
		tc = 0
	} else {
		tc = 1
	}
	l := log{
		typeCode: tc,
		channel:  target,
		name:     config.IrcNick,
		msg:      msg,
	}
	logQueue <- l
}
