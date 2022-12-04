package dao

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"asong.cloud/Golang_Dream/code_demo/Gendry_demo/common"
	"asong.cloud/Golang_Dream/code_demo/Gendry_demo/config"
	"asong.cloud/Golang_Dream/code_demo/Gendry_demo/model"
)

type UserDBTest struct {
	suite.Suite
	db *UserDB
}

func Test_UserDBTest(t *testing.T) {
	suite.Run(t, new(UserDBTest))
}


func (u *UserDBTest) SetupTest() {
	conf := &config.Server{}

	err := conf.Load("../conf/config_local.yml")
	u.Nil(err)
	u.db = NewUserDB(common.MysqlClient(&conf.Mysql))
}

func (u *UserDBTest) Test_GetMethodOne() {
	condition := &model.User{
		ID:       1,
		Username: "asong",
	}
	s, err := u.db.GetMethodOne(context.Background(), condition)
	u.Nil(err)
	u.T().Log(s)
}


func (u *UserDBTest) Test_GetMethodTwo(){
	cond := map[string]interface{}{
		"username": "asong",
		"id": 1,
	}
	s,err := u.db.GetMethodTwo(context.Background(),cond)
	u.Nil(err)
	u.T().Log(s)
}

func (u *UserDBTest) Test_Query()  {
	cond := map[string]interface{}{
		"id in": []int{1,2},
	}
	s,err := u.db.Query(context.Background(),cond)
	u.Nil(err)
	for k,v := range s{
		u.T().Log(k,v)
	}
}

func (u *UserDBTest) Test_Add()  {
	cond := map[string]interface{}{
		"username": "test_add",
		"nickname": "asong",
		"password": "123456",
		"salt": "oooo",
		"avatar": "http://www.baidu.com",
		"uptime": 123,
	}
	s,err := u.db.Add(context.Background(),cond)
	u.Nil(err)
	u.T().Log(s)
}

func (u *UserDBTest) Test_Update()  {
	where := map[string]interface{}{
		"username": "asong",
	}
	data := map[string]interface{}{
		"nickname": "shuai",
	}
	err := u.db.Update(context.Background(),where,data)
	u.Nil(err)
}

func (u *UserDBTest)Test_Delete()  {
	where := map[string]interface{}{
		"username in": []string{"2","test_add"},
	}
	err := u.db.Delete(context.Background(),where)
	u.Nil(err)
}

func (u *UserDBTest) Test_CustomizeGet()  {
	sql := "SELECT * FROM users WHERE username={{username}}"
	data := map[string]interface{}{
		"username": "test_add",
	}
	user,err := u.db.CustomizeGet(context.Background(),sql,data)
	u.Nil(err)
	u.T().Log(user)
}

func (u *UserDBTest) Test_AggregateCount()  {
	where := map[string]interface{}{
		"password": "123456",
	}
	count,err := u.db.AggregateCount(context.Background(),where,"*")
	u.Nil(err)
	u.T().Log(count)
}