package service

import (
	"asong.cloud/Golang_Dream/wire_cron_example/dao"
	"asong.cloud/Golang_Dream/wire_cron_example/model"
)

type UserService struct {
	userDao *dao.UserDB
}

func NewUserService(dao *dao.UserDB) *UserService {
	return &UserService{
		userDao: dao,
	}
}

func (u *UserService)MGet(lastID,size uint64) ([]*model.User,error) {
	return u.userDao.MGet(lastID,size)
}