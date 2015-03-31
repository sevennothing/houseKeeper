
/**
 *  Copyright 2015,
 *  Filename: routers.go
 *  Author: caijun.Li
 *  Date: 2015-03-27
 *  Description:
 *  History:
 *     <author>   <time>   <desc>
 *
 */


 package routers

 import (
 	"../controllers"
	"github.com/astaxie/beego"
 )

 func init(){
 	beego.Router("/",&controllers.MainController{})
	beego.Router("/consumer/login", &controllers.ConsumerController{},"*:ConsumerLogin")
	beego.Router("/consumer/sigin", &controllers.ConsumerController{},"post:ConsumerSigin")

 }
