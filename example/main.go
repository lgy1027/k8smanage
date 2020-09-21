package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	clientset    *kubernetes.Clientset
	dynamiClient dynamic.Interface
	metrics      *versioned.Clientset
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func int32Ptr(i int32) *int32 { return &i }

func init() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Println(err)
	}
	dynamiClient, err = dynamic.NewForConfig(config)
	if err != nil {
		log.Println(err)
	}
	//dynamicClient, err := dynamic.NewForConfig(config)
	metrics, err = versioned.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func main() {

}

func Other() {
	for _, pod := range getPods() {
		if pod.GetNamespace() == "testsd" && pod.GetName() == "test-b8q6sr-7d7566cb49-l5wfg" {
			events := getEvent(pod.GetName(), pod.GetNamespace())
			fmt.Println("Event=============name:", pod.GetName(), "=====namespace:", pod.GetNamespace(), "============event:", events)
			log := getPodLogs(pod.GetNamespace(), pod.GetName())
			fmt.Println(log)
		}
	}

	getNode()
}

func getPodLogs(namespace, name string) string {
	podLogOpts := apiv1.PodLogOptions{}
	req := clientset.CoreV1().Pods(namespace).GetLogs(name, &podLogOpts)
	podLogs, err := req.Stream()
	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()

	return str
}

type EventData struct {
	// 事件事件
	EventTime string
	// 信息
	Messages string
	// 原因
	Reason string
	// 主机Ip
	Host string
}

// 替换时间的T和Z
func ReplaceTime(t string) string {
	t = strings.Replace(t, "T", " ", -1)
	t = strings.Replace(t, "Z", "", -1)
	t = strings.Replace(t, "+0800 CS", "", -1)
	ts := strings.Split(t, ".")
	return ts[0]
}

func getEvent(podName, namespace string) []EventData {
	opt := metav1.ListOptions{}
	opt.FieldSelector = "involvedObject.name=" + podName + ",involvedObject.namespace=" + namespace
	events, err := clientset.CoreV1().Events(namespace).List(opt)
	fmt.Println(err)
	return parseEvents(events.Items)
}

func parseEvents(items []apiv1.Event) []EventData {
	data := []EventData{}
	for _, v := range items {
		t := EventData{}
		t.Reason = v.Reason
		t.Messages = v.Message
		t.EventTime = v.FirstTimestamp.String()
		t.Host = v.Source.Host
		data = append(data, t)
	}
	return data
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func getPods() []apiv1.Pod {
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return pods.Items
}

func testNs() {
	// 通过实现 clientset 的 CoreV1Interface 接口列表中的 NamespacesGetter 接口方法 Namespaces 返回 NamespaceInterface
	// NamespaceInterface 接口拥有操作 Namespace 资源的方法，例如 Create、Update、Get、List 等方法
	name := "lgy-test"
	namespacesClient := clientset.CoreV1().Namespaces()
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Status: apiv1.NamespaceStatus{
			Phase: apiv1.NamespaceActive,
		},
	}

	// 创建一个新的 Namespaces
	fmt.Println("Creating Namespaces...")
	result, err := namespacesClient.Create(namespace)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Namespaces %s on %s\n", result.ObjectMeta.Name, result.ObjectMeta.CreationTimestamp)

	// 获取指定名称的 Namespaces 信息
	fmt.Println("Getting Namespaces...")
	result, err = namespacesClient.Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Name: %s, Status: %s, selfLink: %s, uid: %s\n",
		result.ObjectMeta.Name, result.Status.Phase, result.ObjectMeta.SelfLink, result.ObjectMeta.UID)

	// 删除指定名称的 Namespaces 信息
	fmt.Println("Deleting Namespaces...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := namespacesClient.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Printf("Deleted Namespaces %s\n", name)
}

