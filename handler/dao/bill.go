package dao

type BillRecord struct {
	SampleId    int     `json:"sample_id"`
	UserId      int     `json:"user_id"`
	Type        int     `json:"type"`       // 0 支出 1 收入
	AccountId   int     `json:"account_id"` // 1、微信 2、 支付宝 3、银行卡
	AccountName string  `json:"account_name"`
	Money       float64 `json:"money"`
	LabelId     int     `json:"label_id"`
	LabelName   string  `json:"label_name"`
	LabelImg    string  `json:"label_img"`
	Comment     string  `json:"comment"`
	CreateTime   string  `json:"create_time"`
}

type AccountRestResult struct {
	AccountList []AccountPayment `json:"account_list"`
	TotalRest   float64          `json:"total_rest"`
}

type AccountPayment struct {
	Name string  `json:"name"`
	Img  string  `json:"img"`
	Rest float64 `json:"rest"`
}

type Account struct {
	SampleId    int    `json:"sample_id"`
	AccountName string `json:"account_name"`
	AccountImg  string `json:"account_img"`
}

type GetItemByAccountNameRes struct {
	Rest     float64      `json:"rest"`
	Income   float64      `json:"income"`
	Pay      float64      `json:"pay"`
	ItemList []BillRecord `json:"item_list"`
}

type PieRes struct {
	Name    string  `json:"name"`
	Value   float64 `json:"value"`
	Percent float64 `json:"percent"`
}
