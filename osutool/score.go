package osutool

import (
	"fmt"
	"osu_mp_bot/util"
	"osu_mp_bot/webapi"
)

func SendUserScore(chanel, user, bid string) {
	info := webapi.GetUserInfo(user)
	var rev []map[string]string
	var mapRev []map[string]string
	webapi.GetMapInfo(bid, &mapRev)
	webapi.GetScores(bid, user, &rev)
	if len(rev) == 0 {
		IrcSengFunc(chanel, "获取"+info.Name+"信息失败")
		fmt.Println(rev)
		return
	}
	var title string
	if len(mapRev) == 0 {
		title = ""
		fmt.Println(mapRev)
	} else if mapRev[0]["title_unicode"] == "" {
		title = mapRev[0]["title"]
	} else {
		title = mapRev[0]["title_unicode"]
	}
	IrcSengFunc(user, info.Name+" 的成绩  [osu://b/"+bid+" "+title+"] "+util.GetModName(rev[0]["enabled_mods"])+"Rank: "+rev[0]["rank"]+" score: "+rev[0]["score"]+" pp: "+rev[0]["pp"]+" max combo: "+rev[0]["maxcombo"]+" 300: "+rev[0]["count300"]+" 100: "+rev[0]["count100"]+" 50: "+rev[0]["count50"]+" miss: "+rev[0]["countmiss"])
}

func SendUserRecent(chanel, user string) {
	info := webapi.GetUserInfo(user)
	var rev []map[string]string
	webapi.GetPlayRecent(user, &rev)
	if len(rev) == 0 {
		IrcSengFunc(chanel, "获取"+info.Name+"信息失败")
		fmt.Println(rev)
		return
	}
	var mapRev []map[string]string
	webapi.GetMapInfo(rev[0]["beatmap_id"], &mapRev)
	var title string
	if len(mapRev) == 0 {
		title = ""
		fmt.Println(mapRev)
	} else if mapRev[0]["title_unicode"] == "" {
		title = mapRev[0]["title"]
	} else {
		title = mapRev[0]["title_unicode"]
	}
	IrcSengFunc(chanel, info.Name+" 的成绩 bid: "+rev[0]["beatmap_id"]+" [osu://b/"+rev[0]["beatmap_id"]+" "+title+"] "+util.GetModName(rev[0]["enabled_mods"])+"Rank: "+rev[0]["rank"]+" score: "+rev[0]["score"]+" max combo: "+rev[0]["maxcombo"]+" 300: "+rev[0]["count300"]+" 100: "+rev[0]["count100"]+" 50: "+rev[0]["count50"]+" miss: "+rev[0]["countmiss"])
}
