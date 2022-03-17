package deploy

import (
	"context"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/lgy1027/kubemanage/inital"
	"github.com/lgy1027/kubemanage/inital/client"
	k8s2 "github.com/lgy1027/kubemanage/k8s"
	"github.com/lgy1027/kubemanage/utils"
	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service interface {
	Deploy(ctx context.Context, req *DeployRequest) (*DeploymentResponse, error)
	UploadDeploy(ctx context.Context, req *UploadRequest) (*UploadResponse, error)
	Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error)
	Expansion(ctx context.Context, req *ExpansionRequest) (*ExpansionResponse, error)
	Rollback(ctx context.Context, req *RollbackRequest) (*RollbackResponse, error)
	Stretch(ctx context.Context, req *StretchRequest) (*StretchResponse, error)
	CreateNs(ctx context.Context, req *NamespaceRequest) (*NamespaceResponse, error)
	DeleteNs(ctx context.Context, req *NamespaceRequest) (*NamespaceResponse, error)
}

// NewService return a Service interface
func NewService() Service {
	return &deployService{}
}

type deployService struct{}

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

// @Tags resource
// @Summary  部署资源
// @Produce  json
// @Accept  json
// @Param   kind body DeployRequest true "参数"
// @Success 200 {string} json "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/resource/deploy [post]
func (dm *deployService) Deploy(ctx context.Context, req *DeployRequest) (*DeploymentResponse, error) {
	deploy := &DeploymentResponse{}
	var err error
	switch req.Kind {
	case utils.DEPLOY_DEPLOYMENT:
		_, exist, _ := client.GetBaseClient().Deployment.Exist(req.Namespace, req.Name)
		if exist {
			log.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", "Deployment", req.Name, req.Namespace)
			return nil, k8s2.ErrExist
		}
		deployment := ExpandDeployment(req)
		deploy.Deploy, err = client.GetBaseClient().Deployment.Create(req.Namespace, deployment)
	case utils.DEPLOY_STATEFULSET:
		_, exist, _ := client.GetBaseClient().Sf.Exist(req.Namespace, req.Name)
		if exist {
			log.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", "Deployment", req.Name, req.Namespace)
			return nil, k8s2.ErrExist
		}
		statefulSets := ExpandStatefulSets(req)
		deploy.Deploy, err = client.GetBaseClient().Sf.Create(req.Namespace, statefulSets)
	case utils.DEPLOY_Service:
		_, exist, _ := client.GetBaseClient().Sv.Exist(req.Namespace, req.Name)
		if exist {
			log.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", "Deployment", req.Name, req.Namespace)
			return nil, k8s2.ErrExist
		}
		service := ExpandService(req)
		deploy.Deploy, err = client.GetBaseClient().Sv.Create(req.Namespace, service)
	default:
		log.Error("资源类型不存在，可选 Deployment | StatefulSet | Service")
		return nil, k8s2.ErrInvokerKind
	}
	err = k8s2.PrintErr(err)
	if err != nil {
		log.Error("创建资源失败")
		return nil, k8s2.ErrCreate
	}
	if req.CreateService {
		service := ExpandService(req)
		deploy.Service, err = client.GetBaseClient().Sv.Create(req.Namespace, service)
		err = k8s2.PrintErr(err)
		if err != nil {
			log.Error("创建资源失败")
			return nil, k8s2.ErrCreate
		}
	}
	return deploy, nil
}

// @Tags resource
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
				_, exist, _ := client.GetBaseClient().Deployment.Exist(namespace, name)
				if exist {
					return nil, fmt.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", kind, name, namespace)
				}
				_, errs = client.GetBaseClient().Deployment.DynamicCreateForCustom(namespace, apiVersion, obj)
			case utils.DEPLOY_STATEFULSET:
				_, exist, _ := client.GetBaseClient().Sf.Exist(namespace, name)
				if exist {
					return nil, fmt.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", kind, name, namespace)
				}
				_, errs = client.GetBaseClient().Sf.DynamicCreateForCustom(namespace, apiVersion, obj)
			case utils.DEPLOY_Service:
				_, exist, _ := client.GetBaseClient().Sv.Exist(namespace, name)
				if exist {
					return nil, fmt.Errorf("应用已存在, Kind:%s, Name: %s, namespace:%s", kind, name, namespace)
				}
				_, errs = client.GetBaseClient().Sv.DynamicCreateForCustom(namespace, apiVersion, obj)
			default:
				log.Error("资源类型不存在，可选 Deployment | StatefulSet | Service")
				return nil, k8s2.ErrInvokerKind
			}
			if errs != nil {
				log.Errorf("部署资源失败:%v, Kind:%s, Name: %s, namespace:%s \n", err, kind, name, namespace)
				return nil, k8s2.ErrCreate
			}
		}
	}
	return &UploadResponse{
		File: handler.Filename,
	}, nil
}

