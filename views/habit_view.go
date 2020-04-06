package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"seedHabits/dao"
	"seedHabits/services"
	"strconv"
)

func GetHabitListHandler(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, _ := strconv.Atoi(userIdStr)

	dbpg, _ := dao.ConnectPgDB()
	res, err := services.GetHabitListByUserId(dbpg, userId)
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
}

func PunchHabitHandler(c *gin.Context) {
	var param dao.Detail
	err := c.ShouldBindJSON(&param)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"msg":  err,
			"code": 1,
		})
		return
	}
	fmt.Println("habit_id :", param.HabitId)
	fmt.Println("user_id:", param.UserId)
	dbpg, _ := dao.ConnectPgDB()
	err = services.InserDailyDetail(dbpg, param)
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
}

func UpdatePunchRecordHandler(c *gin.Context) {
	var param dao.Detail
	err := c.BindJSON(&param)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
		return
	}
	dbpg, _ := dao.ConnectPgDB()
	err = services.UpdateDailyDetail(dbpg, param)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"msg":err,
			"code":1,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "update punch info successed!",
		"code": 0,
	})
}

func GetHabitHistoryHandler(c *gin.Context) {
	userIdStr := c.Query("user_id")
	habitIdStr := c.Query("habit_id")
	fmt.Println(userIdStr, habitIdStr)
	user_id, _ := strconv.Atoi(userIdStr)
	habit_id, _ := strconv.Atoi(habitIdStr)

	db, _ := dao.ConnectPgDB()

	res, err := services.GetHistoryByUserIdAndHabitId(db, user_id, habit_id)
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
}

func AddHabitHandler(c *gin.Context) {
	var params dao.AddHabitParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
		return
	}

	db, _ := dao.ConnectPgDB()
	res, err := services.InsertNewHabit(db, params.HabitName, params.Img)
	fmt.Println("habit_id", res)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
		return
	}

	err = services.InsertInfo(db, params, res)
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
}
