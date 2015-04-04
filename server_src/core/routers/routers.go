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

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/application/user/login", &controllers.UserController{}, "*:UserLogin")
	beego.Router("/application/user/sigin", &controllers.UserController{}, "post:UserSigin")
	beego.Router("/application/user/bindMobile", &controllers.UserController{}, "*:UserBindMobile")
	beego.Router("/application/user/resetPassword", &controllers.UserController{}, "*:UserResetPassword")

}
