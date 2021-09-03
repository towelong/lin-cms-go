package model

type Group struct {
	BaseModel
	Name  string `json:"name"`
	Info  string `json:"info"`
	Level string `json:"level"`
}
