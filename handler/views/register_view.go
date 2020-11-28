package views

import (
	"github.com/gin-gonic/gin"
	"seedHabits/handler/dao"
	"seedHabits/handler/services"
)

func HelloHandler(c *gin.Context) {
	c.String(200, "Hello, aogo,打卡成功")
}

func RegisterHandler(c *gin.Context) {
	var params dao.TUsers
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(200, gin.H{
			"msg":  err.Error(),
			"code": 1,
		})
		return
	}
	resName := services.ParamsValid(params.Name)
	resPassword := services.ParamsValid(params.Password)
	if resName == false || resPassword == false {
		c.JSON(200, gin.H{
			"msg":  "输入参数不合法",
			"code": 1,
		})
		return
	}
	res, err := services.UserRegister(params)
	if err != nil {
		c.JSON(200, gin.H{
			"msg":  err.Error(),
			"code": 1,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  res,
		"code": 0,
	})
}
