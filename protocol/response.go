// @Author : liguoyu
// @Date: 2019/10/29 15:42
package protocol

import (
	errors2 "github.com/lgy1027/kubemanage/utils/errors"
	"github.com/pkg/errors"

	"encoding/json"
	"fmt"
)

type extr struct {
	InnerError string `json:"inner_error"`
	ErrorStack string `json:"error_stack"`
}
type OptionsResponse struct {
	Errmsg string `json:"errmsg"`
}

type Response struct {
	// @description 错误编号 -1 失败 0 成功
	Errno int `json:"errno"`
	// @description 错误信息
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
	Extr   *extr       `json:"extr,omitempty"`
}

func NewResponse() *Response {
	return &Response{
		Extr: &extr{},
	}
}

func (r *Response) SetExtr(err error) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	type causer interface {
		Cause() error
	}

	if e1, ok := err.(stackTracer); ok { // pkg error
		st := e1.StackTrace()
		r.Extr.ErrorStack = fmt.Sprintf("%+v", st[:])
	} else if e2, ok := err.(causer); ok { // withMessage or withErrno
		if e3, ok := e2.Cause().(stackTracer); ok {
			st := e3.StackTrace()
			r.Extr.ErrorStack = fmt.Sprintf("%+v", st[:])
		} else if e4, ok := e3.(causer); ok { // withMessage and withErrno
			if e5, ok := e4.Cause().(stackTracer); ok {
				st := e5.StackTrace()
				r.Extr.ErrorStack = fmt.Sprintf("%+v", st[:])
			}
		}
	}

	r.Extr.InnerError = err.Error()
}

func (r *Response) SetErrno(err error) {
	if e, ok := err.(errors2.Errnoer); ok {
		r.Errno = e.Errno()
	} else {
		r.Errno = -1
	}
}

func (r *Response) SetErrmsg(err error) {
	if e, ok := err.(errors2.Tipper); ok {
		r.Errmsg = e.Tip()
	} else {
		r.Errmsg = err.Error()
	}
}

func (r *Response) FillError(err error) {
	r.SetErrno(err)
	r.SetErrmsg(err)
	r.SetExtr(err)
}

func (r *Response) SetData(data interface{}) {
	r.Data = data
}

func (r *Response) String() string {
	bs, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
