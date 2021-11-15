package service

import (
	"github.com/jinzhu/copier"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"github.com/towelong/lin-cms-go/pkg/response"
	"gorm.io/gorm"
	"log"
)

type IPermissionService interface {
	CreateNewPermission(module, permissionName string, mount bool)
	GetPermissions() (permissions []*model.Permission, err error)
	GetStructPermissions() (map[string][]vo.Permission, error)
	UpdatePermission(permission model.Permission) error
	RemoveNotMountPermission(ids []int) error
	DispatchPermission(dto dto.DispatchPermissionDTO) error
	DispatchPermissions(dto dto.DispatchPermissionsDTO) error
	RemovePermissions(dto dto.DispatchPermissionsDTO) error
}

type PermissionService struct {
	DB           *gorm.DB
	GroupService GroupService
}

func (p *PermissionService) CreateNewPermission(module, permissionName string, mount bool) {
	var permission model.Permission
	db := p.DB.Where("module = ? AND name = ?", module, permissionName).First(&permission)
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

func (p *PermissionService) GetPermissionById(id int) (model.Permission, error) {
	var permission model.Permission
	err := p.DB.First(&permission, id).Error
	if err != nil {
		return model.Permission{}, err
	}
	return permission, nil
}

func (p *PermissionService) GetPermissions() (permissions []*model.Permission, err error) {
	db := p.DB.Where("mount <> ?", 0).Find(&permissions)
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
	db := p.DB.Not(ids).Find(&permissions)
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

func (p *PermissionService) DispatchPermission(dto dto.DispatchPermissionDTO) error {
	if _, err := p.GroupService.GetGroupById(dto.GroupId); err != nil {
		return response.NewResponse(10024)
	}
	if _, err := p.GetPermissionById(dto.PermissionId); err != nil {
		return response.NewResponse(10231)
	}
	// 校验所在分组是否存在此权限
	var groupPermissionModel model.GroupPermission
	if res := p.DB.Where("group_id = ? AND permission_id = ?", dto.GroupId, dto.PermissionId).First(&groupPermissionModel); res.RowsAffected > 0 {
		return response.NewResponse(10029)
	}
	groupPermission := model.GroupPermission{GroupID: dto.GroupId, PermissionID: dto.PermissionId}
	create := p.DB.Create(&groupPermission)
	return create.Error
}

func (p *PermissionService) DispatchPermissions(dispatchPermissionsDTO dto.DispatchPermissionsDTO) error {
	for _, permissionId := range dispatchPermissionsDTO.PermissionIds {
		err := p.DispatchPermission(dto.DispatchPermissionDTO{PermissionId: permissionId, GroupId: dispatchPermissionsDTO.GroupId})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PermissionService) RemovePermissions(permissionsDTO dto.DispatchPermissionsDTO) error {
	for _, permissionId := range permissionsDTO.PermissionIds {
		if _, err := p.GetPermissionById(permissionId); err != nil {
			return response.NewResponse(10231)
		}
	}
	if _, err := p.GroupService.GetGroupById(permissionsDTO.GroupId); err != nil {
		return response.NewResponse(10024)
	}
	db := p.DB.Where("group_id = ? AND permission_id IN ?", permissionsDTO.GroupId, permissionsDTO.PermissionIds).Delete(model.GroupPermission{})
	return db.Error
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
