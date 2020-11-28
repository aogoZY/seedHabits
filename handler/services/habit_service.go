package services

import (
	"errors"
	"fmt"
	"seedHabits/handler/dao"
	"time"
)

func GetHabitListByUserId(userId int) (res []dao.UserHabits, err error) {
	HabitList := make([]dao.Info, 0)
	err = dao.DBX.Where("user_id=?", userId).Find(&HabitList)
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

func GetHistoryByUserIdAndHabitId(user_id int, habit_id int) (res []dao.HabitHistoryRes, err error) {
	var habitHistoryInfo []dao.HabitHistoryInfo
	var habitHistoryItem dao.HabitHistoryRes
	err = dao.DBX.Table("detail").Desc("update_time").Where("user_id = ? ", user_id).And("habit_id = ?", habit_id).Find(&habitHistoryInfo)
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

func InsertNewHabit(habitName string, img string) (res int, err error) {

	sql := "insert into habit (habit_img,habit_name) values (?,?)"
	_, err =  dao.DBX.Exec(sql, img, habitName)

	//habit := Habit{HabitName:habitName,HabitImg:img}
	//affected, err := db.Insert(habit).Omit("habit_id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	var habit_id int
	has, err :=  dao.DBX.Table("habit").Cols("habit_id").Where("habit_name=? and habit_img = ?", habitName, img).Get(&habit_id)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	if has {
		return habit_id, nil
	}
	return 0, errors.New("新建习惯失败")
}

func InsertInfo(params dao.AddHabitParams, id int) error {
	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")

	info := dao.Info{UserId: params.UserId, HabitId: id, CreateTime: timeNow, HabitName: params.HabitName, HabitImg: params.Img}
	affected, err :=  dao.DBX.Insert(&info)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected > 0 {
		return nil
	}
	return errors.New("insert failed")
}

func UpdateDailyDetail(params dao.Detail) error {
	//nTime := time.Now()
	//updateTime := nTime.Format("2006-01-02 15:04:05")
	//fmt.Printf("day:%v\n", updateTime)
	//fmt.Println(params)
	//sql := "update detail set create_time = ?, word = ?,img = ? where user_id = ? and habit_id = ? and habit_name = ? and create_time > ?"
	//_, err := db.Exec(sql, params.CreateTime, params.Word, params.Img, params.UserId, params.HabitId, params.HabitName, today)

	//sql := "update detail set word = ?,img = ? where user_id = ? and habit_id = ? and habit_name = ? and create_time > ?"
	//has, err := db.Exec(sql, params.Word, params.Img, params.UserId, params.HabitId, params.HabitName, today)
	//fmt.Printf("has:%v", has)
	//var base64_image_content string = params.Img
	////data:image/jpeg;base64,
	//b, _ := regexp.MatchString(`^data:\s*image\/(\w+);base64,`, base64_image_content)
	//if !b {
	//	return errors.New("image file type wrong")
	//}
	//
	//re, _ := regexp.Compile(`^data:\s*image\/(\w+);base64,`)
	//fmt.Println("re",re)
	//allData := re.FindAllSubmatch([]byte(base64_image_content), 2)
	//fmt.Println("data",allData)
	//fileType := string(allData[0][1]) //png ，jpeg 后缀获取
	//fmt.Println("fileType",fileType)
	//
	//base64Str := re.ReplaceAllString(base64_image_content, "")
	//fmt.Println("base64Str",base64Str)
	//
	//ddd, _ := base64.StdEncoding.DecodeString(base64Str) //成图片文件并把文件写入到buffer
	//fmt.Println("ddd",ddd)
	//err2 := ioutil.WriteFile("./output." + fileType, ddd, 0666)   //buffer输出到jpg文件中（不做处理，直接写到文件）
	//if err2 != nil{
	//	fmt.Println(err2)
	//	return err2
	//}
	path, err := WriteFile(params.Img)
	if err != nil {
		return err
	}
	fmt.Println(path)
	params.Img = path
	fmt.Println(params.Img)

	affected, err :=  dao.DBX.Cols("word", "img").Where("sample_id= ?", params.SampleId).Update(&params)
	fmt.Printf("affected:%v", affected)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected == 1 {
		fmt.Printf("affected:%v", affected)
		return nil
	}
	return errors.New("update failed!")
}


// 新建打卡记录
func InserDailyDetail(params dao.Detail) error {
	path, err := WriteFile(params.Img)
	if err != nil {
		return err
	}
	fmt.Println(path)
	params.Img = path
	fmt.Println(params.Img)

	affected, err :=  dao.DBX.Omit("sample_id").Insert(params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected > 0 {
		fmt.Println("affected:", affected)
		return nil
	}
	return errors.New("insert failed!")
}
