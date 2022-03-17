package cluster

import (
	"context"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/lgy1027/kubemanage/inital"
	"github.com/lgy1027/kubemanage/inital/client"
	k8s2 "github.com/lgy1027/kubemanage/k8s"
	"github.com/lgy1027/kubemanage/model"
	"github.com/lgy1027/kubemanage/pkg/assemble"
	"github.com/lgy1027/kubemanage/utils"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"strconv"
	"sync"
)

var (
	lineReadLimit int64 = 5000
	byteReadLimit int64 = 5000000
)

type Service interface {
	Cluster(ctx context.Context, req *ClusterRequest) (*ClusterResponse, error)
	Nodes(ctx context.Context, req *NodesRequest) (*NodesResponse, error)
	Node(ctx context.Context, req *NodeRequest) (*NodeResponse, error)
	Ns(ctx context.Context, req *NsRequest) (*NsResponse, error)
	NameSpace(ctx context.Context, req *NameSpaceRequest) (*NameSpaceResponse, error)
	PodInfo(ctx context.Context, req *PodInfoRequest) (*PodInfoResponse, error)
	//PodLog(ctx context.Context, req *PodLogRequest) (*PodLogResponse, error)
	Pods(ctx context.Context, req *PodsRequest) (*PodsResponse, error)
	Deployment(ctx context.Context, req *ResourceRequest) (*DeploymentsResponse, error)
	StatefulSet(ctx context.Context, req *ResourceRequest) (*StatefulSetsResponse, error)
	Services(ctx context.Context, req *ResourceRequest) (*ServiceResponse, error)
	GetYaml(ctx context.Context, req *GetYamlRequest) (*GetYamlResponse, error)
	Event(ctx context.Context, req *EventRequest) (*EventResponse, error)
	VersionList(ctx context.Context, req *VersionRequest) (*VersionResponse, error)
}

// NewService return a Service interface
func NewService() Service {
	return &clusterService{}
}

type clusterService struct{}

// @Tags cluster
// @Summary 获取集群信息
// @Produce  json
// @Success 200 {object} protocol.Response{data=ClusterResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/detail [post]
func (cs *clusterService) Cluster(ctx context.Context, req *ClusterRequest) (*ClusterResponse, error) {

	nodes, err := client.GetBaseClient().Node.List("")
	if err != nil {
		log.Debug("获取节点列表失败,err:", err.Error())
		return nil, k8s2.ErrNodeListGet
	}
	cluster := model.Cluster{
		Nodes: make([]model.NodeDetail, 0),
	}
	namespaces, err := client.GetBaseClient().Ns.List("")
	if err != nil {
		log.Debug("获取命名空间列表失败,err:", err.Error())
		return nil, k8s2.ErrNamespaceListGet
	}
	cluster.NameSpaceNum = len(namespaces.Items)
	cluster.NodeNum = len(nodes.Items)
	var (
		cpu, mem, useCpu, useMem float64
	)
	for _, node := range nodes.Items {
		nodeDetail := model.NodeDetail{
			Name:   node.GetName(),
			HostIp: node.Status.Addresses[0].Address,
		}
		list, err := client.GetBaseClient().Pod.List("", metav1.ListOptions{
			FieldSelector: fmt.Sprintf("%s=%s", "spec.nodeName", node.GetName()),
		})
		if err != nil {
			log.Debug("获取容器列表失败,err:", err.Error())
			return nil, k8s2.ErrPodListGet
		}
		cluster.PodNum += node.Status.Allocatable.Pods().Value()
		nodeDetail.PodNum = node.Status.Allocatable.Pods().Value()
		cluster.ActivePodNum += int64(len(list.Items))
		nodeDetail.PodRun = int64(len(list.Items))
		if fmt.Sprintf("%s", node.Status.Conditions[len(node.Status.Conditions)-1].Type) == "Ready" {
			cluster.RunNodeNum++
			metric, err := k8s2.GetNodeMetrics(node.GetName())
			if err == nil {
				//log.Debugf("获取节点指标失败,节点: %s, err: %v\n", node.GetName(), err.Error())
				//return nil, k8s2.ErrMetricsGet
				resource := assemble.AssembleResource(*metric, node)
				nodeDetail.Resource = resource
				cpuNum, _ := strconv.ParseFloat(resource.CpuNum, 64)
				cpuUse, _ := strconv.ParseFloat(resource.CpuUse, 64)
				memSize, _ := strconv.ParseFloat(resource.MemSize, 64)
				memUse, _ := strconv.ParseFloat(resource.MemUse, 64)
				cpu += cpuNum
				mem += memSize
				useCpu += cpuUse
				useMem += memUse
			}
		}
		cluster.Nodes = append(cluster.Nodes, nodeDetail)
	}
	cluster.Resource.MemSize = strconv.FormatFloat(mem, 'f', -1, 64)
	cluster.Resource.MemUse = strconv.FormatFloat(useMem, 'f', -1, 64)
	cluster.Resource.CpuNum = strconv.FormatFloat(cpu, 'f', -1, 64)
	cluster.Resource.CpuUse = strconv.FormatFloat(useCpu, 'f', -1, 64)
	return &ClusterResponse{
		cluster,
	}, nil
}

