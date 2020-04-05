package services

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"seedHabits/dao"
	"time"
)

func GetHabitListByUserId(db *xorm.Engine, userId int) (res []dao.UserHabits, err error) {
	HabitList := make([]dao.Info, 0)
	err = db.Where("user_id=?", userId).Find(&HabitList)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	var userHabit dao.UserHabits
	for _, item := range HabitList {
		userHabit.Id = item.HabitId
		userHabit.Name = item.HabitName
		userHabit.Img = item.HabitImg
		res = append(res, userHabit)
	}
	return res, nil
}

func GetHistoryByUserIdAndHabitId(db *xorm.Engine, user_id int, habit_id int) (res []dao.HabitHistoryRes, err error) {
	var habitHistoryInfo []dao.HabitHistoryInfo
	var habitHistoryItem dao.HabitHistoryRes
	err = db.Table("detail").Desc("update_time").Where("user_id = ? ", user_id).And("habit_id = ?", habit_id).Find(&habitHistoryInfo)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	fmt.Println(habitHistoryInfo)
	length := len(habitHistoryInfo)
	for i, item := range habitHistoryInfo {
		day := length - i
		habitHistoryItem.Day = day
		habitHistoryItem.HabitHistoryInfo = item
		fmt.Println(habitHistoryItem)
		res = append(res, habitHistoryItem)
	}
	return res, nil

}

func InsertNewHabit(db *xorm.Engine, habitName string, img string) (res int, err error) {

	sql := "insert into habit (habit_img,habit_name) values (?,?)"
	_, err = db.Exec(sql, img, habitName)

	//habit := Habit{HabitName:habitName,HabitImg:img}
	//affected, err := db.Insert(habit).Omit("habit_id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	var habit_id int
	has, err := db.Table("habit").Cols("habit_id").Where("habit_name=? and habit_img = ?", habitName, img).Get(&habit_id)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	if has {
		return habit_id, nil
	}
	return 0, errors.New("新建习惯失败")
}

func InsertInfo(db *xorm.Engine, params dao.AddHabitParams, id int) error {
	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")

	info := dao.Info{UserId: params.UserId, HabitId: id, CreateTime: timeNow, HabitName: params.HabitName, HabitImg: params.Img}
	affected, err := db.Insert(&info)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected > 0 {
		return nil
	}
	return errors.New("insert failed")
}

func UpdateDailyDetail(db *xorm.Engine, params dao.Detail) error {
	//nTime := time.Now()
	//updateTime := nTime.Format("2006-01-02 15:04:05")
	//fmt.Printf("day:%v\n", updateTime)
	//fmt.Println(params)
	//sql := "update detail set create_time = ?, word = ?,img = ? where user_id = ? and habit_id = ? and habit_name = ? and create_time > ?"
	//_, err := db.Exec(sql, params.CreateTime, params.Word, params.Img, params.UserId, params.HabitId, params.HabitName, today)

	//sql := "update detail set word = ?,img = ? where user_id = ? and habit_id = ? and habit_name = ? and create_time > ?"
	//has, err := db.Exec(sql, params.Word, params.Img, params.UserId, params.HabitId, params.HabitName, today)
	//fmt.Printf("has:%v", has)

	affected,err := db.Cols("word","img").Where("sample_id= ?",params.SampleId).Update(&params)
	fmt.Printf("affected:%v",affected)
	if err!=nil{
		fmt.Println(err)
		return err
	}
	if affected == 1{
		fmt.Printf("affected:%v",affected)
		return nil
	}
	return errors.New("update failed!")
}

// 新建打卡记录
func InserDailyDetail(db *xorm.Engine, params dao.Detail) error {
	affected, err := db.Omit("sample_id").Insert(params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected > 0{
		fmt.Println("affected:",affected)
		return nil
	}
	return errors.New("insert failed!")
}
