package cluster

import "github.com/go-kit/kit/endpoint"

// Endpoints is a set of endpoint
type Endpoints struct {
	NodesEndpoint       endpoint.Endpoint
	ClusterEndpoint     endpoint.Endpoint
	NodeEndpoint        endpoint.Endpoint
	NameSpaceEndpoint   endpoint.Endpoint
	PodInfoEndpoint     endpoint.Endpoint
	PodsEndpoint        endpoint.Endpoint
	DeploymentEndpoint  endpoint.Endpoint
	StatefulSetEndpoint endpoint.Endpoint
	ServiceEndpoint     endpoint.Endpoint
}

// NewEndpoints return a *Endpoints
func NewEndpoints(svc Service) Endpoints {
	return Endpoints{
		ClusterEndpoint:     MakeClusterEndpoint(svc),
		NodesEndpoint:       MakeNodesEndpoint(svc),
		NodeEndpoint:        MakeNodeEndpoint(svc),
		NameSpaceEndpoint:   MakeNameSpaceEndpoint(svc),
		PodInfoEndpoint:     MakePodInfoEndpoint(svc),
		PodsEndpoint:        MakePodsEndpoint(svc),
		DeploymentEndpoint:  MakeDeploymentEndpoint(svc),
		StatefulSetEndpoint: MakeStatefulSetEndpoint(svc),
		ServiceEndpoint:     MakeServiceEndpoint(svc),
	}
}
