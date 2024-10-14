package database

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

func AddRoom(roomid, name string, modeCode int) {
	_, err := addInfo.Exec(roomid, name, modeCode)
	if err != nil {
		fmt.Println("sql AddInfo err:", err)
	}
}

func AddBeatmap(bid string, playMode int, Star, cs, ar, od, hp float32) {
	_, err := addBeatmap.Exec(bid, playMode, Star, cs, ar, od, hp)
	var driverErr *mysql.MySQLError
	if errors.As(err, &driverErr) {
		if driverErr.Number == 1062 {
			return
		}
	}
	if err != nil {
		fmt.Println("sql AddBeatmap err:", err)
	}
}

func AddPickedMap(roomid, bid, user string) {
	_, err := addPickedMap.Exec(roomid, bid, user)
	if err != nil {
		fmt.Println("sql AddPickedMap err:", err)
	}
}

func AddHostRotate(roomid, user string) {
	var rev map[int]string
	GetHostRotate(roomid, &rev)
	for _, item := range rev {
		if item == user {
			return
		}
	}
	_, err := addHostRotate.Exec(roomid, user)
	if err != nil {
		fmt.Println("sql AddHostRotate err:", err)
	}
}
