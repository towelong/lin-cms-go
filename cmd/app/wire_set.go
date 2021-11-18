package main

import (
	"github.com/google/wire"
	"github.com/towelong/lin-cms-go/api/app"
	"github.com/towelong/lin-cms-go/internal"
	"github.com/towelong/lin-cms-go/internal/extension/file"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/pkg/db"
	"github.com/towelong/lin-cms-go/internal/pkg/jwt"
	"github.com/towelong/lin-cms-go/internal/router"
	"github.com/towelong/lin-cms-go/internal/service"
)

var set = wire.NewSet(
	db.InitDB,
	app.APISet,
	router.Set,
	internal.InjectSet,
	internal.InitEngine,
	jwt.Set,
	service.Set,
	middleware.AuthSet,
	middleware.LogSet,
	file.Set,
)
