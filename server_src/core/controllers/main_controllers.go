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
	//"fmt"
	"github.com/astaxie/beego"
	//	"github.com/astaxie/beego/validation"
	//"../models/errorcode"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Post() {
	this.Ctx.WriteString("hello post")
}

func (this *MainController) Get() {
	this.Ctx.WriteString("hello get")
}