// @Tags cluster
// @Summary 获取所有节点信息
// @Produce  json
// @Success 200 {array} protocol.Response{data=model.NodeDetail} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/nodes [post]
func (cs *clusterService) Nodes(ctx context.Context, req *NodesRequest) (*NodesResponse, error) {
	nodeList := make([]model.NodeDetail, 0)
	nodes, err := client.GetBaseClient().Node.List("")
	if err != nil {
		log.Debug("获取节点列表失败,err:", err.Error())
		return nil, k8s2.ErrNodeListGet
	}

	for _, node := range nodes.Items {
		nodeDetail := model.NodeDetail{
			Name:   node.GetName(),
			HostIp: node.Status.Addresses[0].Address,
		}
		list, err := client.GetBaseClient().Pod.List("", metav1.ListOptions{
			FieldSelector: fmt.Sprintf("%s=%s", "spec.nodeName", node.GetName()),
		})
		if err != nil {
			log.Debug("获取容器列表失败,err:", err.Error())
			return nil, k8s2.ErrPodListGet
		}
		nodeDetail.PodNum = node.Status.Allocatable.Pods().Value()
		nodeDetail.PodRun = int64(len(list.Items))
		if fmt.Sprintf("%s", node.Status.Conditions[len(node.Status.Conditions)-1].Type) == "Ready" {
			metric, err := k8s2.GetNodeMetrics(node.GetName())
			if err == nil {
				resource := assemble.AssembleResource(*metric, node)
				nodeDetail.Resource = resource
			}
		}
		nodeList = append(nodeList, nodeDetail)
	}
	return &NodesResponse{
		NodeList: nodeList,
	}, err
}

// @Tags cluster
// @Summary 获取节点信息
// @Produce  json
// @Accept  json
// @Param   name query string true "节点名"
// @Success 200 {object} protocol.Response{data=NodeResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/node [get]
func (cs *clusterService) Node(ctx context.Context, req *NodeRequest) (*NodeResponse, error) {
	node, err := client.GetBaseClient().Node.Get("", req.Name)
	if k8serror.IsNotFound(err) {
		log.Debug("节点不存在,err:", err.Error())
		return &NodeResponse{
			Exist: false,
		}, nil
	} else if err != nil {
		log.Debug("节点不存在,err:", err.Error())
		return nil, k8s2.ErrNodeGet
	}
	var (
		metrics []v1beta1.NodeMetrics
	)
	if fmt.Sprintf("%s", node.Status.Conditions[len(node.Status.Conditions)-1].Type) == "Ready" {
		metric, err := k8s2.GetNodeMetrics(node.GetName())
		if err == nil {
			//log.Debugf("获取节点指标失败,节点: %s, err: %v\n", node.GetName(), err.Error())
			//return nil, k8s2.ErrMetricsGet
			metrics = append(metrics, *metric)
		}
	}
	nodeDetail := assemble.AssembleNode(*node, nil, metrics, nil)

	return &NodeResponse{
		Exist: true,
		Node:  nodeDetail,
	}, nil
}

