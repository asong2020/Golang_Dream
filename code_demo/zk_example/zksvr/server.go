package zksvr

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func GetConnect(zkList []string) (conn *zk.Conn) {
	conn,_,err := zk.Connect(zkList,10*time.Second)
	if err != nil{
		fmt.Errorf("err: %s",err.Error())
	}
	return conn
}

func AddNode(conn *zk.Conn)  {
	var data = []byte("test value")
	var flags int32 = 0
	acls := zk.WorldACL(zk.PermAll)
	s,err := conn.Create("/test",data,flags,acls)
	if err != nil{
		fmt.Errorf("create fail %d\n",err)
		return
	}
	fmt.Printf("create success %s",s)
}