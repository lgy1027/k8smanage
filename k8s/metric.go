package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"relaper.com/kubemanage/inital"
)

func GetPodListMetrics(namespace string) ([]v1beta1.PodMetrics, error) {
	podMetricsList := make([]v1beta1.PodMetrics, 0)
	podMetrics, err := inital.GetGlobal().GetMetricsClient().MetricsV1beta1().PodMetricses(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	podMetricsList = podMetrics.Items
	return podMetricsList, nil
}

func GetPodMetrics(namespace, podName string) (*v1beta1.PodMetrics, error) {
	podMetric, err := inital.GetGlobal().GetMetricsClient().MetricsV1beta1().PodMetricses(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return podMetric, nil
}

func GetNodeListMetrics() ([]v1beta1.NodeMetrics, error) {
	nodeMetricsList := make([]v1beta1.NodeMetrics, 0)
	podMetrics, err := inital.GetGlobal().GetMetricsClient().MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodeMetricsList = podMetrics.Items
	return nodeMetricsList, nil
}

func GetNodeMetrics(nodeName string) (*v1beta1.NodeMetrics, error) {
	nodeMetric, err := inital.GetGlobal().GetMetricsClient().MetricsV1beta1().NodeMetricses().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return nodeMetric, nil
}
