package main_test

import (
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"relaper.com/kubemanage/example/k8s"
	"relaper.com/kubemanage/example/models"
	"strconv"
	"testing"
)

var (
	clientset  *kubernetes.Clientset
	restClient *rest.RESTClient
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
	config.APIPath = "api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	restClient, err = rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

}

func Test_RestClient(t *testing.T) {
	result := &corev1.NamespaceList{}
	err := clientset.CoreV1().
		RESTClient().
		Get().
		Namespace("lgy").
		Resource("deployments").
		VersionedParams(&metav1.ListOptions{Limit: 500}, scheme.ParameterCodec).Do().Into(result)
	fmt.Println(err)
	fmt.Println(result)
}

func Test_Resource(t *testing.T) {
	nodes := k8s.GetNodes(clientset, "")
	clusterStatus := models.ClusterStatus{}
	for _, item := range nodes[0:4] {
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
	clusterStatus.PodNum = k8s.GetPodsNumber("", clientset)
	clusterStatus.Services = k8s.GetServiceNumber(clientset, "")
	clusterStatus.MemSize = clusterStatus.MemSize / 1024 / 1024 / 1024

	detail := models.CloudClusterDetail{}
	detail.ClusterCpu = clusterStatus.CpuNum
	detail.ClusterMem = clusterStatus.MemSize
	detail.ClusterNode = clusterStatus.Nodes
	detail.ClusterPods = clusterStatus.PodNum

	if detail.ClusterCpu > 0 && detail.ClusterMem > 0 {
		used := k8s.GetClusterUsed(clientset)
		detail.UsedMem = used.UsedMem
		detail.UsedCpu = used.UsedCpu
		floatCpu := (float64(detail.UsedCpu) / float64(detail.ClusterCpu)) * 100
		floatMem := (float64(detail.UsedMem) / float64(detail.ClusterMem)) * 100
		cp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", floatCpu), 64)
		mp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", floatMem), 64)
		detail.CpuUsePercent = cp
		detail.MemUsePercent = mp
		detail.MemFree = detail.ClusterMem - detail.UsedMem
		detail.CpuFree = detail.ClusterCpu - detail.UsedCpu
		detail.Services = used.Services
	}

	fmt.Println(detail)
}
