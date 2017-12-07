package api

import (
	"log"

	//  xxx
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

// 家里环境
const db = "root:123456@tcp(192.168.100.11:3306)/ranzhi?charset=utf8&parseTime=true"

const ServerToken = "88888888888888888888888888888888"

func init() {
	initEngine()
}

// initEngine 由 int 来调用，初始化一个 engine
// 同时 daemon routine 定时检测连接是否通常，如果断了，就重新创建一个
func initEngine() {
	var err error
	engine, err = xorm.NewEngine("mysql", db)
	if err != nil {
		log.Fatal(err)
	}
	engine.Logger().SetLevel(core.LOG_INFO)
	//engine.ShowSQL(true)
}

func PingMysql(v ... interface{}) {
	if err := engine.Ping(); err != nil {
		initEngine()
	}
}

func SetShowSQL(show bool) {
	engine.ShowSQL(show)
}
