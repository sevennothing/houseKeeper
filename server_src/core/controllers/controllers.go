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
	"time"
	//	"strconv"
	"../models/errorcode"
	"../models/tools"
	"../models/usersdb"
	"reflect"
)

var globalSessions *session.Manager
var globalUser *usersdb.MysqlUser
var globalTool *tools.HkTool

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
//	   @ Mobile	  string need
//	   @ Mail	  string option
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
	Mail   string `valid:"Email; MaxSize(100)"` // Email字段需要符合邮箱格式，并且最大长度不能大于100个字符
	Mobile string `valid:"Mobile"`              // Mobile必须为正确的手机号
}

func (u *user) Valid(v *validation.Validation) {
	if strings.Index(u.Name, "admin") != -1 {
		// 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		v.SetError("Name", "名称里不能含有 admin")
	}
}

func (this *UserController) UserSigin() {
	var res = resultSigin{}
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

	var valid = validation.Validation{}
	username := this.GetString("username")
	password := this.GetString("password")
	mobile := this.GetString("mobile")
	mail := this.GetString("mail")
	gender := this.GetString("gender")
	head_photo := this.GetString("head_photo")

	//fmt.Printf("username=%s\npassword=%s\nmail=%s\ngender=%s\nhead_photo=%s\n",username,password,email,gender,head_photo)
	u := user{
		Name:   username,
		Mail:   mail,
		Mobile: mobile,
	}
	_, err = valid.Valid(&u)
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
			Mobile:    "", // 暂时不写入用户数据库，等待验证通过
			Mail:      mail,
			Gender:    gender,
			Password:  password,
			HeadPhoto: head_photo,
			//新用户未激活，激活需要短信验证
			Active: false}
		user.ID = globalTool.GenerateID16() // 生成16位唯一字符串ID

		uinfo, err := globalUser.InsertUser(user)
		res.UserInfo = uinfo
		if err != nil {
			res.Code = errorCode.SiginUsernameIsExist
			res.Message = err.Error()
		} else {
			//something
			authCode := globalTool.GenerateRandNumber(6)     //产生6位数验证码
			acExpire := time.Now().Unix() + (10 * 60 * 1000) // 有效期为10分钟
			sess.Set("username", uinfo.Login)
			sess.Set("mobile", mobile)           // 会话记录绑定手机号码
			sess.Set("authCode", authCode)       //会话记录认证码
			sess.Set("authCodeExpire", acExpire) //设置认证码过期时间
		}
		jsonRes, _ := json.Marshal(res)
		this.Ctx.WriteString(string(jsonRes))
	}
}

//  bind mobile Access
//  params:
//     @ mobile string
//	   @ authCode string
//	return:
//		successed:{}
//		faild:{}
type resultBindMobile struct {
	Code     string       `json:code`
	Message  string       `json:message`
	UserInfo usersdb.User `json:userInfo`
	//Token

}

func (this *UserController) UserBindMobile() {
	var res = resultBindMobile{}
	sess, err := globalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	if err != nil {
		fmt.Println("set error,", err)
		res.Code = errorCode.SessionStartError
		res.Message = err.Error()
		jsonRes, _ := json.Marshal(res)
		this.Ctx.WriteString(string(jsonRes))

		return
	}
	if this.Ctx.Request.Method == "GET" {
		//
		mobile := this.GetString("mobile")
		authCode := globalTool.GenerateRandNumber(6)     //产生6位数验证码
		acExpire := time.Now().Unix() + (10 * 60 * 1000) // 有效期为10分钟
		sess.Set("mobile", mobile)                       // 会话记录绑定手机号码
		sess.Set("authCode", authCode)                   //会话记录认证码
		sess.Set("authCodeExpire", acExpire)             //设置认证码过期时间
		//TODO: send auth code through mobile message

	} else if this.Ctx.Request.Method == "POST" {
		//check authCode vilidate

		current_time := time.Now().Unix()
		if current_time > sess.Get("authCodeExpire").(int64) {
			res.Code = errorCode.AuthFaild
			res.Message = "The Auth Code already expired"
			jsonRes, _ := json.Marshal(res)
			this.Ctx.WriteString(string(jsonRes))
			return
		}
		theAuthCode := this.GetString("authCode")
		if theAuthCode == sess.Get("autchCode") {
			// auth sucessed
			// modify user database
			err = globalUser.BindMobileForUser(sess.Get("username").(string), sess.Get("mobile").(string))
			if err != nil {
				res.Code = errorCode.BindMobileFaild
				res.Message = err.Error()

			}
			jsonRes, _ := json.Marshal(res)
			this.Ctx.WriteString(string(jsonRes))
			return

		} else {
			// auth faild
			res.Code = errorCode.AuthFaild
			res.Message = "The Auth Code is Not Match"
			jsonRes, _ := json.Marshal(res)
			this.Ctx.WriteString(string(jsonRes))
			return
		}
	} else {
		// don't support
		res.Code = errorCode.MethodNotSupport
		res.Message = "Reset password only support GET or POST method"

		jsonRes, _ := json.Marshal(res)
		this.Ctx.WriteString(string(jsonRes))
		return
	}

}

