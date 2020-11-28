package dao

type Info struct {
	UserId     int    `json:"user_id"`
	HabitId    int    `json:"habit_id"`
	CreateTime string `json:"create_time"`
	HabitName  string `json:"habit_name"`
	HabitImg   string `json:"habit_img"`
}

