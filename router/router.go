package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"seedHabits/views"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(Cors())

	router.GET("/", views.HelloHandler) //测试

	router.POST("/user/register", views.RegisterHandler) //注册

	router.POST("/user/login", views.LoginHandler) //登录

	router.GET("/habit/list/:userId", views.GetHabitListHandler) //返回用户习惯列表

	router.POST("/punch", views.PunchHabitHandler) //打卡

	router.GET("habit/history", views.GetHabitHistoryHandler) //获取某习惯的打卡历史记录

	router.POST("/habit/add", views.AddHabitHandler) //新建习惯

	router.GET("/bill/label", views.GetBillLabelListHandler) //查询所有的账户类别

	router.POST("bill/add", views.RecordBillHandler) //新增记账记录

	router.GET("/account/rest", views.GetSumByAccountIdHandler) //查询各账户类型下的收支及总收支

	router.GET("/bill/item", views.GetBillListByMonthHandler) //通过指定账户类型和时间获取该期间的明细列表

	router.GET("/bill/pie", views.GetPieHandler) //获取饼状图统计

	router.POST("bill/item/update", views.UpdateBillRecordHandler) //更新某记账详情记录

	router.GET("bill/item/delete", views.DeleteBillRecordHandler) //删除某记账详情记录

	return router
}

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
