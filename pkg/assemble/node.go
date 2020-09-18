package assemble

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"relaper.com/kubemanage/model"
	"strings"
	"sync"
)

func AssembleNodes(nodes []v1.Node, pods []v1.Pod, nodeMetrics []v1beta1.NodeMetrics, podMetricsList []v1beta1.PodMetrics) []model.NodeDetail {
	role := func(label map[string]string) string {
		role := ""
		for k := range label {
			if index := strings.Index(k, "node-role.kubernetes.io/"); index == -1 {
				continue
			} else {
				slice := strings.Split(k, "/")
				if len(slice) > 1 {
					role = slice[1]
				}
				break
			}
		}
		return role
	}
	var wg sync.WaitGroup
	nodeDetail := make([]model.NodeDetail, 0)
	for _, node := range nodes {
		wg.Add(1)
		var (
			active, total int64
			podsDetail    []model.PodDetail
		)
		if len(pods) > 0 {
			go func() {
				active, total = GetPodNum(node.Name, pods)
				podsDetail = AssemblePod(node.GetName(), pods, podMetricsList)
				wg.Done()
			}()
		} else {
			wg.Done()
		}
		var resource model.ResourceDetail
		wg.Add(1)
		if len(nodeMetrics) > 0 {
			go func() {
				for _, metric := range nodeMetrics {
					if metric.GetName() == node.Name {
						resource = AssembleResource(metric, node)
						break
					}
				}
				wg.Done()
			}()
		} else {
			wg.Done()
		}
		nodeDetial := model.NodeDetail{
			Name:    node.GetName(),
			NodeID:  fmt.Sprintf("%s", node.GetUID()),
			HostIp:  node.Status.Addresses[0].Address,
			Status:  fmt.Sprintf("%s", node.Status.Conditions[len(node.Status.Conditions)-1].Type),
			IsValid: fmt.Sprintf("%s", node.Status.Conditions[len(node.Status.Conditions)-1].Status),
			PodNum:  node.Status.Allocatable.Pods().Value(),
			//PodTotal:          total,
			//PodRun:            active,
			Label:             node.GetLabels(),
			Annotation:        node.GetAnnotations(),
			CreateTime:        node.GetCreationTimestamp().String(),
			ImageNum:          len(node.Status.Images),
			KuBeLetVersion:    node.Status.NodeInfo.KubeletVersion,
			KuProxyVersion:    node.Status.NodeInfo.KubeProxyVersion,
			LastHeartbeatTime: node.Status.Conditions[0].LastHeartbeatTime.String(),
			SystemType:        node.Status.NodeInfo.OperatingSystem,
			SystemOs:          node.Status.NodeInfo.OSImage,
			DockVersion:       node.Status.NodeInfo.ContainerRuntimeVersion,
			KernlVersion:      node.Status.NodeInfo.KernelVersion,
			Role:              role(node.Labels),
			ClusterName:       node.GetClusterName(),
		}
		wg.Wait()
		nodeDetial.Resource = resource
		nodeDetial.PodRun = active
		nodeDetial.PodTotal = total
		nodeDetial.Pods = podsDetail
		nodeDetail = append(nodeDetail, nodeDetial)

	}
	return nodeDetail
}
