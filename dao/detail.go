package dao

import "time"

type Detail struct {
	SampleId   int       `json:"sample_id"`
	CreateTime time.Time `xorm:"create_time created" json:"created_time" description:"创建时间"`
	Word       string    `json:"word"`
	Img        string    `json:"img"`
	UserId     int       `json:"user_id"`
	HabitId    int       `json:"habit_id"`
	UserName   string    `json:"user_name"`
	HabitName  string    `json:"habit_name"`
}

type PunchRequest struct {
	PunchFlag bool `json:"punch_flag"` // flag为true，代表以打卡，新建记录 / false 修改记录
	Detail
}

type HabitHistoryInfo struct {
	CreateTime string `json:"create_time"`
	Word       string `json:"word"`
	Img        string `json:"img"`
}

type HabitHistoryRes struct {
	HabitHistoryInfo
	Day int `json:"day"`
}

type AddHabitParams struct {
	UserId    int    `json:"user_id"`
	HabitName string `json:"habit_name"`
	Img       string `json:"img"`
}

type Habit struct {
	HabitId   int    `json:"habit_id"`
	HabitImg  string `json:"habit_img"`
	HabitName string `json:"habit_name"`
}

type Label struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}
