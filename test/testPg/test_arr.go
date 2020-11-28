package main

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"strconv"
)

var DBX *xorm.Engine

//func connectDB() *sql.DB {
//	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
//		"password=%s dbname=%s sslmode=disable",
//		host, port, user, password, dbname)
//
//	db, err := sql.Open("postgres", psqlInfo)
//	if err != nil {
//		panic(err)
//	}
//
//	err = db.Ping()
//	if err != nil {
//		panic(err)
//	}
//	return db
//}

func init() {
	var err error
	if DBX == nil {
		conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "zhouyang", "", "zhouyang")
		DBX, err = xorm.NewEngine("postgres", conn)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	err = DBX.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("connect pg success")
	DBX.ShowSQL(true)
}

//var err error
//if DBX == nil {
//	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "zhouyang", "", "zhouyang")
//	DBX, err = xorm.NewEngine("postgres", conn)
//}
//
//if err != nil {
//	fmt.Println(err)
//}
//DBX.ShowSQL()
//err = DBX.Ping()
//if err != nil {
//	fmt.Println(err)
//}
//fmt.Println("pong")

//}

type Student struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Class  int    `json:"class"`
	Gender bool   `json:"gender"`
	Age    int    `json:"age"`
}

func main() {
	res, _ := Select(1)
	fmt.Println(res)

}

type SelectRes struct {
	Name   string `json:"name"`
	Class  int    `json:"class"`
	Gender bool   `json:"gender"`
	Age    int    `json:"age"`
}

func Select(class int) (resList []SelectRes, err error) {
	fmt.Println("enter select()")
	fmt.Println(DBX)
	rows, err := DBX.Query("select name,age,class,gender from student where class = ?", class)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("-----\n", rows)
	var res SelectRes
	for _, item := range rows {
		res.Name = string(item["name"])
		res.Age, _ = strconv.Atoi(string(item["age"]))
		res.Class, _ = strconv.Atoi(string(item["class"]))
		res.Gender, _ = strconv.ParseBool(string(item["gender"]))
		resList = append(resList, res)
	}
	return
}
