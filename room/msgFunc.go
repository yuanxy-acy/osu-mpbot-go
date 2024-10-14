package room

import (
	"osu_mp_bot/database"
	"osu_mp_bot/util"
	"strconv"
	"strings"
)

func pick(roomid, bid, user string) {
	title, star, cs, ar, od, hp, ifGot := getMapDiff(bid)
	if !ifGot {
		sendMsg(roomid, "获取谱面信息失败")
	}
	if mapInAble(roomid, star, cs, ar, od, hp) {
		database.AddPickedMap(roomid, bid, user)
		sendMsg(roomid, "选图成功：[osu://b/"+bid+" "+title+"]")
		sendMsg(roomid, "目前还有"+strconv.Itoa(database.GetPickedMapCount(roomid))+"张谱面正在排队")
	} else {
		sendMsg(roomid, "[osu://b/"+bid+" "+title+"] 不符合房间要求")
	}
}
func poolPick(roomId, mod string) bool {
	bid := database.GetPoolMapBid(database.GetMapPool(roomId), strings.ToLower(mod))
	if bid == "" {
		sendMsg(roomId, mod+"不存在")
		return false
	} else {
		sendMsg(roomId, "!mp map "+bid)
		return true
	}
}

func nextMap(roomid string) {
	if database.GetPickedMapCount(roomid) != 0 {
		bid, user := database.GetFirstPickedMap(roomid)
		sendMsg(roomid, "!mp map "+bid)
		sendMsg(roomid, user+" picked")
	} else {
		sendMsg(roomid, "没有正在排队的谱面呢~ 我们来随机一张图吧~~")
		bid := getRandomMap(roomid)
		if bid == "" {
			sendMsg(roomid, "数据库中没有符合条件的曲子 .QAQ.")
			return
		}
		sendMsg(roomid, "!mp map "+bid)
		return
	}
	database.DelFirstPickedMap(roomid)
	sendMsg(roomid, "目前还有"+strconv.Itoa(database.GetPickedMapCount(roomid))+"张谱面正在排队")
}

func skipHost(roomId, user string) {
	if user == database.GetRoomHost(roomId) {
		nextHost(roomId)
	} else {
		database.DelHostRotate(user)
		database.AddHostRotate(roomId, user)
		sendMsg(roomId, "已将"+user+"置于队列末尾")
	}
}

func nextHost(roomid string) {
	host := database.GetRoomHost(roomid)
	if host != "" {
		database.AddHostRotate(roomid, host)
	}
	var rev map[int]string
	database.GetHostRotate(roomid, &rev)
	if len(rev) < 1 {
		return
	}
	host = rev[0]
	database.SetRoomHost(roomid, host)
	sendMsg(roomid, "!mp host "+host)
	var msg = ""
	sendMsg(roomid, "当前房主为：[https://osu.ppy.sh/users/"+host+" "+host+"]")
	for i := 1; i < len(rev); i++ {
		msg += "[https://osu.ppy.sh/users/" + rev[i] + " " + rev[i] + "] "
	}
	sendMsg(roomid, "等待队列："+msg)
	database.DelHostRotate(host)
}

func showHostRotate(roomid string) {
	var rev map[int]string
	database.GetHostRotate(roomid, &rev)
	var msg = ""
	host := database.GetRoomHost(roomid)
	sendMsg(roomid, "当前房主为：[https://osu.ppy.sh/users/"+host+" "+host+"]")
	for i := 0; i < len(rev); i++ {
		msg += "[https://osu.ppy.sh/users/" + rev[i] + " " + rev[i] + "] "
	}
	sendMsg(roomid, "等待队列："+msg)
}

