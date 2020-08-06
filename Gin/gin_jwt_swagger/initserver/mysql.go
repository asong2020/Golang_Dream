package initserver

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
)
func InitMysql()  {
	connInfo := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", global.AsongServer.Mysql.Username, global.AsongServer.Mysql.Password, global.AsongServer.Mysql.Host, global.AsongServer.Mysql.Db)
	var err error
	global.AsongDb, err =sql.Open("mysql",connInfo)
	if err != nil{
		global.AsongLogger.WithFields(logrus.Fields{"err":err}).Error("init mysql err")
	}
	global.AsongDb.SetMaxIdleConns(global.AsongServer.Mysql.Conn.MaxIdle)
	global.AsongDb.SetMaxOpenConns(global.AsongServer.Mysql.Conn.Maxopen)
	global.AsongDb.SetConnMaxLifetime(5 * time.Minute)
}