// @Tags cluster
// @Summary  获取命名空间列表
// @Produce  json
// @Accept  json
// @Param   params body NsRequest false "命名空间名"
// @Success 200 {object} protocol.Response{data=NsResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/ns [post]
func (cs *clusterService) Ns(ctx context.Context, req *NsRequest) (*NsResponse, error) {
	namespaces, err := client.GetBaseClient().Ns.List("")
	if err != nil {
		log.Debug("获取命名空间失败,err:", err.Error())
		return nil, k8s2.ErrNamespaceListGet
	}
	namespacesDetail := make([]model.NamespaceDetail, 0)
	for _, ns := range namespaces.Items {
		namespacesDetail = append(namespacesDetail, model.NamespaceDetail{
			Name:       ns.GetName(),
			CreateTime: ns.GetCreationTimestamp().String(),
			Status:     ns.Status.String(),
		})
	}
	return &NsResponse{
		Num:        len(namespacesDetail),
		Namespaces: namespacesDetail,
	}, nil
}

// @Tags cluster
// @Summary  获取命名空间信息
// @Produce  json
// @Accept  json
// @Param   namespace query string true "命名空间名"
// @Success 200 {object} protocol.Response{data=NameSpaceResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/namespace [get]
func (cs *clusterService) NameSpace(ctx context.Context, req *NameSpaceRequest) (*NameSpaceResponse, error) {
	namespaceDetail := model.NamespaceDetail{}
	namespace, err := client.GetBaseClient().Ns.Get(req.NameSpace, "")
	if k8serror.IsNotFound(err) {
		log.Debug("命名空间不存在,err:", err.Error())
		return &NameSpaceResponse{
			Exist: false,
		}, nil
	} else if err != nil {
		log.Debug("获取命名空间失败,err:", err.Error())
		return nil, k8s2.ErrNamespaceGet
	}
	deploy, err := client.GetBaseClient().Deployment.List(req.NameSpace)
	if err != nil {
		log.Debug("获取部署列表失败,err:", err.Error())
		return nil, k8s2.ErrDeploymentGet
	}
	statefulSet, err := client.GetBaseClient().Sf.List(req.NameSpace, metav1.ListOptions{})
	if err != nil {
		log.Debug("获取状态副本及失败,err:", err.Error())
		return nil, k8s2.ErrStatefulSetGet
	}
	services, err := client.GetBaseClient().Sv.List(req.NameSpace)
	if err != nil {
		log.Debug("获取服务列表失败,err:", err.Error())
		return nil, k8s2.ErrServiceGet
	}

	pods, err := client.GetBaseClient().Pod.List(req.NameSpace, metav1.ListOptions{})
	if err != nil {
		log.Debug("获取服务列表失败,err:", err.Error())
		return nil, k8s2.ErrPodsGet
	}
	namespaceDetail.Name = namespace.GetName()
	namespaceDetail.CreateTime = namespace.GetCreationTimestamp().String()
	namespaceDetail.Status = namespace.Status.String()
	namespaceDetail.DeploymentNum = len(deploy.Items)
	namespaceDetail.StatefulSetNum = len(statefulSet.Items)
	namespaceDetail.ServiceNum = len(services.Items)
	namespaceDetail.PodNum = len(pods.Items)
	return &NameSpaceResponse{Namespaces: namespaceDetail, Exist: true}, nil
}

// @Tags cluster
// @Summary  获取pod信息
// @Produce  json
// @Accept  json
// @Param   params body PodInfoRequest true "命名空间名 名字"
// @Success 200 {object} protocol.Response{data=PodInfoResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/pod [post]
func (cs *clusterService) PodInfo(ctx context.Context, req *PodInfoRequest) (*PodInfoResponse, error) {
	podInfo, err := client.GetBaseClient().Pod.Get(req.NameSpace, req.PodName)
	if k8serror.IsNotFound(err) {
		return &PodInfoResponse{
			Exist: false,
		}, nil
	} else if statusError, isStatus := err.(*k8serror.StatusError); isStatus {
		return nil, errors.New(statusError.ErrStatus.Message)
	} else if err != nil {
		return nil, errors.New("内部错误")
	}
	var podMetric *v1beta1.PodMetrics
	podMetric, err = k8s2.GetPodMetrics(req.NameSpace, req.PodName)
	if err != nil {
		log.Debug("PodInfo: 获取pod指标信息失败。err:", err.Error())
	}
	podDetail := assemble.AssemblePodSummary(*podInfo, podMetric)
	return &PodInfoResponse{
		Pod:   podDetail,
		Exist: true,
	}, nil
}

