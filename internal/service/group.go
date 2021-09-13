package service

import (
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/pkg/response"
	"github.com/towelong/lin-cms-go/pkg/router"
	"gorm.io/gorm"
)

type IGroupService interface {
	GetGroupByLevel(level int) (group *model.Group, err error)
	GetUserHasPermission(useId int, meta router.Meta) bool
	GetUserGroupByUserId(userId int) ([]model.Group, error)
	CheckGroupsValid(ids []int) error
	CheckGroupsExist(ids []int) error
	CheckGroupExistById(id int) error
}

type GroupService struct {
	DB *gorm.DB
}

func (g *GroupService) GetGroupByLevel(level int) (group *model.Group, err error) {
	res := g.DB.Where("level = ?", level).First(&group)
	err = res.Error
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (g *GroupService) GetUserHasPermission(useId int, meta router.Meta) bool {
	var (
		userGroups       []model.UserGroup
		groupPermissions []model.GroupPermission
		permission       model.Permission
	)
	db := g.DB.Where("user_id = ?", useId).Find(&userGroups)
	if db.Error != nil {
		return false
	}
	var groupIds = make([]int, 0)
	for _, userGroup := range userGroups {
		groupIds = append(groupIds, userGroup.GroupID)
	}
	db = g.DB.Where("group_id IN ?", groupIds).Find(&groupPermissions)
	if db.Error != nil {
		return false
	}
	var permissionIds = make([]int, 0)
	for _, groupPermission := range groupPermissions {
		permissionIds = append(permissionIds, groupPermission.PermissionID)
	}
	db = g.DB.Where("name = ? AND mount = ? AND module = ? AND id IN ?", meta.Permission, bool2Int(meta.Mount), meta.Module, permissionIds).First(&permission)
	return db.Error == nil
}

func (g *GroupService) GetUserGroupByUserId(userId int) ([]model.Group, error) {
	var groups []model.Group
	err := g.DB.Raw(`SELECT g.id, g.name, g.info,g.level,
        g.create_time,g.update_time,g.delete_time
        from lin_group AS g
        WHERE
        g.delete_time IS NULL
        AND
        g.id IN
        (
        SELECT ug.group_id
        FROM lin_user AS u
        LEFT JOIN lin_user_group as ug
        ON ug.user_id = u.id
        WHERE u.id = ?
        AND u.delete_time IS NULL
        )`, userId).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *GroupService) CheckGroupExistById(id int) error {
	var group model.Group
	return g.DB.First(&group, id).Error
}

func (g *GroupService) CheckGroupsValid(ids []int) error {
	group, _ := g.GetGroupByLevel(Root)
	for _, id := range ids {
		if id == group.ID {
			return response.NewResponse(10073)
		}
	}
	return nil
}

func (g *GroupService) CheckGroupsExist(ids []int) error {
	for _, id := range ids {
		err := g.CheckGroupExistById(id)
		if err != nil {
			return response.NewResponse(10023)
		}
	}
	return nil
}
