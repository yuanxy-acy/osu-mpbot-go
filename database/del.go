package database

import "fmt"

func DelFirstPickedMap(roomid string) {
	_, err := delFirstPickedMap.Exec(roomid)
	if err != nil {
		fmt.Println("DelFirstPickedMap sql err:", err)
	}
}

func DelHostRotate(user string) {
	_, err := delHostRotate.Exec(user)
	if err != nil {
		fmt.Println("DelHostRotate sql err:", err)
	}
}

func DelRoom(roomid string) {
	/*_, err := delInfo.Exec(roomid)
	if err != nil {
		fmt.Println("DelInfo sql err:", err)
	}
	_, err = delRules.Exec(roomid)
	if err != nil {
		fmt.Println("DelRules sql err:", err)
	}
	_, err = delRoomTeam.Exec(roomid)
	if err != nil {
		fmt.Println("DelRoomTeam sql err:", err)
	}*/
	_, err := delAllPickedMap.Exec(roomid)
	if err != nil {
		fmt.Println("DelPickedMap sql err:", err)
	}
	_, err = delAllHostRotate.Exec(roomid)
	if err != nil {
		fmt.Println("DelAllHostRotate sql err:", err)
	}
}