////@Tags cluster
////@Summary  获取pod日志
////@Produce  json
////@Accept  json
////@Param   namespace query string true "命名空间名 名字"
////@Param   podName query string true "Pod名字"
////@Param   container query string false "容器名"
////@Param   follow query boolean false "是否开启实时日志"
////@Param   previous query boolean false "显示历史日志"
////@Success 200 {object} protocol.Response{data=PodLogResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
////@Router /cluster/v1/podLog [get]
//func (cs *clusterService) PodLog(ctx context.Context, req *PodLogRequest) (*PodLogResponse, error) {
//	pod, err := client.GetBaseClient().Pod.Get(req.NameSpace, req.PodName)
//	if err != nil {
//		log.Error("获取pod失败,err:", err.Error())
//		return nil, k8s2.ErrPodGet
//	}
//	if req.Container == "" {
//		req.Container = pod.Spec.Containers[0].Name
//	}
//	var refLineNum = 0
//	var offsetFrom = 2000000000
//	var offsetTo = 2000000100
//
//	refTimestamp := NewestTimestamp
//	logSelector := DefaultSelection
//
//	logSelector = &Selection{
//		ReferencePoint: LogLineId{
//			LogTimestamp: LogTimestamp(refTimestamp),
//			LineNum:      refLineNum,
//		},
//		OffsetFrom:      offsetFrom,
//		OffsetTo:        offsetTo,
//		LogFilePosition: "end",
//	}
//
//	logOptions := &v1.PodLogOptions{
//		Container:  req.Container,
//		Follow:     req.Follow,
//		Previous:   req.Previous,
//		Timestamps: true,
//	}
//
//	if logSelector.LogFilePosition == Beginning {
//		logOptions.LimitBytes = &byteReadLimit
//	} else {
//		logOptions.TailLines = &lineReadLimit
//	}
//	readCloser, err := inital.GetGlobal().GetClientSet().CoreV1().RESTClient().Get().Namespace(req.NameSpace).
//		Name(req.PodName).Resource("pods").
//		SubResource("log").VersionedParams(logOptions, scheme.ParameterCodec).
//		Stream()
//	//podLog, err := client.GetBaseClient().Pod.Log(req.NameSpace, req.PodName)
//	if err != nil {
//		log.Error("获取日志流失败,err:", err.Error())
//		return nil, k8s2.ErrLogGet
//	}
//	defer func() {
//		_ = readCloser.Close()
//	}()
//	//if req.Follow {
//	//	bufReader := bufio.NewReaderSize(readCloser, 256)
//	//	for {
//	//		line, _, err := bufReader.ReadLine()
//	//		// line = []byte(fmt.Sprintf("%s", string(line)))
//	//		line = utils.ToValidUTF8(line, []byte(""))
//	//		if err != nil {
//	//			if err == io.EOF {
//	//				_, err = .Write(line)
//	//			}
//	//			return err
//	//		}
//	//		// line = append(line, []byte("\r\n")...)
//	//		// line = append(bytes.Trim(line, " "), []byte("\r\n")...)
//	//		_, err = writer.Write(line)
//	//		if err != nil {
//	//			return err
//	//		}
//	//	}
//	//}
//
//	result, err := ioutil.ReadAll(readCloser)
//	if err != nil {
//		log.Error("日志流读取失败,err:", err.Error())
//		return nil, k8s2.ErrLogGet
//	}
//	rawLogs := string(result)
//	parsedLines := ToLogLines(rawLogs)
//	logLines, fromDate, toDate, logSelection, lastPage := parsedLines.SelectLogs(logSelector)
//	readLimitReached := isReadLimitReached(int64(len(rawLogs)), int64(len(parsedLines)), logSelector.LogFilePosition)
//	truncated := readLimitReached && lastPage
//	info := LogInfo{
//		PodName:       req.PodName,
//		ContainerName: req.Container,
//		FromDate:      fromDate,
//		ToDate:        toDate,
//		Truncated:     truncated,
//	}
//
//	return &PodLogResponse{
//		Log: &LogDetails{
//			Info:      info,
//			Selection: logSelection,
//			LogLines:  logLines,
//		},
//	}, err
//}

func isReadLimitReached(bytesLoaded int64, linesLoaded int64, logFilePosition string) bool {
	return (logFilePosition == Beginning && bytesLoaded >= byteReadLimit) ||
		(logFilePosition == End && linesLoaded >= lineReadLimit)
}

