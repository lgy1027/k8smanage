package cluster

import (
	"github.com/lgy1027/kubemanage/inital"
	"github.com/lgy1027/kubemanage/pkg/file"
	"github.com/lgy1027/kubemanage/protocol"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(endpoints Endpoints) http.Handler {
	ctxModel := protocol.CtxModel_Production
	if inital.GetGlobal().GetOptions().Dev {
		ctxModel = protocol.CtxModel_Debug
	}
	options := []httptransport.ServerOption{
		httptransport.ServerBefore(
			httptransport.PopulateRequestContext,
			protocol.MakeJWTTokenToContext(),
			protocol.MakeServerBefore(protocol.CtxOpts{ctxModel})),
	}

	return MakeRouter(endpoints, options...)
}

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

		v1.Handle("/ns", httptransport.NewServer(
			endpoints.NsEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &NsRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/namespace", httptransport.NewServer(
			endpoints.NameSpaceEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &NameSpaceRequest{} })),
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

		//v1.Handle("/podLog", httptransport.NewServer(
		//	endpoints.PodLogEndpoint,
		//	protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &PodLogRequest{} })),
		//	protocol.EncodeHTTPGenericResponse,
		//	opts...,
		//))

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

		v1.Handle("/getYaml", httptransport.NewServer(
			endpoints.GetYamlEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &GetYamlRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/event", httptransport.NewServer(
			endpoints.EventEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &EventRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.Handle("/version", httptransport.NewServer(
			endpoints.VersionEndpoint,
			protocol.LogRequest(protocol.MakeDecodeHTTPRequest(func() interface{} { return &VersionRequest{} })),
			protocol.EncodeHTTPGenericResponse,
			opts...,
		))

		v1.HandleFunc("/downloadYaml", file.HandleDownload)
	}
	return r
}
