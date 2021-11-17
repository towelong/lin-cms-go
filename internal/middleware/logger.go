package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/pkg/log"
	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg/router"
)

var LogSet = wire.NewSet(wire.Struct(new(Logs), "*"))

type Logs struct {
	LogService service.ILogService
}

// ====== 获取response
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ======

// 日志系统
func Log(ctx *gin.Context) {
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

// 行为日志
func (l Logs) Logger(template string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		ctx.Next()
		raw, exit := ctx.Get("currentUser")
		if exit {
			user := raw.(model.User)
			parseString, oldString := parseTemplate(ctx, blw, user, template)
			msg := strings.ReplaceAll(template, oldString, parseString)
			meta, ok := ctx.Get("meta")
			var permission string
			if !ok {
				permission = ""
			} else {
				routeMeta := meta.(router.Meta)
				permission = routeMeta.Permission
			}

			log := model.Log{
				Message:    msg,
				UserID:     user.ID,
				Username:   user.Username,
				StatusCode: blw.Status(),
				Method:     ctx.Copy().Request.Method,
				Path:       ctx.Copy().Request.RequestURI,
				Permission: permission,
			}
			err := l.LogService.CreateLog(log)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
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

func parseTemplate(ctx *gin.Context, resp *bodyLogWriter, user model.User, template string) (parseStr string, oldStr string) {
	r := regexp.MustCompile("{[^}]+}")
	tp := r.FindString(template)
	return getValueByPropName(ctx, user, resp, tp), tp
}

func getValueByPropName(ctx *gin.Context, obj interface{}, resp *bodyLogWriter, prop string) string {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Struct {
		return ""
	}
	s := strings.Split(prop, ".")
	if len(s) != 2 {
		return ""
	}
	object := strings.TrimLeft(s[0], "{")
	propName := strings.TrimRight(s[1], "}")
	v := reflect.ValueOf(obj)
	switch object {
	case "user":
		if _, b := t.FieldByName(FirstUpper(propName)); b {
			return v.FieldByName(FirstUpper(propName)).String()
		}
		return ""
	case "request":
		if propName == "url" {
			return ctx.Copy().Request.RequestURI
		}
		if propName == "ip" {
			return ctx.Copy().ClientIP()
		}
		return ""
	case "response":
		if propName == "status" {
			return fmt.Sprint(resp.Status())
		}
		return ""
	}
	return ""
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
