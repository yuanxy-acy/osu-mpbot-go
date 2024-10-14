package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"osu_mp_bot/config"
	"time"
)

var db *sql.DB

var (
	getInfo,
	getAllActiveRoom,
	setRoomName,
	addInfo,
	delInfo,
	getNowBid,
	setNowBid,
	getRoomHost,
	setRoomHost,
	setPlayMode,
	getPlayMode,
	getMpMode,
	setMpMode,

	getRules,
	setRules,
	delRules,

	addBeatmap,
	getBeatmapIdByRules,

	getFirstPickedMap,
	getPickedMapCount,
	delAllPickedMap,
	delFirstPickedMap,
	addPickedMap,

	getHostRotate,
	addHostRotate,
	delHostRotate,
	delAllHostRotate,

	delRoomTeam,
	getTeamCaptain,
	setTeamCaptain,
	getMatchState,
	setMatchState,
	getMapPool,
	setMapPool,

	getMapPoolIdLike,
	getMapPoolInfo,
	getPoolMapBid *sql.Stmt
)

func init() {
	var err error
	db, err = sql.Open("mysql", config.DatabaseUser+":"+config.DatabasePassword+"@tcp("+config.DatabaseHost+")/"+config.DatabaseName+"?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Second * 600)
	_, err = db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exec("create table if not exists roomlist (roomid varchar(30) NOT NULL PRIMARY KEY, name varchar(60) default '', nowBid varchar(30) default '', host varchar(30) default '', playMode int default 0, mpMode int default 0,minCS float default 0,maxCS float default 11,minAR float default 0,maxAR float default 11,minOD float default 0,maxOD float default 11,minHP float default 0,maxHP float default 11,minStar float default 0,maxStar float default 11,red varchar(30) default '',blue varchar(30) default '',state varchar(30) default '',mapPool int default -1)")
	exec("create table if not exists beatmap (bid varchar(30) NOT NULL PRIMARY KEY,playMode int,Star float,CS float,AR float,OD float,HP float)")
	exec("create table if not exists pickedmap (roomid varchar(30) NOT NULL,bid varchar(30) NOT NULL,user varchar(30) NOT NULL)")
	exec("create table if not exists hostrotate (roomid varchar(30) NOT NULL,user varchar(30) NOT NULL)")
	exec("create table if not exists mappools (poolid int NOT NULL PRIMARY KEY,name varchar(30) default '',avgStar float default 0,NM int default 0,HD int default 0,HR int default 0,DT int default 0,FM int default 0,EX int default 0,TB int default 0)")
	exec("create table if not exists mappoolmap (poolid int NOT NULL,bid varchar(30),mods varchar(30))")
	//exec("create table if not exists scores (user varchar(30) NOT NULL,bid varchar(30) NOT NULL,")

	getInfo = prepare("select name from roomlist where roomid = ?")
	getAllActiveRoom = prepare("select roomid from roomlist where state != 'closed'")
	setRoomName = prepare("update roomlist set name = ? where roomid = ?")
	addInfo = prepare("insert into roomlist (roomid, name, mpMode) values (?, ?, ?)")
	delInfo = prepare("delete from roomlist where roomid = ?")
	getNowBid = prepare("select nowBid from roomlist where roomid = ?")
	setNowBid = prepare("update roomlist set nowBid = ? where roomid = ?")
	getRoomHost = prepare("select host from roomlist where roomid = ?")
	setRoomHost = prepare("update roomlist set host = ? where roomid = ?")
	getPlayMode = prepare("select playMode from roomlist where roomid = ?")
	setPlayMode = prepare("update roomlist set playMode = ? where roomid = ?")
	getMpMode = prepare("select mpMode from roomlist where roomid = ?")
	setMpMode = prepare("update roomlist set mpMode = ? where roomid = ?")

	getRules = prepare("select minCS,maxCS,minAR,maxAR,minOD,maxOD,minHP,maxHP,minStar,maxStar from roomlist where roomid = ?")
	setRules = prepare("update roomlist set minCS = ?, maxCS = ?, minAR = ?, maxAR = ?, minOD = ?, maxOD = ?, minHP = ?, maxHP = ?, minStar = ?, maxStar = ? where roomid = ?")
	delRules = prepare("delete from roomlist where roomid = ?")

	addBeatmap = prepare("insert into beatmap values (?, ?, ?, ?, ?, ?, ?)")
	getBeatmapIdByRules = prepare("select bid from beatmap where playMode = ? and CS >= ? and CS <= ? and AR >= ? and AR <= ? and OD >= ? and OD <=? and HP >= ? and HP <= ? and Star >= ? and Star <= ?")

	getFirstPickedMap = prepare("select bid,user from pickedmap where roomid = ? limit 1")
	getPickedMapCount = prepare("select count(roomid) as mapCount from pickedmap where roomid = ?")
	delAllPickedMap = prepare("delete from pickedmap where roomid = ?")
	delFirstPickedMap = prepare("delete from pickedmap where roomid = ? limit 1")
	addPickedMap = prepare("insert into pickedmap values (?, ?, ?)")

	getHostRotate = prepare("select user from hostrotate where roomid = ?")
	addHostRotate = prepare("insert  into hostrotate values (?, ?)")
	delHostRotate = prepare("delete from hostrotate where user = ?")
	delAllHostRotate = prepare("delete from hostrotate where roomid = ?")

	delRoomTeam = prepare("delete from roomlist where roomid = ?")
	getTeamCaptain = prepare("select red, blue from roomlist where roomid = ?")
	setTeamCaptain = prepare("update roomlist set red = ?, blue = ? where roomid = ?")
	getMatchState = prepare("select state from roomlist where roomid = ?")
	setMatchState = prepare("update roomlist set state = ? where roomid = ?")
	getMapPool = prepare("select mapPool from roomlist where roomid = ?")
	setMapPool = prepare("update roomlist set mapPool = ? where roomid = ?")

	getMapPoolIdLike = prepare("select poolid, name, avgStar from mappools where name like concat('%' ,? ,'%')")
	getMapPoolInfo = prepare("select name, avgStar, NM, HD, HR, DT, FM, EX, TB from mappools where poolid = ?")

	getPoolMapBid = prepare("select bid from mappoolmap where poolid = ? and mods = ?")

	fmt.Println("数据库初始化完成")
}

func exec(sqlContext string) {
	_, err := db.Exec(sqlContext)
	if err != nil {
		fmt.Println(sqlContext + " 失败")
		fmt.Println(err)
	}
}

func prepare(sqlContext string) *sql.Stmt {
	stmt, err := db.Prepare(sqlContext)
	if err != nil {
		fmt.Println(sqlContext + " 失败")
		fmt.Println(err)
	}
	return stmt
}

func rowsToStringMap(rows *sql.Rows, rev *map[int]string) {
	count := 0
	var get string
	*rev = make(map[int]string)
	for rows.Next() {
		err := rows.Scan(&get)
		if err != nil {
			fmt.Println("sql rows scan err:", err)
		}
		(*rev)[count] = get
		count++
	}
	err := rows.Err()
	if err != nil {
		fmt.Println("sql rows err:", err)
	}
	err = rows.Close()
	if err != nil {
		fmt.Println("sql rows close err:", err)
	}
}

func Close() {
	excClose(getInfo)
	excClose(getAllActiveRoom)
	excClose(setRoomName)
	excClose(addInfo)
	excClose(delInfo)
	excClose(getNowBid)
	excClose(setNowBid)
	excClose(getRoomHost)
	excClose(setRoomHost)
	excClose(getPlayMode)
	excClose(setPlayMode)
	excClose(getMpMode)
	excClose(setMpMode)

	excClose(getRules)
	excClose(setRules)
	excClose(delRules)

	excClose(addBeatmap)
	excClose(getBeatmapIdByRules)

	excClose(getFirstPickedMap)
	excClose(getPickedMapCount)
	excClose(delAllPickedMap)
	excClose(delFirstPickedMap)
	excClose(addPickedMap)

	excClose(getHostRotate)
	excClose(addHostRotate)
	excClose(delHostRotate)
	excClose(delAllHostRotate)

	excClose(delRoomTeam)
	excClose(getTeamCaptain)
	excClose(setTeamCaptain)
	excClose(getMatchState)
	excClose(setMatchState)
	excClose(getMapPool)
	excClose(setMapPool)

	excClose(getMapPoolIdLike)
	excClose(getMapPoolInfo)
	excClose(getPoolMapBid)
}

func excClose(stmt *sql.Stmt) {
	err := stmt.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