func getNode() {
	//获取NODE
	fmt.Println("####### 获取node ######")
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	var storge int64
	//var storgeUser int64
	for _, node := range nodes.Items {
		mc, err := metrics.MetricsV1beta1().NodeMetricses().Get(node.GetName(), metav1.GetOptions{})
		if err != nil {
			break
		}
		storgeUser := mc.Usage.StorageEphemeral().String()
		fmt.Println(storgeUser)
		storge += node.Status.Capacity.StorageEphemeral().Value()
	}
	fmt.Println(storge)
	//fmt.Println(storgeUser)
	//node := nodes.Items[1]
	//fmt.Println(node)
	//desc, _ := node.Descriptor()
	//fmt.Println(string(desc))
	//name := nodes.Items[2].Name
	//for _, nds := range nodes.Items {
	//	fmt.Printf("NodeName: %s\n", nds.Name)
	//}
	//
	////获取 指定NODE 的详细信息
	//fmt.Println("\n ####### node详细信息 ######")
	//nodeRel, err := clientset.CoreV1().Nodes().Get(name, metav1.GetOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Name: %s \n", nodeRel.Name)
	//fmt.Printf("CreateTime: %s \n", nodeRel.CreationTimestamp)
	//fmt.Printf("NowTime: %s \n", nodeRel.Status.Conditions[0].LastHeartbeatTime)
	//fmt.Printf("kernelVersion: %s \n", nodeRel.Status.NodeInfo.KernelVersion)
	//fmt.Printf("SystemOs: %s \n", nodeRel.Status.NodeInfo.OSImage)
	//fmt.Printf("Cpu: %s \n", nodeRel.Status.Capacity.Cpu())
	//fmt.Printf("docker: %s \n", nodeRel.Status.NodeInfo.ContainerRuntimeVersion)
	//// fmt.Printf("Status: %s \n", nodeRel.Status.Conditions[len(nodes.Items[0].Status.Conditions)-1].Type)
	//fmt.Printf("Status: %s \n", nodeRel.Status.Conditions[len(nodeRel.Status.Conditions)-1].Type)
	//fmt.Printf("Mem: %s \n", nodeRel.Status.Allocatable.Memory().String())

}

func DeployMent() {
	deploymentsClient := clientset.AppsV1().Deployments("lgy")
	var r apiv1.ResourceRequirements
	//资源分配会遇到无法设置值的问题，故采用json反解析
	j := `{"limits": {"cpu":"200m", "memory": "250Mi"}, "requests": {"cpu":"100m", "memory": "100Mi"}}`
	err := json.Unmarshal([]byte(j), &r)
	var v []apiv1.Volume
	volum := `[{"name":"","hostPath":{"path":""}},{"name":"","hostPath":{"path":""}}]`
	err = json.Unmarshal([]byte(volum), &v)
	maxSurge := intstr.FromInt(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"test": "test",
			},
			Annotations: map[string]string{
				"desc": "nginx test for namespace lgy ",
			},
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{ // 由于replicas为3,则整个升级,pod个数在2-4个之间
					MaxSurge:       &maxSurge, // 滚动升级时会先启动1个pod
					MaxUnavailable: &maxSurge, // 滚动升级时允许的最大Unavailable的pod个数

				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					NodeSelector: map[string]string{
						"node": "three",
					},
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							Resources:       r,
							Command:         []string{},
							Args:            []string{},
							ImagePullPolicy: apiv1.PullIfNotPresent, // PullAlways PullPolicy = "Always"  PullNever PullPolicy = "Never"  PullIfNotPresent PullPolicy = "IfNotPresent"
							Env:             []apiv1.EnvVar{},       // 环境变量
							WorkingDir:      "",                     // 工作目录
							//VolumeMounts: []apiv1.VolumeMount{     //  挂载volumes中定义的磁盘
							//	apiv1.VolumeMount{
							//		Name: "",    // 外部挂载目录
							//		MountPath: "", // 容器内部目录
							//	},
							//},
						},
					},
					// 定义磁盘给上面volumeMounts挂载
					//Volumes: v,
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		if !k8serror.IsAlreadyExists(err) {
			panic(err)
		}
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	CreateServer()
	// Update Deployment
	prompt()
	fmt.Println("Updating deployment...")

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get("demo-deployment", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = int32Ptr(1)                           // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")

	// List Deployments
	prompt()
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// Delete Deployment
	prompt()
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete("demo-deployment", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")

	if err := clientset.CoreV1().Services("lgy").Delete("demo-service", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted service......")
}

func CreateServer() {
	trgetPort := intstr.FromInt(80)
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-service",
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "demo",
			},
			Type: apiv1.ServiceTypeNodePort,
			//ClusterIP:	"NodePort",
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Port:       8000,
					TargetPort: trgetPort,
					Protocol:   apiv1.ProtocolTCP,
				},
			},
		},
	}
	result, err := clientset.CoreV1().Services("lgy").Create(service)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func DynamicDeploy() {
	//deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deploymentYAML := "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  annotations:\n    deployment.kubernetes.io/revision: \"1\"\n    desc: 'nginx test for namespace lgy '\n  creationTimestamp: \"2020-09-03T02:37:42Z\"\n  generation: 1\n  labels:\n    test: test\n  name: demo-deployment\n  namespace: lgy\n  resourceVersion: \"7720167\"\n  selfLink: /apis/apps/v1/namespaces/lgy/deployments/demo-deployment\n  uid: 92a4222a-025f-49eb-b3b6-71e6d1eba368\nspec:\n  progressDeadlineSeconds: 600\n  replicas: 2\n  revisionHistoryLimit: 10\n  selector:\n    matchLabels:\n      app: demo\n  strategy:\n    rollingUpdate:\n      maxSurge: 1\n      maxUnavailable: 1\n    type: RollingUpdate\n  template:\n    metadata:\n      creationTimestamp: null\n      labels:\n        app: demo\n    spec:\n      containers:\n      - image: nginx:1.12\n        imagePullPolicy: IfNotPresent\n        name: web\n        ports:\n        - containerPort: 80\n          name: http\n          protocol: TCP\n        resources:\n          limits:\n            cpu: 200m\n            memory: 250Mi\n          requests:\n            cpu: 100m\n            memory: 100Mi\n        terminationMessagePath: /dev/termination-log\n        terminationMessagePolicy: File\n      dnsPolicy: ClusterFirst\n      nodeSelector:\n        node: three\n      restartPolicy: Always\n      schedulerName: default-scheduler\n      securityContext: {}\n      terminationGracePeriodSeconds: 30\nstatus:\n  availableReplicas: 2\n  conditions:\n  - lastTransitionTime: \"2020-09-03T02:37:43Z\"\n    lastUpdateTime: \"2020-09-03T02:37:43Z\"\n    message: Deployment has minimum availability.\n    reason: MinimumReplicasAvailable\n    status: \"True\"\n    type: Available\n  - lastTransitionTime: \"2020-09-03T02:37:42Z\"\n    lastUpdateTime: \"2020-09-03T02:37:43Z\"\n    message: ReplicaSet \"demo-deployment-dbd8ff5c4\" has successfully progressed.\n    reason: NewReplicaSetAvailable\n    status: \"True\"\n    type: Progressing\n  observedGeneration: 1\n  readyReplicas: 2\n  replicas: 2\n  updatedReplicas: 2"
	obj := &unstructured.Unstructured{}

	// decode YAML into unstructured.Unstructured
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := dec.Decode([]byte(deploymentYAML), nil, obj)

	if err != nil {
		panic(err)
	}
	// Get the common metadata, and show GVK
	fmt.Println(obj.GetName(), gvk.String())

	// encode back to JSON
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	enc.Encode(obj)

}

