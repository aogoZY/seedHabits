package dao

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"seedHabits/conf"
	"seedHabits/sdk/log"
	"seedHabits/sdk/trace"
)

var DBX *xorm.Engine

//var MC *memcache.Client
var Tracer *trace.Service

func Init() {
	var err error
	if DBX == nil {
		connstr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			conf.Config.Database.Host, conf.Config.Database.Port, conf.Config.Database.User, conf.Config.Database.Password, conf.Config.Database.DBName)
		//connstr2:="host=127.0.0.1 port=5432 user=zhouyang password= dbname=zhouyang sslmode=disable"
		fmt.Println(connstr)
		//fmt.Println(connstr2)
		DBX, err = xorm.NewEngine("postgres", connstr)
		if err != nil {
			log.Logger.Error(err)
			DBX.ShowSQL()
			err = DBX.Ping()
			if err != nil || DBX == nil {
				log.Logger.Error("db init failed")
				return
			}
			log.Logger.Info("pong")
		}
		DBX.ShowSQL(true)
	}
}

func TracerInit() {
	Tracer = trace.GetService(conf.Config.Tracer)
	Tracer.Init()
	fmt.Println("tarcer init success")
	fmt.Println(Tracer)
}
