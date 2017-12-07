package model

import (
	"time"
)

type ImChatuser struct {
	Id    int       `json:"id" xorm:"not null pk autoincr MEDIUMINT(8)"`
	Cgid  string    `json:"cgid" xorm:"not null default '' unique(chatuser) index CHAR(40)"`
	User  int       `json:"user" xorm:"not null default 0 index unique(chatuser) MEDIUMINT(8)"`
	Order int       `json:"order" xorm:"not null default 0 index SMALLINT(5)"`
	Star  string    `json:"star" xorm:"not null default '0' index ENUM('0','1')"`
	Hide  string    `json:"hide" xorm:"not null default '0' index ENUM('1','0')"`
	Mute  string    `json:"mute" xorm:"not null default '0' ENUM('0','1')"`
	Join  time.Time `json:"join" xorm:"not null default '0000-00-00 00:00:00' DATETIME"`
	Quit  time.Time `json:"quit" xorm:"not null default '0000-00-00 00:00:00' DATETIME"`
}
