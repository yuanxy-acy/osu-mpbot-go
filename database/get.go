package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func GetInfo(id string) string {
	var name string
	err := getInfo.QueryRow(id).Scan(&name)
	if err != nil {
		fmt.Println("GetInfo sql err:", err)
		return ""
	}
	return name
}

func GetAllActiveRoom(rev *map[int]string) {
	rows, err := getAllActiveRoom.Query()
	if err != nil {
		fmt.Println("GetAllActiveRoom sql err:", err)
		return
	}
	rowsToStringMap(rows, rev)
}

func GetNowBid(roomid string) string {
	var bid string
	err := getNowBid.QueryRow(roomid).Scan(&bid)
	if err != nil {
		fmt.Println("GetNowBid sql err:", err)
		return ""
	}
	return bid
}

func GetPlayMode(roomid string) int {
	var modeCode int
	err := getPlayMode.QueryRow(roomid).Scan(&modeCode)
	if err != nil {
		fmt.Println("GetPlayMode sql err:", err)
		return 0
	}
	return modeCode
}

func GetMpMode(roomid string) int {
	var modeCode int
	err := getMpMode.QueryRow(roomid).Scan(&modeCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			AddRoom(roomid, "", -1)
		}
		fmt.Println("GetMpMode sql err:", err)
		return -1
	}
	return modeCode
}

func GetRules(roomId string) (float32, float32, float32, float32, float32, float32, float32, float32, float32, float32) {
	var minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar float32
	err := getRules.QueryRow(roomId).Scan(&minCS, &maxCS, &minAR, &maxAR, &minOD, &maxOD, &minHP, &maxHP, &minStar, &maxStar)
	if err != nil {
		fmt.Println("GetRules sql err:", err)
		return 0, 11, 0, 11, 0, 11, 0, 11, 0, 11
	}
	return minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar
}

func GetBeatmapIdByRules(playMode int, minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar float32, rev *map[int]string) {
	rows, err := getBeatmapIdByRules.Query(playMode, minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar)
	if err != nil {
		fmt.Println("GetBeatmapIdByRules sql err:", err)
		return
	}
	rowsToStringMap(rows, rev)
}

func GetFirstPickedMap(roomid string) (string, string) {
	var bid, user string
	err := getFirstPickedMap.QueryRow(roomid).Scan(&bid, &user)
	if err != nil {
		fmt.Println("GetFirstPickedMap sql err:", err)
		return "", ""
	}
	return bid, user
}

func GetPickedMapCount(roomid string) int {
	var count int
	err := getPickedMapCount.QueryRow(roomid).Scan(&count)
	if err != nil {
		fmt.Println("GetPickedMapCount sql err:", err)
		return 0
	}
	return count
}

func GetHostRotate(roomid string, rev *map[int]string) {
	rows, err := getHostRotate.Query(roomid)
	if err != nil {
		fmt.Println("GetHostRotate sql err:", err)
		return
	}
	rowsToStringMap(rows, rev)
}

func GetRoomHost(roomid string) string {
	var user string
	err := getRoomHost.QueryRow(roomid).Scan(&user)
	if err != nil {
		fmt.Println("GetRoomHost sql err:", err)
		return ""
	}
	return user
}

func GetTeamCaptain(roomid string) (string, string) {
	var red, blue string
	err := getTeamCaptain.QueryRow(roomid).Scan(&red, &blue)
	if err != nil {
		fmt.Println("GetTeamCaptain sql err:", err)
		return "", ""
	}
	return red, blue
}

func GetMatchState(roomid string) string {
	var state string
	err := getMatchState.QueryRow(roomid).Scan(&state)
	if err != nil {
		fmt.Println("GetMatchState sql err:", err)
		return ""
	}
	return state
}

func GetMapPool(roomid string) int {
	var poolId int
	err := getMapPool.QueryRow(roomid).Scan(&poolId)
	if err != nil {
		fmt.Println("GetMapPool sql err:", err)
		return -1
	}
	return poolId
}

func GetMapPoolIdLike(keyWord string, rev *map[int][]string) {
	rows, err := getMapPoolIdLike.Query(keyWord)
	if err != nil {
		fmt.Println("GetMapPoolIdLike sql err:", err)
		return
	}
	var poolId, name, avgStar string
	*rev = make(map[int][]string)
	count := 0
	for rows.Next() {
		err := rows.Scan(&poolId, &name, &avgStar)
		if err != nil {
			fmt.Println("sql rows scan err:", err)
		}
		(*rev)[count] = make([]string, 3)
		(*rev)[count][0] = poolId
		(*rev)[count][1] = name
		(*rev)[count][2] = avgStar
		count++
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("sql rows err:", err)
	}
	err = rows.Close()
	if err != nil {
		fmt.Println("sql rows close err:", err)
	}
}

func GetMapPoolInfo(poolId int) (string, float32, int, int, int, int, int, int, int) {
	var name string
	var avgStar float32
	var nm, hd, hr, dt, fm, ex, tb int
	err := getMapPoolInfo.QueryRow(poolId).Scan(&name, &avgStar, &nm, &hd, &hr, &dt, &fm, &ex, &tb)
	if err != nil {
		fmt.Println("GetMapPoolInfo sql err:", err)
		return "", 0, 0, 0, 0, 0, 0, 0, 0
	}
	return name, avgStar, nm, hd, hr, dt, fm, ex, tb
}

func GetPoolMapBid(poolId int, mod string) string {
	var bid string
	err := getPoolMapBid.QueryRow(poolId, mod).Scan(&bid)
	if err != nil {
		fmt.Println("GetPoolMapBid sql err:", err)
		return ""
	}
	return bid
}
