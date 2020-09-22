package cluster

import (
	"context"
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	app "relaper.com/kubemanage/cache"
	"relaper.com/kubemanage/inital"
	"relaper.com/kubemanage/inital/client"
	k8s2 "relaper.com/kubemanage/k8s"
	"relaper.com/kubemanage/model"
	"relaper.com/kubemanage/pkg/assemble"
	"relaper.com/kubemanage/utils"
	errors2 "relaper.com/kubemanage/utils/errors"
	"sync"
)

type Service interface {
	Cluster(ctx context.Context, req *ClusterRequest) (*ClusterResponse, error)
	Nodes(ctx context.Context, req *NodesRequest) (*NodesResponse, error)
	Node(ctx context.Context, req *NodeRequest) (*NodeResponse, error)
	NameSpaces(ctx context.Context, req *NameSpacesRequest) (*NameSpacesResponse, error)
	PodInfo(ctx context.Context, req *PodInfoRequest) (*PodInfoResponse, error)
	PodLog(ctx context.Context, req *PodInfoRequest) (*PodLogResponse, error)
	Pods(ctx context.Context, req *PodsRequest) (*PodsResponse, error)
	Deployment(ctx context.Context, req *ResourceRequest) (*DeploymentsResponse, error)
	StatefulSet(ctx context.Context, req *ResourceRequest) (*StatefulSetsResponse, error)
	Services(ctx context.Context, req *ResourceRequest) (*ServiceResponse, error)
	GetYaml(ctx context.Context, req *GetYamlRequest) (*GetYamlResponse, error)
}

// NewService return a Service interface
func NewService() Service {
	return &clusterService{}
}

type clusterService struct{}

// @Summary 获取集群信息
// @Produce  json
// @Success 200 {object} protocol.Response{data=ClusterResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/detail [post]
func (cs *clusterService) Cluster(ctx context.Context, req *ClusterRequest) (*ClusterResponse, error) {
	clusterDetail := &model.Cluster{}
	val, exist, err := inital.GetGlobal().GetCache().Get(utils.CLUSTER_PREFIX_KEY)
	if err != nil {
		log.Debug("从缓存获取信息失败,err:", err.Error())
		clusterDetail = app.GetClusterData()
	} else {
		if exist {
			err = json.Unmarshal([]byte(val), clusterDetail)
			if err != nil {
				log.Debug("序列化失败,err:", err.Error())
			}
			clusterDetail = app.GetClusterData()
		} else {
			clusterDetail = app.GetClusterData()
			app.CacheCluster(clusterDetail)
		}
	}
	return &ClusterResponse{
		*clusterDetail,
	}, nil
}

// @Summary 获取所有节点信息
// @Produce  json
// @Success 200 {array} protocol.Response{data=model.NodeDetail} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/nodes [post]
func (cs *clusterService) Nodes(ctx context.Context, req *NodesRequest) (*NodesResponse, error) {
	nodes, err := client.GetBaseClient().Node.List("")
	if err != nil {
		err = k8s2.PrintErr(err)
		return nil, errors2.WithTipMessage(err, "获取节点列表信息失败")
	}
	nodeList := nodes.Items
	var (
		podsList []v1.Pod
	)
	nodeMetricsList, err := k8s2.GetNodeListMetrics()
	if err != nil {
		log.Debugf("Method [Nodes] = > Get NodeMetrics is err:%v", err.Error())
	}
	pods, err := client.GetBaseClient().Pod.List("")
	if err != nil {
		log.Debugf("Method [Nodes] = > Get pods is err:%v", err.Error())
	} else {
		podsList = pods.Items
	}
	podMetricsList, err := k8s2.GetPodListMetrics("")
	if err != nil {
		log.Debug("获取pod指标失败，err:", err.Error())
	}
	nodeDetail := assemble.AssembleNodes(nodeList, podsList, nodeMetricsList, podMetricsList)
	return &NodesResponse{
		NodeList: nodeDetail,
	}, nil
}

