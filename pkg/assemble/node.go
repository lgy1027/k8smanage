package assemble

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"relaper.com/kubemanage/model"
	"strings"
	"sync"
)

var role = func(label map[string]string) string {
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

func AssembleNodes(nodes []v1.Node, pods []v1.Pod, nodeMetrics []v1beta1.NodeMetrics, podMetricsList []v1beta1.PodMetrics) []model.NodeDetail {
	nodeDetailList := make([]model.NodeDetail, 0)
	var wgC sync.WaitGroup
	for _, node := range nodes {
		if fmt.Sprintf("%s", node.Status.Conditions[len(node.Status.Conditions)-1].Type) != "Ready" {
			continue
		}
		wgC.Add(1)
		go func(node v1.Node) {
			nodeDetail := AssembleNode(node, pods, nodeMetrics, podMetricsList)
			nodeDetailList = append(nodeDetailList, nodeDetail)
			wgC.Done()
		}(node)
	}
	wgC.Wait()
	return nodeDetailList
}

func AssembleNode(node v1.Node, pods []v1.Pod, nodeMetrics []v1beta1.NodeMetrics, podMetricsList []v1beta1.PodMetrics) model.NodeDetail {
	var (
		active, total int64
		podsDetail    []model.PodDetail
	)
	var wg sync.WaitGroup
	if len(pods) > 0 {
		wg.Add(1)
		go func() {
			active, total = GetPodNum(node.Name, pods)
			podsDetail = AssemblePod(node.GetName(), pods, podMetricsList)
			wg.Done()
		}()
	}
	var resource model.ResourceDetail
	if len(nodeMetrics) > 0 {
		wg.Add(1)
		go func() {
			for _, metric := range nodeMetrics {
				if metric.GetName() != node.Name {
					continue
				}
				resource = AssembleResource(metric, node)
				wg.Done()
			}
		}()
	}
	nodeDetail := model.NodeDetail{
		Name:              node.GetName(),
		NodeID:            fmt.Sprintf("%s", node.GetUID()),
		HostIp:            node.Status.Addresses[0].Address,
		Status:            fmt.Sprintf("%s", node.Status.Conditions[len(node.Status.Conditions)-1].Type),
		IsValid:           fmt.Sprintf("%s", node.Status.Conditions[len(node.Status.Conditions)-1].Status),
		PodNum:            node.Status.Allocatable.Pods().Value(),
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
		Conditions:        node.Status.Conditions,
	}
	wg.Wait()
	nodeDetail.Resource = resource
	nodeDetail.PodRun = active
	nodeDetail.PodTotal = total
	nodeDetail.Pods = podsDetail
	return nodeDetail
}