// @Tags resource
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
		err = client.GetBaseClient().Deployment.Delete(req.Namespace, req.Name)
	case utils.DEPLOY_STATEFULSET:
		err = client.GetBaseClient().Sf.Delete(req.Namespace, req.Name)
	case utils.DEPLOY_Service:
		err = client.GetBaseClient().Sv.Delete(req.Namespace, req.Name)
	default:
		return nil, errors.New("资源类型不存在，可选 Deployment | StatefulSet | Service")
	}
	err = k8s2.PrintErr(err)
	if err != nil {
		log.Error("删除失败 err:", err.Error())
		return nil, k8s2.ErrDelete
	}
	return nil, nil
}

// @Tags resource
// @Summary  创建命名空间
// @Produce  json
// @Accept   json
// @Param   params body NamespaceRequest true "命名空间名"
// @Success 200 {string} json "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/namespace/create [post]
func (dm *deployService) CreateNs(ctx context.Context, req *NamespaceRequest) (*NamespaceResponse, error) {
	_, err := client.GetBaseClient().Ns.Create(req.Namespace)
	err = k8s2.PrintErr(err)
	if err != nil {
		log.Error("创建失败 err:", err.Error())
		return nil, k8s2.ErrCreate
	}
	return nil, nil
}

// @Tags resource
// @Summary  删除命名空间
// @Produce  json
// @Accept   json
// @Param   params body NamespaceRequest true "命名空间名"
// @Success 200 {string} json "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/namespace/delete [post]
func (dm *deployService) DeleteNs(ctx context.Context, req *NamespaceRequest) (*NamespaceResponse, error) {
	err := client.GetBaseClient().Ns.Delete(req.Namespace, "")
	err = k8s2.PrintErr(err)
	if err != nil {
		log.Error("删除失败 err:", err.Error())
		return nil, k8s2.ErrDelete
	}
	return nil, nil
}

// @Tags resource
// @Summary  扩容服务，CPU和内存
// @Produce  json
// @Accept   json
// @Param   params body ExpansionRequest true "参数"
// @Success 200 {string} json "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/resource/expansion [post]
func (dm *deployService) Expansion(ctx context.Context, req *ExpansionRequest) (*ExpansionResponse, error) {
	var (
		deploy interface{}
		err    error
		flag   bool
	)
	switch req.Kind {
	case utils.DEPLOY_DEPLOYMENT:
		deploy, flag, err = client.GetBaseClient().Deployment.Exist(req.Namespace, req.Name)
	case utils.DEPLOY_STATEFULSET:
		deploy, flag, err = client.GetBaseClient().Sf.Exist(req.Namespace, req.Name)
	default:
		return nil, k8s2.ErrInvokerKind
	}
	if err != nil {
		return nil, k8s2.ErrDeploymentK8sGet
	}
	if !flag {
		return nil, k8s2.ErrNotFound
	}
	//maxCpu := resource.MustParse(req.MaxCpu)
	//reqCpu := resource.MustParse(req.Cpu)
	//maxMemory := resource.MustParse(req.MaxMemory)
	//reqMemory := resource.MustParse(req.Memory)
	//resources := corev1.ResourceList{}
	//limits := corev1.ResourceList{}

	//if maxMemory.Value() > 0 {
	//	limits[corev1.ResourceMemory] = *resource.NewQuantity(maxMemory.Value(), resource.BinarySI)
	//}
	//
	//if maxCpu.MilliValue() > 100 && maxCpu.MilliValue() < 1000 {
	//	limits[corev1.ResourceCPU] = *resource.NewMilliQuantity(maxCpu.MilliValue(), resource.BinarySI)
	//} else if maxCpu.Value() > 0 {
	//	limits[corev1.ResourceCPU] = *resource.NewQuantity(maxCpu.Value(), resource.BinarySI)
	//}
	//
	//if reqCpu.MilliValue() > 100 && reqCpu.MilliValue() < 1000 {
	//	resources[corev1.ResourceCPU] = *resource.NewMilliQuantity(reqCpu.MilliValue(), resource.BinarySI)
	//} else if reqCpu.Value() > 0 {
	//	resources[corev1.ResourceCPU] = *resource.NewQuantity(reqCpu.Value(), resource.BinarySI)
	//}
	//
	//resources[corev1.ResourceMemory] = *resource.NewQuantity(reqMemory.Value(), resource.BinarySI)

	switch req.Kind {
	case utils.DEPLOY_DEPLOYMENT:
		dep := deploy.(*v1.Deployment)
		dep.Spec.Template.Spec.Containers[0].Resources = *req.Resources
		//dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
		err = client.GetBaseClient().Deployment.Update(req.Namespace, dep)
	case utils.DEPLOY_STATEFULSET:
		dep := deploy.(*v1.StatefulSet)
		dep.Spec.Template.Spec.Containers[0].Resources = *req.Resources
		//dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
		err = client.GetBaseClient().Sf.Update(req.Namespace, dep)
	}
	if err != nil {
		log.Debugf("扩容服务失败 msg:%s\n", err.Error())
		return nil, k8s2.ErrDeploymentK8sUpdate
	}
	return nil, nil
}

