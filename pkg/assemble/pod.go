package assemble

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"relaper.com/kubemanage/model"
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

func AssemblePod(node string, pods []v1.Pod) []model.PodDetail {
	podDetail := make([]model.PodDetail, 0)
	for _, pod := range pods {
		if node != "" {
			if pod.Spec.NodeName != node {
				continue
			}
		}
		resource := model.ResourceDetail{}
		podDetail = append(podDetail, model.PodDetail{
			Name:       pod.GetName(),
			NodeName:   pod.Spec.NodeName,
			Namespace:  pod.GetNamespace(),
			Id:         fmt.Sprintf("%s", pod.GetUID()),
			Status:     fmt.Sprintf("%s", pod.Status.Phase),
			CreateTime: pod.GetCreationTimestamp().String(),
			Label:      pod.GetLabels(),
			Annotation: pod.GetAnnotations(),
			Resource:   resource,
		})
	}
	return podDetail
}
