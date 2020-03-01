package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "zhouyang"
	password = ""
	dbname   = "zhouyang"
)

func main()  {
	db :=ConnecDB()
	fmt.Println(db)
	query(db)
}

func ConnecDB() *sql.DB{
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func query(db *sql.DB){
	var id,name,time string
	sqlInfo:=fmt.Sprintf("select * from person where name=%s","awer")
	rows,err:=db.Query(sqlInfo)

	if err!= nil{
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next(){
		err:= rows.Scan(&id,&name,&time)
		if err!= nil{
			fmt.Println(err)
		}
	}

	err = rows.Err()
	if err!= nil{
		fmt.Println(err)
	}
	fmt.Println(id,name,time)
}


