package app

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"relaper.com/kubemanage/inital"
	k8s2 "relaper.com/kubemanage/k8s"
	"relaper.com/kubemanage/model"
	"relaper.com/kubemanage/pkg/assemble"
	"relaper.com/kubemanage/utils"
	"relaper.com/kubemanage/utils/tools"
	"strconv"
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
		metricsList := make([]v1beta1.NodeMetrics, 0)
		metrics := inital.GetGlobal().GetMetricsClient()
		nodeMetricsList, err := metrics.MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{})
		if err != nil {
			metricsList = nil
			log.Debug("获取节点指标失败，err:", err.Error())
		} else {
			metricsList = nodeMetricsList.Items
		}

		podsList := pods.(*v1.PodList).Items
		nodeDetailList := assemble.AssembleNodes(nodeList, podsList, metricsList)
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
		cpu    float64
		mem    float64
		useCpu float64
		useMem float64
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
