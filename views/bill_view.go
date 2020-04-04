package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"seedHabits/dao"
	"seedHabits/services"
	"strconv"
)

func GetBillLabelListHandler(c *gin.Context) {
	pg, _ := dao.ConnectPgDB()
	res, err := services.GetBillLabelList(pg)
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
}

func RecordBillHandler(c *gin.Context) {
	var Params dao.BillRecord
	err := c.BindJSON(&Params)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  err,
		})
		return
	}
	db, _ := dao.ConnectPgDB()
	err = services.InsertBillRecord(db, Params)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "insert bill record successed !",
	})
}

func GetSumByAccountIdHandler(c *gin.Context) {
	userIdStr := c.Query("user_id")
	pg, _ := dao.ConnectPgDB()
	user_id, err := strconv.Atoi(userIdStr)
	fmt.Println(user_id)
	res, err := services.GetAccountRest(pg, user_id)
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
}

func GetBillListByMonthHandler(c *gin.Context) {
	userIdStr := c.Query("user_id")
	user_id, _ := strconv.Atoi(userIdStr)
	date := c.Query("date")
	accountIdStr := c.DefaultQuery("account_id", "")
	account_id, _ := strconv.Atoi(accountIdStr)
	account_name := c.DefaultQuery("account_name", "")
	fmt.Printf("accpunt_id：%v", account_id)
	fmt.Printf("accpunt_name：%v", account_name)

	db, _ := dao.ConnectPgDB()
	res, err := services.GetTotalAndItemListByMonth(db, user_id, date, account_id, account_name)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "successed!",
		"data": res,
	})
}
func GetPieHandler(c *gin.Context) {
	userIdStr := c.Query("user_id")
	date := c.Query("date")
	searchTypeStr := c.Query("search_type")
	payOrGetStr := c.Query("pay")
	user_id, _ := strconv.Atoi(userIdStr)
	search_type, _ := strconv.Atoi(searchTypeStr)
	PayOrGet, _ := strconv.Atoi(payOrGetStr)

	fmt.Println(user_id, date, search_type, PayOrGet)
	pg, _ := dao.ConnectPgDB()
	res, err := services.GetPieByType(pg, user_id, date, search_type, PayOrGet)
	if err != nil {
		fmt.Println(err)
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
}

func UpdateBillRecordHandler(c *gin.Context) {
	var Params dao.BillRecord
	err := c.BindJSON(&Params)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
		return
	}
	fmt.Println("type", Params.Type)
	pg, _ := dao.ConnectPgDB()
	err = services.UpdateBillItem(pg, Params)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "successed !",
	})
}

func DeleteBillRecordHandler(c *gin.Context) {
	UserIdStr := c.Query("user_id")
	ItemIdStr := c.Query("sample_id")
	user_id, _ := strconv.Atoi(UserIdStr)
	item_id, _ := strconv.Atoi(ItemIdStr)
	pg, _ := dao.ConnectPgDB()
	err := services.DeleteBillItem(pg, user_id, item_id)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "successed!",
	})
}
