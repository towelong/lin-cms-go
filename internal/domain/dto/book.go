package dto

type CreateOrUpdateBookDTO struct {
	Title   string `json:"title" validate:"required,gt=0,lte=50" label:"图书标题"`
	Author  string `json:"author" validate:"required,gt=0,lte=50" label:"图书作者"`
	Summary string `json:"summary" validate:"required,gt=0,lte=1000" label:"图书简介"`
	Image   string `json:"image" validate:"omitempty,gt=0,lte=100" label:"图书封面"`
}
