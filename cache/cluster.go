package app

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	v1 "k8s.io/api/core/v1"
	"relaper.com/kubemanage/inital"
	"relaper.com/kubemanage/inital/client"
	k8s2 "relaper.com/kubemanage/k8s"
	"relaper.com/kubemanage/model"
	"relaper.com/kubemanage/pkg/assemble"
	"relaper.com/kubemanage/utils"
	"relaper.com/kubemanage/utils/tools"
	"strconv"
	"sync"
)

func Cache() {
	cluster := GetClusterData()
	namespaceDetail := GetNamespaceDetail("")
	go CacheCluster(cluster)
	go CacheNamespace(namespaceDetail)
	go CachePods()
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
		nodes, err := client.GetBaseClient().Node.List("")
		if err != nil {
			log.Debug("获取节点列表失败，err:", err.Error())
			return
		}
		nodeList := nodes.Items
		cluster.NodeNum = len(nodeList)
		pods, err := client.GetBaseClient().Pod.List("")
		if err != nil {
			log.Debug("获取Pod列表失败，err:", err.Error())
		}
		nodeMetricsList, err := k8s2.GetNodeListMetrics()
		if err != nil {
			log.Debug("获取节点指标失败，err:", err.Error())
		}

		podMetrics, err := k8s2.GetPodListMetrics("")
		if err != nil {
			log.Debug("获取pod指标失败，err:", err.Error())
		}

		podsList := pods.Items
		nodeDetailList := assemble.AssembleNodes(nodeList, podsList, nodeMetricsList, podMetrics)
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

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		deploys, err := client.GetBaseClient().Deployment.List(name)
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get deployment is err:%v", err.Error())
		} else if len(deploys.Items) > 0 {
			namespaceDetail.DeploymentList = assemble.AssembleDeployment(name, deploys.Items)
		}
		wg.Done()
	}()

	go func() {
		stats, err := client.GetBaseClient().Sf.List(name)
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get statefulSet is err:%v", err.Error())
		} else if len(stats.Items) > 0 {
			namespaceDetail.StatefulSetList = assemble.AssembleStatefulSet(name, stats.Items)
		}
		wg.Done()
	}()

	go func() {
		svcs, err := client.GetBaseClient().Sv.List(name)
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get service is err:%v", err.Error())
		} else if len(svcs.Items) > 0 {
			namespaceDetail.ServiceList = assemble.AssembleService(name, svcs.Items)
		}
		wg.Done()
	}()
	wg.Wait()
	return namespaceDetail
}

func CachePods() {
	list, err := client.GetBaseClient().Pod.List("")
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
