package dao

type TUsers struct {
	Id int    `json:"id"`
	Name     string `json:"name" `
	Password string `json:"password"`
	Slogan   string `json:"slogan"`
	Img      string `json:"img"`
}

func (d *TUsers) TableName() string {
	return "users"
}

