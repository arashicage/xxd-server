package daemon

import (
	"time"
	"log"
)

const (
	// ping and create log 10 second
	pingMysql = 10 * time.Second
)

type DaemonTask func(v ...interface{})

func (dt DaemonTask) Do(v ...interface{}) {
	dt(v...)
}

func (dt DaemonTask)Tick() {
	go func() {
		ticker := time.NewTicker(pingMysql)
		defer ticker.Stop()

		/*	或者来一组
		defer func() {
			mysqlPingTicker.Stop()
		}()
		*/

		for {
			select {
			case <-ticker.C:
				dt.Do(nil)
				log.Println(".")
			}
		}
	}()
}
