package service

import (
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"gorm.io/gorm"
)

type IFileService interface {
	GetFileByMD5(md5 string) (model.File, error)
	CreateFile(file model.File) (model.File, error)
}

type FileService struct {
	DB *gorm.DB
}

func (f FileService) GetFileByMD5(md5 string) (model.File, error) {
	var file model.File
	db := f.DB.Where("md5 = ?", md5).First(&file)
	return file, db.Error
}

func (f FileService) CreateFile(file model.File) (model.File, error) {
	db := f.DB.Omit("create_time", "update_time").Create(&file)
	return file, db.Error
}
