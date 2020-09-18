package cluster

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
	{
		v1 := r.PathPrefix("/v1").Subrouter()

		v1.Handle("/detail", httptransport.NewServer(
			endpoints.ClusterEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &ClusterRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/nodes", httptransport.NewServer(
			endpoints.NodesEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &NodesRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/node", httptransport.NewServer(
			endpoints.NodeEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &NodeRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/namespace", httptransport.NewServer(
			endpoints.NameSpaceEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &NameSpacesRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/pod", httptransport.NewServer(
			endpoints.PodInfoEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &PodInfoRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/pods", httptransport.NewServer(
			endpoints.PodsEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &PodsRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/podLog", httptransport.NewServer(
			endpoints.PodLogEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &PodInfoRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/deployment", httptransport.NewServer(
			endpoints.DeploymentEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &ResourceRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/statefulSet", httptransport.NewServer(
			endpoints.StatefulSetEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &ResourceRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/service", httptransport.NewServer(
			endpoints.ServiceEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &ResourceRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))
	}
	return r
}
