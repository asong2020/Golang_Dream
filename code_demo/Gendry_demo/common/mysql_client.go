package common

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"

	"asong.cloud/Golang_Dream/code_demo/Gendry_demo/config"
)

func MysqlClient(conf *config.Mysql) *sql.DB {

	db, err := manager.
		New(conf.Db,conf.Username,conf.Password,conf.Host).Set(
		manager.SetCharset("utf8"),
		manager.SetAllowCleartextPasswords(true),
		manager.SetInterpolateParams(true),
		manager.SetTimeout(1 * time.Second),
		manager.SetReadTimeout(1 * time.Second),
			).Port(conf.Port).Open(true)

	if err != nil {
		fmt.Printf("init mysql err %v\n", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("ping mysql err: %v", err)
	}
	db.SetMaxIdleConns(conf.Conn.MaxIdle)
	db.SetMaxOpenConns(conf.Conn.Maxopen)
	db.SetConnMaxLifetime(5 * time.Minute)
	//scanner.SetTagName("json")  // 全局设置，只允许设置一次
	fmt.Println("init mysql successc")
	return db
}