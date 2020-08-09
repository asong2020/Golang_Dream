package resp

type LoginResponse struct {
	User      ResponseUser `json:"user"`
	Token     string        `json:"token"`
	ExpiresAt int64         `json:"expiresAt"`
}


type ResponseUser struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar  string `json:"avatar"`
}