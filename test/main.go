package main

import (
	"github.com/gin-gonic/gin"
	"time"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)



func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, 打卡成功")
	})
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		fmt.Println(name)
		tNow:=time.Now()
		timeNow := tNow.Format("2006-01-02 15:04:05")
		db:=connectDB()
		queryName,queryTime,err :=query(db,name)
		if err!=nil{
			c.String(200,"err:%s",err)
		}
		if queryName == ""{
			insert(db,name,timeNow)
			c.String(200, "Hello %s login success", name)
		}else{
			c.String(200,"%s hava signed at %s",queryName,queryTime)
		}
		
	})
	r.POST("/seed/user/register",func(c *gin.Context){
		nickname:=c.PostForm("nickName")
		password :=c.PostForm("password")
		db:=connectDB()
		res:=insertRegister(db,nickname,password)
		fmt.Println(res)
		if res!=nil{
			c.JSON(200, gin.H{
				"msg": "error",
				"code": 1,
			})
		}else{
			c.JSON(200,gin.H{
				"msg":"success",
				"code":0,
			})
		} 
	})


	r.Run() // listen and serve on 0.0.0.0:8080
}


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


func query(db *sql.DB,insign string)(name string,time string,err error){
	rows,err:=db.Query(" select name,time from person where name=$1",insign)

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

func insert(db *sql.DB,name string,time string)error{
	stmt,err := db.Prepare("insert into person(name,time) values($1,$2)")
	if err != nil {
		fmt.Println(err)
	}
	_,err = stmt.Exec(name,time)
	if err != nil {
		fmt.Println(err)
		return err
	}else {
		fmt.Println("insert success")
		return nil
	}
}

func insertRegister(db *sql.DB,name string,pwd string)error{
	stmt,err:=db.Prepare("insert into register(name,password) values($1,$2)")
	if err!=nil{
		fmt.Println(err)
		return err
	}
	_,err = stmt.Exec(name,pwd)
	if err!=nil{
		return err
	}
	return nil
}

func queryRegister(db *sql.DB,name string)(bool,error){
	rows,err:=db.Query(" select name register where name=$1",name)

	if err!= nil{
		return false,err
	}
	defer rows.Close()

	for rows.Next(){
		err:= rows.Scan(&name)
		if err!= nil{
			return false,err
		}
	}

	err = rows.Err()
	if err!= nil{
		return false,err
	}

	return true,nil
}



