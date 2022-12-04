package main



func main()  {
	u := &UserSvr{}
	u.init() // 初始化配置文件

	u.Run() //启动server
}
