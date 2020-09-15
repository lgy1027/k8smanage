// @Author : liguoyu
// @Date: 2019/10/29 15:42
package deploy

import (
	"context"

	"errors"

	"github.com/go-kit/kit/endpoint"
	tipErrors "relaper.com/kubemanage/utils/errors"
)

func (s *Endpoints) Deploy(ctx context.Context, request *DeployRequest) (*DeploymentResponse, error) {
	if resp, err := s.DeployEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*DeploymentResponse), nil
	}
}

func MakeDeployEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*DeployRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeDeployEndpoint"), "内部错误")
		} else {
			return svc.Deploy(ctx, req)
		}
	}
}

func (s *Endpoints) UploadDeploy(ctx context.Context, request *UploadRequest) (*UploadResponse, error) {
	if resp, err := s.UploadEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*UploadResponse), nil
	}
}

func MakeUploadDeployEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*UploadRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeUploadDeployEndpoint"), "内部错误")
		} else {
			return svc.UploadDeploy(ctx, req)
		}
	}
}

func (s *Endpoints) Delete(ctx context.Context, request *DeleteRequest) (*DeleteResponse, error) {
	if resp, err := s.DeployEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*DeleteResponse), nil
	}
}

func MakeDeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*DeleteRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeDeleteEndpoint"), "内部错误")
		} else {
			return svc.Delete(ctx, req)
		}
	}
}

func (s *Endpoints) CreateNs(ctx context.Context, request *NamespaceRequest) (*NamespaceResponse, error) {
	if resp, err := s.DeployEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*NamespaceResponse), nil
	}
}

func MakeCreateNsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*NamespaceRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeCreateNsEndpoint"), "内部错误")
		} else {
			return svc.CreateNs(ctx, req)
		}
	}
}

func (s *Endpoints) DeleteNs(ctx context.Context, request *NamespaceRequest) (*NamespaceResponse, error) {
	if resp, err := s.DeployEndpoint(ctx, request); err != nil {
		return nil, err
	} else {
		return resp.(*NamespaceResponse), nil
	}
}

func MakeDeleteNsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if v, ok := request.(Validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}

		if req, ok := request.(*NamespaceRequest); !ok {
			return nil, tipErrors.WithTipMessage(errors.New("MakeDeleteNsEndpoint"), "内部错误")
		} else {
			return svc.DeleteNs(ctx, req)
		}
	}
}
