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
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/session"
	"github.com/astaxie/beego/validation"
	"strings"
	//	"strconv"
	"../models/errorcode"
	"../models/usersdb"
	"reflect"
)

var globalSessions *session.Manager
var globalUser *usersdb.MysqlUser

type sessionConfig struct {
	CookieName        string `json:cookieName`
	EnableSetCookie   bool   `json:enableSetCookie,omitempty`
	Gclifetime        int    `json:gclifetime`
	MaxLifetime       int    `json:maxLifetime`
	Secure            bool   `json:secure`
	SessionIDHashFunc string `json:sessionIDHashFunc`
	SessionIDHashKey  string `json:sessionIDHashKey`
	CookieLifeTime    int    `json:cookieLifeTime`
	ProviderConfig    string `json:providerConfig`
}

func init() {
	fmt.Println("Beego Seesion init")
	cookieName := beego.AppConfig.String("cookiename")
	sessionIDHashKey := beego.AppConfig.String("sessionIDHashKey")
	jsonConfig := sessionConfig{CookieName: cookieName,
		EnableSetCookie:   true,
		Gclifetime:        3600,
		MaxLifetime:       3600,
		Secure:            false,
		SessionIDHashFunc: "sha1",
		SessionIDHashKey:  sessionIDHashKey,
		CookieLifeTime:    3600,
		ProviderConfig:    ""}
	b, _ := json.Marshal(jsonConfig)
	confStr := string(b)
	fmt.Println(confStr)
	confStr = strings.Replace(confStr, "EnableSetCookie", "EnableSetCookie,omitempty", -1)
	fmt.Println(confStr)
	//memory mysql redis or file
	globalSessions, _ = session.NewManager("memory", confStr)
	go globalSessions.GC()

	//users
	globalUser = &usersdb.MysqlUser{DBPath: beego.AppConfig.String("mysqldb")}

}

type MainController struct {
	beego.Controller
}

type UserController struct {
	beego.Controller
}

func (this *MainController) Post() {
	this.Ctx.WriteString("hello post")
}

func (this *MainController) Get() {
	this.Ctx.WriteString("hello get")
}

func (this *UserController) Get() {
	this.Ctx.WriteString("GET: hello consumer login")
}

//  login Access
//  params:
//     @ username string
//	   @ password string
//	return:
//		successed:{sessionID:X, faildCnt:0,isLogin:true,UserInfo:{}}
//		faild:{sessionID:x, faildCnt:x,isLogin:false,code:x,message:x,UserInfo:{}}
type resultLogin struct {
	SessionID string       `json:sessionID`
	FaildCnt  int          `json:faildCnt`
	IsLogin   bool         `json:isLogin`
	Code      string       `json:code`
	Message   string       `json:message`
	UserInfo  usersdb.User `json:userInfo`
	//Token

}

func (this *UserController) UserLogin() {
	var res = resultLogin{FaildCnt: 0, IsLogin: false}
	sess, err := globalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	if err != nil {
		fmt.Println("set error,", err)
		res.Code = errorCode.SessionStartError
		res.Message = err.Error()
		jsonRes, _ := json.Marshal(res)
		// exit
		this.Ctx.WriteString(string(jsonRes))

		return
	}

	//ckeck user login info
	req_username := this.GetString("username")
	req_passwd := this.GetString("password")
	// find in mysql database with user
	uinfo, err := globalUser.FindUser(req_username, "")
	res.UserInfo = uinfo
	var isAuth bool = false
	if err != nil {
		res.Code = errorCode.LoginFaild
		res.Message = err.Error()
		jsonRes, _ := json.Marshal(res)
		this.Ctx.WriteString(string(jsonRes))
		return
	}

	if uinfo.Password == req_passwd {
		// password access
		isAuth = true
	} else {
		// password auth faild
		isAuth = false
	}

	var faildCnt int = 0

	if reflect.TypeOf(sess.Get("faildCnt")) != nil {
		faildCnt = sess.Get("faildCnt").(int)
	}

	//fmt.Println(username)
	if this.Ctx.Request.Method == "GET" {
		fmt.Println("======")
		// logout
		sess.Set("username", "")
		sess.Set("isLogin", false)

	} else {
		sess.Set("username", req_username)
		if !isAuth {
			sess.Set("faildCnt", faildCnt+1)
			res.Code = errorCode.LoginFaild
			res.Message = "password author error"
		} else {
			sess.Set("faildCnt", 0)
		}
		sess.Set("isLogin", isAuth)
	}

	res.SessionID = sess.SessionID()
	res.IsLogin = sess.Get("isLogin").(bool)
	res.FaildCnt = sess.Get("faildCnt").(int)
	jsonRes, _ := json.Marshal(res)
	this.Ctx.WriteString(string(jsonRes))
}

//  sigin Access
//  params:
//     @ username string
//	   @ password string
//	   @ Mail	  string
//	   @ Gender	  string
//	   @ head_photo	string	-->URL
//	return:
//		successed:{UserInfo:{}}
//		faild:{code:x,message:x,UserInfo:{}}
type resultSigin struct {
	Code     string       `json:code`
	Message  string       `json:message`
	UserInfo usersdb.User `json:userInfo`
}

type user struct {
	Name string `valid:"Required; Range(5,10)"` // Name 不能为空并且以Bee开头
	//Age    int    `valid:"Range(1, 140)"` // 1 <= Age <= 140，超出此范围即为不合法
	Mail string `valid:"Email; MaxSize(100)"` // Email字段需要符合邮箱格式，并且最大长度不能大于100个字符
	//	Mobile string `valid:"Mobile"` // Mobile必须为正确的手机号
}

func (u *user) Valid(v *validation.Validation) {
	if strings.Index(u.Name, "admin") != -1 {
		// 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("Name", "名称里不能含有 admin")
	}
}

func (this *UserController) UserSigin() {
	var res = resultSigin{}
	var valid = validation.Validation{}
	username := this.GetString("username")
	password := this.GetString("password")
	mail := this.GetString("mail")
	gender := this.GetString("gender")
	head_photo := this.GetString("head_photo")

	//fmt.Printf("username=%s\npassword=%s\nmail=%s\ngender=%s\nhead_photo=%s\n",username,password,email,gender,head_photo)
	u := user{
		Name: username,
		Mail: mail,
	}
	_, err := valid.Valid(&u)
	if err != nil {
		res.Code = errorCode.MalformedArgument
		res.Message = err.Error()
		jsonRes, _ := json.Marshal(res)
		this.Ctx.WriteString(string(jsonRes))
		return
	} else {
		var user = usersdb.User{Login: username,
			FirstName: "",
			LastName:  "",
			Mail:      mail,
			Gender:    gender,
			Password:  password,
			HeadPhoto: head_photo}
		uinfo, err := globalUser.InsertUser(user)
		res.UserInfo = uinfo
		if err != nil {
			res.Code = errorCode.SiginUsernameIsExist
			res.Message = err.Error()
		} else {
			//something
		}
		jsonRes, _ := json.Marshal(res)
		this.Ctx.WriteString(string(jsonRes))
	}
}
