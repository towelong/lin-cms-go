package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/pkg/log"
)

func Logger(ctx *gin.Context) {
	start := time.Now()
	data, _ := ctx.GetRawData()
	query := ctx.Copy().Request.URL.RawQuery
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 重新赋值
	ctx.Next()
	costs := time.Since(start)
	body := string(data)
	if body == "" {
		body = "{}"
	}
	logger := log.NewCustomerLogger()
	defer logger.Sync()
	msg := fmt.Sprintf(`%s [%s] -> [%s] from: %s costs: %dms
data: {
	params: %s,
	body: %s
}`,
		time.Now().Format("2006-01-02 15:04:05"),
		ctx.Request.Method,
		ctx.Request.RequestURI,
		ctx.ClientIP(),
		costs.Milliseconds(),
		fomatterQuery(query),
		body,
	)
	logger.Info(msg)
}

func fomatterQuery(query string) string {
	if query != "" {
		queryParts := strings.Split(query, "&")
		var temp string
		for _, part := range queryParts {
			s := strings.Split(part, "=")
			temp = temp + fmt.Sprintf(`"%s":"%s",`, s[0], s[1])
		}
		return "{" + strings.TrimRight(temp, ",") + "}"
	}
	return "{}"
}
