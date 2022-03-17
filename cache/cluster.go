package app

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/lgy1027/kubemanage/inital"
	"github.com/lgy1027/kubemanage/inital/client"
	k8s2 "github.com/lgy1027/kubemanage/k8s"
	"github.com/lgy1027/kubemanage/model"
	"github.com/lgy1027/kubemanage/pkg/assemble"
	"github.com/lgy1027/kubemanage/utils"
	"github.com/lgy1027/kubemanage/utils/tools"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"strconv"
	"sync"
)

func Cache() {
	cluster := GetClusterData()
	namespaceDetail := GetNamespaceDetail("")
	go CacheCluster(cluster)
	go CacheNamespace(namespaceDetail)
	go CachePods()
	go CacheNode(cluster.Nodes)
}

func CacheNode(nodes []model.NodeDetail) {
	for _, node := range nodes {
		nodeJson, err := json.Marshal(node)
		if err == nil {
			err = inital.GetGlobal().GetCache().Set(utils.NODE_PREFIX_KEY+node.Name, nodeJson, utils.NODE_TIME)
			if err != nil {
				log.Debugf("缓存集群数据失败，err:%v", err.Error())
			}
		} else {
			log.Debugf("集群信息json转换失败，err:%v", err.Error())
		}
	}
	nodesJson, err := json.Marshal(nodes)
	if err == nil {
		err = inital.GetGlobal().GetCache().Set(utils.NODE_PREFIX_KEY, nodesJson, utils.NODE_TIME)
		if err != nil {
			log.Debugf("缓存集群数据失败，err:%v", err.Error())
		}
	} else {
		log.Debugf("集群信息json转换失败，err:%v", err.Error())
	}
}

func CacheCluster(cluster *model.Cluster) {
	clusterJson, err := json.Marshal(cluster)
	if err == nil {
		err = inital.GetGlobal().GetCache().Set(utils.CLUSTER_PREFIX_KEY, clusterJson, utils.CLUSTER_DETAIL_TIME)
		if err != nil {
			log.Debugf("缓存集群数据失败，err:%v", err.Error())
		}
	} else {
		log.Debugf("集群信息json转换失败，err:%v", err.Error())
	}
}

func CacheNamespace(namespaceDetail []model.NamespaceDetail) {
	for _, ns := range namespaceDetail {
		data, err := json.Marshal(ns)
		if err == nil {
			err = inital.GetGlobal().GetCache().HSet(utils.NAMESPACE_PREFIX_KEY, utils.NAMESPACE_PREFIX_KEY+ns.Name, data, utils.NAMESPACE_TIME)
			if err != nil {
				log.Debugf("缓存命名空间数据失败，err:%v, Data:%v", err.Error(), ns.Name)
			}
		} else {
			log.Debugf("命名空间json转换失败，err:%v, Data:%v", err.Error(), ns.Name)
		}
	}
}

func GetClusterData() *model.Cluster {
	cluster := &model.Cluster{
		Nodes: make([]model.NodeDetail, 0),
	}
	var wg tools.WaitGroupWrapper
	wg.Wrap(func() {
		nodeDetailList, _ := GetNodes("")
		cluster.NodeNum = len(nodeDetailList)
		for _, node := range nodeDetailList {
			if node.Status == "Ready" {
				cluster.RunNodeNum++
			}
			cluster.PodNum += node.PodNum
			cluster.ActivePodNum += node.PodRun
		}
		cluster.Nodes = nodeDetailList
	})
	wg.Wrap(func() {
		namespaces, err := client.GetBaseClient().Ns.List("")
		if err != nil {
			log.Debug("获取命名空间失败，err:", err.Error())
			return
		}
		cluster.NameSpaceNum = len(namespaces.Items)
	})
	wg.Wait()
	var (
		cpu, mem, useCpu, useMem float64
	)
	if len(cluster.Nodes) > 0 {
		for _, node := range cluster.Nodes {
			cpuNum, _ := strconv.ParseFloat(node.Resource.CpuNum, 64)
			cpuUse, _ := strconv.ParseFloat(node.Resource.CpuUse, 64)
			memSize, _ := strconv.ParseFloat(node.Resource.MemSize, 64)
			memUse, _ := strconv.ParseFloat(node.Resource.MemUse, 64)
			cpu += cpuNum
			mem += memSize
			useCpu += cpuUse
			useMem += memUse

		}
		cluster.Resource.MemSize = strconv.FormatFloat(mem, 'f', -1, 64)
		cluster.Resource.MemUse = strconv.FormatFloat(useMem, 'f', -1, 64)
		cluster.Resource.CpuNum = strconv.FormatFloat(cpu, 'f', -1, 64)
		cluster.Resource.CpuUse = strconv.FormatFloat(useCpu, 'f', -1, 64)
	}
	return cluster
}

