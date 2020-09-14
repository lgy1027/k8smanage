package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"relaper.com/kubemanage/inital"
)

type Pod struct{}

func NewPod() Base {
	return &Pod{}
}

func (pos *Pod) Get(namespace, name string) (interface{}, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
}

func (pos *Pod) List(namespace string) (interface{}, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Pods(namespace).List(metav1.ListOptions{})
}

// 删除某个pod后自动重建
func (pos *Pod) Delete(namespace, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return inital.GetGlobal().GetClientSet().CoreV1().Pods(namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}
