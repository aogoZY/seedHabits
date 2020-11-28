package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"seedHabits/handler/dao"
	"seedHabits/handler/services"
	"strconv"
)

func GetBillLabelListHandler(c *gin.Context) {
	res, err := services.GetBillLabelList()
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"msg":  err.Error(),
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
			"msg":  err.Error(),
		})
		return
	}
	err = services.InsertBillRecord(Params)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  err.Error(),
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
	user_id, err := strconv.Atoi(userIdStr)
	fmt.Println(user_id)
	res, err := services.GetAccountRest(user_id)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
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

	res, err := services.GetTotalAndItemListByMonth(user_id, date, account_id, account_name)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
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
	res, err := services.GetPieByType(user_id, date, search_type, PayOrGet)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
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
			"msg":  err.Error(),
		})
		return
	}
	fmt.Println("type", Params.Type)
	err = services.UpdateBillItem(Params)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
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
	err := services.DeleteBillItem(user_id, item_id)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "successed!",
	})
}
