package dao

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"asong.cloud/Golang_Dream/wire_cron_example/config"
)

type ClientDB struct {
	client *sql.DB
}

func NewClientDB(m *config.Mysql) *sql.DB {
	connInfo := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", m.Username, m.Password, m.Host, m.Db)

	client, err := sql.Open("mysql", connInfo)
	if err != nil {
		log.Println("init mysql err")
	}
	err = client.Ping()
	if err != nil {
		log.Println("ping mysql err")
	}
	client.SetMaxIdleConns(m.Conn.MaxIdle)
	client.SetMaxOpenConns(m.Conn.Maxopen)
	client.SetConnMaxLifetime(5 * time.Minute)
	fmt.Println("init mysql success")
return client
}