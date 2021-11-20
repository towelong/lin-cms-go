package response

import (
	"fmt"
	"github.com/towelong/lin-cms-go/pkg"
)

type Response struct {
	Code     int64       `json:"code"`
	Message  interface{} `json:"message"`
	Request  string      `json:"request"`
	HttpCode int         `json:"-"`
}

func NewResponse(code int64) *Response {
	return &Response{
		Code:    code,
		Message: pkg.Code2Message(code),
		Request: "",
	}
}

func New(code int64, httpCode int) *Response {
	return &Response{
		Code:     code,
		Message:  pkg.Code2Message(code),
		Request:  "",
		HttpCode: httpCode,
	}
}

func (r *Response) SetCode(code int64) {
	r.Code = code
	// code 变了， message也要跟着变
	r.Message = pkg.Code2Message(code)
}

func (r *Response) SetMessage(message interface{}) {
	r.Message = message
}

func (r *Response) Error() string {
	switch m := r.Message.(type) {
	case string:
		return m
	case map[string]string:
		var msg string
		for k, v := range m {
			msg += fmt.Sprintf("%s : %s", k, v)
		}
		return msg
	default:
		return ""
	}
}

func (r *Response) SetRequest(request string) {
	r.Request = request
}
