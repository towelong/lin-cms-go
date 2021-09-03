package service

import (
	"github.com/jinzhu/copier"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"gorm.io/gorm"
	"log"
)

type IPermissionService interface {
	CreateNewPermission(module, permissionName string, mount bool)
	GetPermissions() (permissions []*model.Permission, err error)
	GetStructPermissions() (map[string][]vo.Permission, error)
	UpdatePermission(permission model.Permission) error
	RemoveNotMountPermission(ids []int) error
}

type PermissionService struct {
	DB *gorm.DB
}

func (p *PermissionService) CreateNewPermission(module, permissionName string, mount bool) {
	var permission model.Permission
	db := p.DB.Where("module = ? AND name = ?  AND delete_time is null", module, permissionName).First(&permission)
	if db.RowsAffected == 0 {
		permission = model.Permission{
			Module: module,
			Name:   permissionName,
			Mount:  bool2Int(mount),
		}
		p.DB.Select("Module", "Name", "Mount").Create(&permission)
	}
	if db.RowsAffected > 0 && (int2Bool(permission.Mount) != mount) {
		permission.Mount = bool2Int(mount)
		p.DB.Save(&permission)
	}
}

func (p *PermissionService) GetPermissions() (permissions []*model.Permission, err error) {
	db := p.DB.Where("delete_time is null").Find(&permissions)
	if db.RowsAffected > 0 {
		return permissions, nil
	}
	return nil, db.Error
}

func (p *PermissionService) GetStructPermissions() (map[string][]vo.Permission, error) {
	permissions, err := p.GetPermissions()
	structPermission := make(map[string][]vo.Permission)
	for _, item := range permissions {
		_, ok := structPermission[item.Module]
		if ok {
			var permission vo.Permission
			err = copier.Copy(&permission, item)
			structPermission[item.Module] = append(structPermission[item.Module], permission)
		} else {
			permissions := make([]vo.Permission, 0)
			var permission vo.Permission
			err = copier.Copy(&permission, item)
			permissions = append(permissions, permission)
			structPermission[item.Module] = append(structPermission[item.Module], permissions...)
		}
	}
	return structPermission, err
}

func (p *PermissionService) UpdatePermission(permission model.Permission) error {
	db := p.DB.Save(&permission)
	return db.Error
}

func (p *PermissionService) RemoveNotMountPermission(ids []int) error {
	var permissions []model.Permission
	db := p.DB.Where("delete_time is null").Not(ids).Find(&permissions)
	if db.Error != nil {
		return db.Error
	}
	for _, permission := range permissions {
		permission.Mount = 0
		err := p.UpdatePermission(permission)
		if err != nil {
			log.Printf("removeNotMountPermission err is %v\n", err)
			return err
		}
	}
	return nil
}

func bool2Int(x bool) int {
	if x {
		return 1
	}
	return 0
}

func int2Bool(x int) bool {
	return x != 0
}
