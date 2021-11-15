package main

import (
	"context"
	"fmt"
	"github.com/towelong/lin-cms-go/internal/config"
	"github.com/towelong/lin-cms-go/pkg/validator"
	"log"
)

func main() {
	validator.InitValidator()
	configs, err := config.LoadConfig("../../config")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	// pprof
	Pprof(fmt.Sprintf(":%d", configs.Server.PprofPort))

	// 初始化服务
	injector, err := NewInjector()
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	ctx := context.Background()
	initServer(ctx, configs, injector)
}
