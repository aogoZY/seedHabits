package services

import (
	"errors"
	"fmt"
	"seedHabits/handler/dao"
)

type userService struct {
	Dao dao.TUsers
}

func GetUserService() *userService {
	return &userService{}
}

func (d *userService) QueryByNameAndPassword(name string, password string) (res *dao.TUsers, err error) {
	var userDao dao.TUsers
	has, err := dao.DBX.Where("name=? and password=?",name,password).Get(&userDao)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if has {
		return &userDao, nil
	}
	return nil, nil
}

func (d *userService) Insert(insertVo *dao.TUsers) error {
	has, err := dao.DBX.Table(d.Dao.TableName()).Omit("id").Insert(insertVo)
	fmt.Println(has)
	if err != nil {
		return err
	}
	if has > 0 {
		return nil
	}
	return errors.New("insert users failed")
}
