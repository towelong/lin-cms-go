package dto

type SearchLogDTO struct {
	BasePage
	Start   string `json:"start" form:"start" validate:"omitempty,gte=1" tag:"start"`
	End     string `json:"end" form:"end" validate:"omitempty,gte=1" tag:"end"`
	Name    string `json:"name" form:"name" validate:"omitempty,gte=1" tag:"name"`
	Keyword string `json:"keyword" form:"keyword" validate:"omitempty,gte=1" tag:"keyword"`
}
