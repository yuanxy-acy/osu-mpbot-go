package webapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"osu_mp_bot/config"
	"strings"
)

var token = config.WebapiToken
var host = "https://osu.ppy.sh/"
var contentType = "application/x-www-form-urlencoded"

func osuApiV1(path string, data string) []byte {
	resp, err := http.Post(host+path, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("osuApiV1 failed, err:%v\n", err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed,err:%v\n", err)
		return nil
	}
	return b
}

func GetMapInfo(bid string, rev *[]map[string]string) {
	err := json.Unmarshal(osuApiV1("api/get_beatmaps", "k="+token+"&b="+bid), rev)
	if err != nil {
		fmt.Println("JSON decode error!")
		fmt.Println(err)
	}
}

func GetPlayRecent(userId string, rev *[]map[string]string) {
	err := json.Unmarshal(osuApiV1("api/get_user_recent", "k="+token+"&u="+userId+"&type=string&limit=1"), rev)
	if err != nil {
		fmt.Println("JSON decode error!")
		fmt.Println(err)
	}
}

func GetScores(bid, userId string, rev *[]map[string]string) {
	err := json.Unmarshal(osuApiV1("api/get_scores", "k="+token+"&u="+userId+"&b="+bid+"&type=string&limit=1"), rev)
	if err != nil {
		fmt.Println("JSON decode error!")
		fmt.Println(err)
	}
}

type UserInfo struct {
	Id   string `json:"user_id"`
	Name string `json:"username"`
}

func GetUserInfo(user string) UserInfo {
	var info []UserInfo
	err := json.Unmarshal(osuApiV1("api/get_user", "k="+token+"&u="+user+"&type=string&limit=1"), &info)
	if err != nil {
		fmt.Println("JSON decode error!")
		fmt.Println(err)
	}
	if len(info) < 1 {
		return UserInfo{
			Id:   "",
			Name: "",
		}
	}
	return info[0]
}
