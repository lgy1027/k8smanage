package k8s

import (
	"github.com/pkg/errors"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
)

var (
	ErrDeploymentK8sGet    = errors.New("Kubernetes获取错误,请查询是否存在")
	ErrDeploymentK8sUpdate = errors.New("Kubernetes更新错误,请联系管理员")
	ErrDeploymentK8sScale  = errors.New("Kubernetes伸缩错误,请联系管理员")
	ErrProjectPodsList     = errors.New("该项目的pods列表获取错误，请查看是否存在")
	ErrMetricsGet          = errors.New("获取指标失败")
	ErrNodeListGet         = errors.New("获取节点列表失败")
	ErrPodListGet          = errors.New("获取容器列表失败")
	ErrNamespaceListGet    = errors.New("获取命名空间列表失败")
	ErrNamespaceGet        = errors.New("获取命名空间失败")
	ErrNodeGet             = errors.New("获取节点失败")
	ErrDeploymentGet       = errors.New("获取无状态服务失败")
	ErrStatefulSetGet      = errors.New("获取有状态服务失败")
	ErrServiceGet          = errors.New("获取服务列表失败")
	ErrPodGet              = errors.New("获取容器列表失败")
	ErrReplicaSetGet       = errors.New("副本集获取失败")
	ErrUpdate              = errors.New("更新失败")
	ErrInvokerKind         = errors.New("未知资源类型")
	ErrNotFound            = errors.New("应用不存在")
	ErrDelete              = errors.New("删除失败")
	ErrCreate              = errors.New("创建失败")
	ErrExist               = errors.New("资源已存在")
)

func PrintErr(err error) error {
	if err != nil {
		if k8serror.IsAlreadyExists(err) {
			return errors.New("资源已存在")
		}
		if k8serror.IsNotFound(err) {
			return errors.New("资源不存在")
		}
		statusError, isStatus := err.(*k8serror.StatusError)
		if isStatus {
			return errors.New(statusError.ErrStatus.Message)
		}
		return err
	}
	return nil
}