func Otherer() {
	//mc, err := metricsv.NewForConfig(config)
	//nodeMetrics, err := mc.MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{})
	////podInfo,err := mc.MetricsV1beta1().PodMetricses("lgy").Get("nginx-v1", metav1.GetOptions{})
	////mc.MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{})
	//nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	//for _, info := range nodeMetrics.Items {
	//	for _, nodeRel := range nodes.Items {
	//		if nodeRel.GetName() == info.GetName() {
	//			fmt.Println("=========================================")
	//			fmt.Println("node name:", nodeRel.GetName())
	//			//fmt.Println(nodeRel.Status.Allocatable.Memory().Value())
	//			//fmt.Println(nodeRel.Status.Capacity.Memory().Value())
	//			//memNum := nodeRel.Status.Capacity.Memory().Value()
	//			//memNumValue := decimal.NewFromInt(memNum)
	//			//memNumValue = memNumValue.Div(decimal.NewFromInt(1024)).Div(decimal.NewFromInt(1024)).DivRound(decimal.NewFromInt(1024),2)
	//			//fmt.Println(memNumValue.Float64())
	//			//fmt.Println(info.Usage.Memory().Value())
	//			//memNum = info.Usage.Memory().Value()
	//			//umemNumValue := decimal.NewFromInt(memNum)
	//			//umemNumValue = umemNumValue.Div(decimal.NewFromInt(1024)).Div(decimal.NewFromInt(1024)).DivRound(decimal.NewFromInt(1024),2)
	//			//fmt.Println(umemNumValue.Float64())
	//			//fmt.Println(umemNumValue.DivRound(memNumValue,2).Float64())
	//			//fmt.Println("=========================================")
	//			fmt.Println(nodeRel.Status.Allocatable.Cpu().MilliValue())
	//			fmt.Println(nodeRel.Status.Allocatable.Cpu().Value())
	//			fmt.Println(nodeRel.Status.Capacity.Cpu().MilliValue())
	//			fmt.Println(nodeRel.Status.Capacity.Cpu().Value())
	//			fmt.Println(info.Usage.Cpu().MilliValue())
	//			fmt.Println(info.Usage.Cpu().Value())
	//			fmt.Println()
	//			fmt.Println()
	//			fmt.Println()
	//			break
	//		}
	//	}
	//}

	//podsMetrics, err := mc.MetricsV1beta1().PodMetricses("istio-system").Get("jaeger-operator-bdbb4954b-9zmnt", metav1.GetOptions{})
	//fmt.Println(err)
	//fmt.Println(podsMetrics.Containers[0].Usage.Memory().Value() / 1024 / 1024)
	//fmt.Println(podsMetrics.Containers[0].Usage.Cpu().MilliValue())
	//fmt.Println()

	//mc.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(metav1.ListOptions{})
	//mc.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).Get("your pod name", metav1.GetOptions{})

	//DeployMent()
	//DynamicDeploy()
	// 通过实现 clientset 的 CoreV1Interface 接口列表中的 PodsGetter 接口方法 Pods(namespace string)返回 PodInterface
	// PodInterface 接口拥有操作 Pod 资源的方法，例如 Create、Update、Get、List 等方法
	// 注意：Pods() 方法中 namespace 不指定则获取 Cluster 所有 Pod 列表

	//nodes,err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	//fmt.Println(err)
	//fmt.Println(nodes)
	//pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Printf("There are %d pods in the k8s cluster\n", len(pods.Items))

	// 获取指定 namespace 中的 Pod 列表信息
	//namespace := "kubesphere-monitoring-system"
	//pods, err = clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("\nThere are %d pods in namespaces %s\n", len(pods.Items), namespace)
	//for _, pod := range pods.Items {
	//	fmt.Printf("Name: %s, Status: %s, CreateTime: %s\n", pod.ObjectMeta.Name, pod.Status.Phase, pod.ObjectMeta.CreationTimestamp)
	//}
	//time.Sleep(10 * time.Second)

	//for {
	//	// 通过实现 clientset 的 CoreV1Interface 接口列表中的 PodsGetter 接口方法 Pods(namespace string)返回 PodInterface
	//	// PodInterface 接口拥有操作 Pod 资源的方法，例如 Create、Update、Get、List 等方法
	//	// 注意：Pods() 方法中 namespace 不指定则获取 Cluster 所有 Pod 列表
	//	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//	fmt.Printf("There are %d pods in the k8s cluster\n", len(pods.Items))
	//	var (
	//		node1 int
	//		node2 int
	//		node3 int
	//		master int
	//		total int
	//		failed int
	//	)
	//	for _, pod := range pods.Items {
	//		if pod.Status.Phase != apiv1.PodRunning{
	//			fmt.Println("node:",pod.Spec.NodeName," podName:",pod.Name, " podStatus:",pod.Status.Phase)
	//			failed++
	//			continue
	//		}
	//		switch pod.Spec.NodeName {
	//		case "node1":
	//			node1++
	//		case "node2":
	//			node2++
	//		case "node3":
	//			node3++
	//		case "master":
	//			master++
	//		}
	//		total++
	//	}
	//	fmt.Printf("master: %d node1: %d node2: %d node3: %d total: %d failed: %d",master,node1,node2,node3,total,failed)
	//
	//	// 获取指定 namespace 中的 Pod 列表信息
	//	//namespce := "kubesphere-monitoring-system"
	//	//pods, err = clientset.CoreV1().Pods(namespce).List(metav1.ListOptions{})
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//	//fmt.Printf("\nThere are %d pods in namespaces %s\n", len(pods.Items), namespce)
	//	//for _, pod := range pods.Items {
	//	//	fmt.Printf("Name: %s, Status: %s, CreateTime: %s\n", pod.ObjectMeta.Name, pod.Status.Phase, pod.ObjectMeta.CreationTimestamp)
	//	//}
	//	//
	//	//// 获取所有的 Namespaces 列表信息
	//	//ns, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//	//nss := ns.Items
	//	//fmt.Printf("\nThere are %d namespaces in cluster\n", len(nss))
	//	//for _, ns := range nss {
	//	//	fmt.Printf("Name: %s, Status: %s, CreateTime: %s\n", ns.ObjectMeta.Name, ns.Status.Phase, ns.CreationTimestamp)
	//	//}
	//
	//	//time.Sleep(10 * time.Second)
	//	break
	//}

	//testNs()
	//getNode()
}
