package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"relaper.com/kubemanage/inital"
)

type Node struct{}

func NewNode() Base {
	return &Node{}
}

func (node *Node) Get(namespace, name string) (interface{}, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Nodes().Get(name, metav1.GetOptions{})
}

func (node *Node) List(namespace string) (interface{}, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Nodes().List(metav1.ListOptions{})
}

func (node *Node) Delete(namespace, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return inital.GetGlobal().GetClientSet().CoreV1().Nodes().Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}