// @Tags resource
// @Summary  容器伸缩
// @Produce  json
// @Accept   json
// @Param   params body StretchRequest true "参数"
// @Success 200 {string} json "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/resource/stretch [post]
func (dm *deployService) Stretch(ctx context.Context, req *StretchRequest) (*StretchResponse, error) {
	var err error
	switch req.Kind {
	case utils.DEPLOY_DEPLOYMENT:
		err = client.GetBaseClient().Deployment.Scale(req.Name, req.Namespace, req.Replicas)
	case utils.DEPLOY_STATEFULSET:
		err = client.GetBaseClient().Sf.Scale(req.Name, req.Namespace, req.Replicas)
	default:
		log.Debugf("容器伸缩失败，kind: %s, namespace: %s, name :%s 错误信息：%v \n", req.Kind, req.Namespace, req.Name, err)
		return nil, k8s2.ErrInvokerKind
	}
	if err != nil {
		log.Debugf("容器伸缩失败，kind: %s, namespace: %s, name :%s 错误信息：%v \n", req.Kind, req.Namespace, req.Name, err)
		return nil, k8s2.ErrDeploymentK8sScale
	}
	return nil, nil
}

// @Tags resource
// @Summary  版本回滚
// @Produce  json
// @Accept   json
// @Param   params body RollbackRequest true "参数"
// @Success 200 {string} json "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /resource/v1/resource/rollback [post]
func (dm *deployService) Rollback(ctx context.Context, req *RollbackRequest) (*RollbackResponse, error) {
	rollback, err := inital.GetGlobal().GetClientSet().AppsV1().ReplicaSets(req.Namespace).Get(req.VersionName, metav1.GetOptions{})
	if k8serror.IsNotFound(err) {
		return nil, errors.New("回滚版本不存在")
	} else if err != nil {
		return nil, k8s2.ErrReplicaSetGet
	}

	switch req.Kind {
	case utils.DEPLOY_DEPLOYMENT:
		deploy, err := client.GetBaseClient().Deployment.Get(req.Namespace, req.Name)
		if err != nil {
			return nil, k8s2.ErrDeploymentGet
		}
		deploy.Spec.Template = rollback.Spec.Template
		err = client.GetBaseClient().Deployment.Update(req.Namespace, deploy)
		if err != nil {
			return nil, k8s2.ErrUpdate
		}
	case utils.DEPLOY_STATEFULSET:
		deploy, err := client.GetBaseClient().Sf.Get(req.Namespace, req.Name)
		if err != nil {
			return nil, k8s2.ErrDeploymentGet
		}
		deploy.Spec.Template = rollback.Spec.Template
		err = client.GetBaseClient().Sf.Update(req.Namespace, deploy)
		if err != nil {
			return nil, k8s2.ErrUpdate
		}
	default:
		return nil, k8s2.ErrInvokerKind
	}

	return nil, nil

}
