package service

import (
	"fmt"
	"github.com/jianfengye/collection"
	"github.com/jinzhu/copier"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"github.com/towelong/lin-cms-go/pkg"
	"github.com/towelong/lin-cms-go/pkg/response"
	"gorm.io/gorm"
)

type IUserService interface {
	GetUserById(id int) (model.User, error)
	GetUserPageByGroupId(groupId int, page int, count int) (*vo.Page, error)
	IsAdmin(id int) (bool, error)
	VerifyUser(username, password string) (model.User, error)
	GetRootUserId() int
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
	err := u.DB.Where("delete_time is null AND level = ?", Root).First(&group).Error
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
		user.Group = append(user.Group, groupsVo...)
	}
	p.SetItems(usersVo)
	p.SetTotal(int(db.RowsAffected))
	return p, nil
}

func (u UserService) GetUserById(id int) (model.User, error) {
	var user model.User
	res := u.DB.First(&user, "id = ? AND delete_time is null", id)
	if res.RowsAffected > 0 {
		return user, nil
	}
	return user, res.Error
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
	db := u.DB.Where("delete_time is null AND identity_type = ? AND identifier = ?", UserPassword.String(), username).First(&userIdentity)
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
