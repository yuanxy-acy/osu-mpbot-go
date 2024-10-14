package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type conf struct {
	WebapiToken      string `json:"webapi_token"`
	DatabaseHost     string `json:"database_host"`
	DatabaseName     string `json:"database_name"`
	DatabaseUser     string `json:"database_user"`
	DatabasePassword string `json:"database_password"`
	IrcNick          string `json:"irc_nick"`
	IrcPass          string `json:"irc_pass"`
}

var (
	DatabaseHost     string
	DatabaseUser     string
	WebapiToken      string
	DatabaseName     string
	DatabasePassword string
	IrcNick          string
	IrcPass          string
)

func init() {
	configPath := flag.String("c", "config.json", "config file path")
	flag.Parse()
	configJson, err := os.ReadFile(*configPath)
	if err != nil {
		fmt.Println("读取配置文件失败")
		fmt.Println(err)
		os.Exit(1)
	}
	var config = conf{}
	err = json.Unmarshal(configJson, &config)
	if err != nil {
		fmt.Println("解析配置文件失败")
		fmt.Println(err)
		os.Exit(1)
	}
	WebapiToken = config.WebapiToken
	DatabaseHost = config.DatabaseHost
	DatabaseName = config.DatabaseName
	DatabaseUser = config.DatabaseUser
	DatabasePassword = config.DatabasePassword
	IrcNick = config.IrcNick
	IrcPass = config.IrcPass
}
