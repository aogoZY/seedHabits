package services

import (
	"fmt"
	"github.com/go-xorm/xorm"
)

func QueryLoginIn(db *xorm.Engine, name string, password string) (int, error) {
	var sampleid int
	_, err := db.Sql("select sampleid from users where name=? and password=?", name, password).Get(&sampleid)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return sampleid, nil
}
