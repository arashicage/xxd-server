package model

import (
	"time"
)

type SysUser struct {
	Id       int       `json:"id" xorm:"not null pk autoincr MEDIUMINT(8)"`
	Dept     int       `json:"dept" xorm:"not null index MEDIUMINT(8)"`
	Account  string    `json:"account" xorm:"not null default '' unique index(accountPassword) CHAR(30)"`
	Password string    `json:"-" field:"password" xorm:"not null default '' index(accountPassword) CHAR(32)"`
	Realname string    `json:"realname" xorm:"not null default '' CHAR(30)"`
	Role     string    `json:"role" xorm:"not null CHAR(30)"`
	Nickname string    `json:"-" field:"nickname" xorm:"not null default '' CHAR(60)"`
	Admin    string    `json:"admin" xorm:"not null default 'no' index ENUM('no','common','super')"`
	Avatar   string    `json:"avatar" xorm:"not null default '' VARCHAR(255)"`
	Birthday time.Time `json:"-" field:"birthday" xorm:"not null DATE"`
	Gender   string    `json:"gender" xorm:"not null default 'u' ENUM('u','f','m')"`
	Email    string    `json:"email" xorm:"not null default '' CHAR(90)"`
	Skype    string    `json:"-" field:"skype" xorm:"not null CHAR(90)"`
	Qq       string    `json:"-" field:"qq" xorm:"not null default '' CHAR(20)"`
	Yahoo    string    `json:"-" field:"yahoo" xorm:"not null default '' CHAR(90)"`
	Gtalk    string    `json:"-" field:"gtalk" xorm:"not null default '' CHAR(90)"`
	Wangwang string    `json:"-" field:"wangwang" xorm:"not null default '' CHAR(90)"`
	Site     string    `json:"site" xorm:"not null VARCHAR(100)"`
	Mobile   string    `json:"mobile" xorm:"not null default '' CHAR(11)"`
	Phone    string    `json:"phone" xorm:"not null default '' CHAR(20)"`
	Address  string    `json:"-" field:"address" xorm:"not null default '' CHAR(120)"`
	Zipcode  string    `json:"-" field:"zipcode" xorm:"not null default '' CHAR(10)"`
	Visits   int       `json:"-" field:"visits" xorm:"not null default 0 MEDIUMINT(8)"`
	Ip       string    `json:"-" field:"ip" xorm:"not null default '' CHAR(50)"`
	Last     time.Time `json:"-" field:"last" xorm:"not null DATETIME"`
	Ping     time.Time `json:"-" field:"ping" xorm:"not null DATETIME"`
	Fails    int       `json:"-" field:"fails" xorm:"not null default 0 TINYINT(3)"`
	Join     time.Time `json:"-" field:"join" xorm:"not null DATETIME"`
	Locked   time.Time `json:"-" field:"locked" xorm:"not null DATETIME"`
	Deleted  string    `json:"-" field:"deleted" xorm:"not null ENUM('0','1')"`
	Status   string    `json:"status" xorm:"not null default 'offline' ENUM('online','away','busy','offline')"`
	Settings string    `json:"-" field:"settings" xorm:"TEXT"`
}

/*
Password string    `json:"-"`
Nickname string    `json:"-"`
Birthday time.Time `json:"-"`
Skype    string    `json:"-"`
Qq       string    `json:"-"`
Yahoo    string    `json:"-"`
Gtalk    string    `json:"-"`
Wangwang string    `json:"-"`
Address  string    `json:"-"`
Zipcode  string    `json:"-"`
Visits   string    `json:"-"`
Ip       string    `json:"-"`
Last     string    `json:"-"`
Ping     string    `json:"-"`
Fails    string    `json:"-"`
Join     time.Time `json:"-"`
Locked   time.Time `json:"-"`
Deleted  string    `json:"-"`
*/
