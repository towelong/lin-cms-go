package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UnifyResponse(code int64, ctx *gin.Context) *Response{
	response := NewResponse(code)
	response.SetRequest(ctx.Request.RequestURI)
	return response
}

func ParmeterInvalid(ctx *gin.Context, code int64, message string) *Response {
	response := NewResponse(code)
	response.SetRequest(ctx.Request.RequestURI)
	response.SetMessage(message)
	return response
}

func WrapResponse(ctx *gin.Context, response *Response) *Response{
	response.SetRequest(ctx.Request.RequestURI)
	return response
}

func Success(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, UnifyResponse(0, ctx))
}

func CreatedVO(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, UnifyResponse(1, ctx))
}

func UpdatedVO(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, UnifyResponse(2, ctx))
}

func DeletedVO(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, UnifyResponse(3, ctx))
}

func NotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, UnifyResponse(10020, ctx))
}

func AuthFail(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, UnifyResponse(10041, ctx))
}

func ServerFail(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, UnifyResponse(9999, ctx))
}
