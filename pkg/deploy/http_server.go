package deploy

import (
	"net/http"
	"relaper.com/kubemanage/inital"
	"relaper.com/kubemanage/protocol"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPHandler return a http.Handler
func NewHTTPHandler(endpoints Endpoints) http.Handler {
	ctxModel := protocol.CtxModel_Production
	if inital.GetGlobal().GetOptions().Dev {
		ctxModel = protocol.CtxModel_Debug
	}
	options := []httptransport.ServerOption{
		//httptransport.ServerErrorLogger(logger),
		httptransport.ServerBefore(
			httptransport.PopulateRequestContext,
			protocol.MakeJWTTokenToContext(),
			protocol.MakeServerBefore(protocol.CtxOpts{ctxModel})),
	}

	return MakeRouter(endpoints, options...)
}

// MakeRouter return a http.Handler.
// 仅保持最基本的错误处理
func MakeRouter(endpoints Endpoints, options ...httptransport.ServerOption) http.Handler {
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(protocol.MakeErrorEncoder()),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	opts = append(opts, options...)
	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()
	resource := v1.PathPrefix("/resource").Subrouter()
	{
		resource.Handle("/deploy", httptransport.NewServer(
			endpoints.DeployEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &DeployRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		resource.Handle("/uploadDeploy", httptransport.NewServer(
			endpoints.UploadEndpoint,
			protocol.MakeUploadRequest(&UploadRequest{}),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		resource.Handle("/delete", httptransport.NewServer(
			endpoints.DeleteEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &DeleteRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		resource.Handle("/expansion", httptransport.NewServer(
			endpoints.ExpansionEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &ExpansionRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))
	}

	namespace := v1.PathPrefix("/namespace").Subrouter()
	{
		namespace.Handle("/create", httptransport.NewServer(
			endpoints.CreateNsEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &NamespaceRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))
		namespace.Handle("/delete", httptransport.NewServer(
			endpoints.DeleteNsEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &NamespaceRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))
	}

	return r
}
