package k8s

import (
	log "github.com/cihub/seelog"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"relaper.com/kubemanage/example/models"
)

// 获取pod数据
//fmt.Println(p.Status.HostIP)
//{"metadata":{"name":"zhaoyun1-rc-28fp6","generateName":"zhaoyun1-rc-","namespace":"default","selfLink":"/api/v1/namespaces/default/pods/zhaoyun1-rc-28fp6","uid":"29676a19-dbbd-11e7-a7e2-0894ef37b2d2","resourceVersion":"287211","creationTimestamp":"2017-12-08T02:11:54Z","labels":{"app":"www-gg-com","max-scale":"3","min-scale":"3"},"annotations":{"kubernetes.io/created-by":"{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicationController\",\"namespace\":\"default\",\"name\":\"zhaoyun1-rc\",\"uid\":\"2966d58b-dbbd-11e7-a7e2-0894ef37b2d2\",\"apiVersion\":\"v1\",\"resourceVersion\":\"287184\"}}\n"},"ownerReferences":[{"apiVersion":"v1","kind":"ReplicationController","name":"zhaoyun1-rc","uid":"2966d58b-dbbd-11e7-a7e2-0894ef37b2d2","controller":true,"blockOwnerDeletion":true}]},"spec":{"containers":[{"name":"zhaoyun1","image":"nginx:1.11","ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{"limits":{"cpu":"1","memory":"0"},"requests":{"cpu":"1","memory":"0"}},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","nodeName":"10.16.55.103","securityContext":{},"schedulerName":"default-scheduler"},"status":{"phase":"Running","conditions":[{"type":"Initialized","status":"True","lastProbeTime":null,"lastTransitionTime":"2017-12-08T02:10:19Z"},{"type":"Ready","status":"True","lastProbeTime":null,"lastTransitionTime":"2017-12-08T02:10:23Z"},{"type":"PodScheduled","status":"True","lastProbeTime":null,"lastTransitionTime":"2017-12-08T02:11:54Z"}],"hostIP":"10.16.55.103","podIP":"172.16.8.15","startTime":"2017-12-08T02:10:19Z","containerStatuses":[{"name":"zhaoyun1","state":{"running":{"startedAt":"2017-12-08T02:10:22Z"}},"lastState":{},"ready":true,"restartCount":0,"image":"nginx:1.11","imageID":"docker-pullable://nginx@sha256:e6693c20186f837fc393390135d8a598a96a833917917789d63766cab6c59582","containerID":"docker://c0a1cae85d6146d415996252750add084373c0f0e90c68fb129e3aa440262645"}],"qosClass":"Burstable"}}
func GetPods(namespace string, clientset *kubernetes.Clientset) []v1.Pod {
	opt := metav1.ListOptions{}
	pods, err := clientset.CoreV1().Pods(namespace).List(opt)
	if err != nil {
		log.Error("获取Pods错误", err.Error())
		return make([]v1.Pod, 0)
	}
	return pods.Items
}

// 获取资源使用情况,cpu，内存
func GetClusterUsed(clientset *kubernetes.Clientset) models.ClusterResources {
	clusterResouces := models.ClusterResources{}
	resources := GetPods("", clientset)
	var cpu int64
	var memory int64
	for _, item := range resources {
		containers := item.Spec.Containers
		for _, container := range containers {
			cpu += container.Resources.Limits.Cpu().Value()
			memory += container.Resources.Limits.Memory().Value()
		}
	}
	clusterResouces.Services = GetServiceNumber(clientset, "")
	clusterResouces.UsedCpu = cpu
	clusterResouces.UsedMem = memory / 1024 / 1024 / 1024
	return clusterResouces
}

// 获取nodes
func GetNodes(clientset *kubernetes.Clientset, labels string) []v1.Node {
	opt := metav1.ListOptions{}
	if labels != "" {
		opt.LabelSelector = labels
	}
	nodes, err := clientset.CoreV1().Nodes().List(opt)
	if err != nil {
		log.Error("获取Nodes错误", err.Error())
		return make([]v1.Node, 0)
	}
	return nodes.Items
}

// 获取pods数量
func GetPodsNumber(namespace string, clientset *kubernetes.Clientset) int {
	opt := metav1.ListOptions{}

	pods, err := clientset.CoreV1().Pods(namespace).List(opt)
	if err != nil {
		log.Error("获取k8s Pods失败", err.Error())
		return 0
	}
	return len(pods.Items)
}

// 获取某个集群服务的数量
func GetServiceNumber(clientset *kubernetes.Clientset, namespace string) int {
	data, _ := GetServices(clientset, namespace)
	return len(data)
}

// 获取某个集群的服务信息
func GetServices(clientset *kubernetes.Clientset, namespace string) ([]v1.Service, error) {
	opt := metav1.ListOptions{}
	data, err := clientset.CoreV1().Services(namespace).List(opt)
	if err != nil {
		log.Error("获取service失败啦", err)
		return make([]v1.Service, 0), err
	}
	return data.Items, err
}

func GetNodeFromCluster(clientset *kubernetes.Clientset) models.ClusterStatus {
	nodes := GetNodes(clientset, "")
	clusterStatus := models.ClusterStatus{}
	for _, item := range nodes {
		if clusterStatus.MemSize == 0 && clusterStatus.CpuNum == 0 {
			clusterStatus.CpuNum = item.Status.Capacity.Cpu().Value()
			clusterStatus.MemSize = item.Status.Capacity.Memory().Value()
			clusterStatus.Nodes = 1
		} else {
			clusterStatus.CpuNum = clusterStatus.CpuNum + item.Status.Capacity.Cpu().Value()
			clusterStatus.MemSize = clusterStatus.MemSize + item.Status.Capacity.Memory().Value()
			clusterStatus.Nodes = clusterStatus.Nodes + 1
		}
	}
	clusterStatus.PodNum = GetPodsNumber("", clientset)
	clusterStatus.Services = GetServiceNumber(clientset, "")
	clusterStatus.MemSize = clusterStatus.MemSize / 1024 / 1024 / 1024
	return clusterStatus
}
