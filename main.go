package models

import (
	_ "database/sql/driver"
	"fmt"
	db "newland/database"
)

type Person struct {
	Id   int    `json:"id" form:"id"`
	Name string `json:"name"`
	Time string `json:"time"`
}

func (p *Person) AddPerson() (successed bool,err error) {
	_, err = db.SqlDB.Exec("INSERT INTO person(id,name,time) VALUES (?, ?)", p.Id, p.Name,p.Time)
	if err != nil {
		return false,err
	}
	return true,nil
}

func (p *Person) GetPersons() (existed bool, err error) {
	persons = make([]Person, 0)
	sql := fmt.Sprintf("select * from person where name = %s",p.Name)
	rows, err := db.SqlDB.Exec(sql)
	defer rows.Close()

	if err != nil {
		return false,err
	}
	if rows<1 {
		return false,nil
	}
	return true,nil
}