func GetNodes(name string) ([]model.NodeDetail, error) {
	var (
		err      error
		nodeList []v1.Node
	)
	if name != "" {
		node, err := client.GetBaseClient().Node.Get("", name)
		if err != nil {
			log.Debug("获取节点信息失败,NodeName:", name)
			return nil, errors.New("获取节点信息失败")
		}
		nodeList = []v1.Node{*node}
	} else {
		nodes, err := client.GetBaseClient().Node.List("")
		if err != nil {
			log.Debug("获取节点列表失败，err:", err.Error())
			return nil, errors.New("获取节点列表失败")
		}
		nodeList = nodes.Items
	}

	var (
		pods            *v1.PodList
		nodeMetricsList []v1beta1.NodeMetrics
		podMetrics      []v1beta1.PodMetrics
	)

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		opts := metav1.ListOptions{}
		if name != "" {
			opts = metav1.ListOptions{
				FieldSelector: fmt.Sprintf("%s=%s", "spec.nodeName", name),
			}
		}
		pods, err = client.GetBaseClient().Pod.List("", opts)
		if err != nil {
			log.Debug("获取Pod列表失败，err:", err.Error())
		}
		wg.Done()
	}()

	go func() {
		nodeMetricsList, err = k8s2.GetNodeListMetrics()
		if err != nil {
			log.Debug("获取节点指标失败，err:", err.Error())
		}
		wg.Done()
	}()

	go func() {
		podMetrics, err = k8s2.GetPodListMetrics("", metav1.ListOptions{})
		if err != nil {
			log.Debug("获取pod指标失败，err:", err.Error())
		}
		wg.Done()
	}()
	wg.Wait()
	var nodeDetailList []model.NodeDetail
	if name != "" {
		nodeDetail := assemble.AssembleNode(nodeList[0], pods.Items, nodeMetricsList, podMetrics)
		nodeDetailList = append(nodeDetailList, nodeDetail)
	} else {
		nodeDetailList = assemble.AssembleNodes(nodeList, pods.Items, nodeMetricsList, podMetrics)
	}
	return nodeDetailList, nil
}

func GetNamespaceDetail(name string) []model.NamespaceDetail {
	namespaceDetail := make([]model.NamespaceDetail, 0)
	if name == "" {
		namespaces, err := client.GetBaseClient().Ns.List("")
		if err != nil {
			log.Debug("获取命名空间失败，err:", err.Error())
			return nil
		}
		items := namespaces.Items
		var wg sync.WaitGroup
		for _, ns := range items {
			wg.Add(1)
			go func(namespace v1.Namespace) {
				nsDetail := GetDetailForRange(namespace)
				namespaceDetail = append(namespaceDetail, nsDetail)
				wg.Done()
			}(ns)
		}
		wg.Wait()
	} else {
		namespace, err := client.GetBaseClient().Ns.Get(name, "")
		if err != nil {
			log.Debug("获取命名空间失败，err:", err.Error())
			return nil
		}
		namespaceDetail = append(namespaceDetail, GetDetailForRange(*namespace))
	}
	return namespaceDetail
}

func GetDetailForRange(namespace v1.Namespace) model.NamespaceDetail {
	namespaceDetail := model.NamespaceDetail{
		Name:       namespace.GetName(),
		CreateTime: namespace.GetCreationTimestamp().String(),
		Status:     string(namespace.Status.Phase),
	}
	name := namespace.GetName()

	var wg tools.WaitGroupWrapper
	wg.Wrap(func() {
		deploys, err := client.GetBaseClient().Deployment.List(name)
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get deployment is err:%v", err.Error())
		} else if deploys != nil && len(deploys.Items) > 0 {
			namespaceDetail.DeploymentList = assemble.AssembleDeployment(name, deploys.Items)
		}
	})

	wg.Wrap(func() {
		stats, err := client.GetBaseClient().Sf.List(name, metav1.ListOptions{})
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get statefulSet is err:%v", err.Error())
		} else if stats != nil && len(stats.Items) > 0 {
			namespaceDetail.StatefulSetList = assemble.AssembleStatefulSet(name, stats.Items)
		}
	})

	wg.Wrap(func() {
		svcs, err := client.GetBaseClient().Sv.List(name)
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get service is err:%v", err.Error())
		} else if svcs != nil && len(svcs.Items) > 0 {
			namespaceDetail.ServiceList = assemble.AssembleService(name, svcs.Items)
		}
	})
	wg.Wait()
	return namespaceDetail
}

func CachePods() {
	list, err := client.GetBaseClient().Pod.List("", metav1.ListOptions{})
	if err != nil {
		log.Debugf("Method [CachePods] = > Get pod list err:%v", err.Error())
		return
	}
	podList := list.Items
	for _, pod := range podList {
		podJson, err := json.Marshal(pod)
		if err == nil {
			err = inital.GetGlobal().GetCache().Set(utils.POD_PREFIX_KEY+pod.GetNamespace()+":"+pod.GetName(), podJson, utils.POD_TIME)
			if err != nil {
				log.Debugf("缓存pod数据失败，err:%v", err.Error())
			}
		} else {
			log.Debugf("pod信息json转换失败，err:%v", err.Error())
		}
	}
	podListJson, err := json.Marshal(podList)
	if err == nil {
		err = inital.GetGlobal().GetCache().Set(utils.POD_PREFIX_KEY, podListJson, utils.POD_TIME)
		if err != nil {
			log.Debugf("缓存pod数据失败，err:%v", err.Error())
		}
	} else {
		log.Debugf("pod信息json转换失败，err:%v", err.Error())
	}
}
