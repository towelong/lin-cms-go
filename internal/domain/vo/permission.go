package vo

type Permission struct {
	ID     int   `json:"id"`
	Name   string `json:"name"`
	Module string `json:"module"`
}