// @Summary 获取节点信息
// @Produce  json
// @Accept  json
// @Param   name query string true "节点名"
// @Success 200 {object} protocol.Response{data=NodeResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/node [get]
func (cs *clusterService) Node(ctx context.Context, req *NodeRequest) (*NodeResponse, error) {
	nodeInfo, err := client.GetBaseClient().Node.Get("", req.Name)
	if k8serror.IsNotFound(err) {
		return &NodeResponse{
			Exist: false,
		}, nil
	} else if statusError, isStatus := err.(*k8serror.StatusError); isStatus {
		return nil, errors.New(statusError.ErrStatus.Message)
	} else if err != nil {
		return nil, errors.New("内部错误")
	}
	var (
		podsList       []v1.Pod
		metrics        []v1beta1.NodeMetrics
		podMetricsList []v1beta1.PodMetrics
	)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		nodeMetric, err := k8s2.GetNodeMetrics(req.Name)
		if err != nil {
			log.Debugf("Method [node] = > Get NodeMetrics is err:%v", err.Error())
		} else {
			metrics = append(metrics, *nodeMetric)
		}
		wg.Done()
	}()
	go func() {
		pods, err := client.GetBaseClient().Pod.List("")
		if err != nil {
			log.Debugf("Method [Node] = > Get pods is err:%v", err.Error())
		} else {
			podsList = pods.Items
		}
		wg.Done()
	}()
	go func() {
		podMetricsList, err = k8s2.GetPodListMetrics("")
		if err != nil {
			log.Debugf("Method [Node] = > GetPodListMetrics is err:%v", err.Error())
		}
		wg.Done()
	}()
	wg.Wait()
	nodeDetail := assemble.AssembleNodes([]v1.Node{*nodeInfo}, podsList, metrics, podMetricsList)
	return &NodeResponse{
		Exist: true,
		Node:  nodeDetail[0],
	}, nil
}

// @Summary  获取命名空间信息
// @Produce  json
// @Accept  json
// @Param   params body NameSpacesRequest false "命名空间名"
// @Success 200 {object} protocol.Response{data=NameSpacesResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/namespace [post]
func (cs *clusterService) NameSpaces(ctx context.Context, req *NameSpacesRequest) (*NameSpacesResponse, error) {
	resp := &NameSpacesResponse{}
	namespaceDetail := make([]model.NamespaceDetail, 0)
	var (
		results []interface{}
		result  string
		err     error
		exist   bool
	)
	if req.Namespace == "" {
		results, exist, err = inital.GetGlobal().GetCache().HVals(utils.NAMESPACE_PREFIX_KEY)
	} else {
		result, exist, err = inital.GetGlobal().GetCache().HGet(utils.NAMESPACE_PREFIX_KEY, utils.NAMESPACE_PREFIX_KEY+req.Namespace)
	}
	if err != nil {
		log.Debug("从缓存获取信息失败,err:", err.Error())
		resp.Namespaces = app.GetNamespaceDetail(req.Namespace)
	} else {
		if exist {
			if req.Namespace == "" {
				for _, v := range results {
					var ns model.NamespaceDetail
					err = json.Unmarshal(v.([]byte), &ns)
					if err != nil {
						log.Debug("json转换失败,err:", err.Error())
					} else {
						namespaceDetail = append(namespaceDetail, ns)
					}
				}
			} else {
				ns := model.NamespaceDetail{}
				err = json.Unmarshal([]byte(result), &ns)
				if err == nil {
					namespaceDetail = append(namespaceDetail, ns)
				} else {
					log.Debug("json转换失败,err:", err.Error())
					namespaceDetail = app.GetNamespaceDetail(req.Namespace)
				}
			}
		} else {
			namespaceDetail = app.GetNamespaceDetail(req.Namespace)
			if len(namespaceDetail) > 0 {
				app.CacheNamespace(namespaceDetail)
			}
		}
	}
	//namespaceDetail := make([]model.NamespaceDetail, 0)
	//if req.Name == "" {
	//	namespaces, err := cs.namespace.List("")
	//	if err != nil {
	//		return nil, errors2.WithTipMessage(err, "获取命名空间列表失败")
	//	}
	//	items := namespaces.(*v1.NamespaceList).Items
	//	wg := sync.WaitGroup{}
	//	for _, ns := range items {
	//		wg.Add(1)
	//		go func(ns v1.Namespace) {
	//			namespaceDetail = append(namespaceDetail, cs.GetDetailForRange(ns))
	//			wg.Done()
	//		}(ns)
	//		//namespaceDetail = append(namespaceDetail,cs.GetDetailForRange(ns))
	//	}
	//	wg.Wait()
	//	return &NameSpacesResponse{
	//		Exist:      true,
	//		Namespaces: namespaceDetail,
	//	}, nil
	//}
	//namespace, err := cs.namespace.Get(req.Name, "")
	//if k8serror.IsNotFound(err) {
	//	return &NameSpacesResponse{
	//		Exist: false,
	//	}, nil
	//} else if statusError, isStatus := err.(*k8serror.StatusError); isStatus {
	//	return nil, errors.New(statusError.ErrStatus.Message)
	//} else if err != nil {
	//	return nil, errors.New("内部错误")
	//}
	//ns := namespace.(*v1.Namespace)
	//namespaceDetail = append(namespaceDetail, cs.GetDetailForRange(*ns))
	return &NameSpacesResponse{
		Namespaces: namespaceDetail,
	}, nil
}

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
	podListMetric := make([]v1beta1.PodMetrics, 0)
	podMetric, err := k8s2.GetPodMetrics(req.NameSpace, req.PodName)
	if err != nil {
		log.Debug("PodInfo: 获取pod指标信息失败。err:", err.Error())
	} else {
		podListMetric = append(podListMetric, *podMetric)
	}
	podDetail := assemble.AssemblePod(podInfo.Spec.NodeName, []v1.Pod{*podInfo}, podListMetric)
	return &PodInfoResponse{
		Pod:   podDetail[0],
		Exist: true,
	}, nil
}

