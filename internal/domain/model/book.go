package model

type Book struct {
	BaseModel
	Title   string `json:"title"`
	Author  string `json:"author"`
	Summary string `json:"summary"`
	Img     string `json:"img"`
}

// TableName 覆盖gorm的命名策略
func (Book) TableName() string {
	return "book"
}
