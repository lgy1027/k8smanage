// @Author : liguoyu
// @Date: 2019/10/29 15:42
package protocol

import (
	"bytes"
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func EncodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	switch r.Method {
	case http.MethodPost:
		r.Body = ioutil.NopCloser(&buf)
	case http.MethodGet:
		//todo 这里需要将结构体数据解析为 url中的参数格式
	}
	return nil
}

func MakeDecodeHTTPResponse(maker func() interface{}) httptransport.DecodeResponseFunc {
	return func(_ context.Context, r *http.Response) (interface{}, error) {
		data := maker()

		if r.StatusCode != http.StatusOK {
			return nil, errors.New(r.Status)
		}

		resp := NewResponse()
		resp.Data = data
		if err := json.NewDecoder(r.Body).Decode(resp); err != nil {
			return nil, err
		}

		if resp.Errno != 0 {
			return nil, errors.New("resp.Extr.InnerError:" + resp.Errmsg)
		}

		return resp.Data, nil
	}
}
