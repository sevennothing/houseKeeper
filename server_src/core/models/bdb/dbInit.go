/*
 * Copyright 2015
 * Filename:income_expend.go
 * Author:caijun.Li  Date:2015-04-08
 * Description:
 * 记录用户的收支数据，数据库为mongoDB
 * History:
 *         <author>      <time>       <desc>
 *
 */

package bdb

import (
	"../tools"
	"fmt"
	"gopkg.in/mgo.v2"
)

var dbConfigPath = "./conf/db.conf"
var session *mgo.Session
var businessDb *mgo.Database

func init() {
	//Init MongoDB
	fmt.Println("MongoDB Session init")

	config, err := tools.ReadFile(dbConfigPath)
	if err != nil {
		fmt.Println("readFile: ", err.Error())
		panic(err.Error())
	}
	//fmt.Println(config)
	//fmt.Println(config["bdbUrl"])
	session, err = mgo.Dial(config["bdbUrl"])
	if err != nil {
		fmt.Println("connect Business DataBase: ", err.Error())
		panic(err.Error())
	}

	//fmt.Println(session)
	businessDb = session.DB(config["businessDBName"])

}

//TODO: close session
