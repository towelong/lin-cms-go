package service

import (
	"github.com/jinzhu/copier"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"github.com/towelong/lin-cms-go/pkg/response"
	"github.com/towelong/lin-cms-go/pkg/router"
	"gorm.io/gorm"
	"time"
)

type IGroupService interface {
	GetRootGroup() (group model.Group)
	GetGroupByLevel(level int) (group *model.Group, err error)
	GetGroupByName(name string) (group model.Group, err error)
	GetUserHasPermission(useId int, meta router.Meta) bool
	GetUserGroupByUserId(userId int) ([]model.Group, error)
	CheckGroupsValid(ids []int) error
	CheckGroupsExist(ids []int) error
	CheckGroupExistById(id int) error
	GetPageGroups(page dto.BasePage) *vo.Page
	GetAllGroups() []vo.Group
	GetGroupById(id int) (groupInfo vo.GroupInfo, err error)
	CreateGroup(groupDTO dto.NewGroupDTO) error
	UpdateGroup(id int, groupDTO dto.UpdateGroupDTO) error
	DeleteGroup(id int) error
}

type GroupService struct {
	DB *gorm.DB
}

func (g *GroupService) GetRootGroup() (group model.Group) {
	err := g.DB.Where("level = ?", Root).First(&group).Error
	if err != nil {
		return model.Group{}
	}
	return group
}

func (g *GroupService) GetGroupByLevel(level int) (group *model.Group, err error) {
	res := g.DB.Where("level = ?", level).First(&group)
	err = res.Error
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (g *GroupService) GetGroupByName(name string) (group model.Group, err error) {
	err = g.DB.Where("name = ?", name).First(&group).Error
	return group, err
}

func (g *GroupService) GetGroupById(id int) (groupInfo vo.GroupInfo, err error) {
	var group model.Group
	res := g.DB.First(&group, id)
	err = res.Error
	if err != nil {
		return vo.GroupInfo{}, response.NewResponse(10024)
	}
	var groupPermissions []model.GroupPermission
	if err = g.DB.Where("group_id = ?", id).Find(&groupPermissions).Error; err != nil {
		groupInfo.Permissions = make([]vo.Permission, 0)
	}
	var ids []int
	for _, groupPermission := range groupPermissions {
		ids = append(ids, groupPermission.PermissionID)
	}
	var permissions = make([]model.Permission, 0)
	if len(ids) > 0 {
		g.DB.Find(&permissions, ids)
	}
	copier.Copy(&groupInfo.Permissions, &permissions)
	copier.Copy(&groupInfo, &group)
	return groupInfo, nil
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

func (g *GroupService) GetPageGroups(page dto.BasePage) *vo.Page {
	var groups []model.Group
	var newGroups = make([]vo.Group, 0)
	newPage := vo.NewPage(page.Page, page.Count)
	rootLevel, _ := g.GetGroupByLevel(Root)
	db := g.DB.Limit(page.Count).Offset(page.Page*page.Count).Where("id <> ?", rootLevel.ID).Find(&groups)
	newPage.Total = int(db.RowsAffected)
	copier.CopyWithOption(&newGroups, &groups, copier.Option{IgnoreEmpty: true})
	newPage.Items = newGroups
	return newPage
}

func (g *GroupService) GetAllGroups() []vo.Group {
	var groups []model.Group
	var newGroups = make([]vo.Group, 0)
	rootGroup, _ := g.GetGroupByLevel(Root)
	g.DB.Where("id <> ?", rootGroup.ID).Find(&groups)
	copier.CopyWithOption(&newGroups, &groups, copier.Option{IgnoreEmpty: true})
	return newGroups
}

func (g *GroupService) CreateGroup(groupDTO dto.NewGroupDTO) error {
	var group model.Group
	copier.Copy(&group, &groupDTO)
	create := g.DB.Select("Name", "Info").Create(&group)
	if create.Error != nil {
		return response.NewResponse(10200)
	}
	if len(groupDTO.PermissionIds) > 0 {
		for _, permissionId := range groupDTO.PermissionIds {
			var permission model.Permission
			if err := g.DB.First(&permission, permissionId).Error; err != nil {
				return response.NewResponse(10231)
			}
			groupPermission := model.GroupPermission{
				GroupID:      group.ID,
				PermissionID: permissionId,
			}
			g.DB.Create(&groupPermission)
		}
	}
	return nil
}

func (g *GroupService) UpdateGroup(id int, groupDTO dto.UpdateGroupDTO) error {
	var group model.Group
	err := g.DB.First(&group, id).Error
	if err != nil {
		return response.NewResponse(10024)
	}
	// 若分组名称被修改则校验名称是否重复
	if group.Name != groupDTO.Name {
		_, err := g.GetGroupByName(groupDTO.Name)
		if err != nil {
			return response.NewResponse(10072)
		}
	}
	copier.CopyWithOption(&group, &groupDTO, copier.Option{IgnoreEmpty: true})
	group.UpdateTime = time.Now()
	err = g.DB.Save(&group).Error
	return err
}

func (g *GroupService) DeleteGroup(id int) error {
	_, err := g.GetGroupById(id)
	if err != nil {
		return response.NewResponse(10024)
	}
	root, _ := g.GetGroupByLevel(Root)
	guest, _ := g.GetGroupByLevel(Guest)
	if root.ID == id {
		return response.NewResponse(10074)
	}
	if guest.ID == id {
		return response.NewResponse(10075)
	}
	var userGroups []model.UserGroup
	g.DB.Where("group_id = ?", id).Find(&userGroups)
	if len(userGroups) > 0 {
		return response.NewResponse(10027)
	}
	// TODO: 关联的权限还未删除
	return g.DB.Delete(&model.Group{}, id).Error
}
