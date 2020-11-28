package dao

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"seedHabits/conf"
	"seedHabits/sdk/log"
)

//
//const (
//	host     = "127.0.0.1"
//	port     = 5432
//	user     = "postgres"
//	password = ""
//	dbname   = "postgres"
//)

var DBX *xorm.Engine
//var MC *memcache.Client

func Init() {
	var err error
	if DBX == nil {
		connstr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			conf.Config.Database.Host, conf.Config.Database.Port, conf.Config.Database.User, conf.Config.Database.Password, conf.Config.Database.DBName)
		DBX, err = xorm.NewEngine("postgres", connstr)
		if err!=nil{
			log.Logger.Error(err)
			DBX.ShowSQL()
			err=DBX.Ping()
			if err!= nil || DBX==nil{
				log.Logger.Error("db init failed")
				return
			}
			log.Logger.Info("pong")
		}
	}
}
