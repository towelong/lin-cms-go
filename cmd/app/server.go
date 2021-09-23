package main

import (
	"context"
	"fmt"
	"github.com/towelong/lin-cms-go/internal"
	"github.com/towelong/lin-cms-go/internal/config"
	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func initServer(ctx context.Context, configs config.Config, injector *internal.Injector) {
	addr := fmt.Sprintf(":%d", configs.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      injector.Engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	// 开启服务
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Println("shutdown server")
		}
	}()
	go func() {
		permissionHandleListener(injector.PermissionService)
		removePermissionListener(injector.PermissionService)
	}()
	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

// 监听新增权限
func permissionHandleListener(service service.IPermissionService) {
	for _, v := range router.RouteMetaInfo {
		if v.Mount {
			service.CreateNewPermission(v.Module, v.Permission, v.Mount)
		}
	}
}

// 监听被移除的权限
func removePermissionListener(service service.IPermissionService) {
	permissions, err := service.GetPermissions()
	if err != nil {
		// 数据库中无权限则直接退出当前逻辑
		return
	}
	// 已经被挂载的权限
	var matchID = make([]int, 0)
	for _, permission := range permissions {
		for _, routeMeta := range router.RouteMetaInfo {
			var stayedInMeta = permission.Module == routeMeta.Module && permission.Name == routeMeta.Permission
			if stayedInMeta {
				matchID = append(matchID, permission.ID)
				break
			}
		}
	}
	// 移除没有被挂载的权限，设置为未挂载
	removeErr := service.RemoveNotMountPermission(matchID)
	if removeErr != nil {
		log.Println(removeErr)
	}
}
