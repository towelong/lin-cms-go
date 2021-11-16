package service

import (
	"fmt"
	"github.com/jianfengye/collection"
	"github.com/jinzhu/copier"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"github.com/towelong/lin-cms-go/pkg"
	"github.com/towelong/lin-cms-go/pkg/response"
	"gorm.io/gorm"
	"time"
)

type IUserService interface {
	GetUserById(id int) (model.User, error)
	GetUserByUsername(username string) (model.User, error)
	GetUserPageByGroupId(groupId int, page int, count int) (*vo.Page, error)
	IsAdmin(id int) (bool, error)
	VerifyUser(username, password string) (model.User, error)
	GetRootUserId() int
	// ChangeUserPassword 用于超级管理员修改用户密码
	ChangeUserPassword(id int, newPassword string) error
	DeleteUser(id int) error
	CreateUser(dto dto.RegisterDTO) error
	CreateUsernamePasswordIdentity(userId int, username, password string) error
	// ChangePassword 用于用户本身修改密码
	ChangePassword(id int, passwordDTO dto.ChangePasswordDTO) error
	UpdateProfile(id int, dto dto.UpdateInfoDTO) error
	GetUserGroupByUserId(id int) (groups []model.Group)
	GetUserPermissionsInfo(id int) (userPermissions vo.UserPermissionInfo, err error)
	UpdateUserInfo(id int, dto dto.UpdateGroupsDTO) error
}

type UserService struct {
	DB           *gorm.DB
	GroupService GroupService
}

func (u UserService) GetRootUserId() int {
	var (
		group     model.Group
		userGroup model.UserGroup
	)
	err := u.DB.Where("level = ?", Root).First(&group).Error
	if err != nil {
		return 0
	}
	err = u.DB.Where("group_id = ?", group.ID).First(&userGroup).Error
	if err != nil {
		return 0
	}
	return userGroup.UserID
}

func (u UserService) GetUserPageByGroupId(groupId int, page int, count int) (*vo.Page, error) {
	var (
		userGroups []model.UserGroup
		users      []model.User
		usersVo    []*vo.User
	)
	p := vo.NewPage(page, count)
	rootId := u.GetRootUserId()
	// groupId = 0 返回所有的分页用户
	if groupId == 0 {
		if err := u.DB.Where("user_id <> ?", rootId).Find(&userGroups).Error; err != nil {
			return p, err
		}
	} else {
		if err := u.DB.Where("user_id <> ? AND group_id = ?", rootId, groupId).Find(&userGroups).Error; err != nil {
			return p, err
		}
	}
	userGroupCollection := collection.NewObjCollection(userGroups)
	userIds, err := userGroupCollection.Pluck("UserID").ToInts()
	if err != nil {
		fmt.Println(err)
	}
	// 若非root用户数量为0，直接返回
	if len(userIds) == 0 {
		p.SetTotal(0)
		users = make([]model.User, 0)
		p.SetItems(users)
		return p, nil
	}
	db := u.DB.Limit(count).Offset(page*count).Find(&users, userIds)
	if db.Error != nil {
		return p, err
	}
	err = copier.Copy(&usersVo, &users)
	if err != nil {
		fmt.Println(err)
	}
	for _, user := range usersVo {
		groups, err := u.GroupService.GetUserGroupByUserId(user.ID)
		if err != nil {
			continue
		}
		var groupsVo []vo.Group
		err = copier.Copy(&groupsVo, &groups)
		if err != nil {
			fmt.Println(err)
		}
		user.Groups = append(user.Groups, groupsVo...)
	}
	p.SetItems(usersVo)
	p.SetTotal(int(db.RowsAffected))
	return p, nil
}

func (u UserService) GetUserById(id int) (model.User, error) {
	var user model.User
	res := u.DB.First(&user, "id = ?", id)
	if res.RowsAffected > 0 {
		return user, nil
	}
	return user, res.Error
}

func (u UserService) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := u.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func (u UserService) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	err := u.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

func (u UserService) IsAdmin(id int) (bool, error) {
	// 先判断用户是否存在
	user, err := u.GetUserById(id)
	if err != nil {
		return false, err
	}
	// 查找root用户的分组id
	group, groupErr := u.GroupService.GetGroupByLevel(Root)
	if groupErr != nil {
		return false, groupErr
	}
	// 查询用户分组表中是否存在记录
	var userGroup model.UserGroup
	res := u.DB.Where("user_id = ? AND group_id = ?", user.ID, group.ID).First(&userGroup)
	if res.RowsAffected > 0 {
		return true, nil
	}
	return false, res.Error
}

