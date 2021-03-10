package kube

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	//"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"os"
	"path/filepath"
)

var (
	clientset *kubernetes.Clientset
	metrics   *versioned.Clientset
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

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

func NodeList() ([]v1.Node, error) {
	labelSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"node-role.kubernetes.io/edge":  "",
			"node-role.kubernetes.io/agent": "",
		},
	}
	labelMap, err := metav1.LabelSelectorAsMap(labelSelector)

	//获取NODE
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelMap).String(),
	})
	if err != nil {
		return nil, err
	}
	return nodes.Items, err
}

func GetNode(name string) (*v1.Node, error) {
	//获取NODE
	node, err := clientset.CoreV1().Nodes().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return node, nil
}

type Pod struct {
	NameSpace  string
	Name       string
	NodeName   string
	Containers []v1beta1.ContainerMetrics
}

func GetPodMetrics() []Pod {
	nodes, err := NodeList()
	if err != nil {
		return nil
	}
	podMetricsList := make([]Pod, 0)
	for _, node := range nodes {
		podList, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
			FieldSelector: fmt.Sprintf("%s=%s", "spec.nodeName", node.GetName()),
		})
		if err == nil {
			if len(podList.Items) > 0 {
				for _, pod := range podList.Items {
					podMetrics, err := metrics.MetricsV1beta1().PodMetricses(pod.GetNamespace()).Get(pod.GetName(), metav1.GetOptions{})
					if err == nil {
						po := Pod{
							Name:       pod.GetName(),
							NameSpace:  pod.GetNamespace(),
							NodeName:   node.GetName(),
							Containers: podMetrics.Containers,
						}
						podMetricsList = append(podMetricsList, po)
					}
				}
			}
		}
	}
	return podMetricsList
}
