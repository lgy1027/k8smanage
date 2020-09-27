package k8s

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"relaper.com/kubemanage/inital"
)

type Node struct {
	clientSet kubernetes.Interface
}

func NewNode(clientSet kubernetes.Interface) *Node {
	return &Node{clientSet}
}

func (node *Node) Get(namespace, name string) (*v1.Node, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Nodes().Get(name, metav1.GetOptions{})
}

func (node *Node) List(namespace string) (*v1.NodeList, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Nodes().List(metav1.ListOptions{})
}

func (node *Node) Delete(namespace, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return inital.GetGlobal().GetClientSet().CoreV1().Nodes().Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}
