package dao

type TUserHabit struct {
	Id string `json:"id"`
	UserId string 	`json:"user_id"`
	HabitId string `json:"habit_id"`
	CreateTime string `json:"create_time" xorm:"created"`

}
