package model


type User struct {
	ID       uint64 `json:"id,omitempty" mapstructure:"id,omitempty" ddb:"id"`
	Username string `json:"username,omitempty" mapstructure:"username,omitempty" ddb:"username"`
	Nickname string `json:"nickname,omitempty" mapstructure:"nickname,omitempty" ddb:"nickname"`
	Password string `json:"password,omitempty" mapstructure:"password,omitempty" ddb:"password"`
	Salt     string `json:"salt,omitempty" mapstructure:"salt,omitempty" ddb:"salt"`
	Avatar   string `json:"avatar,omitempty"mapstructure:"avatar,omitempty" ddb:"avatar"`
	Uptime   uint64 `json:"uptime,omitempty"mapstructure:"uptime,omitempty" ddb:"uptime"`
}

func NewEmptyUser() *User {
	return &User{}
}