package tools

import (
	"context"
	"encoding/json"
	"net/http"
)

type Response struct {
	Errno  int         `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func NewResponse(errNo int, errMsg string) *Response {
	return &Response{
		Errno:  errNo,
		Errmsg: errMsg,
	}
}

func EncodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, resp *Response) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	encoder := json.NewEncoder(w)
	return encoder.Encode(resp)
}
