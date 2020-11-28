package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"seedHabits/handler/dao"
	"seedHabits/handler/services"
	"strconv"
)

func LoginHandler(c *gin.Context) {
	var params dao.TUsers
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
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
	queryLoginRes, err := services.QueryLoginIn(params.Name, params.Password)
	if err != nil {
		c.JSON(200, gin.H{
			"msg":  err.Error(),
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

func AddUserInfoHandler(c *gin.Context) {
	var param dao.TUsers
	err := c.ShouldBind(&param)
	if err != nil {
		c.JSON(200, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	fmt.Println(param)
	err = services.AddUserInfo(param)
	if err != nil {
		c.JSON(200, gin.H{"code": 1, "mes": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}

func GetUserInfoHandler(c *gin.Context) {
	sample_id := c.Query("sample_id")
	fmt.Println("sample_id", sample_id)
	sample_id_int, _ := strconv.Atoi(sample_id)
	res, err := services.GetUserInfo(sample_id_int)
	if err != nil {
		c.JSON(200, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": res})

}
