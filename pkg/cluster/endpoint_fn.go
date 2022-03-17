// @Author : liguoyu
// @Date: 2019/10/29 15:42
package cluster

import (
	"context"

	"errors"

	"github.com/go-kit/kit/endpoint"
	tipErrors "github.com/lgy1027/kubemanage/utils/errors"
)

func (s *Endpoints) Cluster(ctx context.Context, request *ClusterRequest) (*ClusterResponse, error) {
	if resp, err := s.ClusterEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*ClusterResponse), nil
	}
}

func MakeClusterEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*ClusterRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeClusterEndpoint"), "内部错误")
		} else {
			return svc.Cluster(ctx, req)
		}
	}
}

func (s *Endpoints) Nodes(ctx context.Context, request *NodesRequest) (*NodesResponse, error) {
	if resp, err := s.NodesEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*NodesResponse), nil
	}
}

func MakeNodesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*NodesRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeNodesEndpoint"), "内部错误")
		} else {
			return svc.Nodes(ctx, req)
		}
	}
}

func (s *Endpoints) Node(ctx context.Context, request *NodeRequest) (*NodeResponse, error) {
	if resp, err := s.NodeEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*NodeResponse), nil
	}
}

func MakeNodeEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*NodeRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeNodeEndpoint"), "内部错误")
		} else {
			return svc.Node(ctx, req)
		}
	}
}

func (s *Endpoints) Ns(ctx context.Context, request *NsRequest) (*NsResponse, error) {
	if resp, err := s.NsEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*NsResponse), nil
	}
}

func MakeNsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*NsRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeNsEndpoint"), "内部错误")
		} else {
			return svc.Ns(ctx, req)
		}
	}
}

func (s *Endpoints) NameSpace(ctx context.Context, request *NameSpaceRequest) (*NameSpaceResponse, error) {
	if resp, err := s.NameSpaceEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*NameSpaceResponse), nil
	}
}

func MakeNameSpaceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*NameSpaceRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeNameSpaceEndpoint"), "内部错误")
		} else {
			return svc.NameSpace(ctx, req)
		}
	}
}

func (s *Endpoints) PodInfo(ctx context.Context, request *PodInfoRequest) (*PodInfoResponse, error) {
	if resp, err := s.PodInfoEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*PodInfoResponse), nil
	}
}

func MakePodInfoEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*PodInfoRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakePodInfoEndpoint"), "内部错误")
		} else {
			return svc.PodInfo(ctx, req)
		}
	}
}

//func (s *Endpoints) PodLog(ctx context.Context, request *PodLogRequest) (*PodLogResponse, error) {
//	if resp, err := s.PodInfoEndpoint(ctx, request); err != nil {
//		return nil, err
//	} else {
//		return resp.(*PodLogResponse), nil
//	}
//}
//
//func MakePodLogEndpoint(svc Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		if v, ok := request.(Validator); ok {
//			if err = v.Validate(); err != nil {
//				return nil, err
//			}
//		}
//
//		if req, ok := request.(*PodLogRequest); !ok {
//			return nil, tipErrors.WithTipMessage(errors.New("MakePodLogEndpoint"), "内部错误")
//		} else {
//			return svc.PodLog(ctx, req)
//		}
//	}
//}

func (s *Endpoints) Pods(ctx context.Context, request *PodsRequest) (*PodsResponse, error) {
	if resp, err := s.PodsEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*PodsResponse), nil
	}
}

func MakePodsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*PodsRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakePodsEndpoint"), "内部错误")
		} else {
			return svc.Pods(ctx, req)
		}
	}
}

func (s *Endpoints) Deployment(ctx context.Context, request *ResourceRequest) (*DeploymentsResponse, error) {
	if resp, err := s.DeploymentEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*DeploymentsResponse), nil
	}
}

func MakeDeploymentEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*ResourceRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeDeploymentEndpoint"), "内部错误")
		} else {
			return svc.Deployment(ctx, req)
		}
	}
}

func (s *Endpoints) StatefulSet(ctx context.Context, request *ResourceRequest) (*StatefulSetsResponse, error) {
	if resp, err := s.StatefulSetEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*StatefulSetsResponse), nil
	}
}

func MakeStatefulSetEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*ResourceRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeStatefulSetEndpoint"), "内部错误")
		} else {
			return svc.StatefulSet(ctx, req)
		}
	}
}

func (s *Endpoints) Services(ctx context.Context, request *ResourceRequest) (*ServiceResponse, error) {
	if resp, err := s.ServiceEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*ServiceResponse), nil
	}
}

func MakeServiceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*ResourceRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeServiceEndpoint"), "内部错误")
		} else {
			return svc.Services(ctx, req)
		}
	}
}

func (s *Endpoints) GetYaml(ctx context.Context, request *GetYamlRequest) (*GetYamlResponse, error) {
	if resp, err := s.GetYamlEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*GetYamlResponse), nil
	}
}

func MakeGetYamlEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*GetYamlRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeGetYamlEndpoint"), "内部错误")
		} else {
			return svc.GetYaml(ctx, req)
		}
	}
}

func (s *Endpoints) Event(ctx context.Context, request *EventRequest) (*EventResponse, error) {
	if resp, err := s.EventEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*EventResponse), nil
	}
}

func MakeEventEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*EventRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeEventEndpoint"), "内部错误")
		} else {
			return svc.Event(ctx, req)
		}
	}
}

func (s *Endpoints) VersionList(ctx context.Context, request *VersionRequest) (*VersionResponse, error) {
	if resp, err := s.VersionEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*VersionResponse), nil
	}
}

func MakeVersionEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*VersionRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeVersionEndpoint"), "内部错误")
		} else {
			return svc.VersionList(ctx, req)
		}
	}
}
