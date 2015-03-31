/**
*  Copyright 2015,
*  Filename: users_db.go
*  Author: caijun.Li
*  Date: 2015-03-31
*  Description:
*  need create table as sql:
	CREATE TABLE `users` (
		`id` char(64) NOT NULL,
		`login` varchar(10) NOT NULL UNIQUE,
		`firstname`	varchar(4),
		`lastname`	varchar(10),
		`mail`	varchar(40),
		`gender`	varchar(10),
		`password`	varchar(40),
		`type`		varchar(255),
		`family`	char(64),
		`last_login_on`	datetime,
		`created_on`	datetime,
		`updated_on`	datetime,
		PRIMARY KEY (`id`)
	) ENGINE=MyISAM DEFAULT CHARSET=utf8;
*
*  History:
*     <author>   <time>   <desc>
*
*/


package usersdb

import (
	"database/sql"
	"time"
	"crypto/rand"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

const uidLength = 16

type User struct {
	ID					string `json:id`
	Login				string	`json:login`	// login name
	FirstName		string	`json:firstName`
	LastName		string	`json:lastName`
	Mail				string	`json:mail`
	Gender			string	`json:gender`
	Password		string	`json:password`	//hashed password
	Type				string	`json:Type`	//user type: family, user
	Family			string	`json:Type` //familyID
	LastLoginOn	string	`json:lastLoginOn`
	CreatedOn		string	`json:createdOn`
	UpdatedOn		string	`json:updatedOn`

}



type family struct {
	ID				string	`json:id`
	Name			string	`json:name`
	CreatedOn	string	`json:createdOn`
	UpdatedOn	string	`json:updatedOn`
}

// mysql session provider
type MysqlUser struct {
	DBPath    string
}


// connect to mysql
func (mu *MysqlUser) connectInit() *sql.DB {
	db, e := sql.Open("mysql", mu.DBPath)
	if e != nil {
		fmt.Println(e)
		return nil
	}
	return db
}


func  generateUid() (string, error){
	b := make([]byte, uidLength)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", fmt.Errorf("Could not successfully read from the system CSPRNG.")
	}
	return hex.EncodeToString(b), nil
}

func (mu *MysqlUser) InsertUser(uinfo User) (User, error){
	c := mu.connectInit();
	defer c.Close()
	uinfo.ID,_ = generateUid()
	lastLoginOn := time.Now()
	updatedOn := time.Now()
	createdOn := time.Now()
	_,err := c.Exec("INSERT INTO users VALUES(?,?,?,?,?,?,?,?,?,?,?,?)",
						uinfo.ID,
						uinfo.Login,
						uinfo.FirstName,
						uinfo.LastName,
						uinfo.Mail,
						uinfo.Gender,
						uinfo.Password,
						uinfo.Type,
						uinfo.Family,
						lastLoginOn,
						createdOn,
						updatedOn)
	if err != nil {
		fmt.Println(err)
		return uinfo,err
	}

	return uinfo, nil
}


func (mu *MysqlUser) FindUser(login string,uid string) ([]byte, error){
	c := mu.connectInit()
	defer c.Close()
	byLogin, err := c.Prepare("SELECT * FROM users WHERE login = ?")
	if err != nil {
		panic(err.Error())
	}

	defer byLogin.Close()


	byUid, err := c.Prepare("SELECT * FROM users WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	defer byUid.Close()

	var rowdata []byte
	if login != "" {
		byLogin.QueryRow(login).Scan(&rowdata)
	}else if uid != "" {
		byUid.QueryRow(uid).Scan(&rowdata)
	}else{
		return nil,fmt.Errorf("username and uid must be define")
	}
	fmt.Println(rowdata)
	return rowdata, nil
}

func (mu *MysqlUser) UpdateUser(uid string,firstname string,lastname string,mail string,gender string, lastlogin string, password string) (error) {
	c := mu.connectInit()
	defer c.Close()

	updateLastLogin,err := c.Prepare("UPDATE users SET last_login_on=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateLastLogin.Close()

	updateFirstName,err := c.Prepare("UPDATE users SET firstname =? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateFirstName.Close()

	updateLastName,err := c.Prepare("UPDATE users SET lastname=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateLastName.Close()

	updateMail,err := c.Prepare("UPDATE users SET mail=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateMail.Close()

	updateGender,err := c.Prepare("UPDATE users SET gender=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateGender.Close()


	updatePassword,err := c.Prepare("UPDATE users SET password=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updatePassword.Close()

	updateUpdatedOn,err := c.Prepare("UPDATE users SET updated_on=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateUpdatedOn.Close()

	if firstname != "" {
		updateFirstName.Exec(firstname,uid)
	}else if lastname != "" {
		updateLastName.Exec(lastname,uid)
	}else if mail != "" {
		updateMail.Exec(mail,uid)
	}else if gender != "" {
		updateGender.Exec(mail,uid)
	}else if password != "" {
		updatePassword.Exec(password,uid)
	}else if	lastlogin != "" {
		updateLastLogin.Exec(lastlogin,uid)
	}else {
		fmt.Println("UpdateUser need params")
		return errors.New("UpdateUser need params")
		return fmt.Errorf("UpdateUser need params")
	}
	_,err := updateUpdatedOn.Exec(time.now(),uid)

}

