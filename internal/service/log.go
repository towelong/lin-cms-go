package service

import (
	"github.com/jinzhu/copier"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"gorm.io/gorm"
)

type ILogService interface {
	GetLogs(dto dto.BasePage) (page *vo.Page, err error)
	SearchLogs(dto dto.SearchLogDTO) (page *vo.Page, err error)
	GetUsers(dto dto.BasePage) (page *vo.Page, err error)
	CreateLog(dto model.Log) error
}

type LogService struct {
	DB *gorm.DB
}

func (l *LogService) CreateLog(dto model.Log) error {
	return l.DB.Omit("create_time", "update_time").Create(&dto).Error
}

func (l *LogService) GetLogs(dto dto.BasePage) (page *vo.Page, err error) {
	var logs = make([]model.Log, 0)
	db := l.DB.Order("create_time desc").Limit(dto.Count).Offset(dto.Page * dto.Count).Find(&logs)
	page = vo.NewPage(dto.Page, dto.Count)
	page.SetItems(logs)
	page.SetTotal(int(db.RowsAffected))
	return page, db.Error
}

func (l *LogService) SearchLogs(dto dto.SearchLogDTO) (page *vo.Page, err error) {
	var logs = make([]model.Log, 0)
	db := l.DB.Order("create_time desc").Limit(dto.Count).Offset(dto.Page * dto.Count)
	if dto.Keyword != "" {
		keyword := "%" + dto.Keyword + "%"
		db.Where("message like ?", keyword)
	}
	if dto.Name != "" {
		db.Where("username = ?", dto.Name)
	}
	if dto.Start != "" && dto.End != "" {
		db.Where("create_time BETWEEN ? AND ?", dto.Start, dto.End)
	}
	db.Find(&logs)
	page = vo.NewPage(dto.Page, dto.Count)
	var logVos []vo.LogVo
	copier.CopyWithOption(&logVos, &logs, copier.Option{IgnoreEmpty: true})
	page.SetItems(logVos)
	page.SetTotal(int(db.RowsAffected))
	return page, db.Error
}

func (l *LogService) GetUsers(dto dto.BasePage) (page *vo.Page, err error) {
	var (
		logs  []model.Log
		users = make([]string, 0)
	)
	db := l.DB.Model(&logs).
		Distinct("username").
		Limit(dto.Count).Offset(dto.Page*dto.Count).
		Pluck("username", &users)
	page = vo.NewPage(dto.Page, dto.Count)
	page.SetItems(users)
	page.SetTotal(int(db.RowsAffected))
	return page, db.Error
}
