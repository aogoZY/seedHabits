package services

import (
	"seedHabits/handler/dao"
	"testing"
)

func TestAuthenticationUser(t *testing.T) {
	service := GetUserService()
	res, err := service.QueryByNameAndPassword("aogo","123")
	t.Log(res)
	t.Error(err)
}

func TestInsert(t *testing.T) {
	service := GetUserService()
	insertVo := new(dao.TUsers)
	insertVo.Password="12"
	insertVo.Name="fansilin"
	err := service.Insert(insertVo)
	t.Log(err)
}