// @Tags cluster
// @Summary  获取pod列表
// @Produce  json
// @Accept  json
// @Param   params body PodsRequest false "命名空间名"
// @Success 200 {object} protocol.Response{data=PodsResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/pods [post]
func (cs *clusterService) Pods(ctx context.Context, req *PodsRequest) (*PodsResponse, error) {
	var (
		pods *v1.PodList
		err  error
	)
	podDetail := make([]model.PodDetail, 0)
	if req.NameSpace != "" {
		pods, err = client.GetBaseClient().Pod.List(req.NameSpace, metav1.ListOptions{})

	} else if req.NodeName != "" {
		pods, err = client.GetBaseClient().Pod.List("", metav1.ListOptions{
			FieldSelector: fmt.Sprintf("%s=%s", "spec.nodeName", req.NodeName),
		})
	}
	if k8serror.IsNotFound(err) {
		return &PodsResponse{}, nil
	} else if err != nil {
		return nil, k8s2.ErrPodListGet
	}
	var wg sync.WaitGroup
	for _, pod := range pods.Items {
		wg.Add(1)
		go func(pod v1.Pod) {
			var podMetric *v1beta1.PodMetrics
			podMetric, err = k8s2.GetPodMetrics(pod.GetNamespace(), pod.GetName())
			if err != nil {
				log.Debug("PodInfo: 获取pod指标信息失败。err:", err.Error())
			}
			podDetail = append(podDetail, assemble.AssemblePodSummary(pod, podMetric))
			wg.Done()
		}(pod)
	}
	wg.Wait()
	return &PodsResponse{
		Pods: podDetail,
	}, nil
}

