package deploy

import (
	"context"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
	k8s2 "relaper.com/kubemanage/k8s"
	"relaper.com/kubemanage/utils"
)

type Service interface {
	Deploy(ctx context.Context, req *DeployRequest) (*DeploymentResponse, error)
	UploadDeploy(ctx context.Context, request *UploadRequest) (*UploadResponse, error)
	Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error)
	CreateNs(ctx context.Context, req *NamespaceRequest) (*NamespaceResponse, error)
	DeleteNs(ctx context.Context, req *NamespaceRequest) (*NamespaceResponse, error)
}

// NewService return a Service interface
func NewService() Service {
	return &deployService{
		Deployment:  k8s2.NewDeploy(),
		Sv:          k8s2.NewSv(),
		StateFulSet: k8s2.NewStateFulSet(),
		Namespace:   k8s2.NewNs(),
	}
}

type deployService struct {
	Deployment  *k8s2.Deployment
	Sv          *k8s2.Sv
	StateFulSet *k8s2.Sf
	Namespace   *k8s2.Ns
}

/*
{
   "kind":"Deployment",
   "namespace":"lgy",
   "name":"nginx",
   "objectMetaLabels":{"test":"test-nginx"},
   "annotations":{"desc":"nginx test for namespace lgy "},
   "replicas": 3,
   "matchLabels" : {"app": "demo"},
   "maxSurge":1,
   "maxUnavailable":1,
   "templateLabels":{"app": "demo"},
   "nodeSelector":{"node":"three"},
   "podName":"nginx-post",
   "image":"nginx:1.12",
   "podPort":[{
       "name":"http",
       "protocol":"TCP",
       "containerPort":80
   }],
   "resources":{"limits": {"cpu":"200m", "memory": "250Mi"}, "requests": {"cpu":"100m", "memory": "100Mi"}},
   "imagePullPolicy":"IfNotPresent",
   "createService":      true,
   "serviceName":       "nginx-service",
   "servicePorts":    [{
       "port":8080,
       "targetPort":80,
       "protocol":"TCP"
   }],
   "serviceType":"NodePort"
}

*/

// @Summary  部署资源
// @Produce  json
// @Accept  json
// @Param   kind body DeployRequest true "参数"
// @Success 200 {string} json "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/resource/deploy [post]
func (dm *deployService) Deploy(ctx context.Context, req *DeployRequest) (*DeploymentResponse, error) {
	var err error
	switch req.Kind {
	case utils.DEPLOY_DEPLOYMENT:
		_, exist, _ := dm.Deployment.Exist(req.Namespace, req.Name)
		if exist {
			return nil, fmt.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", "Deployment", req.Name, req.Namespace)
		}
		deployment := ExpandDeployment(req)
		_, err = dm.Deployment.Create(req.Namespace, deployment)
	case utils.DEPLOY_STATEFULSET:
		_, exist, _ := dm.StateFulSet.Exist(req.Namespace, req.Name)
		if exist {
			return nil, fmt.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", "StatefulSet", req.Name, req.Namespace)
		}
		statefulSets := ExpandStatefulSets(req)
		_, err = dm.StateFulSet.Create(req.Namespace, statefulSets)
	case utils.DEPLOY_Service:
		_, exist, _ := dm.Sv.Exist(req.Namespace, req.Name)
		if exist {
			return nil, fmt.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", "Service", req.Name, req.Namespace)
		}
		service := ExpandService(req)
		_, err = dm.Sv.Create(req.Namespace, service)
	default:
		return nil, errors.New("资源类型不存在，可选 Deployment | StatefulSet | Service")
	}
	if err != nil {
		err = k8s2.PrintErr(err)
		return nil, err
	}
	if req.CreateService {
		service := ExpandService(req)
		_, err := dm.Sv.Create(req.Namespace, service)
		if err != nil {
			err = k8s2.PrintErr(err)
			return nil, err
		}
	}
	return nil, nil
}

