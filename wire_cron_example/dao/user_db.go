package dao

import (
	"database/sql"
	"fmt"
	"log"

	"asong.cloud/Golang_Dream/wire_cron_example/model"
)

type UserDB struct {
	client *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{
		client: db,
	}
}

func (u *UserDB)MGet(lastID , size uint64)  ([]*model.User,error){
	table := "users"
	query := `
SELECT
id,username,nickname,password,salt,avatar,uptime
FROM
%s
WHERE id > ?
ORDER BY id LIMIT ?
`
	query = fmt.Sprintf(query,table)
	rows,err := u.client.Query(query,lastID,size)
	if err != nil{
		log.Println(err.Error())
		return nil, err
	}
	defer func() {
		if rows != nil{
			err = rows.Close()
			if err != nil{
				log.Println(err.Error())
			}
		}
	}()
	users := make([]*model.User,0)
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID,&user.Username,&user.Nickname,&user.Password,&user.Salt,
			&user.Avatar,&user.Uptime)
		if err != nil{
			log.Println(err.Error())
			continue
		}
		users = append(users,&user)
	}
	return users,nil
}