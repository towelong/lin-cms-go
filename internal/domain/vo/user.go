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

type PurePermission struct {
	Name   string `json:"name"`
	Module string `json:"module"`
}

type UserPermissionInfo struct {
	ID          int                           `json:"id"`
	Nickname    string                        `json:"nickname"`
	Avatar      string                        `json:"avatar"`
	Email       string                        `json:"email"`
	Admin       bool                          `json:"admin"`
	Permissions []map[string][]PurePermission `json:"permissions"`
}
