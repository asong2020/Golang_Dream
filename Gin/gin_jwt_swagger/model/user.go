package model

import (
	"time"
)

type User struct {
	Username string    `json:"username"`
	Nickname string    `json:"nickname"`
	Password string    `json:"password"`
	Salt     string    `json:"salt"`
	Avatar   string    `json:"avatar"`
	Uptime   time.Time `json:"uptime"`
}
