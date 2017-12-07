package model

type ImUsermessage struct {
	Id      int    `json:"id" xorm:"not null pk autoincr MEDIUMINT(8)"`
	Level   int    `json:"level" xorm:"not null default 3 SMALLINT(5)"`
	User    int    `json:"user" xorm:"not null default 0 index MEDIUMINT(8)"`
	Message string `json:"message" xorm:"not null TEXT"`
}
