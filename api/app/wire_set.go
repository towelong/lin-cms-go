package app

import (
	"github.com/google/wire"
	"github.com/towelong/lin-cms-go/api/app/cms"
	v1 "github.com/towelong/lin-cms-go/api/app/v1"
)

var APISet = wire.NewSet(
	wire.Struct(new(cms.AdminAPI), "*"),
	wire.Struct(new(cms.UserAPI), "*"),
	wire.Struct(new(v1.BookAPI), "*"),
)
