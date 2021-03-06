package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"seedHabits/handler/dao"
	"seedHabits/handler/views"
	"seedHabits/sdk/log"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(Cors())

	router.GET("/", views.HelloHandler) //测试

	router.POST("/user/register", views.RegisterHandler) //注册

	router.POST("/user/login", views.LoginHandler) //登录

	router.POST("/user/info/upload", views.AddUserInfoHandler) //新增用户信息

	router.GET("/user/info/get", views.GetUserInfoHandler) //查询用户信息

	router.GET("/habit/list/:userId", views.GetHabitListHandler) //返回用户习惯列表

	router.POST("/punch", views.PunchHabitHandler) //打卡

	router.POST("/punch/update", views.UpdatePunchRecordHandler) //更新当日打卡内容

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

func registerUserService(group *gin.RouterGroup) {
	user := group.Group("/user")
	{
		user.POST("/login", views.LoginHandler)             //登录
		user.POST("/register", views.RegisterHandler)       //注册
		user.POST("/info/upload", views.AddUserInfoHandler) //新增用户信息
		user.GET("/info/get", views.GetUserInfoHandler)     //查询用户信息

	}
}

func applyApiRoutes(engine *gin.Engine) {
	api := engine.Group("api/v1")
	log.Logger.Info(dao.DBX)
	log.Logger.Info(dao.Tracer)
	api.Use(dao.Tracer.Trace())

	registerUserService(api)
}

func setupRoute(g *gin.Engine) {
	applyApiRoutes(g)
	g.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	g.GET("/_healthy_check", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})
	url := ginSwagger.URL("http://localhost:8001/swagger/doc.json")
	g.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
