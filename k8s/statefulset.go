package k8s

import (
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

type Sf struct {
	clientSet kubernetes.Interface
}

func NewStateFulSet(clientSet kubernetes.Interface) *Sf {
	return &Sf{clientSet}
}

func (ssf *Sf) Get(namespace, name string) (*appsv1.StatefulSet, error) {
	return inital.GetGlobal().GetClientSet().AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
}

func (ssf *Sf) Delete(namespace, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return inital.GetGlobal().GetClientSet().AppsV1().StatefulSets(namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func (ssf *Sf) List(namespace string, opt metav1.ListOptions) (*appsv1.StatefulSetList, error) {
	return inital.GetGlobal().GetClientSet().AppsV1().StatefulSets(namespace).List(opt)
}

func (ssf *Sf) Create(namespace string, yaml *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return inital.GetGlobal().GetClientSet().AppsV1().StatefulSets(namespace).Create(yaml)
}

func (ssf *Sf) DynamicGet(namespace string, name string) (*unstructured.Unstructured, error) {
	return inital.GetGlobal().GetResourceClient(utils.STATEFULSET).Namespace(namespace).Get(name, metav1.GetOptions{})
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

func (ssf *Sf) Update(namespace string, result *appsv1.StatefulSet) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := inital.GetGlobal().GetClientSet().AppsV1().StatefulSets(namespace).Update(result)
		return updateErr
	})
	return PrintErr(retryErr)
}

func (ssf *Sf) GetPods(name, namespace string) (*corev1.PodList, error) {
	deploy, err := ssf.Get(name, namespace)
	if err != nil {
		return nil, err
	}
	labelSelector, err := metav1.LabelSelectorAsSelector(deploy.Spec.Selector)
	if err != nil {
		return nil, err
	}
	opt := metav1.ListOptions{LabelSelector: labelSelector.String()}
	podList, err := ssf.clientSet.CoreV1().Pods(namespace).List(opt)
	return podList, err
}

func (ssf *Sf) Scale(name, namespace string, replicas int32) error {
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		scale, err := ssf.clientSet.AppsV1().StatefulSets(namespace).GetScale(name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		scale.Spec.Replicas = replicas
		_, err = ssf.clientSet.AppsV1().StatefulSets(namespace).UpdateScale(name, scale)
		return err
	}); err != nil {
		return err
	}
	return nil
}
