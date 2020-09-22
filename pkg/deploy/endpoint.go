package deploy

import "github.com/go-kit/kit/endpoint"

// Endpoints is a set of endpoint
type Endpoints struct {
	DeployEndpoint    endpoint.Endpoint
	UploadEndpoint    endpoint.Endpoint
	DeleteEndpoint    endpoint.Endpoint
	ExpansionEndpoint endpoint.Endpoint
	CreateNsEndpoint  endpoint.Endpoint
	DeleteNsEndpoint  endpoint.Endpoint
}

// NewEndpoints return a *Endpoints
func NewEndpoints(svc Service) Endpoints {
	return Endpoints{
		DeployEndpoint:    MakeDeployEndpoint(svc),
		UploadEndpoint:    MakeUploadDeployEndpoint(svc),
		DeleteEndpoint:    MakeDeleteEndpoint(svc),
		ExpansionEndpoint: MakeExpansionEndpoint(svc),
		CreateNsEndpoint:  MakeCreateNsEndpoint(svc),
		DeleteNsEndpoint:  MakeDeleteNsEndpoint(svc),
	}
}
