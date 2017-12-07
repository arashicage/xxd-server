package model

import "time"

type SysFile struct {
	Id          int    `json:"id" xorm:"not null pk autoincr MEDIUMINT(8)"`
	Pathname    string `json:"pathName" xorm:"not null CHAR(100)"`
	Title       string `json:"title" xorm:"not null CHAR(90)"`
	Extension   string `json:"extension" xorm:"not null CHAR(30)"`
	Size        int    `json:"size" xorm:"not null default 0 MEDIUMINT(8)"`
	Objecttype  string `json:"objectType" xorm:"not null index(object) CHAR(30)"`
	Objectid    int    `json:"objectId" xorm:"not null index(object) MEDIUMINT(8)"`
	Createdby   string `json:"createdBy" xorm:"not null default '' CHAR(30)"`
	Createddate time.Time   `json:"createdDate" xorm:"not null DATETIME"`
	Editor      string `json:"editor" xorm:"not null default '0' ENUM('1','0')"`
	Primary     string `json:"primary" xorm:"default '0' ENUM('1','0')"`
	Public      string `json:"public" xorm:"not null default '1' ENUM('1','0')"`
	Downloads   int    `json:"downloads" xorm:"not null default 0 MEDIUMINT(8)"`
	Extra       string `json:"extra" xorm:"not null VARCHAR(255)"`
}
