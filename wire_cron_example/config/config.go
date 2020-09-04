package config

type Server struct {
	Mysql Mysql `json:"mysql" yaml:"mysql"`
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
