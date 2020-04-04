package dao

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "zhouyang"
	password = ""
	dbname   = "zhouyang"
)

var engine *xorm.Engine

func ConnectPgDB() (*xorm.Engine, error) {
	sql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	engine, err := xorm.NewEngine("postgres", sql)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	engine.ShowSQL()
	err = engine.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("pong")
	return engine, nil
}