//  user reset password Access
//  params:
//     @ username string
//	   @ newPassword string
//	return:
//		successed:{}
//		faild:{}
type resultResetPwd struct {
	Code     string       `json:code`
	Message  string       `json:message`
	UserInfo usersdb.User `json:userInfo`
	//Token

}

func (this *UserController) UserResetPassword() {
	var res = resultResetPwd{}
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

	if this.Ctx.Request.Method == "GET" {
		//check username
		username := this.GetString("username")
		newPassword := this.GetString("newPassword")

		if username == "" || newPassword == "" {
			res.Code = errorCode.MalformedArgument
			res.Message = err.Error()
			jsonRes, _ := json.Marshal(res)
			this.Ctx.WriteString(string(jsonRes))
			return

		}
		// Find This User
		uinfo, err := globalUser.FindUser(username, "")
		if err != nil {
			res.Code = errorCode.FindUserFaild
			res.Message = err.Error()
			jsonRes, _ := json.Marshal(res)
			this.Ctx.WriteString(string(jsonRes))

			return
		}
		mobile := uinfo.Mobile
		uid := uinfo.ID
		if mobile == "" {
			res.Code = errorCode.ResetPasswordFaild
			res.Message = "this user is not bind mobile"
			jsonRes, _ := json.Marshal(res)
			this.Ctx.WriteString(string(jsonRes))

			return
		}
		// authrization user
		// TODO :  生成临时的认证资源，或者验证码，下发至认证手机或邮箱
		authCode := globalTool.GenerateRandNumber(6)     //产生6位数验证码
		acExpire := time.Now().Unix() + (10 * 60 * 1000) // 有效期为10分钟
		//sess.Set("mobile", mobile)                       // 会话记录绑定手机号码
		sess.Set("uid", uid)
		sess.Set("username", username)       //会话记录登陆名
		sess.Set("newPassword", newPassword) //会话记录新密码
		sess.Set("authCode", authCode)       //会话记录认证码
		sess.Set("authCodeExpire", acExpire) //设置认证码过期时间
		//TODO: send auth code through mobile message

		jsonRes, _ := json.Marshal(res)
		this.Ctx.WriteString(string(jsonRes))
		return
	} else if this.Ctx.Request.Method == "POST" {
		//
		theAuthCode := this.GetString("authCode")
		if theAuthCode == sess.Get("authCode") {
			//auth successed,and update password
			err = globalUser.UpdateUser(sess.Get("uid").(string), "", "", "", "", "", sess.Get("newPassword").(string))
			if err != nil {
				res.Code = errorCode.ResetPasswordFaild
				res.Message = err.Error()
			} else {
				//TODO something
			}

		} else {
			//auth faild
			res.Code = errorCode.ResetPasswordFaild
			res.Message = "The Auth Code is Not Match"

		}
		jsonRes, _ := json.Marshal(res)
		this.Ctx.WriteString(string(jsonRes))
	} else {
		// don't support
		res.Code = errorCode.MethodNotSupport
		res.Message = "Reset password only support GET or POST method"

		jsonRes, _ := json.Marshal(res)
		// exit
		this.Ctx.WriteString(string(jsonRes))
		return
	}

}
