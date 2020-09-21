package k8s

import (
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/retry"
	"relaper.com/kubemanage/inital"
	"relaper.com/kubemanage/utils"
)

type Deployment struct{}

func NewDeploy() *Deployment {
	return &Deployment{}
}

func (deployment *Deployment) List(namespace string) (interface{}, error) {
	return inital.GetGlobal().GetClientSet().AppsV1().Deployments(namespace).List(metav1.ListOptions{})
}

func (deployment *Deployment) Get(namespace, name string) (interface{}, error) {
	return inital.GetGlobal().GetClientSet().AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
}

func (deployment *Deployment) Delete(namespace, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return inital.GetGlobal().GetClientSet().AppsV1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (deployment *Deployment) Exist(namespace, name string) (*appsv1.Deployment, bool, error) {
	deploy, err := inital.GetGlobal().GetClientSet().AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
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
	return inital.GetGlobal().GetClientSet().AppsV1().Deployments(namespace).Create(yaml)
}

func (deployment *Deployment) Update(namespace string, result *appsv1.Deployment) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := inital.GetGlobal().GetClientSet().AppsV1().Deployments(namespace).Update(result)
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
