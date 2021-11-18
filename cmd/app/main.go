package main

import (
	"context"
	"fmt"

	"github.com/towelong/lin-cms-go/internal/config"
	"github.com/towelong/lin-cms-go/internal/pkg/log"
	"github.com/towelong/lin-cms-go/pkg/validator"
)

func main() {
	validator.InitValidator()
	log.NewCustomerLogger()
	configs := config.LoadConfig()
	// pprof
	Pprof(fmt.Sprintf(":%d", configs.Server.PprofPort))

	// 初始化服务
	injector, err := NewInjector()
	if err != nil {
		log.Logger.Error(err.Error())
	}
	ctx := context.Background()
	initServer(ctx, configs, injector)
}
