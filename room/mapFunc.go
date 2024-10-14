package room

import (
	"fmt"
	"math/rand"
	"osu_mp_bot/database"
	"osu_mp_bot/util"
	"osu_mp_bot/webapi"
	"strconv"
)

func sendMapInfo(roomId string, bid string) {
	var rev []map[string]string
	webapi.GetMapInfo(bid, &rev)
	if len(rev) == 0 {
		fmt.Println("获取谱面信息失败")
		sendMsg(roomId, "获取谱面信息失败")
		fmt.Println(rev)
		return
	}
	s, err := strconv.ParseInt(rev[0]["hit_length"], 10, 64)
	var time string
	if err != nil {
		fmt.Println("计算谱面时间失败 ", err)
		fmt.Println(rev[0]["hit_length"])
		time = ""
	} else {
		m := s / 60
		h := m / 60
		s -= m * 60
		m -= h * 60
		if h != 0 {
			time += strconv.FormatInt(h, 10) + ":"
		}
		if m < 10 {
			time += "0" + strconv.FormatInt(m, 10)
		} else {
			time += strconv.FormatInt(m, 10)
		}
		if s < 10 {
			time += ":0" + strconv.FormatInt(s, 10)
		} else {
			time += ":" + strconv.FormatInt(s, 10)
		}
	}
	var title string
	if rev[0]["title_unicode"] == "" {
		title = rev[0]["title"]
	} else {
		title = rev[0]["title_unicode"]
	}
	var approved string
	switch rev[0]["approved"] {
	case "-2":
		approved = "graveyard"
	case "-1":
		approved = "WIP"
	case "0":
		approved = "pending"
	case "1":
		approved = "ranked"
	case "2":
		approved = "approved"
	case "3":
		approved = "qualified"
	case "4":
		approved = "loved"
	}
	var mode string
	switch rev[0]["mode"] {
	case "0":
		mode = "osu!"
	case "1":
		mode = "taiko"
	case "2":
		mode = "catch"
	case "3":
		mode = "mania"
	}
	sendMsg(roomId, approved+" "+mode+" [osu://b/"+rev[0]["beatmap_id"]+" "+title+"] rating: "+rev[0]["difficultyrating"]+" time: "+time+" bpm: "+rev[0]["bpm"]+" CS: "+rev[0]["diff_size"]+" AR: "+rev[0]["diff_approach"]+" OD: "+rev[0]["diff_overall"]+" HP: "+rev[0]["diff_drain"])
}

func sendMapLink(roomId string, bid string) {
	go func() {
		var rev []map[string]string
		webapi.GetMapInfo(bid, &rev)
		if len(rev) == 0 {
			sendMsg(roomId, "获取谱面信息失败")
			return
		}
		sendMsg(roomId, "这边是小夜下载链接[https://dl.sayobot.cn/beatmaps/download/full/"+rev[0]["beatmapset_id"]+" 完整版]  [https://dl.sayobot.cn/beatmaps/download/novideo/"+rev[0]["beatmapset_id"]+" 无视频]")
	}()
}
func mapInAble(roomid string, star, cs, ar, od, hp float32) bool {
	minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar := database.GetRules(roomid)
	return cs >= minCS && cs <= maxCS && ar >= minAR && ar <= maxAR && od >= minOD && od <= maxOD && hp >= minHP && hp <= maxHP && star >= minStar && star <= maxStar
}

func getRandomMap(roomid string) string {
	var rev map[int]string
	minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar := database.GetRules(roomid)
	mode := database.GetPlayMode(roomid)
	switch mode {
	case 1:
		minCS = -1
		minAR = -1
	case 3:
		minAR = -1
	}
	database.GetBeatmapIdByRules(mode, minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar, &rev)
	if len(rev) == 0 {
		return ""
	}
	return rev[rand.Intn(len(rev))]
}

func getMapDiff(bid string) (string, float32, float32, float32, float32, float32, bool) {
	var rev []map[string]string
	webapi.GetMapInfo(bid, &rev)
	if len(rev) == 0 {
		fmt.Println("获取谱面信息失败")
		fmt.Println(rev)
		return "", 0, 0, 0, 0, 0, false
	}
	var title string
	if rev[0]["title_unicode"] == "" {
		title = rev[0]["title"]
	} else {
		title = rev[0]["title_unicode"]
	}
	star := util.StringToFloat(rev[0]["difficultyrating"])
	cs := util.StringToFloat(rev[0]["diff_size"])
	ar := util.StringToFloat(rev[0]["diff_approach"])
	od := util.StringToFloat(rev[0]["diff_overall"])
	hp := util.StringToFloat(rev[0]["diff_drain"])
	mode := util.StringToInt(rev[0]["mode"])
	switch mode {
	case 1:
		cs = -1
		ar = -1
	case 3:
		ar = -1
	}
	go database.AddBeatmap(bid, mode, star, cs, ar, od, hp)
	return title, star, cs, ar, od, hp, true
}
