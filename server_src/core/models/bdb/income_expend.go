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
	//"fmt"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type IncomeAndExpend struct {
	ID  bson.ObjectId `bson:"_id"`
	Uid string        `bson:"uid"`
}

func (this *IncomeAndExpend) Insert_doc(doc IncomeAndExpend) error {
	var inExCollection = businessDb.C("income_expend")

	doc.ID = bson.NewObjectId()

	err := inExCollection.Insert(&doc)

	if err != nil {
		panic(err)
	}

	return nil
}
