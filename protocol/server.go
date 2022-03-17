// @Author : liguoyu
// @Date: 2019/10/29 15:42
package protocol

import (
	"bytes"
	"context"
	"encoding/json"
	log "github.com/cihub/seelog"
	httptransport "github.com/go-kit/kit/transport/http"
	uerrors "github.com/lgy1027/kubemanage/utils/errors"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type contextKey string

const (
	logLength                     = 8192
	RequestLogger      contextKey = "LOGGER"
	CtxModelContextKey contextKey = "CtxModel"

	CtxModel_Debug      = 0
	CtxModel_Production = 1
)

const ErrMsg = "options http method"

func MakeErrorEncoder() httptransport.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// TODO 打印记录合适的错误信息
		resp := NewResponse()
		resp.FillError(err)

		// 输出客户端IP信息
		//clientIp, ok := ctx.Value(httptransport.ContextKeyRequestXForwardedFor).(string)
		//if !ok {
		//	clientIp = "missing"
		//}
		//loggerFunc.Log("clientIp", clientIp, "resp", resp)

		// 非debug模式输出不包含extr信息
		if ctxModel, ok := ctx.Value(CtxModelContextKey).(int); ok {
			if ctxModel != CtxModel_Debug {
				//resp.Extr = nil
			}
		}
		// TODO 上线可以注释 resp.Extr = nil
		resp.Extr = nil
		encoder := json.NewEncoder(w)
		if err.Error() == ErrMsg {
			op := OptionsResponse{
				Errmsg: "http method options ok",
			}
			encoder.Encode(op)
		} else {
			encoder.Encode(resp)
		}
	}
}

type CtxOpts struct {
	//Logger log.Logger
	Model int
}

func MakeServerBefore(opts CtxOpts) httptransport.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		ctx = context.WithValue(ctx, CtxModelContextKey, opts.Model)
		//ctx = context.WithValue(ctx, RequestLogger, opts.Logger)
		return ctx
	}
}

type ReadCloser struct {
	closed bool
	r      *bytes.Reader
}

func (rc *ReadCloser) Read(p []byte) (int, error) {
	if rc.closed == true {
		return 0, io.EOF
	}
	return rc.r.Read(p)
}

func (rc *ReadCloser) Close() error {
	rc.closed = true
	return nil
}

func LogRequest(next httptransport.DecodeRequestFunc) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		var body []byte
		var err error
		header, _ := json.Marshal(r.Header)
		if r.GetBody == nil {
			if body, err = ioutil.ReadAll(io.LimitReader(r.Body, logLength)); err == nil {
				r.Body.Close()
				r.Body = &ReadCloser{r: bytes.NewReader(body)}
			}
		} else {
			if bodyReader, err := r.GetBody(); err == nil {
				defer bodyReader.Close()
				body, _ = ioutil.ReadAll(io.LimitReader(bodyReader, logLength))
			}
		}
		log.Info("requestUri:", r.RequestURI, " header:", string(header), " body:", string(body))
		return next(ctx, r)
	}
}

func MakeDecodeHTTPRequest(maker func() interface{}) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		req := maker()

		if req == nil {
			return nil, nil
		}

		var err error
		switch r.Method {
		case http.MethodGet:
			err = r.ParseForm()
			if err != nil {
				err = errors.Wrap(err, "ParseForm failure")
				break
			}

			// TODO query 提为常量
			err = BindData(req, r.Form, "query")
			if err != nil {
				err = errors.Wrap(err, "query string unmarshal failure")
			}
		case http.MethodPost:
			body, _ := ioutil.ReadAll(r.Body)
			if len(body) == 0 {
				body = []byte("{}")
			}
			r.Body = &ReadCloser{r: bytes.NewReader(body)}
			err = json.NewDecoder(r.Body).Decode(req)
			if err != nil {
				err = errors.Wrap(err, "json unmarshal failure")
			}
		case http.MethodOptions:
			err = errors.New(ErrMsg)
			err = uerrors.WithTipMessage(err, ErrMsg)
			return req, err
		default:
			err = errors.New("unsupported http method")
		}

		err = uerrors.WithTipMessage(err, "系统内部错误")
		return req, err
	}
}

func EncodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	resp := NewResponse()
	resp.Data = response
	encoder := json.NewEncoder(w)
	return encoder.Encode(resp)
}

func TimeoutHandle(next http.Handler, d time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancer := context.WithTimeout(ctx, d)
		defer cancer()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

type Requestor interface {
	SetRequest(r *http.Request)
}

func MakeUploadRequest(req interface{}) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		if iReq, ok := req.(Requestor); !ok {
			return nil, errors.New("内部错误")
		} else {
			iReq.SetRequest(r)
			return req, nil
		}
	}
}