func setCaptain(roomId, user, red, blue, old, team string, argLen int) {
	if old == "" {
		database.SetTeamCaptain(roomId, red, blue)
		checkCaptain(roomId, user, blue)
	} else if old == user {
		if argLen == 3 {
			sendMsg(roomId, "缺少<user>选项，请使用完整!set captain <red/blue> <user>指令转让队长")
		} else {
			database.SetTeamCaptain(roomId, red, blue)
			sendMsg(roomId, "转让队长成功")
		}
	} else {
		switch team {
		case "red":
			sendMsg(roomId, "红队已存在队长,请由队长输入指令")
		case "blue":
			sendMsg(roomId, "蓝队已存在队长,请由队长输入指令")
		}
		return
	}
	switch team {
	case "red":
		sendMsg(roomId, "!mp team "+red+" red")
	case "blue":
		sendMsg(roomId, "!mp team "+blue+" blue")
	}
}

func checkCaptain(roomId, red, blue string) {
	if database.GetMatchState(roomId) == "setting captain" && blue != "" && red != "" {
		sendMsg(roomId, "双方队长已选出，接下来进行红蓝轮流pick图游玩")
		sendMsg(roomId, "首先由红队队长进行pick")
		database.SetMatchState(roomId, "red pick")
	}
}

func setRules(roomid string, args []string) {
	var index = 0
	minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar := database.GetRules(roomid)
	for index < len(args) {
		switch strings.ToLower(args[index]) {
		case "mincs":
			minCS = util.StringToFloat(args[index+1])
			checkMinRule(&minCS, &maxCS)
		case "maxcs":
			maxCS = util.StringToFloat(args[index+1])
			checkMaxRule(&minCS, &maxCS)
		case "minar":
			minAR = util.StringToFloat(args[index+1])
			checkMinRule(&minAR, &maxAR)
		case "maxar":
			maxAR = util.StringToFloat(args[index+1])
			checkMaxRule(&minAR, &maxAR)
		case "minod":
			minOD = util.StringToFloat(args[index+1])
			checkMinRule(&minOD, &maxOD)
		case "maxod":
			maxOD = util.StringToFloat(args[index+1])
			checkMaxRule(&minOD, &maxOD)
		case "minhp":
			minHP = util.StringToFloat(args[index+1])
			checkMinRule(&minHP, &maxHP)
		case "maxhp":
			maxHP = util.StringToFloat(args[index+1])
			checkMaxRule(&minHP, &maxHP)
		case "minstar":
			minStar = util.StringToFloat(args[index+1])
			checkMinRule(&minStar, &maxStar)
		case "maxstar":
			maxStar = util.StringToFloat(args[index+1])
			checkMaxRule(&minStar, &maxStar)
		}
		index += 2
	}
	database.SetRules(roomid, minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar)
	sendMsg(roomid, "设置成功")
}

func checkMinRule(minVal, maxVal *float32) {
	if *minVal > *maxVal {
		*minVal = *maxVal
	}
}

func checkMaxRule(minVal, maxVal *float32) {
	if *minVal > *maxVal {
		*maxVal = *minVal
	}
}

func sendMapPoolSetting(roomId string, poolId int) {
	if poolId == -1 {
		sendMsg(roomId, "当前未选定图池，自由轮流选图")
		database.SetMatchState(roomId, "free pick")
	} else {
		name, avgStar, nm, hd, hr, dt, fm, ex, tb := database.GetMapPoolInfo(poolId)
		var m = ""
		if nm != 0 {
			m += "NM: " + strconv.Itoa(nm)
		}
		if hd != 0 {
			m += "HD: " + strconv.Itoa(hd)
		}
		if hr != 0 {
			m += "HR: " + strconv.Itoa(hr)
		}
		if dt != 0 {
			m += "DT: " + strconv.Itoa(dt)
		}
		if fm != 0 {
			m += "FM: " + strconv.Itoa(fm)
		}
		if ex != 0 {
			m += "EX: " + strconv.Itoa(ex)
		}
		if tb != 0 {
			m += "TB: " + strconv.Itoa(tb)
		}
		sendMsg(roomId, "已选图池："+name+" 平均星数"+util.FloatToString(avgStar)+"图池配置 "+m)
	}
}