// @Tags cluster
// @Summary  获取deployment
// @Produce  json
// @Accept  json
// @Param  params body ResourceRequest true "参数列表"
// @Success 200 {object} protocol.Response{data=DeploymentsResponse} "{"errno":0,"errmsg":"","data":{"items":[]},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/deployment [post]
func (cs *clusterService) Deployment(ctx context.Context, req *ResourceRequest) (*DeploymentsResponse, error) {
	var (
		err   error
		items []model.DeploymentDetail
	)
	if req.Name == "" {
		deploy, err := client.GetBaseClient().Deployment.List(req.NameSpace)
		err = k8s2.PrintErr(err)
		if err == nil {
			if len(deploy.Items) > 0 {
				var selectorKey, selectorVal string
				for _, deploy := range deploy.Items {
					for key, val := range deploy.Spec.Selector.MatchLabels {
						selectorKey = key
						selectorVal = val
					}
					pods, err := client.GetBaseClient().Pod.List(req.NameSpace, metav1.ListOptions{
						LabelSelector: fmt.Sprintf("%s=%s", selectorKey, selectorVal)})
					if err != nil {
						return nil, k8s2.ErrProjectPodsList
					}
					item := assemble.AssembleDeploymentSimple(deploy)
					var metrics []v1beta1.PodMetrics
					metrics, err = k8s2.GetPodListMetrics(req.NameSpace, metav1.ListOptions{})
					if err != nil {
						log.Debugf("获取指标失败")
					}
					podDetail := assemble.AssemblePod("", pods.Items, metrics)
					item.PodDetail = podDetail
					items = append(items, item)
				}
			}
		}
	} else {
		deploy, err := client.GetBaseClient().Deployment.Get(req.NameSpace, req.Name)
		err = k8s2.PrintErr(err)
		if err == nil {
			var selectorKey, selectorVal string
			for key, val := range deploy.Spec.Selector.MatchLabels {
				selectorKey = key
				selectorVal = val
			}
			pods, err := client.GetBaseClient().Pod.List(req.NameSpace, metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%s=%s", selectorKey, selectorVal)})
			if err != nil {
				return nil, k8s2.ErrProjectPodsList
			}
			item := assemble.AssembleDeploymentSimple(*deploy)
			var metrics []v1beta1.PodMetrics
			metrics, err = k8s2.GetPodListMetrics(req.NameSpace, metav1.ListOptions{})
			if err != nil {
				log.Debugf("获取指标失败")
			}
			podDetail := assemble.AssemblePod("", pods.Items, metrics)
			item.PodDetail = podDetail
			items = append(items, item)
		}
	}
	return &DeploymentsResponse{
		items,
	}, err
}

// @Tags cluster
// @Summary  获取statefulSet
// @Produce  json
// @Accept  json
// @Param  params body ResourceRequest true "参数列表"
// @Success 200 {object} protocol.Response{data=StatefulSetsResponse} "{"errno":0,"errmsg":"","data":{"items":[]},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/statefulSet [post]
func (cs *clusterService) StatefulSet(ctx context.Context, req *ResourceRequest) (*StatefulSetsResponse, error) {
	var (
		err   error
		items []model.StatefulSetDetail
	)
	if req.Name == "" {
		deploy, err := client.GetBaseClient().Sf.List(req.NameSpace, metav1.ListOptions{})
		err = k8s2.PrintErr(err)
		if err == nil {
			if len(deploy.Items) > 0 {
				var selectorKey, selectorVal string
				for _, deploy := range deploy.Items {
					for key, val := range deploy.Spec.Selector.MatchLabels {
						selectorKey = key
						selectorVal = val
					}
					pods, err := client.GetBaseClient().Pod.List(req.NameSpace, metav1.ListOptions{
						LabelSelector: fmt.Sprintf("%s=%s", selectorKey, selectorVal)})
					if err != nil {
						return nil, k8s2.ErrProjectPodsList
					}
					item := assemble.AssembleStatefulSetSimple(deploy)
					var metrics []v1beta1.PodMetrics
					metrics, err = k8s2.GetPodListMetrics(req.NameSpace, metav1.ListOptions{})
					if err != nil {
						log.Debugf("获取指标失败")
					}
					podDetail := assemble.AssemblePod("", pods.Items, metrics)
					item.PodDetail = podDetail
					items = append(items, item)
				}
			}

		}
	} else {
		deploy, err := client.GetBaseClient().Sf.Get(req.NameSpace, req.Name)
		err = k8s2.PrintErr(err)
		if err == nil {
			var selectorKey, selectorVal string
			for key, val := range deploy.Spec.Selector.MatchLabels {
				selectorKey = key
				selectorVal = val
			}
			pods, err := client.GetBaseClient().Pod.List(req.NameSpace, metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%s=%s", selectorKey, selectorVal)})
			if err != nil {
				return nil, k8s2.ErrProjectPodsList
			}
			item := assemble.AssembleStatefulSetSimple(*deploy)
			var metrics []v1beta1.PodMetrics
			metrics, err = k8s2.GetPodListMetrics(req.NameSpace, metav1.ListOptions{})
			if err != nil {
				log.Debugf("获取指标失败")
			}
			podDetail := assemble.AssemblePod("", pods.Items, metrics)
			item.PodDetail = podDetail
			items = append(items, item)
		}
	}
	return &StatefulSetsResponse{
		items,
	}, err
}

// @Tags cluster
// @Summary  获取service
// @Produce  json
// @Accept  json
// @Param  params body ResourceRequest true "参数列表"
// @Success 200 {object} protocol.Response{data=ServiceResponse} "{"errno":0,"errmsg":"","data":{"items":[]},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/service [post]
func (cs *clusterService) Services(ctx context.Context, req *ResourceRequest) (*ServiceResponse, error) {
	var (
		err    error
		deploy interface{}
		items  []model.ServiceDetail
	)
	if req.Name == "" {
		deploy, err = client.GetBaseClient().Sv.List(req.NameSpace)
		err = k8s2.PrintErr(err)
		if err == nil {
			svcs := deploy.(*v1.ServiceList).Items
			items = assemble.AssembleService(req.NameSpace, svcs)
		}
	} else {
		deploy, err = client.GetBaseClient().Sv.Get(req.NameSpace, req.Name)
		err = k8s2.PrintErr(err)
		if err == nil {
			dp := deploy.(*v1.Service)
			items = assemble.AssembleService(req.NameSpace, []v1.Service{*dp})
		}
	}
	return &ServiceResponse{
		items,
	}, err
}

func (cs *clusterService) GetDetailForRange(namespace v1.Namespace) model.NamespaceDetail {
	namespaceDetail := model.NamespaceDetail{
		Name:       namespace.GetName(),
		CreateTime: namespace.GetCreationTimestamp().String(),
	}
	name := namespace.GetName()
	//var wg tools.WaitGroupWrapper
	//wg.Wrap(func() {
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
	//})
	//wg.Wrap(func() {
	go func() {
		stats, err := client.GetBaseClient().Sf.List(name, metav1.ListOptions{})
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get statefulSet is err:%v", err.Error())
		} else if len(stats.Items) > 0 {
			namespaceDetail.StatefulSetList = assemble.AssembleStatefulSet(name, stats.Items)
		}
		wg.Done()
	}()
	//})
	//wg.Wrap(func() {
	go func() {
		svcs, err := client.GetBaseClient().Sv.List(name)
		if err != nil {
			log.Debugf("Method [GetDetailForRange] = > Get service is err:%v\n", err.Error())
		} else if len(svcs.Items) > 0 {
			namespaceDetail.ServiceList = assemble.AssembleService(name, svcs.Items)
		}
		wg.Done()
	}()
	//})
	wg.Wait()
	return namespaceDetail
}

// @Tags cluster
// @Summary  获取资源详细配置
// @Produce  json
// @Accept  json
// @Param   namespace query string true "命名空间名 名字"
// @Param   kind query string true "资源类型"
// @Param   name query string true "资源名"
// @Success 200 {object} protocol.Response{data=GetYamlResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/getYaml [get]
func (cs *clusterService) GetYaml(ctx context.Context, req *GetYamlRequest) (*GetYamlResponse, error) {

	var (
		err error
		obj *unstructured.Unstructured
	)
	switch req.Kind {
	case utils.DEPLOY_DEPLOYMENT:
		obj, err = client.GetBaseClient().Deployment.DynamicGet(req.Namespace, req.Name)
	case utils.DEPLOY_STATEFULSET:
		obj, err = client.GetBaseClient().Sf.DynamicGet(req.Namespace, req.Name)
	case utils.DEPLOY_Service:
		obj, err = client.GetBaseClient().Sv.DynamicGet(req.Namespace, req.Name)
	}

	err = k8s2.PrintErr(err)
	if err != nil {
		return nil, err
	}
	delete(obj.Object, "status")
	delete(obj.Object["metadata"].(map[string]interface{}), "creationTimestamp")
	delete(obj.Object["metadata"].(map[string]interface{}), "generation")
	delete(obj.Object["metadata"].(map[string]interface{}), "resourceVersion")
	delete(obj.Object["metadata"].(map[string]interface{}), "selfLink")
	delete(obj.Object["metadata"].(map[string]interface{}), "uid")
	delete(obj.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{}), "terminationGracePeriodSeconds")
	delete(obj.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["metadata"].(map[string]interface{}), "creationTimestamp")
	cont := obj.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})
	delete(cont, "terminationMessagePath")
	delete(cont, "terminationMessagePolicy")
	return &GetYamlResponse{Yaml: obj.Object}, nil
}

