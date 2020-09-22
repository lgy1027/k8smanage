package k8s

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"relaper.com/kubemanage/inital"
)

type Ns struct{}

func NewNs() *Ns {
	return &Ns{}
}

func (ns *Ns) List(namespace string) (*apiv1.NamespaceList, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Namespaces().List(metav1.ListOptions{})
}

func (ns *Ns) Get(namespace, name string) (*apiv1.Namespace, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
}

func (ns *Ns) Create(name string) (*apiv1.Namespace, error) {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Status: apiv1.NamespaceStatus{
			Phase: apiv1.NamespaceActive,
		},
	}
	return inital.GetGlobal().GetClientSet().CoreV1().Namespaces().Create(namespace)
}

func (ns *Ns) Delete(namespace, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return inital.GetGlobal().GetClientSet().CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}
