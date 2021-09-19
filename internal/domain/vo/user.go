package vo

type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Nickname string  `json:"nickname"`
	Avatar   string  `json:"avatar"`
	Email    string  `json:"email"`
	Group    []Group `json:"group"`
}

type UserInfo struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Nickname string  `json:"nickname"`
	Avatar   string  `json:"avatar"`
	Email    string  `json:"email"`
	Groups   []Group `json:"groups"`
}
