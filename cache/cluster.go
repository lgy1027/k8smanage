package app

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"relaper.com/kubemanage/inital"
	k8s2 "relaper.com/kubemanage/k8s"
	"relaper.com/kubemanage/model"
	"relaper.com/kubemanage/pkg/assemble"
	"relaper.com/kubemanage/utils"
	"relaper.com/kubemanage/utils/tools"
	"strconv"
	"sync"
)

var (
	node        k8s2.Base
	namespace   *k8s2.Ns
	pod         k8s2.Base
	deployment  *k8s2.Deployment
	statefulSet *k8s2.Sf
	service     *k8s2.Sv
)

func init() {
	node = k8s2.NewNode()
	namespace = k8s2.NewNs()
	pod = k8s2.NewPod()
	deployment = k8s2.NewDeploy()
	statefulSet = k8s2.NewStateFulSet()
	service = k8s2.NewSv()
}

func CacheCluster() {
	cluster := GetClusterData()
	namespaceDetail := GetNamespaceDetail("")
	for _, ns := range namespaceDetail {
		data, err := json.Marshal(ns)
		if err == nil {
			err = inital.GetGlobal().GetCache().Set(utils.NAMESPACE_PREFIX_KEY+ns.Name, data, utils.NAMESPACE_TIME)
			if err != nil {
				log.Debugf("缓存命名空间数据失败，err:%v, Data:%v", err.Error(), ns.Name)
			}
		} else {
			log.Debugf("命名空间json转换失败，err:%v, Data:%v", err.Error(), ns.Name)
		}
	}
	namespaceJson, err := json.Marshal(namespaceDetail)
	if err == nil {
		err = inital.GetGlobal().GetCache().Set(utils.NAMESPACE_PREFIX_KEY, namespaceJson, utils.NAMESPACE_TIME)
		if err != nil {
			log.Debugf("缓存命名空间数据失败，err:%v, Data:%v", err.Error(), cluster.NameSpaceNum)
		}
	} else {
		log.Debugf("命名空间json转换失败，err:%v, Data:%v", err.Error(), cluster.NameSpaceNum)
	}
	clusterJson, err := json.Marshal(cluster)
	if err == nil {
		err = inital.GetGlobal().GetCache().Set(utils.CLUSTER_PREFIX_KEY, clusterJson, utils.CLUSTER_DETAIL_TIME)
		if err != nil {
			log.Debugf("缓存集群数据失败，err:%v, Data:%v", err.Error(), cluster)
		}
	} else {
		log.Debugf("集群信息json转换失败，err:%v, Data:%v", err.Error(), cluster)
	}
}

func GetClusterData() *model.Cluster {
	cluster := &model.Cluster{
		Nodes: make([]model.NodeDetail, 0),
	}
	var wg tools.WaitGroupWrapper
	wg.Wrap(func() {
		nodes, err := node.List("")
		if err != nil {
			log.Debug("获取节点列表失败，err:", err.Error())
			return
		}
		nodeList := nodes.(*v1.NodeList).Items
		cluster.NodeNum = len(nodeList)
		pods, err := pod.List("")
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

		podsList := pods.(*v1.PodList).Items
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
		namespaces, err := namespace.List("")
		if err != nil {
			log.Debug("获取命名空间失败，err:", err.Error())
			return
		}
		items := namespaces.(*v1.NamespaceList).Items
		cluster.NameSpaceNum = len(items)
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
		namespaces, err := namespace.List("")
		if err != nil {
			log.Debug("获取命名空间失败，err:", err.Error())
			return nil
		}
		items := namespaces.(*v1.NamespaceList).Items
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
		namespace, err := namespace.Get(name, "")
		if err != nil {
			log.Debug("获取命名空间失败，err:", err.Error())
			return nil
		}
		ns := namespace.(*v1.Namespace)
		namespaceDetail = append(namespaceDetail, GetDetailForRange(*ns))
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
	//var wg tools.WaitGroupWrapper
	//wg.Wrap(func() {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		deploys, err := deployment.List(name)
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get deployment is err:%v", err.Error())
		} else if len(deploys.(*appsv1.DeploymentList).Items) > 0 {
			namespaceDetail.DeploymentList = assemble.AssembleDeployment(name, deploys.(*appsv1.DeploymentList).Items)
		}
		wg.Done()
	}()
	//})
	//wg.Wrap(func() {
	go func() {
		stats, err := statefulSet.List(name)
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get statefulSet is err:%v", err.Error())
		} else if len(stats.(*appsv1.StatefulSetList).Items) > 0 {
			namespaceDetail.StatefulSetList = assemble.AssembleStatefulSet(name, stats.(*appsv1.StatefulSetList).Items)
		}
		wg.Done()
	}()
	//})
	//wg.Wrap(func() {
	go func() {
		svcs, err := service.List(name)
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get service is err:%v", err.Error())
		} else if len(svcs.(*v1.ServiceList).Items) > 0 {
			namespaceDetail.ServiceList = assemble.AssembleService(name, svcs.(*v1.ServiceList).Items)
		}
		wg.Done()
	}()
	//})
	wg.Wait()
	return namespaceDetail
}
