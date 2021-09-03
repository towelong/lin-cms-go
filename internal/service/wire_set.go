package service

import "github.com/google/wire"

var Set = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),

	wire.Struct(new(GroupService), "*"),
	wire.Bind(new(IGroupService), new(*GroupService)),

	wire.Struct(new(PermissionService), "*"),
	wire.Bind(new(IPermissionService), new(*PermissionService)),
)
