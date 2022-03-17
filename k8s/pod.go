package k8s

import (
	"bytes"
	"github.com/lgy1027/kubemanage/inital"
	"github.com/pkg/errors"
	"io"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Pod struct {
	clientSet kubernetes.Interface
}

func NewPod(clientSet kubernetes.Interface) *Pod {
	return &Pod{clientSet}
}

func (pos *Pod) Get(namespace, name string) (*apiv1.Pod, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
}

func (pos *Pod) List(namespace string, opt metav1.ListOptions) (*apiv1.PodList, error) {
	return inital.GetGlobal().GetClientSet().CoreV1().Pods(namespace).List(opt)
}

// 删除某个pod后自动重建
func (pos *Pod) Delete(namespace, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return inital.GetGlobal().GetClientSet().CoreV1().Pods(namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

// 删除某个pod后自动重建
func (pos *Pod) Log(namespace, name string) (string, error) {
	req := inital.GetGlobal().GetClientSet().CoreV1().Pods(namespace).GetLogs(name, &apiv1.PodLogOptions{})
	podLogs, err := req.Stream()
	if err != nil {
		return "", errors.New("查看日志失败")
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", errors.New("内部错误")
	}
	str := buf.String()

	return str, nil
}
