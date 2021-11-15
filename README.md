<h1 align="center">
  <br>Lin-CMS-Go<br>
</h1>
<h4 align="center">Lin-CMS的go语言实现, 简要介绍下LinCMS: Lin CMS 是一个前后端分离的 CMS 解决方案。其他信息请移步官方文档</h4>

## LinCMS文档资料
### 文档地址

- [http://doc.cms.talelin.com/](http://doc.cms.talelin.com/)

### 线上 Demo

- [http://face.cms.talelin.com/](http://face.cms.talelin.com/)

### 案例

- [http://sleeve.talelin.com/](http://sleeve.talelin.com/)


***
## Lin-CMS-Go

### 安装依赖

Go >= 1.16

```golang
$ go mod tidy
```

### 目录结构


```sh
.
├── Makefile
├── README.md
├── api
│   └── app  ---- 存放API文件
│       ├── cms
│       │   ├── admin.go
│       │   ├── file.go
│       │   ├── user.go
│       │   └── user_test.go
│       ├── v1
│       │   └── book.go
│       └── wire_set.go   --- 提供对接口的依赖注入
├── cmd
│   └── app
│       ├── main.go -- 服务启动的入口文件
│       ├── pprof.go  -- go的pprof服务
│       ├── server.go -- 服务启动所需的依赖
│       ├── wire.go -- 注入器
│       ├── wire_gen.go
│       └── wire_set.go -- 为注入器提供依赖
├── config ----配置文件
│   ├── config.dev.yaml
│   ├── config.prod.yaml
│   └── config.yaml
├── go.mod
├── go.sum
├── internal
│   ├── config -- 解析配置文件
│   │   ├── config.go
│   │   └── config_test.go
│   ├── domain  -- 业务模型等等
│   │   ├── dto
│   │   │   ├── group.go
│   │   │   ├── permission.go
│   │   │   └── user.go
│   │   ├── model
│   │   │   ├── base_model.go
│   │   │   ├── book.go
│   │   │   ├── file.go
│   │   │   ├── group.go
│   │   │   ├── group_permission.go
│   │   │   ├── log.go
│   │   │   ├── permission.go
│   │   │   ├── user.go
│   │   │   ├── user_group.go
│   │   │   └── user_identity.go
│   │   └── vo
│   │       ├── group.go
│   │       ├── page.go
│   │       ├── permission.go
│   │       └── user.go
│   ├── injector.go
│   ├── middleware  ---中间件
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── error.go
│   ├── pkg
│   │   ├── db
│   │   │   └── gorm.go
│   │   └── jwt
│   │       ├── jwt_maker.go
│   │       └── wire_set.go
│   ├── router --- 路由注册
│   │   └── router.go
│   ├── service --- 接口服务
│   │   ├── constants.go
│   │   ├── file.go
│   │   ├── group.go
│   │   ├── mock
│   │   │   ├── group.go
│   │   │   └── user.go
│   │   ├── permission.go
│   │   ├── user.go
│   │   └── wire_set.go
│   └── web.go
└── pkg --- 一些常用的包
    ├── code.go --- 业务状态码
    ├── crypt.go --- 加密解密模块
    ├── crypt_test.go
    ├── response  --- 通用的相应封装
    │   ├── common.go
    │   └── response.go
    ├── router
    │   └── lin_router.go
    ├── token
    │   ├── constants.go
    │   ├── jwt.go
    │   ├── jwt_test.go
    │   └── payload.go
    └── validator -- 参数校验
        └── validator.go
```

## 如何使用？ 

### 使用

对于本程序来说有一点的学习成本，由于在这个项目中用到了 Google的 `wire`来实现依赖的注入，因此，大家需要对`wire`有一定的了解. 不懂也没有关系，通过本文档进行简单的了解之后，也可以帮助大家进行使用。

### 开发API的基本步骤
以 `v1/book/login`这个接口为例。
首先API放置在`api/app` 这个目录下，`cms`目录对应的是cms相关的API, `v1`目录放置自身的业务API。

1. 声明一个结构体
```go
type BookAPI struct {
	
}
```

2. 为这个结构体声明一个方法

```go
func (book *BookAPI) GetBookList(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"book": "list",
	})
}
```

3. 最终声明一个注册路由的方法
```go
func (book *BookAPI) RegisterServer(routerGroup *gin.RouterGroup) {
	bookRouter := router.NewLinRouter("/book", "图书", routerGroup)
	bookRouter.GET("/list", book.GetBookList)
	bookRouter.LinDELETE("DeleteBook",
		"/",
		bookRouter.Permission("删除图书", true),
		book.Auth.GroupRequired,
		book.DeleteBook,
	)
}
```
4. 为图书的API注册到主路由上
私有的一些模块放置到了`internal`目录里面，而router也在其中。

internal/router/router.go

```go
type Router struct {
	BookAPI  *v1.BookAPI
  .....
}
```
在RegisterAPI方法中, 声明一个v1的路由组

```go
  v1Group := app.Group("/v1")
	{
		r.BookAPI.RegisterServer(v1Group)
	}
```

到此为止，一个简单的API就写完了。


## 进阶使用

在刚才的基础教程中，并没有实际的业务逻辑，那么在这里通过实际的业务逻辑来进行实战操作。
根据刚才的目录结构可以知道，开发一个业务需要在service中进行。

本次以 cms/user/login API 为例

service/user.go
```go
type IUserService interface {
	VerifyUser(username, password string) (model.User, error)
}

type UserService struct {
	DB           *gorm.DB
}

func (u UserService) VerifyUser(username, password string) (model.User, error) {
	var (
		userIdentity model.UserIdentity
		user         model.User
	)
	db := u.DB.Where("identity_type = ? AND identifier = ?", UserPassword.String(), username).First(&userIdentity)
	if db.Error != nil {
		return user, response.NewResponse(10031)
	}
	verifyPsw := pkg.VerifyPsw(password, userIdentity.Credential)
	if verifyPsw {
		err := u.DB.Where("username = ?", username).First(&user).Error
		if err != nil {
			return user, response.NewResponse(10031)
		}
		return user, nil
	}
	return user, response.NewResponse(10032)
}

```
可以看到定义了一个接口`IUserService`，以及一个结构体`UserService`，结构体中有DB，那么这个`DB（gorm）`就是拿来操作数据库的。

最后定义了一个结构体方法`VerifyUser`，里面就书写了一系列的业务逻辑。到这里并没有结束！

那么大家可能会好奇，这个DB并没有通过什么方法进行传递赋值，为什么可以直接拿来用呢？
这个就是wire这个库在起作用，他是一个依赖注入的库，通过自动生成来处理我们模块与模块之间的依赖。

写到这里大家就差不多能理解了。

但是`wire`并没有`SpringBoot`那么隐式的帮我们注入依赖，而是通过了一个`wire_set`来帮我们管理依赖。在service目录下有个`wire_set.go`文件，里面放的就是需要被wire管理的依赖。
```go
var Set = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),
)
```
这样我们就完成了整个`service`的开发。

回到控制器 `api/app/cms/user.go`
```go
type UserAPI struct {
	JWT         token.IToken
	UserService service.IUserService
}
```
在这里面直接注入`IUserService`即可。

最后路由是这样的

```go
func (u UserAPI) RegisterServer(routerGroup *gin.RouterGroup) {
	userRouter := router.NewLinRouter("/user", "用户", routerGroup)
	userRouter.POST("/login", u.Login)
}
```

除此之外这个`UserAPI`也需要被wire托管。
同样的在`api/app/wire_set.go`中，进行注册即可。
```go
var APISet = wire.NewSet(
	wire.Struct(new(cms.UserAPI), "*"),
)
```

> 最重要的一点：wire_set要生效，必须 使用 `go generate ./... `来进行生成wire帮助我们处理依赖的文件 `cmd/app/wire_gen.go`

但其实看起来跟我们手动去管理依赖的方式是一样的。它看起来像这样：
```go

// Injectors from wire.go:

func NewInjector() (*internal.Injector, error) {
	gormDB, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	userService := &service.UserService{
		DB:           gormDB,
		GroupService: groupService,
	}
	routerRouter := &router.Router{
		AdminAPI: adminAPI,
		UserAPI:  userAPI,
		BookAPI:  bookAPI,
	}
  ....  // 此处进行省略
	return injector, nil
}
```

## LinRouter

路由的使用：使用过koa的小伙伴应该知道，go的gin框架也是属于类koa的一个库，在gin自带的路由基础上进行了进一步的封装。

用起来像这样：
```go
 Auth的引入都是通过注入来进行引入，直接放在结构体中就可以。
type UserAPI struct {
	UserService service.IUserService
	Auth        middleware.Auth
}

// 权限路由 参考 api/app/cms/admin.go
adminRouter := router.NewLinRouter("/admin", "管理员", routerGroup)
	adminRouter.LinGET(
		"GetAllPermissions",
		"/permission",
		adminRouter.Permission("查询所有可分配的权限", false),
		admin.Auth.AdminRequired,
		admin.GetAllPermissions,
	)
// 不需要挂载权限的路由
userRouter := router.NewLinRouter("/user", "用户", routerGroup)
userRouter.POST("/register", u.Auth.AdminRequired, u.Register)
// 普通路由
bookRouter := router.NewLinRouter("/book", "图书", routerGroup)
	bookRouter.GET("/list", book.GetBookList)
```


## License

MIT