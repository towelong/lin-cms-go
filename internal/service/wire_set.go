package service

import "github.com/google/wire"

var Set = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),

	wire.Struct(new(GroupService), "*"),
	wire.Bind(new(IGroupService), new(*GroupService)),

	wire.Struct(new(PermissionService), "*"),
	wire.Bind(new(IPermissionService), new(*PermissionService)),

	wire.Struct(new(LogService), "*"),
	wire.Bind(new(ILogService), new(*LogService)),

	wire.Struct(new(FileService), "*"),
	wire.Bind(new(IFileService), new(*FileService)),
)
