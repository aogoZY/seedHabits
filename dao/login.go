package dao

type Users struct {
	Name     string `json:"name" `
	Password string `json:"password"`
}

type UserHabits struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

