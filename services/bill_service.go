package services

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"seedHabits/dao"
	"strconv"
)

var PaymentList = []int{1, 2, 3} // 1、微信 2、支付宝 3、银行卡

const (
	Income int = 1 //收入
	Pay    int = 0 //支出
)

func GetBillLabelList(db *xorm.Engine) (res []dao.Label, err error) {
	err = db.Find(&res)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	return res, nil
}

func InsertBillRecord(db *xorm.Engine, params dao.BillRecord) (err error) {
	affected, err := db.Omit("sample_id").Insert(&params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected > 0 {
		return nil
	}
	return errors.New("insert failed!")
}

func GetAccountRest(db *xorm.Engine, user_id int) (res dao.AccountRestResult, err error) {
	billRecord := new(dao.BillRecord)
	fmt.Printf("billRecord:%+v\n", billRecord)
	var accountPayment dao.AccountPayment
	var accountPaymentList []dao.AccountPayment
	var total float64
	for _, paymentItem := range PaymentList {
		fmt.Println(paymentItem)
		GetMoney, err := db.Where("account_id = ? and type = ? and user_id = ?", paymentItem, Income, user_id).Sum(billRecord, "money")
		fmt.Println(GetMoney)
		if err != nil {
			fmt.Println(err)
			return res, err
		}
		PayMoney, err := db.Where("account_id = ? and type = ? and user_id = ?", paymentItem, Pay, user_id).Sum(billRecord, "money")
		fmt.Println(PayMoney)
		RestbyPaymentItem := GetMoney - PayMoney
		fmt.Println(RestbyPaymentItem)
		accountPayment.Rest = RestbyPaymentItem
		account := &dao.Account{SampleId: paymentItem}
		_, err = db.Get(account)
		if err != nil {
			fmt.Println(err)
			return res, err
		}

		accountPayment.Img = account.AccountImg
		accountPayment.Name = account.AccountName
		accountPaymentList = append(accountPaymentList, accountPayment)
		total += RestbyPaymentItem
	}
	res.AccountList = accountPaymentList
	res.TotalRest = total
	return res, nil
}

func GetTotalAndItemListByMonth(db *xorm.Engine, user_id int, date string, account_id int, account_name string) (res dao.GetItemByAccountNameRes, err error) {
	billRecord := make([]dao.BillRecord, 0)
	startDate, endDate, _ := GetStartDayAndEndDayByMonth(date)
	session := db.Desc("create_time").Where("user_id =? and create_time > ? and create_time < ?", user_id, startDate, endDate)
	if account_id != 0 && account_name != "" {
		session = session.Where("account_id = ? and account_name = ?", account_id, account_name)
	}
	err = session.Find(&billRecord)
	if err != nil {
		fmt.Println(err)
		return res, err
	}

	fmt.Printf(" billRecord: %+v\n", billRecord)
	fmt.Println("----------分割线------------")
	res.ItemList = billRecord
	bill := new(dao.BillRecord)

	sessionPay := db.Where("user_id =? and create_time > ? and create_time < ?  and type = ?", user_id, startDate, endDate, Pay)
	if account_id != 0 && account_name != "" {
		sessionPay = sessionPay.Where("account_id = ? ", account_id)
	}
	pay, err := sessionPay.Sum(bill, "money")
	fmt.Printf("pay: %v\n", pay)

	sessionIncome := db.Where("user_id =? and create_time > ? and create_time < ?  and type = ?", user_id, startDate, endDate, Income)
	if account_id != 0 && account_name != "" {
		sessionIncome = sessionIncome.Where("account_id=?", account_id)
	}
	income, err := sessionIncome.Sum(bill, "money")
	fmt.Printf("income:%v\n", income)
	res.Income = income
	res.Pay = pay
	res.Rest = income - pay
	fmt.Printf("rest:%v", res.Rest)
	return res, nil
}

func GetRearchWayById(id int) (value string) {
	if id == 0 {
		return "account_name"
	} else if id == 1 {
		return "label_name"
	} else if id == 2 {
		return "comment"
	} else {
		return "not supported now!"
	}
}

func GetPieByType(pg *xorm.Engine, user_id int, date string, search_type int, PayOrGet int) (res []dao.PieRes, err error) {
	startDate, endDate, _ := GetStartDayAndEndDayByMonth(date)

	search_field := GetRearchWayById(search_type)
	fmt.Println(search_field)
	billRecord := new(dao.BillRecord)
	pay, err := pg.Where("user_id =? and type = ? and create_time < ? and create_time > ?", user_id, PayOrGet, endDate, startDate).Sum(billRecord, "money")
	fmt.Printf("all total:%v\n", pay)

	sql := fmt.Sprintf("select %s,sum(money) from bill_record where user_id =%v and type = %v and create_time < '%s' and create_time>'%s' group by(%s)", search_field, user_id, PayOrGet, endDate, startDate, search_field)
	fmt.Println(sql)
	results, err := pg.SQL(sql).QueryString()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(results)
	var pie dao.PieRes
	var PieList []dao.PieRes
	var total float64

	for _, v := range results {
		pie.Name = v[search_field]
		sum := v["sum"]
		sumFloat, _ := strconv.ParseFloat(sum, 64)
		pie.Value = sumFloat
		Percent := (sumFloat / pay) * 100
		pie.Percent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", Percent), 64)
		fmt.Printf("percent:%v", pie.Percent)
		PieList = append(PieList, pie)
		total += sumFloat
	}
	fmt.Println(total)
	fmt.Println(PieList)
	res = PieList
	return res, nil

	//sum, err := pg.Table("bill_record").GroupBy("account_id").Where("user_id =?", user_id).Sum(billRecord, "money")
	//fmt.Printf("sum:%v",sum)

	//select account_name,sum(money) from  bill_record where user_id = 1 and type = 0 group by(account_name);
	//
	//	pg.GroupBy("account_id").Find(&billRecord)
}

func UpdateBillItem(pg *xorm.Engine, Params dao.BillRecord) error {
	fmt.Println("sapmleId", Params.SampleId)
	sample_id := Params.SampleId
	affected, err := pg.Cols("type", "label_id", "label_name", "label_img", "money", "account_id", "account_name", "comment").In("sample_id", sample_id).Update(&Params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected > 0 {
		return nil
	}
	return errors.New("update failed")
}

func DeleteBillItem(pg *xorm.Engine, user_id int, sample_id int) (err error) {
	billRecord := &dao.BillRecord{SampleId: sample_id, UserId: user_id}
	affected, err := pg.Delete(billRecord)
	fmt.Printf("affected:%v", affected)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if affected == 1 {
		return nil
	}
	return errors.New("delete failed")
}
