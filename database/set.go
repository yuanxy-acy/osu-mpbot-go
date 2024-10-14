package database

import (
	"fmt"
)

func SetNowBid(roomid, bid string) {
	_, err := setNowBid.Exec(bid, roomid)
	if err != nil {
		fmt.Println("SetNowBid sql err:", err)
	}
}

func SetPlayMode(roomid string, modeCode int) {
	_, err := setPlayMode.Exec(modeCode, roomid)
	if err != nil {
		fmt.Println("SetPlayMode sql err:", err)
	}
}

func SetMpMode(roomid string, modeCode int) {
	_, err := setMpMode.Exec(modeCode, roomid)
	if err != nil {
		fmt.Println("SetMpMode sql err:", err)
	}
}

func SetRules(roomid string, minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar float32) {
	_, err := setRules.Exec(minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar, roomid)
	if err != nil {
		fmt.Println("SetRules sql err:", err)
	}
}

func SetRoomName(roomid, name string) {
	_, err := setRoomName.Exec(name, roomid)
	if err != nil {
		fmt.Println("SetPlayMode sql err:", err)
	}
}

func SetRoomHost(roomid, host string) {
	_, err := setRoomHost.Exec(host, roomid)
	if err != nil {
		fmt.Println("SetPlayMode sql err:", err)
	}
}

func SetTeamCaptain(roomid, red, blue string) {
	_, err := setTeamCaptain.Exec(red, blue, roomid)
	if err != nil {
		fmt.Println("SetTeamCaptain sql err:", err)
	}
}

func SetMatchState(roomid, state string) {
	_, err := setMatchState.Exec(state, roomid)
	if err != nil {
		fmt.Println("SetMatchState sql err:", err)
	}
}

func SetMapPool(roomid string, poolId int) {
	_, err := setMapPool.Exec(poolId, roomid)
	if err != nil {
		fmt.Println("SetMapPool sql err:", err)
	}
}
