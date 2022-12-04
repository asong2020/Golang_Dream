package task

import (
	"fmt"

	"asong.cloud/Golang_Dream/wire_cron_example/service"
)

type Scanner struct {
	lastID uint64
	user *service.UserService
}

const  (
	ScannerSize = 10
)

func NewScanner(user *service.UserService)  *Scanner{
	return &Scanner{
		user: user,
	}
}

func (s *Scanner)Run()  {
	err := s.scannerDB()
	if err != nil{
		fmt.Errorf(err.Error())
	}
}

func (s *Scanner)scannerDB()  error{
	s.reset()
	flag := false
	for {
		users,err:=s.user.MGet(s.lastID,ScannerSize)
		if err != nil{
			return err
		}
		if len(users) < ScannerSize{
			flag = true
		}
		s.lastID = users[len(users) - 1].ID
		for k,v := range users{
			fmt.Println(k,v)
		}
		if flag{
			return nil
		}
	}
}

func (s *Scanner)reset()  {
	s.lastID = 0
}