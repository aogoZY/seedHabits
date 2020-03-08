package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "zhouyang"
	password = ""
	dbname   = "zhouyang"
)

func connectDB() *sql.DB{
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
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


func query(db *sql.DB)(name string,time string,err error){
	rows,err:=db.Query(" select name,time from person where name=$1","gjkj")

	if err!= nil{
		return "","",err
	}
	defer rows.Close()
		for rows.Next(){
			err:= rows.Scan(&name,&time)
			if err!= nil{
				return "","",err
			}
		}
	
		err = rows.Err()
		if err!= nil{
			return "","",err
		}
	
		return name,time,nil
	}
	
}

func insert(db *sql.DB){
	stmt,err := db.Prepare("insert into person(name,time) values($1,$2)")
	if err != nil {
		fmt.Println(err)
	}
	_,err = stmt.Exec("amao","2020-3-2 11:23")
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println("insert success")
	}
}

func main()  {
	db:=connectDB()
	name,time,err :=query(db)
	if err!=nil{
		fmt.Printf("err:%s",err)
	}
	fmt.Printf("%s is already signed at %s",name,time)
	// insert(db)
	
}


