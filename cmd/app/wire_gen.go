// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/towelong/lin-cms-go/api/app/cms"
	"github.com/towelong/lin-cms-go/api/app/v1"
	"github.com/towelong/lin-cms-go/internal"
	"github.com/towelong/lin-cms-go/internal/extension/file"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/pkg/db"
	"github.com/towelong/lin-cms-go/internal/pkg/jwt"
	"github.com/towelong/lin-cms-go/internal/router"
	"github.com/towelong/lin-cms-go/internal/service"
)

// Injectors from wire.go:

func NewInjector() (*internal.Injector, error) {
	gormDB, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	groupService := service.GroupService{
		DB: gormDB,
	}
	permissionService := &service.PermissionService{
		DB:           gormDB,
		GroupService: groupService,
	}
	userService := &service.UserService{
		DB:           gormDB,
		GroupService: groupService,
	}
	serviceGroupService := &service.GroupService{
		DB: gormDB,
	}
	iToken := jwt.NewJWTMaker()
	auth := middleware.Auth{
		JWT:          iToken,
		UserService:  userService,
		GroupService: serviceGroupService,
	}
	logService := &service.LogService{
		DB: gormDB,
	}
	logs := middleware.Logs{
		LogService: logService,
	}
	adminAPI := &cms.AdminAPI{
		PermissionService: permissionService,
		UserService:       userService,
		GroupService:      serviceGroupService,
		Auth:              auth,
		Logs:              logs,
	}
	userAPI := &cms.UserAPI{
		JWT:          iToken,
		UserService:  userService,
		GroupService: serviceGroupService,
		Auth:         auth,
	}
	logAPI := &cms.LogAPI{
		LogService: logService,
		Auth:       auth,
	}
	fileService := &service.FileService{
		DB: gormDB,
	}
	localUploader := &file.LocalUploader{
		FileService: fileService,
	}
	fileAPI := &cms.FileAPI{
		LocalUploader: localUploader,
		Auth:          auth,
	}
	bookAPI := &v1.BookAPI{
		Auth: auth,
	}
	routerRouter := &router.Router{
		AdminAPI: adminAPI,
		UserAPI:  userAPI,
		LogAPI:   logAPI,
		FileAPI:  fileAPI,
		BookAPI:  bookAPI,
	}
	engine := internal.InitEngine(routerRouter)
	injector := &internal.Injector{
		Engine:            engine,
		PermissionService: permissionService,
	}
	return injector, nil
}
