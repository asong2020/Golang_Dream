package config

type Server struct {
	System System `json:"system" yaml:"system"`
	Mysql  Mysql  `json:"mysql" yaml:"mysql"`
	Redis  Redis  `json:"redis" yaml:"redis"`
	Log    Log    `json:"log" yaml:"log"`
	Jwt    Jwt    `json:"jwt" yaml:"jwt"`
}

type System struct {
	Port int `json:"port" yaml:"port"`
}

type Mysql struct {
	Host     string `json:"host" yaml:"host"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Db       string `json:"db" yaml:"db"`
	Conn     struct {
		MaxIdle int `json:"maxidle" yaml:"maxidle"`
		Maxopen int `json:"maxopen" yaml:"maxopen"`
	}
}

type Redis struct {
	Addr     string `json:"addr" yaml:"addr"`
	Db       int    `json:"db" yaml:"db"`
	Password string `json:"password" yaml:"password"`
	Poolsize int    `json:"poolsize" yaml:"poolsize"`
	Cache    struct {
		Tokenexpired int `json:"tokenexpired" yaml:"tokenexpired"`
	}
}

type Log struct {
	Prefix  string `json:"prefix" yaml:"prefix"`
	LogFile bool   `json:"log_file" yaml:"log_file"`
	Stdout  string `json:"stdout" yaml:"stdout"`
	File    string `json:"file" yaml:"file"`
}

type Jwt struct {
	Signkey string `json:"signkey" yaml:"signkey"`
}
