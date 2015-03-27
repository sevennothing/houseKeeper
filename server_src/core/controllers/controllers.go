/**
 *  Copyright 2015,
 *  Filename: controllers.go
 *  Author: caijun.Li
 *  Date: 2015-03-27
 *  Description:
 *  History:
 *     <author>   <time>   <desc>
 *
 */

package controllers

import (
	//"encoding/json"
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/example/beeapi/models"
)

type MainController struct {
	beego.Controller
}

type ConsumerController struct {
	beego.Controller
}

func (this *MainController) Post() {
	this.Ctx.WriteString("hello post")
}

func (this *MainController) Get() {
	this.Ctx.WriteString("hello get")
}


func (this *ConsumerController) Get() {
	this.Ctx.WriteString("hello consumer login")
}
func (this *ConsumerController) ConsumerLogin() {
	this.Ctx.WriteString("hello post consumer login")
}
