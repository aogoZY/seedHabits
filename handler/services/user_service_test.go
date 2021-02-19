package services

import (
	"seedHabits/conf"
	"seedHabits/handler/dao"
	"seedHabits/sdk/jsonutil"
	"seedHabits/sdk/log"
	"testing"
)

func init() {
	conf.Init("/Users/zhouyang/go/src/seedHabits/conf/conf.toml")
	log.Init(conf.Config.LogConfig)
	dao.Init()
}

func TestAuthenticationUser(t *testing.T) {
	service := GetUserService()
	res, err := service.QueryByNameAndPassword("aogo", "123")
	t.Log(res)
	log.Logger.Info(jsonutil.GetIndexJson(res))
	t.Error(err)
}

func TestInsert(t *testing.T) {
	service := GetUserService()
	insertVo := new(dao.TUsers)
	insertVo.Password = "12"
	insertVo.Name = "fansilin"
	err := service.Insert(insertVo)
	t.Log(err)
}
