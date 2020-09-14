package k8s

import (
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"relaper.com/kubemanage/inital"
	"relaper.com/kubemanage/utils"
)

type Sf struct{}

func NewStateFulSet() *Sf {
	return &Sf{}
}

func (ssf *Sf) Get(namespace, name string) (interface{}, error) {
	return inital.GetGlobal().GetClientSet().AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
}

func (ssf *Sf) Delete(namespace, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return inital.GetGlobal().GetClientSet().AppsV1().StatefulSets(namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (ssf *Sf) List(namespace string) (interface{}, error) {
	return inital.GetGlobal().GetClientSet().AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
}

func (ssf *Sf) Create(namespace string, yaml *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return inital.GetGlobal().GetClientSet().AppsV1().StatefulSets(namespace).Create(yaml)
}

func (ssf *Sf) DynamicCreate(namespace string, yaml *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	return inital.GetGlobal().GetResourceClient(utils.STATEFULSET).Namespace(namespace).Create(yaml, metav1.CreateOptions{})
}

func (ssf *Sf) DynamicCreateForCustom(namespace, apiVersion string, yaml *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	return inital.GetGlobal().GetResourceClientForCustom(utils.STATEFULSET, apiVersion).Namespace(namespace).Create(yaml, metav1.CreateOptions{})
}

func (ssf *Sf) Exist(namespace, name string) (*appsv1.StatefulSet, bool, error) {
	deploy, err := inital.GetGlobal().GetClientSet().AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
	if k8serror.IsNotFound(err) {
		return nil, false, nil
	} else if statusError, isStatus := err.(*k8serror.StatusError); isStatus {
		return nil, false, errors.New(statusError.ErrStatus.Message)
	} else if err != nil {
		return nil, false, errors.New(err.Error())
	}
	return deploy, true, nil
}
