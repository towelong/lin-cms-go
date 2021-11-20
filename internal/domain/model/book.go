package model

type Book struct {
	BaseModel
	Title   string `json:"title"`
	Author  string `json:"author"`
	Summary string `json:"summary"`
	Image   string `json:"image"`
}

// TableName 覆盖gorm的命名策略
func (Book) TableName() string {
	return "book"
}
