package dal

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sqlx.DB
}

type DatabaseConfig struct {
	Driver      string
	Host        string
	Port        int
	User        string
	Password    string
	DBName      string
	MaxOpenConn int
}

func connDB(cfg *DatabaseConfig) (*Database, error) {
	var Dbinstance = Database{}
	connStr := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=disable", cfg.DBName, cfg.Host, cfg.User, cfg.Password, cfg.Port)
	db, err := sqlx.Open(cfg.Driver, connStr)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	if err != nil || db == nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil || db == nil {
		return nil, err
	}
	Dbinstance.DB = db
	return &Dbinstance, nil
}