func (u UserService) VerifyUser(username, password string) (model.User, error) {
	var (
		userIdentity model.UserIdentity
		user         model.User
	)
	db := u.DB.Where("identity_type = ? AND identifier = ?", UserPassword.String(), username).First(&userIdentity)
	if db.Error != nil {
		return user, response.NewResponse(10031)
	}
	verifyPsw := pkg.VerifyPsw(password, userIdentity.Credential)
	if verifyPsw {
		err := u.DB.Where("username = ?", username).First(&user).Error
		if err != nil {
			return user, response.NewResponse(10031)
		}
		return user, nil
	}
	return user, response.NewResponse(10032)
}

func (u UserService) ChangeUserPassword(id int, newPassword string) error {
	user, err := u.GetUserById(id)
	if err != nil {
		return response.NewResponse(10021)
	}
	var userIdentity model.UserIdentity
	db := u.DB.Where("user_id = ?", user.ID).First(&userIdentity)
	password := pkg.EncodePassword(newPassword)
	save := db.Model(&userIdentity).Update("credential", password)
	return save.Error
}

func (u UserService) DeleteUser(id int) error {
	user, err := u.GetUserById(id)
	if err != nil {
		return response.NewResponse(10021)
	}
	if u.GetRootUserId() == id {
		return response.NewResponse(10079)
	}
	// 1. 软删除user表中的数据
	u.DB.Delete(&user)
	// 2. 软删除user—identity表中的数据
	var userIdentity model.UserIdentity
	u.DB.Where("user_id = ?", user.ID).Delete(&userIdentity)
	// 3. 软删除user-group表中的数据
	var userGroup model.UserGroup
	update := u.DB.Where("user_id = ?", user.ID).Delete(&userGroup)
	return update.Error
}

