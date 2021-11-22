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
	configs := config.LoadConfig()
	log.NewCustomerLogger()
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
