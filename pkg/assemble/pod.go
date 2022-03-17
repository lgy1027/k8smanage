package assemble

import (
	"fmt"
	"github.com/lgy1027/kubemanage/model"
	"github.com/shopspring/decimal"
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"strconv"
	"sync"
)

func GetPodNum(node string, pods []v1.Pod) (int64, int64) {
	var (
		active int64
		total  int64
	)
	for _, pod := range pods {
		if pod.Spec.NodeName != node {
			continue
		}
		total++
		if pod.Status.Phase != v1.PodRunning {
			continue
		}
		active++
	}
	return active, total
}

func AssemblePodSummary(pod v1.Pod, metric *v1beta1.PodMetrics) model.PodDetail {
	podInfo := model.PodDetail{
		Name:         pod.GetName(),
		NodeName:     pod.Spec.NodeName,
		Namespace:    pod.GetNamespace(),
		Id:           fmt.Sprintf("%s", pod.GetUID()),
		Status:       fmt.Sprintf("%s", pod.Status.Phase),
		CreateTime:   pod.GetCreationTimestamp().String(),
		Label:        pod.GetLabels(),
		Annotation:   pod.GetAnnotations(),
		RestartCount: pod.Status.ContainerStatuses[0].RestartCount,
		HostIp:       pod.Status.HostIP,
		PodIp:        pod.Status.PodIP,
	}
	if metric != nil && len(metric.Containers) > 0 {
		var resource model.ResourceDetail
		memUseValue := decimal.NewFromInt(metric.Containers[0].Usage.Memory().Value())
		memUseValue = memUseValue.Div(decimal.NewFromInt(1024)).DivRound(decimal.NewFromInt(1024), 2)
		resource.MemUse = memUseValue.String()
		resource.CpuUse = strconv.FormatInt(metric.Containers[0].Usage.Cpu().MilliValue(), 10)
		podInfo.Resource = resource
	}

	return podInfo
}

func AssemblePod(node string, pods []v1.Pod, podMetricsList []v1beta1.PodMetrics) []model.PodDetail {
	podDetail := make([]model.PodDetail, 0)
	var wg sync.WaitGroup
	for _, pod := range pods {
		if node != "" {
			if pod.Spec.NodeName != node {
				continue
			}
		}
		wg.Add(1)
		go func(pod v1.Pod) {
			podInfo := AssemblePodSimple(pod, podMetricsList)
			podDetail = append(podDetail, podInfo)
			wg.Done()
		}(pod)
	}
	wg.Wait()
	return podDetail
}

func AssemblePodSimple(pod v1.Pod, podMetricsList []v1beta1.PodMetrics) model.PodDetail {
	var wg sync.WaitGroup
	var resource model.ResourceDetail
	if len(podMetricsList) > 0 {
		wg.Add(1)
		go func() {
			for _, metric := range podMetricsList {
				if metric.GetNamespace() == pod.GetNamespace() && metric.GetName() == pod.GetName() {
					if len(metric.Containers) > 0 {
						memUseValue := decimal.NewFromInt(metric.Containers[0].Usage.Memory().Value())
						memUseValue = memUseValue.Div(decimal.NewFromInt(1024)).DivRound(decimal.NewFromInt(1024), 2)
						resource.MemUse = memUseValue.String()
						resource.CpuUse = strconv.FormatInt(metric.Containers[0].Usage.Cpu().MilliValue(), 10)
					}
					break
				}
			}
			wg.Done()
		}()
	}
	podInfo := model.PodDetail{
		Name:         pod.GetName(),
		NodeName:     pod.Spec.NodeName,
		Namespace:    pod.GetNamespace(),
		Id:           fmt.Sprintf("%s", pod.GetUID()),
		Status:       fmt.Sprintf("%s", pod.Status.Phase),
		CreateTime:   pod.GetCreationTimestamp().String(),
		Label:        pod.GetLabels(),
		Annotation:   pod.GetAnnotations(),
		RestartCount: pod.Status.ContainerStatuses[0].RestartCount,
		HostIp:       pod.Status.HostIP,
		PodIp:        pod.Status.PodIP,
	}
	wg.Wait()
	podInfo.Resource = resource
	//podInfo.EventData = event
	return podInfo
}
