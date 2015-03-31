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
	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego/session"
	"encoding/json"
	"fmt"
	"strings"
	"../models"
)

//type usersdb models.usersdb

var globalSessions *session.Manager
var globalUser	*usersdb.MysqlUser


type sessionConfig	struct {
	CookieName  string `json:cookieName`
//	EnableSetCookie,omitempty   bool `json:enableSetCookie,omitempty`
//	EnableSetCookiei\,omitempty   bool `json:enableSetCookie,omitempty`
	EnableSetCookie   bool `json:enableSetCookie,omitempty`
	Gclifetime  int `json:gclifetime`
	MaxLifetime int `json:maxLifetime`
	Secure  bool `json:secure`
	SessionIDHashFunc string `json:sessionIDHashFunc`
	SessionIDHashKey string `json:sessionIDHashKey`
	CookieLifeTime int `json:cookieLifeTime`
	ProviderConfig string `json:providerConfig`
}

func init(){
	fmt.Println("Beego Seesion init");
	cookieName := beego.AppConfig.String("cookiename")
	sessionIDHashKey := beego.AppConfig.String("sessionIDHashKey")
	jsonConfig := sessionConfig{CookieName:cookieName, 
					//"EnableSetCookie,omitempty": true, 
					//EnableSetCookie\,omitempty: true, 
					EnableSetCookie: true, 
					Gclifetime:3600, 
					MaxLifetime: 3600, 
					Secure: false, 
					SessionIDHashFunc: "sha1", 
					SessionIDHashKey: sessionIDHashKey, 
					CookieLifeTime: 3600, 
					ProviderConfig: ""}
	b,_ :=json.Marshal(jsonConfig)
	confStr := string(b)
	fmt.Println(confStr)
	confStr = strings.Replace(confStr,"EnableSetCookie","EnableSetCookie,omitempty",-1)
	fmt.Println(confStr)
	//memory mysql redis or file
	globalSessions, _ = session.NewManager("memory", confStr)
	go globalSessions.GC()

	//users
	globalUser = &usersdb.MysqlUser{DBPath:beego.AppConfig.String("mysqldb")}

}

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
	this.Ctx.WriteString("GET: hello consumer login")
}

func (this *ConsumerController) ConsumerLogin() {
	sess,err := globalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	if err != nil {
		fmt.Println("set error,", err)
	}

	//ckeck user login info
	req_username := this.GetString("username")
	//req_passwd := this.GetString("password")
	// find in mysql database with user


	username := sess.Get("username")
	//isLogin := sess.Get("isLogin")

	fmt.Println(username)
	if this.Ctx.Request.Method == "GET" {
		fmt.Println("======")
	} else {
		sess.Set("username", this.GetString("username"))
		sess.Set("isLogin", true)
	}


	uinfo,err := globalUser.FindUser(req_username,"")
	if err != nil {
		this.Ctx.WriteString(err.Error());
	}


	//this.Ctx.WriteString("consumer login:" + this.GetString("username") + "databaseInfo:" + uinfo)
	this.Ctx.WriteString("consumer login:" + this.GetString("username") + "\ndatabaseInfo:" + string(uinfo))
}

type user struct {
	Id     int
	Name   string `valid:"Required; Range(5,10)"` // Name 不能为空并且以Bee开头
	//Age    int    `valid:"Range(1, 140)"` // 1 <= Age <= 140，超出此范围即为不合法
	Email  string `valid:"Email; MaxSize(100)"` // Email字段需要符合邮箱格式，并且最大长度不能大于100个字符
	Mobile string `valid:"Mobile"` // Mobile必须为正确的手机号
	//	IP     string `valid:"IP"` // IP必须为一个正确的IPv4地址
}
func (u *user) Valid(v *validation.Validation) {
	if strings.Index(u.Name, "admin") != -1 {
		// 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("Name", "名称里不能含有 admin")
	}
}

func (this *ConsumerController) ConsumerSigin() {
	var valid = validation.Validation{}
	username := this.GetString("username")
	password := this.GetString("password")
	email := this.GetString("email")
	gender := this.GetString("gender")
	head_photo := this.GetString("head_photo")

	fmt.Printf("username=%s\npassword=%s\nemail=%s\ngender=%s\nhead_photo=%s\n",username,password,email,gender,head_photo)
	u := user{Id:1,Name:username,Email:email,Mobile:"12322222322"}
	_,err :=  valid.Valid(&u)
	if err != nil{
		fmt.Println(err)
		this.Ctx.WriteString("error .....")
	}else{
		var user = usersdb.User{Login:username,FirstName:"",Mail:email,Gender:gender,Password:password}
		uinfo,err  := globalUser.InsertUser(user)
		if err != nil {
			this.Ctx.WriteString("user sigin error:" + err.Error())

		}else{
			this.Ctx.WriteString("consumer sigin ok,uid=" + uinfo.ID)
		}
	}
}

