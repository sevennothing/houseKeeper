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
		`mobile`	char(11),
		`mail`	varchar(40),
		`gender`	varchar(10),
		`password`	varchar(40),
		`type`		varchar(255),
		`family`	char(64),
		`last_login_on`	datetime,
		`created_on`	datetime,
		`updated_on`	datetime,
		`head_photo`	varchar(255),
		`active`		bool,
		PRIMARY KEY (`id`)
	) ENGINE=MyISAM DEFAULT CHARSET=utf8;
*
*  History:
*     <author>   <time>   <desc>
*
*/

package usersdb

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const uidLength = 16

type User struct {
	ID          string `json:id`
	Login       string `json:login` // login name
	FirstName   string `json:firstName`
	LastName    string `json:lastName`
	Mobile      string `json:mobile`
	Mail        string `json:mail`
	Gender      string `json:gender`
	Password    string `json:password` //hashed password
	Type        string `json:Type`     //user type: family, user
	Family      string `json:Type`     //familyID
	LastLoginOn string `json:lastLoginOn`
	CreatedOn   string `json:createdOn`
	UpdatedOn   string `json:updatedOn`
	HeadPhoto   string `json:head_photo`
	Active      bool   `json:active`
}

type family struct {
	ID        string `json:id`
	Name      string `json:name`
	CreatedOn string `json:createdOn`
	UpdatedOn string `json:updatedOn`
}

// mysql session provider
type MysqlUser struct {
	DBPath string
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

func generateUid() (string, error) {
	b := make([]byte, uidLength)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", fmt.Errorf("Could not successfully read from the system CSPRNG.")
	}
	return hex.EncodeToString(b), nil
}

func (mu *MysqlUser) InsertUser(uinfo User) (User, error) {
	c := mu.connectInit()
	defer c.Close()
	uinfo.ID, _ = generateUid()
	lastLoginOn := time.Now()
	updatedOn := time.Now()
	createdOn := time.Now()
	_, err := c.Exec("INSERT INTO users VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		uinfo.ID,
		uinfo.Login,
		uinfo.FirstName,
		uinfo.LastName,
		uinfo.Mobile,
		uinfo.Mail,
		uinfo.Gender,
		uinfo.Password,
		uinfo.Type,
		uinfo.Family,
		lastLoginOn,
		createdOn,
		updatedOn,
		uinfo.HeadPhoto,
		uinfo.Active)
	if err != nil {
		fmt.Println(err)
		return uinfo, err
	}

	return uinfo, nil
}

func (mu *MysqlUser) FindUser(login string, uid string) (User, error) {
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

	var rowdata = User{ID: "",
		Login:       "",
		FirstName:   "",
		LastName:    "",
		Mobile:      "",
		Mail:        "",
		Gender:      "",
		Password:    "",
		Type:        "",
		Family:      "",
		LastLoginOn: "",
		CreatedOn:   "",
		UpdatedOn:   "",
		HeadPhoto:   "",
		Active:      false}

	if login != "" {
		err := byLogin.QueryRow(login).Scan(&rowdata.ID, &rowdata.Login, &rowdata.FirstName, &rowdata.LastName, &rowdata.Mobile, &rowdata.Mail, &rowdata.Gender, &rowdata.Password, &rowdata.Type, &rowdata.Family, &rowdata.LastLoginOn, &rowdata.CreatedOn, &rowdata.UpdatedOn, &rowdata.HeadPhoto, &rowdata.Active)
		if err != nil {
			fmt.Println("Error on FindUser: ", err.Error())
		}
	} else if uid != "" {
		err := byUid.QueryRow(login).Scan(&rowdata.ID, &rowdata.Login, &rowdata.FirstName, &rowdata.LastName, &rowdata.Mobile, &rowdata.Mail, &rowdata.Gender, &rowdata.Password, &rowdata.Type, &rowdata.Family, &rowdata.LastLoginOn, &rowdata.CreatedOn, &rowdata.UpdatedOn, &rowdata.HeadPhoto, &rowdata.Active)
		if err != nil {
			fmt.Println("Error on FindUser: ", err.Error())
		}
	} else {
		//return User{},fmt.Errorf("username and uid must be define")
		return User{}, errors.New("must be specify at least one parameter between uname and uid")
	}
	//fmt.Println("rowdata =  ",rowdata)

	return rowdata, nil
}

func (mu *MysqlUser) UpdateUser(uid string, firstname string, lastname string, mail string, gender string, lastlogin string, password string) error {
	// argument check

	c := mu.connectInit()
	defer c.Close()

	updateLastLogin, err := c.Prepare("UPDATE users SET last_login_on=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateLastLogin.Close()

	updateFirstName, err := c.Prepare("UPDATE users SET firstname =? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateFirstName.Close()

	updateLastName, err := c.Prepare("UPDATE users SET lastname=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateLastName.Close()

	updateMail, err := c.Prepare("UPDATE users SET mail=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateMail.Close()

	updateGender, err := c.Prepare("UPDATE users SET gender=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateGender.Close()

	updatePassword, err := c.Prepare("UPDATE users SET password=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updatePassword.Close()
	/*
		setActive, err := c.Prepare("UPDATE users SET active=? where id=?")
		if err != nil {
			panic(err.Error())
		}
		defer setActive.Close()
	*/
	updateUpdatedOn, err := c.Prepare("UPDATE users SET updated_on=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateUpdatedOn.Close()

	if firstname != "" {
		updateFirstName.Exec(firstname, uid)
	}
	if lastname != "" {
		updateLastName.Exec(lastname, uid)
	}
	if mail != "" {
		updateMail.Exec(mail, uid)
	}
	if gender != "" {
		updateGender.Exec(mail, uid)
	}
	if password != "" {
		updatePassword.Exec(password, uid)
	}
	if lastlogin != "" {
		updateLastLogin.Exec(lastlogin, uid)
	}

	updateUpdatedOn.Exec(time.Now(), uid)

	return nil

}

func (mu *MysqlUser) BindMobileForUser(login string, mobile string) error {
	c := mu.connectInit()
	defer c.Close()

	updateMobile, err := c.Prepare("UPDATE users SET mobile=?,active=?,updated_on=? where login=?")
	if err != nil {
		panic(err.Error())
	}
	defer updateMobile.Close()
	updateMobile.Exec(mobile, true, login, time.Now())

	return nil

}
