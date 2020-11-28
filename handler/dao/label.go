package dao

type TLabel struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	LabelImg   string `json:"label_img"`
}

func (d *TLabel) TableName() string {
	return "label"
}