// @Tags cluster
// @Summary  获取事件
// @Produce  json
// @Accept  json
// @Param   kind query string true "0 node 1 Deployment 2 StatefulSet 3 Service 4 pod"
// @Param   name query string true "名"
// @Param   namespace query string false "命名空间名 除node外必填"
// @Success 200 {object} protocol.Response{data=EventResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/event [get]
func (cs *clusterService) Event(ctx context.Context, req *EventRequest) (*EventResponse, error) {
	event := k8s2.GetEvents(req.Namespace, req.Name, req.Kind)
	return &EventResponse{
		Event: event,
	}, nil
}

// @Tags cluster
// @Summary  版本号列表
// @Produce  json
// @Accept  json
// @Param   namespace query string true "命名空间名 名字"
// @Param   name query string true "资源名"
// @Param   label query string true "资源唯一标签"
// @Success 200 {object} protocol.Response{data=VersionResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/version [get]
func (cs *clusterService) VersionList(ctx context.Context, req *VersionRequest) (*VersionResponse, error) {
	rs, err := inital.GetGlobal().GetClientSet().AppsV1().ReplicaSets(req.Namespace).List(metav1.ListOptions{
		LabelSelector: req.Label,
	})
	if err != nil {
		log.Debugf("副本级获取失败 err:%v\n", err.Error())
		return nil, k8s2.ErrReplicaSetGet
	}
	var versionList []model.Versions
	if rs != nil && len(rs.Items) > 0 {
		for _, version := range rs.Items {
			for _, owner := range version.OwnerReferences {
				if owner.Name == req.Name {
					versionList = append(versionList, model.Versions{
						Version:     version.Annotations["deployment.kubernetes.io/revision"],
						VersionName: version.GetName(),
					})
					break
				}
			}
		}
	}
	return &VersionResponse{
		VersionList: versionList,
	}, nil
}