// @Summary  文件部署资源
// @Produce  json
// @Param resource formData file true "yaml文件"
// @Success 200 {object} protocol.Response{data=UploadResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/resource/uploadDeploy [post]
func (dm *deployService) UploadDeploy(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	file, handler, err := req.FormFile("resource")
	if err != nil {
		log.Debug("[UploadDeploy]: err:", err.Error())
		return nil, errors.New("文件上传失败")
	}
	defer file.Close()

	objs, err := ExpandMultiYamlFileToObject(file)
	if err != nil {
		return nil, err
	}
	if objs != nil && len(objs) > 0 {
		for _, obj := range objs {
			apiVersion, kind, name, namespace := obj.GetAPIVersion(), obj.GetKind(), obj.GetName(), obj.GetNamespace()
			if namespace == "" {
				namespace = utils.DEFAULTNS
			}
			var errs error
			switch kind {
			case utils.DEPLOY_DEPLOYMENT:
				_, exist, _ := dm.Deployment.Exist(namespace, name)
				if exist {
					return nil, fmt.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", kind, name, namespace)
				}
				_, errs = dm.Deployment.DynamicCreateForCustom(namespace, apiVersion, obj)
			case utils.DEPLOY_STATEFULSET:
				_, exist, _ := dm.StateFulSet.Exist(namespace, name)
				if exist {
					return nil, fmt.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", kind, name, namespace)
				}
				_, errs = dm.StateFulSet.DynamicCreateForCustom(namespace, apiVersion, obj)
			case utils.DEPLOY_Service:
				_, exist, _ := dm.Sv.Exist(namespace, name)
				if exist {
					return nil, fmt.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", kind, name, namespace)
				}
				_, errs = dm.Sv.DynamicCreateForCustom(namespace, apiVersion, obj)
			default:
				return nil, errors.New("资源类型不存在，可选 Deployment | StatefulSet | Service")
			}
			if errs != nil {
				return nil, fmt.Errorf("部署资源失败:%v, Kind:%s, Name: %s, namespace:%s", err, kind, name, namespace)
			}
		}
	}
	return &UploadResponse{
		File: handler.Filename,
	}, nil
}

// @Summary  删除资源
// @Produce  json
// @Accept   json
// @Param   params body DeleteRequest true "资源对象 可选 Deployment | StatefulSet | Service "
// @Success 200 {string} json "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/resource/delete [post]
func (dm *deployService) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	var err error
	switch req.Kind {
	case utils.DEPLOY_DEPLOYMENT:
		err = dm.Deployment.Delete(req.Namespace, req.Name)
	case utils.DEPLOY_STATEFULSET:
		err = dm.StateFulSet.Delete(req.Namespace, req.Name)
	case utils.DEPLOY_Service:
		err = dm.Sv.Delete(req.Namespace, req.Name)
	default:
		return nil, errors.New("资源类型不存在，可选 Deployment | StatefulSet | Service")
	}
	err = k8s2.PrintErr(err)
	return &DeleteResponse{}, err
}

// @Summary  创建命名空间
// @Produce  json
// @Accept   json
// @Param   params body NamespaceRequest true "命名空间名"
// @Success 200 {string} json "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/namespace/create [post]
func (dm *deployService) CreateNs(ctx context.Context, req *NamespaceRequest) (*NamespaceResponse, error) {
	_, err := dm.Namespace.Create(req.Namespace)
	err = k8s2.PrintErr(err)
	return nil, err
}

// @Summary  删除命名空间
// @Produce  json
// @Accept   json
// @Param   params body NamespaceRequest true "命名空间名"
// @Success 200 {string} json "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/namespace/delete [post]
func (dm *deployService) DeleteNs(ctx context.Context, req *NamespaceRequest) (*NamespaceResponse, error) {
	err := dm.Namespace.Delete(req.Namespace, "")
	err = k8s2.PrintErr(err)
	return nil, err
}
