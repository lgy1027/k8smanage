package k8s

import (
	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"relaper.com/kubemanage/inital"
	"relaper.com/kubemanage/utils"
)

type Sv struct {
	clientSet kubernetes.Interface
}

func NewSv(clientSet kubernetes.Interface) *Sv {
	return &Sv{clientSet}
}

func (sv *Sv) List(namespace string) (*apiv1.ServiceList, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Services(namespace).List(metav1.ListOptions{})
}

func (sv *Sv) Get(namespace, name string) (*apiv1.Service, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
}

func (sv *Sv) Delete(namespace, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return inital.GetGlobal().GetClientSet().CoreV1().Services(namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (sv *Sv) Create(namespace string, yaml *apiv1.Service) (*apiv1.Service, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Services(namespace).Create(yaml)
}

func (sv *Sv) DynamicGet(namespace string, name string) (*unstructured.Unstructured, error) {
	return inital.GetGlobal().GetResourceClient(utils.SERVICE).Namespace(namespace).Get(name, metav1.GetOptions{})
}

func (sv *Sv) DynamicCreate(namespace string, yaml *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	return inital.GetGlobal().GetResourceClient(utils.SERVICE).Namespace(namespace).Create(yaml, metav1.CreateOptions{})
}

func (sv *Sv) DynamicCreateForCustom(namespace, apiVersion string, yaml *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	return inital.GetGlobal().GetResourceClientForCustom(utils.SERVICE, apiVersion).Namespace(namespace).Create(yaml, metav1.CreateOptions{})
}

func (sv *Sv) Exist(namespace, name string) (*apiv1.Service, bool, error) {
	service, err := inital.GetGlobal().GetClientSet().CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	if k8serror.IsNotFound(err) {
		return nil, false, nil
	} else if statusError, isStatus := err.(*k8serror.StatusError); isStatus {
		return nil, false, errors.New(statusError.ErrStatus.Message)
	} else if err != nil {
		return nil, false, errors.New(err.Error())
	}
	return service, true, err
}
