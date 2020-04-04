package services

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"seedHabits/dao"
	"time"
)

func UserRegister(params dao.Users) (msg string, err error) {
	dbpg, _ := dao.ConnectPgDB()
	registerFlag, err := queryRegister(dbpg, params.Name)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if registerFlag {
		msg = "u have registered!"
		return msg, nil
	}
	res, UserId, err := insertRegister(dbpg, params.Name, params.Password)
	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")
	err = insertUserHabitInfo(dbpg, UserId, timeNow)
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

func insertRegister(db *xorm.Engine, name string, pwd string) (bool, int, error) {
	var registerUser dao.Users
	var sampleid int
	registerUser.Name = name
	registerUser.Password = pwd
	affected, err := db.Insert(&registerUser)
	if err != nil {
		fmt.Println(err)
		return false, 0, err
	}
	if affected > 0 {
		_, err = db.Sql("select sampleid from users where name=?", name).Get(&sampleid)
		if err != nil {
			fmt.Println(err)
			return false, 0, err
		}
		return true, sampleid, nil
	}
	return false, 0, nil
}

func queryRegister(db *xorm.Engine, name string) (bool, error) {
	has, err := db.Table("users").Where("name=?", name).Exist()
	if err != nil {
		return false, err
	}
	if has {
		return true, nil
	}
	return false, nil
}

func insertUserHabitInfo(db *xorm.Engine, id int, time string) error {
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

	affected, err := db.Insert(&userHabitInfos)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected == 3 {
		return nil
	}
	return errors.New("自动插入三条数据失败")
}
