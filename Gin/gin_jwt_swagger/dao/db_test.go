package dao

import (
	"os"
	"testing"
	"time"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/config"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/initserver"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/model"
)

func TestInsertUser(t *testing.T) {
	dir,err := os.Getwd()
	if err != nil{
		t.Error(err.Error())
	}
	config.CofParse(dir+"/config.yaml",&global.AsongServer)
	initserver.InitMysql()
	type args struct {
		user *model.User
	}
	user := &model.User{
		Username: "asong",
		Nickname: "Golang梦工厂",
		Password: "123456",
		Salt: "test",
		Avatar: "default",
		Uptime: time.Now(),
	}
	arg := args{
		user: user,
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{"asong",arg,1},// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InsertUser(tt.args.user); tt.wantErr {
				t.Errorf("InsertUser() = %v", got)
			}
		})
	}
}

