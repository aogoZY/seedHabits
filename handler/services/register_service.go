package services

import (
	"errors"
	"fmt"
	"seedHabits/handler/dao"
	"time"
)

func UserRegister(params dao.TUsers) (msg string, err error) {
	registerFlag, err := queryRegister(params.Name)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if registerFlag {
		msg = "u have registered!"
		return msg, nil
	}
	res, UserId, err := insertRegister(params.Name, params.Password)
	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")
	err = insertUserHabitInfo(UserId, timeNow)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if res {
		msg = "register success!"
		return msg, nil
	}
	return "", nil
}

func insertRegister(name string, pwd string) (bool, int, error) {
	var registerUser dao.TUsers
	var sampleid int
	registerUser.Name = name
	registerUser.Password = pwd
	affected, err :=  dao.DBX.Insert(&registerUser)
	if err != nil {
		fmt.Println(err)
		return false, 0, err
	}
	if affected > 0 {
		_, err =  dao.DBX.Where("name = ?",name).Cols("sampleid").Get(&sampleid)
		if err != nil {
			fmt.Println(err)
			return false, 0, err
		}
		return true, sampleid, nil
	}
	return false, 0, nil
}

func queryRegister(name string) (bool, error) {
	has, err :=  dao.DBX.Table("users").Where("name=?", name).Exist()
	if err != nil {
		return false, err
	}
	if has {
		return true, nil
	}
	return false, nil
}

func insertUserHabitInfo(id int, time string) error {
	userHabitInfos := make([]dao.Info, 3)
	userHabitInfos[0].HabitId = 4
	userHabitInfos[0].HabitName = "记账"
	userHabitInfos[0].UserId = id
	userHabitInfos[0].CreateTime = time

	userHabitInfos[1].HabitId = 5
	userHabitInfos[1].HabitName = "打代码"
	userHabitInfos[1].UserId = id
	userHabitInfos[1].CreateTime = time

	userHabitInfos[2].HabitId = 7
	userHabitInfos[2].HabitName = "读书"
	userHabitInfos[2].UserId = id
	userHabitInfos[2].CreateTime = time

	affected, err := dao.DBX.Insert(&userHabitInfos)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected == 3 {
		return nil
	}
	return errors.New("自动插入三条数据失败")
}
