package dao

type HabitDetail struct {
	Id          string `json:"id"`
	UserHabitId string `json:"user_habit_id"`
	CreateTime  string `json:"create_time" xorm:"created"`
	UpdateTime  string `json:"update_time" xorm:"updated"`
	Word        string `json:"word"`
	Img         string `json:"img"`
}

func (d *HabitDetail) TableName() string {
	return "habit_detail"
}
