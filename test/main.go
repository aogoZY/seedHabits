package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
		tNow := time.Now()
		timeNow := tNow.Format("2006-01-02 15:04:05")
		db := connectDB()
		queryName, queryTime, err := query(db, name)
		if err != nil {
			c.String(200, "err:%s", err)
		}
		if queryName == "" {
			insert(db, name, timeNow)
			c.String(200, "Hello %s login success", name)
		} else {
			c.String(200, "%s hava signed at %s", queryName, queryTime)
		}

	})
	r.POST("/user/register", func(c *gin.Context) {
		var params Users
		err := c.ShouldBindJSON(&params)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "绑定失败",
				"code": 1,
			})
			return
		}
		fmt.Printf("nickname:%s\n", params.Name)
		fmt.Printf("pwd:%s\n", params.Password)

		nickname := params.Name
		nickname = strings.Replace(nickname, " ", "", -1)
		if len(nickname) == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户名不能为空",
				"code": 1,
			})
			return
		}
		password := params.Password
		password = strings.Replace(password, " ", "", -1)
		if len(password) == 0 {
			c.JSON(200, gin.H{
				"msg":  "密码不能为空",
				"code": 1,
			})
			return
		}
		dbpg, _ := connectPgDB()
		registerFlag, err := queryRegister(dbpg, nickname)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  err,
				"code": 1,
			})
			return
		}
		if registerFlag {
			c.JSON(200, gin.H{
				"msg":  "u have registered!",
				"code": 1,
			})
			return
		}
		res, UserId, err := insertRegister(dbpg, nickname, password)
		fmt.Println("inser reigster res:", res)
		fmt.Println("insert register UserId:", UserId)
		fmt.Println("inser register err:", err)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  err,
				"code": 1,
			})
			return
		}
		tNow := time.Now()
		timeNow := tNow.Format("2006-01-02 15:04:05")
		err = insertUserHabitInfo(dbpg, UserId, timeNow)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  err,
				"code": 1,
			})
			return
		}
		if res {
			c.JSON(200, gin.H{
				"msg":  "register success!",
				"code": 0,
			})
		}
	})

	r.POST("/user/login", func(c *gin.Context) {
		var params Users
		err := c.ShouldBindJSON(&params)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  err,
			})
		}
		name := strings.Replace(params.Name, " ", "", -1)
		if name == "" {
			c.JSON(200, gin.H{
				"msg":  "用户名不能为空",
				"code": 1,
			})
			return
		}
		password := strings.Replace(params.Password, " ", "", -1)
		if password == "" {
			c.JSON(200, gin.H{
				"msg":  "密码不能为空",
				"code": 1,
			})
			return
		}
		dbpg, _ := connectPgDB()
		queryLoginRes, err := queryLoginIn(dbpg, name, password)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  err,
				"code": 1,
			})
			return
		}
		if queryLoginRes > 0 {
			c.JSON(200, gin.H{
				"msg":  "登录成功",
				"code": 0,
				"data": gin.H{"userId": queryLoginRes},
			})
			return
		}
		c.JSON(200, gin.H{
			"msg":  "登录失败,用户名或密码不正确",
			"code": 1,
		})

	})
	r.GET("/habit/list/:userId", func(c *gin.Context) {
		userIdStr := c.Param("userId")
		fmt.Println("userID:", userIdStr)
		userId, _ := strconv.Atoi(userIdStr)

		dbpg, _ := connectPgDB()
		res, err := getHabitListByUserId(dbpg, userId)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  err,
				"code": 1,
			})
			return
		}
		c.JSON(200, gin.H{
			"msg":  "successed",
			"code": 0,
			"data": res,
		})
	})

	r.POST("/punch", func(c *gin.Context) {
		var param PunchRequest
		err := c.ShouldBindJSON(&param)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"msg":  err,
				"code": 1,
			})
			return
		}
		fmt.Println("habit_id :", param.Detail.HabitId)
		fmt.Println("user_id:", param.Detail.UserId)
		dbpg, _ := connectPgDB()
		NTime := time.Now()
		punchTime := NTime.Format("2006-01-02 15:04:05")
		param.Detail.CreateTime = punchTime
		fmt.Printf("%+v", param)
		fmt.Println(param.PunchFlag)

		if param.PunchFlag {
			err := UpdateDailyDetail(dbpg, param.Detail)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.JSON(200, gin.H{
				"msg":  "update punch info successed!",
				"code": 0,
			})
			return
		}
		err = InserDailyDetail(dbpg, param.Detail)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"msg":  err,
				"code": 1,
			})
			return
		}
		c.JSON(200, gin.H{
			"msg":  "insert punch info successed! Keep moving on!",
			"code": 0,
		})
	})
	r.GET("habit/history", func(c *gin.Context) {
		userIdStr := c.Query("user_id")
		habitIdStr := c.Query("habit_id")
		fmt.Println(userIdStr, habitIdStr)
		user_id, _ := strconv.Atoi(userIdStr)
		habit_id, _ := strconv.Atoi(habitIdStr)

		db, _ := connectPgDB()

		res, err := GetHistoryByUserIdAndHabitId(db, user_id, habit_id)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"msg":  err,
				"code": 1,
				"data": res,
			})
			return
		}
		c.JSON(200, gin.H{
			"msg":  "你已经打卡了好多哇！",
			"code": 0,
			"data": res,
		})
	})
	r.POST("/habit/add", func(c *gin.Context) {
		var params AddHabitParams
		err := c.ShouldBindJSON(&params)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  err,
			})
			return
		}
		db, _ := connectPgDB()
		res, err := InsertNewHabit(db, params.HabitName, params.Img)
		fmt.Println("habit_id", res)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  err,
			})
			return
		}
		err = InsertInfo(db, params, res)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  err,
			})
			return
		}
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "insert new habit successed!",
		})
		return
	})
	r.GET("/bill/label", func(c *gin.Context) {
		pg, _ := connectPgDB()
		res, err := GetLableList(pg)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"msg":  err,
				"code": 1,
			})
			return
		}
		c.JSON(200, gin.H{
			"msg":  "ok",
			"code": 0,
			"data": res,
		})

	})
	r.POST("bill/add", func(c *gin.Context) {
		var Params BillRecord
		err := c.BindJSON(&Params)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  err,
			})
			return
		}
		db, _ := connectPgDB()
		err = InsertBillRecord(db, Params)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  err,
			})
			return
		}
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "insert bill record successed !",
		})

	})
	r.GET("/account/rest", func(c *gin.Context) {
		userIdStr := c.Query("user_id")
		pg, _ := connectPgDB()
		user_id, err := strconv.Atoi(userIdStr)
		fmt.Println(user_id)
		res, err := GetAccountRest(pg, user_id)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  err,
			})
			return
		}
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "successed",
			"data": res,
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

