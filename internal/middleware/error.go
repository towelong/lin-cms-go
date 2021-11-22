package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/towelong/lin-cms-go/pkg/response"
	validator2 "github.com/towelong/lin-cms-go/pkg/validator"
	"net/http"
	"strings"
)

func ErrorHandler(ctx *gin.Context) {
	ctx.Next()
	length := len(ctx.Errors)
	if length > 0 {
		e := ctx.Errors[length-1]
		switch err := e.Err.(type) {
		case validator.ValidationErrors:
			wrapError(ctx, err)
		case *validator.ValidationErrors:
			wrapError(ctx, *err)
		case *response.Response:
			err.SetRequest(ctx.Request.RequestURI)
			httpCode := err.HttpCode
			if httpCode == 0 {
				httpCode = 400
			}
			ctx.JSON(httpCode, err)
		default:
			response.ServerFail(ctx)
		}
	}
}

func wrapError(ctx *gin.Context, err validator.ValidationErrors) {
	mapErrors := make(map[string]string)
	var (
		errString string
		r         *response.Response
	)
	for _, v := range err {
		errString = v.Translate(validator2.Trans)
		filedName := strings.ToLower(v.StructField())
		mapErrors[filedName] = errString
	}
	r = response.UnifyResponse(10030, ctx)
	if len(err) > 1 {
		r.SetMessage(mapErrors)
	} else {
		r.SetMessage(errString)
	}
	httpCode := r.HttpCode
	if httpCode == 0 {
		httpCode = http.StatusBadRequest
	}
	ctx.JSON(http.StatusBadRequest, r)
}
