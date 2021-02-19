package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"seedHabits/handler/dao"
	"seedHabits/handler/services"
	"seedHabits/sdk/log"
	"strconv"
)

//@summary 用户模块
//@description 用户登陆功能
//@param LoginReq body dao.LoginReq true "用户名与密码"
//@Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
//@router /api/v1/user/login [post]
func LoginHandler(c *gin.Context) {
	var params dao.LoginReq
	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
	}
	log.Logger.Info("hello")
	log.Logger.Info(params)
	nameValidRes := services.ParamsValid(params.Name)
	passwordValidRes := services.ParamsValid(params.Password)
	if nameValidRes == false || passwordValidRes == false {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "输入参数不合法",
		})
		log.Logger.Error("aaaaa")
		return
	}
	userService := services.GetUserService()
	queryLoginRes, err := userService.QueryByNameAndPassword(params.Name, params.Password)
	if err != nil {
		c.JSON(200, gin.H{
			"msg":  err.Error(),
			"code": 1,
		})
		log.Logger.Error("bbbbb")
		return
	}
	if queryLoginRes != nil {
		c.JSON(200, gin.H{
			"msg":  "登录成功",
			"code": 0,
			"data": gin.H{"userId": queryLoginRes.Id},
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

//@summary 用户模块
//@description 查询用户信息功能
//@param id query int true "用户id"
//@Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
//@router /api/v1/user/info/get [get]
func GetUserInfoHandler(c *gin.Context) {
	id := c.Query("id")
	fmt.Println("id", id)
	id_int, _ := strconv.Atoi(id)
	res, err := services.GetUserInfo(id_int)
	if err != nil {
		c.JSON(200, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": res})

}
