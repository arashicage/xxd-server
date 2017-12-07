package model

import (
	"time"
	"strconv"
	"strings"
)

type Time time.Time

type ImChat struct {
	Id     int    `json:"id" xorm:"not null pk autoincr MEDIUMINT(8)"`
	Gid    string `json:"gid" xorm:"not null default '' index CHAR(40)"`
	Name   string `json:"name" xorm:"not null default '' index VARCHAR(60)"`
	Type   string `json:"type" xorm:"not null default 'group' index VARCHAR(20)"`
	Admins string `json:"admins" xorm:"not null default '' VARCHAR(255)"`
	//Committers     string    `json:"committers" xorm:"not null default '' VARCHAR(255)"`
	Committers     string    `json:"-" field:"committers" xorm:"not null default '' VARCHAR(255)"`
	Subject        int       `json:"subject" xorm:"not null default 0 MEDIUMINT(8)"`
	Public         string    `json:"public" xorm:"not null default '0' index ENUM('0','1')"`
	Createdby      string    `json:"createdBy" xorm:"not null default '' index VARCHAR(30)"`
	Createddate    Time `json:"createdDate" xorm:"not null default '0000-00-00 00:00:00' DATETIME"`
	Editedby       string    `json:"editedBy" xorm:"not null default '' index VARCHAR(30)"`
	Editeddate     Time `json:"editedDate" xorm:"not null default '0000-00-00 00:00:00' DATETIME"`
	Lastactivetime Time `json:"lastActiveTime" xorm:"not null default '0000-00-00 00:00:00' DATETIME"`
}

func (t Time) MarshalJSON() ([]byte, error) {
	//$chat->editedDate == '0000-00-00 00:00:00' ? '' : strtotime($chat->editedDate);  //todo
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *Time) UnmarshalJSON(s []byte) (err error) {
	r := strings.Replace(string(s), `"`, ``, -1)

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q/1000, 0)
	return
}

func (t Time) String() string { return time.Time(t).String() }
