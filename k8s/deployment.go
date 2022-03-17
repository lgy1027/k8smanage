package k8s

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/lgy1027/kubemanage/inital"
	"github.com/lgy1027/kubemanage/utils"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type Deployment struct {
	clientSet kubernetes.Interface
}

func NewDeploy(clientSet kubernetes.Interface) *Deployment {
	return &Deployment{clientSet}
}

func (deployment *Deployment) List(namespace string) (*appsv1.DeploymentList, error) {
	return deployment.clientSet.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
}

func (deployment *Deployment) Get(namespace, name string) (*appsv1.Deployment, error) {
	return deployment.clientSet.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
}

func (deployment *Deployment) Delete(namespace, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return deployment.clientSet.AppsV1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (deployment *Deployment) Exist(namespace, name string) (*appsv1.Deployment, bool, error) {
	deploy, err := deployment.clientSet.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if k8serror.IsNotFound(err) {
		return nil, false, nil
	} else if statusError, isStatus := err.(*k8serror.StatusError); isStatus {
		return nil, false, errors.New(statusError.ErrStatus.Message)
	} else if err != nil {
		return nil, false, errors.New(err.Error())
	}
	return deploy, true, nil
}

func (deployment *Deployment) Create(namespace string, yaml *appsv1.Deployment) (*appsv1.Deployment, error) {
	return deployment.clientSet.AppsV1().Deployments(namespace).Create(yaml)
}

func (deployment *Deployment) Update(namespace string, result *appsv1.Deployment) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := deployment.clientSet.AppsV1().Deployments(namespace).Update(result)
		return updateErr
	})
	return PrintErr(retryErr)
}

func (deployment *Deployment) DynamicGet(namespace, name string) (*unstructured.Unstructured, error) {
	return inital.GetGlobal().GetResourceClient(utils.DEPLOYMENT).Namespace(namespace).Get(name, metav1.GetOptions{})
}

func (deployment *Deployment) DynamicCreate(namespace string, yaml *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	return inital.GetGlobal().GetResourceClient(utils.DEPLOYMENT).Namespace(namespace).Create(yaml, metav1.CreateOptions{})
}

func (deployment *Deployment) DynamicCreateForCustom(namespace, apiVersion string, yaml *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	dc := inital.GetGlobal().GetResourceClientForCustom(utils.DEPLOYMENT, apiVersion).Namespace(namespace)
	return dc.Create(yaml, metav1.CreateOptions{})
}

func (deployment *Deployment) GetPods(name, namespace string) (*corev1.PodList, error) {
	deploy, err := deployment.Get(name, namespace)
	if err != nil {
		return nil, errors.New("获取资源失败")
	}
	labelSelector, err := metav1.LabelSelectorAsSelector(deploy.Spec.Selector)
	if err != nil {
		return nil, errors.New("获取标签失败")
	}
	opt := metav1.ListOptions{LabelSelector: labelSelector.String()}
	podList, err := deployment.clientSet.CoreV1().Pods(namespace).List(opt)
	return podList, err
}

func (deployment *Deployment) GetLatestReplicaSet(name, namespace string) (*appsv1.Deployment, string, error) {
	deploy, err := deployment.Get(name, namespace)
	if err != nil {
		log.Debug(err.Error())
		return nil, "", errors.New("查找资源出错")
	}
	revision := deploy.Annotations["deployment.kubernetes.io/revision"]
	labelSelector, err := metav1.LabelSelectorAsSelector(deploy.Spec.Selector)
	if err != nil {
		log.Debug(err.Error())
		return nil, "", errors.New("标签解析错误")
	}

	opt := metav1.ListOptions{LabelSelector: labelSelector.String()}
	replicasets, err := deployment.clientSet.AppsV1().ReplicaSets(namespace).List(opt)
	if err != nil {
		log.Debug(err.Error())
		return nil, "", errors.New("查找副本集错误")
	}
	for _, rs := range replicasets.Items {
		if rs.Annotations["deployment.kubernetes.io/revision"] == revision {
			return deploy, rs.Name, nil
		}
	}
	return nil, "", fmt.Errorf("尚未创建最新副本")
}

func (deployment *Deployment) Scale(name, namespace string, replicas int32) error {
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		scale, err := deployment.clientSet.AppsV1().Deployments(namespace).GetScale(name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		scale.Spec.Replicas = replicas
		_, err = deployment.clientSet.AppsV1().Deployments(namespace).UpdateScale(name, scale)
		return err
	}); err != nil {
		return err
	}
	return nil
}