// @Summary  获取pod日志
// @Produce  json
// @Accept  json
// @Param   namespace query string true "命名空间名 名字"
// @Param   podName query string true "Pod名字"
// @Success 200 {object} protocol.Response{data=PodLogResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/podLog [get]
func (cs *clusterService) PodLog(ctx context.Context, req *PodInfoRequest) (*PodLogResponse, error) {
	podLog, err := client.GetBaseClient().Pod.Log(req.NameSpace, req.PodName)
	return &PodLogResponse{
		Log: podLog,
	}, err
}

// @Summary  获取pod列表
// @Produce  json
// @Accept  json
// @Param   params body PodsRequest false "命名空间名"
// @Success 200 {object} protocol.Response{data=PodsResponse} "{"errno":0,"errmsg":"","data":{},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/pods [post]
func (cs *clusterService) Pods(ctx context.Context, req *PodsRequest) (*PodsResponse, error) {
	pods, err := client.GetBaseClient().Pod.List(req.NameSpace)
	if k8serror.IsNotFound(err) {
		return &PodsResponse{}, nil
	} else if statusError, isStatus := err.(*k8serror.StatusError); isStatus {
		return nil, errors.New(statusError.ErrStatus.Message)
	} else if err != nil {
		return nil, errors.New("内部错误")
	}
	items := pods.Items
	podListMetric, err := k8s2.GetPodListMetrics(req.NameSpace)
	if err != nil {
		log.Debug("获取pod监控资源失败，err:", err.Error())
	}
	podDetail := assemble.AssemblePod("", items, podListMetric)
	return &PodsResponse{
		Pods: podDetail,
	}, nil
}

// @Summary  获取deployment
// @Produce  json
// @Accept  json
// @Param  params body ResourceRequest true "参数列表"
// @Success 200 {object} protocol.Response{data=DeploymentsResponse} "{"errno":0,"errmsg":"","data":{"items":[]},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/deployment [post]
func (cs *clusterService) Deployment(ctx context.Context, req *ResourceRequest) (*DeploymentsResponse, error) {
	var (
		err    error
		deploy interface{}
		items  []model.DeploymentDetail
	)
	if req.Name == "" {
		deploy, err = client.GetBaseClient().Deployment.List(req.NameSpace)
		err = k8s2.PrintErr(err)
		if err == nil {
			deploys := deploy.(*appsv1.DeploymentList).Items
			items = assemble.AssembleDeployment(req.NameSpace, deploys)
		}
	} else {
		deploy, err = client.GetBaseClient().Deployment.Get(req.NameSpace, req.Name)
		err = k8s2.PrintErr(err)
		if err == nil {
			dp := deploy.(*appsv1.Deployment)
			items = assemble.AssembleDeployment(req.NameSpace, []appsv1.Deployment{*dp})
		}
	}
	return &DeploymentsResponse{
		items,
	}, err
}

// @Summary  获取statefulSet
// @Produce  json
// @Accept  json
// @Param  params body ResourceRequest true "参数列表"
// @Success 200 {object} protocol.Response{data=StatefulSetsResponse} "{"errno":0,"errmsg":"","data":{"items":[]},"extr":{"inner_error":"","error_stack":""}}"
// @Router /cluster/v1/statefulSet [post]
func (cs *clusterService) StatefulSet(ctx context.Context, req *ResourceRequest) (*StatefulSetsResponse, error) {
	var (
		err    error
		deploy interface{}
		items  []model.StatefulSetDetail
	)
	if req.Name == "" {
		deploy, err = client.GetBaseClient().Sf.List(req.NameSpace)
		err = k8s2.PrintErr(err)
		if err == nil {
			stats := deploy.(*appsv1.StatefulSetList).Items
			items = assemble.AssembleStatefulSet(req.NameSpace, stats)
		}
	} else {
		deploy, err = client.GetBaseClient().Sf.Get(req.NameSpace, req.Name)
		err = k8s2.PrintErr(err)
		if err == nil {
			dp := deploy.(*appsv1.StatefulSet)
			items = assemble.AssembleStatefulSet(req.NameSpace, []appsv1.StatefulSet{*dp})
		}
	}
	return &StatefulSetsResponse{
		items,
	}, err
}

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
		stats, err := client.GetBaseClient().Sf.List(name)
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
			log.Debugf("Method [GetDetailForRange] = > Get service is err:%v", err.Error())
		} else if len(svcs.Items) > 0 {
			namespaceDetail.ServiceList = assemble.AssembleService(name, svcs.Items)
		}
		wg.Done()
	}()
	//})
	wg.Wait()
	return namespaceDetail
}

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
