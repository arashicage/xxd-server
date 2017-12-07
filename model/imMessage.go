package model

type ImMessage struct {
	Id          int    `json:"id" xorm:"not null pk autoincr MEDIUMINT(8)"`
	Gid         string `json:"gid" xorm:"not null default '' index CHAR(40)"`
	Cgid        string `json:"cgid" xorm:"not null default '' index CHAR(40)"`
	User        string `json:"user" xorm:"not null default '' index VARCHAR(30)"`
	Date        Time   `json:"date" xorm:"not null default '0000-00-00 00:00:00' DATETIME"`
	Type        string `json:"type" xorm:"not null default 'normal' index ENUM('normal','broadcast')"`
	Content     string `json:"content" xorm:"not null TEXT"`
	Contenttype string `json:"contentType" xorm:"not null default 'text' ENUM('emotion','image','file','object','text')"`
}
