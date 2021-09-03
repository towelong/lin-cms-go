package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/towelong/lin-cms-go/internal/service"
)

var InjectSet = wire.NewSet(wire.Struct(new(Injector), "*"))

type Injector struct {
	Engine *gin.Engine
	PermissionService service.IPermissionService
}