func (u *UserService) CreateUser(dto dto.RegisterDTO) error {
	user, _ := u.GetUserByUsername(dto.Username)
	// 若记录存在
	if user.ID > 0 {
		return response.NewResponse(10071)
	}
	if dto.Email != "" {
		userByEmail, _ := u.GetUserByEmail(dto.Email)
		// 若记录存在
		if userByEmail.ID > 0 {
			return response.NewResponse(10076)
		}
	}
	// 开启事务
	err := u.DB.Transaction(func(tx *gorm.DB) error {
		var user model.User
		copier.Copy(&user, &dto)
		if err := tx.Select("Username", "Email").Create(&user).Error; err != nil {
			return err
		}
		// 若指定了权限分组
		if dto.GroupIds != nil && len(dto.GroupIds) > 0 {
			if err := u.GroupService.CheckGroupsExist(dto.GroupIds); err != nil {
				return err
			}
			if err := u.GroupService.CheckGroupsValid(dto.GroupIds); err != nil {
				return err
			}
			var (
				userGroups = make([]model.UserGroup, 0)
				userGroup  model.UserGroup
			)
			for _, groupId := range dto.GroupIds {
				userGroups = append(userGroups, model.UserGroup{
					UserID:  user.ID,
					GroupID: groupId,
				})
			}
			if err := tx.Model(&userGroup).Create(userGroups).Error; err != nil {
				return err
			}
		} else {
			// 未指定分组则默认Guest
			guest, _ := u.GroupService.GetGroupByLevel(Guest)
			group := model.UserGroup{
				UserID:  user.ID,
				GroupID: guest.ID,
			}
			if err := tx.Create(&group).Error; err != nil {
				return err
			}
		}
		if err := u.CreateUsernamePasswordIdentity(user.ID, dto.Username, dto.Password); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return response.NewResponse(10200)
	}
	return nil
}

func (u *UserService) CreateUsernamePasswordIdentity(userId int, username, password string) error {
	userIdentity := model.UserIdentity{
		UserID:       userId,
		Identifier:   username,
		Credential:   pkg.EncodePassword(password),
		IdentityType: UserPassword.String(),
	}
	return u.DB.Select("UserID", "Identifier", "Credential", "IdentityType").Create(&userIdentity).Error
}

func (u *UserService) ChangePassword(id int, passwordDTO dto.ChangePasswordDTO) error {
	userIdentity, err := u.GetUserIdentityById(id)
	if err != nil {
		return err
	}
	if !pkg.VerifyPsw(passwordDTO.OldPassword, userIdentity.Credential) {
		return response.NewResponse(10032)
	}
	password := pkg.EncodePassword(passwordDTO.NewPassword)
	update := u.DB.Model(&userIdentity).Update("credential", password)
	return update.Error
}

func (u *UserService) GetUserIdentityById(id int) (model.UserIdentity, error) {
	_, err := u.GetUserById(id)
	if err != nil {
		return model.UserIdentity{}, response.NewResponse(10021)
	}
	var userIdentity model.UserIdentity
	u.DB.Where("user_id = ?", id).First(&userIdentity)
	return userIdentity, nil
}

func (u *UserService) UpdateProfile(id int, dto dto.UpdateInfoDTO) error {
	user, err := u.GetUserById(id)
	if err != nil {
		return response.NewResponse(10021)
	}
	copier.CopyWithOption(&user, &dto, copier.Option{IgnoreEmpty: true})
	user.UpdateTime = time.Now()
	err = u.DB.Save(&user).Error
	if err != nil {
		return response.NewResponse(10200)
	}
	return nil
}

func (u *UserService) GetUserGroupByUserId(id int) (groups []model.Group) {
	groups, _ = u.GroupService.GetUserGroupByUserId(id)
	if len(groups) == 0 {
		groups = make([]model.Group, 0)
	}
	for i, group := range groups {
		if group.ID == u.GroupService.GetRootGroup().ID {
			groups = append(groups[:i], groups[i+1:]...)
		}
	}
	return groups
}

func (u *UserService) GetUserPermissionsInfo(id int) (userPermissions vo.UserPermissionInfo, err error) {
	user, err := u.GetUserById(id)
	if err != nil {
		return vo.UserPermissionInfo{}, response.NewResponse(10021)
	}
	userPermissions.Permissions = make([]map[string][]vo.PurePermission, 0)
	copier.Copy(&userPermissions, &user)
	if id == u.GetRootUserId() {
		userPermissions.Admin = true
	} else {
		userPermissions.Admin = false
	}
	groups := u.GetUserGroupByUserId(user.ID)
	groupIds := make([]int, 0)
	for _, group := range groups {
		groupIds = append(groupIds, group.ID)
	}
	var (
		structMap  = make(map[string][]vo.PurePermission)
		structMaps = make([]map[string][]vo.PurePermission, 0)
	)
	if len(groupIds) > 0 {
		var groupPermissions []model.GroupPermission
		u.DB.Where("group_id in ?", groupIds).Find(&groupPermissions)
		if len(groupPermissions) == 0 {
			return userPermissions, nil
		}
		var permissionIds = make([]int, 0)
		for _, groupPermission := range groupPermissions {
			permissionIds = append(permissionIds, groupPermission.PermissionID)
		}
		var permissions []model.Permission
		u.DB.Find(&permissions, permissionIds)
		for _, permission := range permissions {
			if _, ok := structMap[permission.Module]; ok {
				var purePermission vo.PurePermission
				copier.Copy(&purePermission, &permission)
				structMap[permission.Module] = append(structMap[permission.Module], purePermission)
			} else {
				var purePermissions = make([]vo.PurePermission, 0)
				var purePermission vo.PurePermission
				copier.Copy(&purePermission, &permission)
				purePermissions = append(purePermissions, purePermission)
				structMap[permission.Module] = append(structMap[permission.Module], purePermissions...)
			}
		}
		for key, value := range structMap {
			var newMap = make(map[string][]vo.PurePermission)
			newMap[key] = value
			structMaps = append(structMaps, newMap)
		}
	}
	userPermissions.Permissions = structMaps
	return userPermissions, nil
}

func (u *UserService) UpdateUserInfo(id int, dto dto.UpdateGroupsDTO) error {
	newGroupIds := dto.GroupIds
	if _, err := u.GetUserById(id); err != nil {
		return response.NewResponse(10021)
	}
	rootLevel, _ := u.GroupService.GetGroupByLevel(Root)
	for _, groupId := range newGroupIds {
		// 校验分组是否为非Root分组
		if groupId == rootLevel.ID {
			return response.NewResponse(10073)
		}
		// 校验分组是否存在
		if _, err := u.GroupService.GetGroupById(groupId); err != nil {
			return err
		}
	}
	var existGroups []model.UserGroup
	u.DB.Where("user_id = ?", id).Find(&existGroups)
	// 将existGroupIds取出来
	existGroupIds, _ := collection.NewObjCollection(existGroups).Pluck("GroupID").ToInts()
	// 创建两个int集合，一个为数据库中存在的group_id集合， 一个为前端传递的group_id集合
	existColl := collection.NewIntCollection(existGroupIds)
	newColl := collection.NewIntCollection(newGroupIds)
	// existGroupIds没有的 即为新增的
	addIds, _ := newColl.Filter(func(obj interface{}, index int) bool {
		val := obj.(int)
		return !existColl.Contains(val)
	}).ToInts()
	// newGroupIds没有的 即为删除的
	deleteIds, _ := existColl.Filter(func(obj interface{}, index int) bool {
		val := obj.(int)
		return !newColl.Contains(val)
	}).ToInts()
	// 删除existGroupIds有，而newGroupIds没有的
	if len(deleteIds) > 0 {
		u.DB.Where("group_id IN ?", deleteIds).Delete(&existGroups)
	}
	// 添加newGroupIds有，而existGroupIds没有的
	if len(addIds) > 0 {
		var userGroups []model.UserGroup
		for _, addId := range addIds {
			userGroups = append(userGroups, model.UserGroup{
				UserID:  id,
				GroupID: addId,
			})
		}
		create := u.DB.Create(&userGroups)
		return create.Error
	}
	return nil
}
