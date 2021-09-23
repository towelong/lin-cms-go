package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/pkg/router"
)

type BookAPI struct {
	Auth middleware.Auth
}

func (book *BookAPI) GetBookList(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"book": "list",
	})
}

func (book *BookAPI) DeleteBook(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"book": "list",
	})
}

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
