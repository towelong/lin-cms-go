package model

type Permission struct {
	BaseModel
	Name   string `json:"name"`
	Module string `json:"module"`
	Mount  int    `json:"mount"`
}