func connectDB() *sql.DB {
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

func connectPgDB() (*xorm.Engine, error) {
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

func query(db *sql.DB, insign string) (name string, time string, err error) {
	rows, err := db.Query(" select name,time from person where name=$1", insign)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &time)
		if err != nil {
			return "", "", err
		}
	}

	err = rows.Err()
	if err != nil {
		return "", "", err
	}

	return name, time, nil
}

func insert(db *sql.DB, name string, time string) error {
	stmt, err := db.Prepare("insert into person(name,time) values($1,$2)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(name, time)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
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

type Users struct {
	Name     string `json:"name" `
	Password string `json:"password"`
}

func insertRegister(db *xorm.Engine, name string, pwd string) (bool, int, error) {
	var registerUser Users
	var sampleid int
	registerUser.Name = name
	registerUser.Password = pwd
	affected, err := db.Insert(&registerUser)
	if err != nil {
		fmt.Println(err)
		return false, 0, err
	}
	if affected > 0 {
		_, err = db.Sql("select sampleid from users where name=?", name).Get(&sampleid)
		if err != nil {
			fmt.Println(err)
			return false, 0, err
		}
		return true, sampleid, nil
	}
	return false, 0, nil
}

func queryRegister(db *xorm.Engine, name string) (bool, error) {
	has, err := db.Table("users").Where("name=?", name).Exist()
	if err != nil {
		return false, err
	}
	if has {
		return true, nil
	}
	return false, nil

}

type Info struct {
	UserId     int    `json:"user_id"`
	HabitId    int    `json:"habit_id"`
	CreateTime string `json:"create_time"`
	HabitName  string `json:"habit_name"`
	HabitImg   string `json:"habit_img"`
}

func insertUserHabitInfo(db *xorm.Engine, id int, time string) error {
	userHabitInfos := make([]Info, 3)
	userHabitInfos[0].HabitId = 4
	userHabitInfos[0].HabitName = "记账"
	userHabitInfos[0].UserId = id
	userHabitInfos[0].CreateTime = time
	userHabitInfos[0].HabitImg = "accounting"

	userHabitInfos[1].HabitId = 5
	userHabitInfos[1].HabitName = "打代码"
	userHabitInfos[1].UserId = id
	userHabitInfos[1].CreateTime = time
	userHabitInfos[1].HabitImg = "coding"

	userHabitInfos[2].HabitId = 7
	userHabitInfos[2].HabitName = "读书"
	userHabitInfos[2].UserId = id
	userHabitInfos[2].CreateTime = time
	userHabitInfos[2].HabitImg = "reading"

	affected, err := db.Insert(&userHabitInfos)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected == 3 {
		return nil
	}
	return errors.New("自动插入三条数据失败")
}

func queryLoginIn(db *xorm.Engine, name string, password string) (int, error) {
	var sampleid int
	_, err := db.Sql("select sampleid from users where name=? and password=?", name, password).Get(&sampleid)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return sampleid, nil

}

type UserHabits struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

func getHabitListByUserId(db *xorm.Engine, userId int) (res []UserHabits, err error) {
	HabitList := make([]Info, 0)
	err = db.Where("user_id=?", userId).Find(&HabitList)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	var userHabit UserHabits
	for _, item := range HabitList {
		userHabit.Id = item.HabitId
		userHabit.Name = item.HabitName
		userHabit.Img = item.HabitImg
		res = append(res, userHabit)
	}
	return res, nil
}

type Detail struct {
	SampleId   int    `json:"sample_id"`
	CreateTime string `json:"create_time"`
	Word       string `json:"word"`
	Img        string `json:"img"`
	UserId     int    `json:"user_id"`
	HabitId    int    `json:"habit_id"`
	UserName   string `json:"user_name"`
	HabitName  string `json:"habit_name"`
}

type PunchRequest struct {
	PunchFlag bool `json:"punch_flag"`
	Detail
}

type PunchInfo struct {
	UserId    int    `json:"user_id"`
	HabitId   int    `json:"habit_id"`
	Word      string `json:"word"`
	Img       string `json:"img"`
	HabitName string `json:"habit_name"`
}

// 新建打卡记录
func InserDailyDetail(db *xorm.Engine, params Detail) error {
	// detail := new(Detail)
	// detail.Word = params.Word
	// detail.Img = params.Img
	// detail.UserId=params.UserId
	// detail.HabitId=params.HabitId
	// detail.HabitName =params.HabitName

	// detail.CreateTime=punchTime

	sql := "insert into detail(create_time,word,img,user_id,habit_id,habit_name) values (?, ?, ?, ?, ?, ?)"
	res, err := db.Exec(sql, params.CreateTime, params.Word, params.Img, params.UserId, params.HabitId, params.HabitName)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(res)
	return nil
}

//修改打卡记录
func UpdateDailyDetail(db *xorm.Engine, params Detail) error {
	//sql ="update user set age = ? where name = ?"
	//res, err := engine.Exec(sql, 1, "xorm")
	//
	sql := "update detail set create_time = ?, word = ?,img = ? where user_id = ? and habit_id = ? and habit_name = ?"
	res, err := db.Exec(sql, params.CreateTime, params.Word, params.Img, params.UserId, params.HabitId, params.HabitName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("update detail res:%s", res)
	return nil
}

type HabitHistoryInfo struct {
	CreateTime string `json:"create_time"`
	Word       string `json:"word"`
	Img        string `json:"img"`
}

type HabitHistoryRes struct {
	HabitHistoryInfo
	Day int `json:"day"`
}

func GetHistoryByUserIdAndHabitId(db *xorm.Engine, user_id int, habit_id int) (res []HabitHistoryRes, err error) {
	var habitHistoryInfo []HabitHistoryInfo
	var habitHistoryItem HabitHistoryRes
	//var habitHistoryRes []HabitHistoryRes
	err = db.Table("detail").Desc("create_time").Where("user_id = ? ", user_id).And("habit_id = ?", habit_id).Find(&habitHistoryInfo)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	length := len(habitHistoryInfo)
	for i, item := range habitHistoryInfo {
		day := length - i
		habitHistoryItem.Day = day
		habitHistoryItem.HabitHistoryInfo = item
		fmt.Println(habitHistoryItem)
		res = append(res, habitHistoryItem)
	}
	return res, nil

}

type AddHabitParams struct {
	UserId    int    `json:"user_id"`
	HabitName string `json:"habit_name"`
	Img       string `json:"img"`
}

type Habit struct {
	HabitId   int    `json:"habit_id"`
	HabitImg  string `json:"habit_img"`
	HabitName string `json:"habit_name"`
}

func InsertNewHabit(db *xorm.Engine, habitName string, img string) (res int, err error) {

	sql := "insert into habit (habit_img,habit_name) values (?,?)"
	_, err = db.Exec(sql, img, habitName)

	//habit := Habit{HabitName:habitName,HabitImg:img}
	//affected, err := db.Insert(habit).Omit("habit_id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	var habit_id int
	has, err := db.Table("habit").Cols("habit_id").Where("habit_name=? and habit_img = ?", habitName, img).Get(&habit_id)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	if has {
		return habit_id, nil
	}
	return 0, errors.New("新建习惯失败")
}

func InsertInfo(db *xorm.Engine, params AddHabitParams, id int) error {
	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")
	info := Info{UserId: params.UserId, HabitId: id, CreateTime: timeNow, HabitName: params.HabitName, HabitImg: params.Img}
	affected, err := db.Insert(&info)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected > 0 {
		return nil
	}
	return errors.New("insert failed")
}

type Label struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

func GetLableList(db *xorm.Engine) (res []Label, err error) {
	err = db.Find(&res)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	return res, nil

}

type BillRecord struct {
	UserId      int       `json:"user_id"`
	Type        int       `json:"type"`       // 0 支出 1 收入
	AccountId   int       `json:"account_id"` // 1、微信 2、 支付宝 3、银行卡
	AccountName string    `json:"account_name"`
	Money       float64   `json:"money"`
	LabelId     int       `json:"label_id"`
	LabelName   string    `json:"label_name"`
	Comment     string    `json:"comment"`
	CreatTime   time.Time `xorm:"create_time created" json:"creat_time" description:"创建时间"`
}

func InsertBillRecord(db *xorm.Engine, params BillRecord) (err error) {
	affected, err := db.Insert(&params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected > 0 {
		return nil
	}
	return errors.New("insert not successed!")

}

type AccountRestResult struct {
	AccountList []AccountPayment `json:"account_list"`
	TotalRest   float64          `json:"total_rest"`
}

type AccountPayment struct {
	Name string  `json:"name"`
	Img  string  `json:"img"`
	Rest float64 `json:"rest"`
}

var PaymentList = []int{1, 2, 3} // 1、微信 2、支付宝 3、银行卡

const (
	Pay    int = 1
	Income int = 0
)

type Account struct {
	SampleId    int    `json:"sample_id"`
	AccountName string `json:"account_name"`
	AccountImg  string `json:"account_img"`
}

func GetAccountRest(db *xorm.Engine, user_id int) (res AccountRestResult, err error) {
	billRecord := new(BillRecord)
	fmt.Printf("billRecord:%+v\n", billRecord)
	var accountPayment AccountPayment
	var accountPaymentList []AccountPayment
	var total float64
	for _, paymentItem := range PaymentList {
		fmt.Println(paymentItem)
		PayMoney, err := db.Where("account_id = ? and type = ? and user_id = ?", paymentItem, Income, user_id).Sum(billRecord, "money")
		fmt.Println(PayMoney)
		if err != nil {
			fmt.Println(err)
			return res, err
		}
		GetMoney, err := db.Where("account_id = ? and type = ? and user_id = ?", paymentItem, Pay, user_id).Sum(billRecord, "money")
		fmt.Println(GetMoney)
		RestbyPaymentItem := GetMoney - PayMoney
		fmt.Println(RestbyPaymentItem)
		accountPayment.Rest = RestbyPaymentItem
		account := &Account{SampleId: paymentItem}
		_, err = db.Get(account)
		if err != nil {
			fmt.Println(err)
			return res, err
		}

		accountPayment.Img = account.AccountImg
		accountPayment.Name = account.AccountName
		accountPaymentList = append(accountPaymentList, accountPayment)
		total += RestbyPaymentItem
	}
	res.AccountList = accountPaymentList
	res.TotalRest = total
	return res, nil
}
