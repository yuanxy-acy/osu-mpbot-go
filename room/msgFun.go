package room

import (
	"osu_mp_bot/database"
	"osu_mp_bot/util"
	"osu_mp_bot/webapi"
	"strconv"
	"strings"
)

func msgFun(roomId, user string, _args []string) {
	if len(_args[0]) < 3 {
		return
	}
	_args[0] = _args[0][1:]
	if user == "BanchoBot" {
		banchoMsg(roomId, _args)
		return
	}
	_args[0] = strings.Replace(_args[0], "！", "!", 1)
	var args []string
	for _, value := range _args {
		if value == "" {
			args = append(args, value)
		}
	}
	if database.GetMpMode(roomId) == -1 {
		args[0] = args[0][1:]
		baseCmdFun(roomId, user, args)
		return
	}
	if args[0][0] == '!' {
		args[0] = args[0][1:]
		switch strings.ToLower(args[0]) {
		case "ping":
			sendMsg(roomId, "pong! ! ! .≥v≤.")
		case "pick":
			if len(args) < 2 {
				return
			}
			switch database.GetMpMode(roomId) {
			case 0:
				if len(args) < 2 {
					sendMsg(roomId, "请输入 bid")
					return
				}
				switch args[1] {
				case "pool":
					if len(args) == 2 {
						return
					}
					bid := database.GetPoolMapBid(database.GetMapPool(roomId), strings.ToLower(args[2]))
					if bid == "" {
						sendMsg(roomId, args[2]+"不存在")
						return
					} else {
						database.AddPickedMap(roomId, bid, args[2]+" "+user)
						sendMsg(roomId, "目前还有"+strconv.Itoa(database.GetPickedMapCount(roomId))+"张谱面正在排队")
						return
					}
				case "host":
					info := webapi.GetUserInfo(user)
					if database.GetRoomHost(roomId) == "" {
						sendMsg(roomId, "!mp host "+info.Name)
					} else {
						sendMsg(roomId, "请等待当前房主选图")
					}
				default:
					info := webapi.GetUserInfo(user)
					pick(roomId, args[1], info.Name)
				}
			case 2:
				if len(args) < 2 || database.GetMapPool(roomId) == -1 {
					return
				}
				args[1] = strings.ToLower(args[1])
				red, blue := database.GetTeamCaptain(roomId)
				switch database.GetMatchState(roomId) {
				case "red pick":
					if user == red {
						if !poolPick(roomId, args[1]) {
							database.SetMatchState(roomId, "red picked")
						}
					}
				case "blue pick":
					if user == blue {
						if !poolPick(roomId, args[1]) {
							database.SetMatchState(roomId, "blue picked")
						}
					}
				}
			}
		case "skip":
			switch database.GetMpMode(roomId) {
			case 0:
				nextMap(roomId)
			case 1, 3:
				skipHost(roomId, user)
			case 2:
				if database.GetMapPool(roomId) == -1 {
					skipHost(roomId, user)
				}
			}
		case "start":
			sendMsg(roomId, "!mp start 10")
		case "abort":
			sendMsg(roomId, "!mp abort")
		case "set":
			if len(args) < 2 {
				sendMsg(roomId, "set 参数不足")
				return
			}
			switch strings.ToLower(args[1]) {
			case "v1":
				switch database.GetMpMode(roomId) {
				case 0, 1:
					sendMsg(roomId, "!mp set 0 0 16")
				case 2:
					sendMsg(roomId, "!mp set 2 0 16")
				case 3:
					sendMsg(roomId, "!mp set 1 0 16")
				}
			case "v2":
				switch database.GetMpMode(roomId) {
				case 0, 1:
					sendMsg(roomId, "!mp set 0 3 16")
				case 2:
					sendMsg(roomId, "!mp set 2 3 16")
				case 3:
					sendMsg(roomId, "!mp set 1 3 16")
				}
			case "dt":
				if database.GetMpMode(roomId) == 2 {
					return
				}
				sendMsg(roomId, "!mp mods dt freemod")
			case "nm":
				if database.GetMpMode(roomId) == 2 {
					return
				}
				sendMsg(roomId, "!mp mods none")
			case "fm":
				if database.GetMpMode(roomId) == 2 {
					return
				}
				sendMsg(roomId, "!mp mods freemod")
			case "rule":
				if len(args) < 3 {
					sendMsg(roomId, "!set rule 参数不足")
					return
				}
				setRules(roomId, args[2:])
			case "playmode":
				if len(args) < 3 {
					sendMsg(roomId, "set playmode 参数不足")
					return
				}
				database.SetPlayMode(roomId, util.StringToInt(args[2]))
			case "mpmode", "mp":
				if len(args) < 3 {
					sendMsg(roomId, "set mpmode 参数不足")
					return
				}
				setMpMode(roomId, user, args[2])
			case "captain", "cap":
				if database.GetMpMode(roomId) != 3 {
					return
				}
				if len(args) < 3 {
					sendMsg(roomId, "set captain 参数不足")
					return
				}
				red, blue := database.GetTeamCaptain(roomId)
				info := webapi.GetUserInfo(user)
				switch strings.ToLower(args[2]) {
				case "red":
					if len(args) == 3 {
						setCaptain(roomId, info.Name, info.Name, blue, red, "red", 3)
					} else {
						setInfo := webapi.GetUserInfo(args[3])
						setCaptain(roomId, info.Name, setInfo.Name, blue, red, "red", 4)
					}
				case "blue":
					if len(args) == 3 {
						setCaptain(roomId, info.Name, red, info.Name, blue, "blue", 3)
					} else {
						setInfo := webapi.GetUserInfo(args[3])
						setCaptain(roomId, info.Name, red, setInfo.Name, blue, "blue", 4)
					}
				}
			case "mappool":
				if len(args) < 3 {
					sendMsg(roomId, "!set mappool <poolid> 命令缺少参数<poolid>")
				}
				name, _, _, _, _, _, _, _, _ := database.GetMapPoolInfo(util.StringToInt(args[2]))
				if name == "" {
					sendMsg(roomId, "poolid:"+args[2]+"不存在")
				} else {
					database.SetMapPool(roomId, util.StringToInt(args[2]))
					sendMapPoolSetting(roomId, database.GetMapPool(roomId))
					if database.GetMpMode(roomId) == 2 && database.GetMatchState(roomId) == "free pick" {
						sendMsg(roomId, "!mp clearhost")
						sendMsg(roomId, "请双方指定队长，使用!set captain <red/blue> <user>指定或转让，当不存在队长时可省略<user>选项指定自己为队长")
						database.SetMatchState(roomId, "setting captain")
					}
				}
			}
		case "get":
			if len(args) < 2 {
				sendMsg(roomId, "get 参数不足")
				return
			}
			switch strings.ToLower(args[1]) {
			case "rules":
				minCS, maxCS, minAR, maxAR, minOD, maxOD, minHP, maxHP, minStar, maxStar := database.GetRules(roomId)
				sendMsg(roomId, "minStar "+util.FloatToString(minStar)+" maxStar "+util.FloatToString(maxStar))
				sendMsg(roomId, "minCS "+util.FloatToString(minCS)+" maxCS "+util.FloatToString(maxCS))
				sendMsg(roomId, "minAR "+util.FloatToString(minAR)+" maxAR "+util.FloatToString(maxAR))
				sendMsg(roomId, "minOD "+util.FloatToString(minOD)+" maxOD "+util.FloatToString(maxOD))
				sendMsg(roomId, "minHP "+util.FloatToString(minHP)+" maxHP "+util.FloatToString(maxHP))
			case "playmode", "pm":
				switch database.GetMpMode(roomId) {
				case 0:
					sendMsg(roomId, "当前游玩模式为：pick&random map")
				case 1:
					sendMsg(roomId, "当前游玩模式为：auto host rotate")
				case 2:
					sendMsg(roomId, "当前游玩模式为：read vs blue")
				case 3:
					sendMsg(roomId, "当前游戏模式为：relay mode - auto host rotate")
				}
			case "mappool":
				if len(args) == 2 {
					sendMapPoolSetting(roomId, database.GetMapPool(roomId))
				} else {
					var rev map[int][]string
					database.GetMapPoolIdLike(strings.Join(args[2:], " "), &rev)
					for _, item := range rev {
						sendMsg(roomId, " poolid:"+item[0]+" name:"+item[1]+" avgStar:"+item[2])
					}
				}
			}
		default:
			baseCmdFun(roomId, user, args)
		}
		return
	}
}
