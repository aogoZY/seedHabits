package services

import (
	"errors"
	"fmt"
	"seedHabits/handler/dao"
)

type UserService struct {
	Dao dao.TUsers
}

func getUserService() *UserService {
	return &UserService{}
}

func (d *UserService)QueryLoginIn(name string, password string) (int, error) {
	var sampleid int
	_, err := dao.DBX.Table(d.Dao.TableName()).Where("where use_name=? and password=?", name, password).Cols("sample").Get(&sampleid)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return sampleid, nil
}

func AddUserInfo(params dao.TUsers) error {
	fmt.Println(params.Img)
	path, err := WriteFile(params.Img)
	if err != nil {
		fmt.Println(err)
		return err
	}
	params.Img = path
	//affected,err := dao.DBX.Where("sample_id=?",params.SampleId).Omit("password").Omit("sample_id").Update(&params)
	if err != nil {
		return err
	}
	//if affected == 1 {
	//	return nil
	//}
	return errors.New("insert user info failed")
}

func GetUserInfo(id int) (res dao.TUsers, err error) {
	var user dao.TUsers
	has, err := dao.DBX.Where("id=?", id).Omit("password").Get(&user)
	if err != nil {
		return res, err
	}
	if has {
		fmt.Println("has:", has)
		return user, nil
	}
	return res, errors.New("get user info failed")
}
