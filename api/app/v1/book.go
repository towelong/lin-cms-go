package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg/response"
	"github.com/towelong/lin-cms-go/pkg/router"
)

type BookAPI struct {
	Auth        middleware.Auth
	BookService service.IBookService
}

func (book *BookAPI) GetBookList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, book.BookService.GetBookList())
}

func (book *BookAPI) GetBookById(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if id <= 0 {
		ctx.Error(response.ParmeterInvalid(ctx, 10030, "id必须是正整数"))
		return
	}
	b, err := book.BookService.FindBookById(id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, b)
}

func (book *BookAPI) CreateBook(ctx *gin.Context) {
	var bookDto dto.CreateOrUpdateBookDTO
	if err := ctx.ShouldBindJSON(&bookDto); err != nil {
		ctx.Error(err)
		return
	}
	book.BookService.CreateBook(bookDto)
	response.CreatedVO(ctx, 12)
}

func (book BookAPI) ChangeBookById(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if id <= 0 {
		ctx.Error(response.ParmeterInvalid(ctx, 10030, "id必须是正整数"))
		return
	}
	var bookDto dto.CreateOrUpdateBookDTO
	if err := ctx.ShouldBindJSON(&bookDto); err != nil {
		ctx.Error(err)
		return
	}
	book.BookService.ChangeBookById(id, bookDto)
	response.UpdatedVO(ctx, 13)
}

func (book *BookAPI) DeleteBook(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if id <= 0 {
		ctx.Error(response.ParmeterInvalid(ctx, 10030, "id必须是正整数"))
		return
	}
	book.BookService.DeleteBookById(id)
	response.DeletedVO(ctx, 14)
}

func (book *BookAPI) RegisterServer(routerGroup *gin.RouterGroup) {
	bookRouter := router.NewLinRouter("/book", "图书", routerGroup)
	bookRouter.GET("", book.GetBookList)
	bookRouter.LinDELETE("DeleteBook",
		"/:id",
		bookRouter.Permission("删除图书", true),
		book.Auth.GroupRequired,
		book.DeleteBook,
	)
	bookRouter.POST("", book.CreateBook)
	bookRouter.GET("/:id", book.GetBookById)
	bookRouter.PUT("/:id", book.ChangeBookById)
}
