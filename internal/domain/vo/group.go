package vo

type Group struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Info  string `json:"info"`
}

type GroupInfo struct {
	Group
	Permissions []Permission `json:"permissions"`
}
