/**
 *  Copyright 2015,
 *  Filename: main.go
 *  Author: caijun.Li
 *  Date: 2015-03-27
 *  Description:
 *  History:
 *     <author>   <time>   <desc>
 *
 */
package main

import (
	_ "./routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
