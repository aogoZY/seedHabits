package main

import (
	"github.com/gin-gonic/gin"
	"time"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/go-xorm/xorm"
	"log"
	"net/http"
	"strings"
	"errors"
	"strconv"

)



// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
   return func(c *gin.Context) {
      method := c.Request.Method
      fmt.Println(method)
      c.Header("Access-Control-Allow-Origin", "*")
      c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
      c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
      c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
      c.Header("Access-Control-Allow-Credentials", "true")

      // 放行所有OPTIONS方法，因为有的模板是要请求两次的
      if method == "OPTIONS" {
         c.AbortWithStatus(http.StatusNoContent)
      }

      // 处理请求
      c.Next()
   }
}


func main() {
	// r := gin.Default()

	r := gin.New()
	r.Use(Cors())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, aogo,打卡成功")
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
	r.POST("/user/register",func(c *gin.Context){
		var params Users
		err := c.ShouldBindJSON(&params)
		if err!=nil{
			c.JSON(200, gin.H{
				"msg": "绑定失败",
				"code": 1,
			})
			return
		}
		fmt.Printf("nickname:%s\n",params.Name)
		fmt.Printf("pwd:%s\n",params.Password)

		nickname:=params.Name
		nickname=strings.Replace(nickname," ","",-1)
		if len(nickname) == 0{
			c.JSON(200, gin.H{
				"msg": "用户名不能为空",
				"code": 1,
			})
			return
		}
		password:=params.Password
		password=strings.Replace(password," ","",-1)
		if len(password) == 0{
			c.JSON(200, gin.H{
				"msg": "密码不能为空",
				"code": 1,
			})
			return
		}
		dbpg,_:=connectPgDB()
		registerFlag,err:=queryRegister(dbpg,nickname)
		if err!=nil{
			c.JSON(200, gin.H{
				"msg": err,
				"code": 1,
			})
			return
		}
		if registerFlag{
			c.JSON(200, gin.H{
				"msg": "u have registered!",
				"code": 1,
			})
			return
		}
		res,UserId,err:=insertRegister(dbpg,nickname,password)
		fmt.Println("inser reigster res:",res)
		fmt.Println("insert register UserId:",UserId)
		fmt.Println("inser register err:",err)
		
		tNow:=time.Now()
		timeNow := tNow.Format("2006-01-02 15:04:05")
		err=insertUserHabitInfo(dbpg,UserId,timeNow)
		if err!=nil{
			c.JSON(200, gin.H{
				"msg": err,
				"code": 1,
			})
			return
		}
		if res{
			c.JSON(200,gin.H{
				"msg":"register success!",
				"code":0,
			})
		} 
	})

	r.POST("/user/login",func(c *gin.Context){
		var params Users
		err:=c.ShouldBindJSON(&params)
		if err!=nil{
			c.JSON(200,gin.H{
				"code": 1,
				"msg": err,
			})
		}
		name:=strings.Replace(params.Name," ","",-1)
		if name == ""{
			c.JSON(200, gin.H{
				"msg": "用户名不能为空",
				"code": 1,
			})
			return
		}
		password:=strings.Replace(params.Password," ","",-1)
		if password == ""{
			c.JSON(200, gin.H{
				"msg": "密码不能为空",
				"code": 1,
			})
			return
		}
		dbpg,_:=connectPgDB()
		queryLoginRes,err:=queryLoginIn(dbpg,name,password)
		if err!=nil{
			c.JSON(200,gin.H{
				"msg":err,
				"code":1,
			})
			return
		}
		if queryLoginRes>0{
			c.JSON(200,gin.H{
				"msg":"登录成功",
				"code":0,
				"data": gin.H{"userId":queryLoginRes},
			})
			return
		}
		c.JSON(200,gin.H{
			"msg":"登录失败,用户名或密码不正确",
			"code":1,
	
		})


	})
	r.GET("/habit/list/:userId", func(c *gin.Context) {
		userIdStr :=c.Param("userId")
		fmt.Println("userID:",userIdStr)
		userId, _ := strconv.Atoi(userIdStr)

		dbpg,_:=connectPgDB()
		res,err:=getHabitListByUserId(dbpg,userId)
		if err!=nil{
			c.JSON(200,gin.H{
				"msg":err,
				"code":1,
			})
			return 
		}
		c.JSON(200,gin.H{
				"msg":"successed",
				"code":0,
				"data":res,
		})
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


func connectPgDB()(*xorm.Engine,error){
	sql:=fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host,port,user,password,dbname)
	engine,err:=xorm.NewEngine("postgres",sql)
	if err!=nil{
		log.Fatal(err)
		return nil,err
	}
	engine.ShowSQL()
	err=engine.Ping()
	if err!=nil{
		return nil,err
	}
	fmt.Println("pong")
	return engine,nil

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

// func insertRegister(db *sql.DB,name string,pwd string)error{
// 	stmt,err:=db.Prepare("insert into register(name,password) values($1,$2)")
// 	if err!=nil{
// 		fmt.Println(err)
// 		return err
// 	}
// 	_,err = stmt.Exec(name,pwd)
// 	if err!=nil{
// 		return err
// 	}
// 	return nil
// }


type Users struct{
	Name string `json:"name" `
	Password string `json:"password"`

}

func insertRegister(db *xorm.Engine,name string,pwd string)(bool,int,error){
	var registerUser Users
	var sampleid int
	registerUser.Name=name
	registerUser.Password=pwd
	affected,err:=db.Insert(&registerUser)
    if err!=nil{
		fmt.Println(err)
		return false,0,err
	}
	if affected>0{
		_, err = db.Sql("select sampleid from users where name=?", name).Get(&sampleid)
		if err!=nil{
			fmt.Println(err)
			return false,0,err
		}
		return true,sampleid,nil
	}
	return false,0,nil
}




func queryRegister(db *xorm.Engine,name string)(bool,error){
	has,err:=db.Table("users").Where("name=?",name).Exist()
	if err!=nil{
		return false,err
	}
	if has{
		return true,nil
	}
	return false,nil

}

type Info struct{
	UserId int `json:"user_id"`
	HabitId int `json:"habit_id"`
	CreateTime string `json:"create_time"`
	HabitName string 	`json:"habit_name"`
}

func insertUserHabitInfo(db *xorm.Engine,id int,time string)error{
	userHabitInfos :=make([]Info,3)
	userHabitInfos[0].HabitId = 4
	userHabitInfos[0].HabitName = "记账"
	userHabitInfos[0].UserId=id
	userHabitInfos[0].CreateTime=time

	userHabitInfos[1].HabitId = 5
	userHabitInfos[1].HabitName = "打代码"
	userHabitInfos[1].UserId=id
	userHabitInfos[1].CreateTime=time


	userHabitInfos[2].HabitId = 7
	userHabitInfos[2].HabitName = "读书"
	userHabitInfos[2].UserId=id
	userHabitInfos[2].CreateTime=time

	
	affected, err := db.Insert(&userHabitInfos)
	if err!=nil{
		fmt.Println(err)
		return err
	}
	if affected == 3{
		return nil
	}
	return errors.New("自动插入三条数据失败")
}

func queryLoginIn(db *xorm.Engine,name string,password string)(int,error){
	var sampleid int
	_, err := db.Sql("select sampleid from users where name=? and password=?", name,password).Get(&sampleid)
	if err!=nil{
		fmt.Println(err)
		return 0,err
	}
	return sampleid,nil

}

type UserHabits struct{
	Id int `json:"id"`
	Name string 	`json:"name"`
	img string 	`json:"img"`
}

func getHabitListByUserId(db *xorm.Engine,userId int)(res []UserHabits,err error){
	HabitList :=make([]Info,0)
	err=db.Where("user_id=?",userId).Find(&HabitList)
	if err!=nil{
		fmt.Println(err)
		return res,err
	}
	var userHabit UserHabits
	for _,item:=range HabitList{
		userHabit.Id=item.HabitId
		userHabit.Name =item.HabitName
		res=append(res,userHabit)
	}
	return res,nil
}