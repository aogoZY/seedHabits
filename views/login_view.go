package views

import (
	"github.com/gin-gonic/gin"
	"seedHabits/dao"
	"seedHabits/services"
)

func LoginHandler(c *gin.Context) {
	var params dao.Users
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
	}
	nameValidRes := services.ParamsValid(params.Name)
	passwordValidRes := services.ParamsValid(params.Password)
	if nameValidRes == false || passwordValidRes == false {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "输入参数不合法",
		})
		return
	}
	dbpg, _ := dao.ConnectPgDB()
	queryLoginRes, err := services.QueryLoginIn(dbpg, params.Name, params.Password)
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
}
