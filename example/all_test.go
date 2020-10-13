package main_test

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"relaper.com/kubemanage/example/k8s"
	"relaper.com/kubemanage/example/models"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	clientset  *kubernetes.Clientset
	restClient *rest.RESTClient
)

func Test_Arr(t *testing.T) {
	// 切片连接 查询重复值
	arr1 := []int{1, 2, 3, 4, 5, 6}
	arr2 := []int{5, 6, 7, 8, 9, 0}
	fmt.Println(ArrayIntersection(arr1, arr2))

	// 寻找最长含有不重复字符的字串长度大小
	fmt.Println(lengthOfNonRepeatingSubStr("hello 世界!"))

	head := &Node{
		value: 1,
		next: &Node{
			value: 2,
			next: &Node{
				value: 3,
				next: &Node{
					value: 4,
					next:  nil,
				},
			},
		},
	}
	//head = reverse(head)
	//printNode(head)

	fmt.Println(hasCycle(head))

	add("123000", "1230001")
}

func singleNumber(arr []int) int {
	count := len(arr)
	ones := 0
	twos := 0
	xthrees := 0
	for i := 0; i < count; i++ {
		twos |= ones & arr[i]
		ones ^= arr[i]
		xthrees = ^(ones & twos)
		ones &= xthrees
		twos &= xthrees
	}
	return ones
}

func stringReverse(str string) string {
	reverse := []rune(str)
	strLen := len(str)
	for i, j := 0, strLen-1; i < j; i, j = i+1, j-1 {
		reverse[i], reverse[j] = reverse[j], reverse[i]
	}
	return string(reverse)
}

func add(str1 string, str2 string) string {

	if len(str1) < len(str2) {
		str1 = strings.Repeat("0", len(str2)-len(str1)) + str1
	} else if len(str1) > len(str2) {
		str2 = strings.Repeat("0", len(str1)-len(str2)) + str2
	}
	str1 = stringReverse(str1)
	str2 = stringReverse(str2)

	count := len(str1)
	nums := make([]byte, count)
	carry := false

	for i := 0; i < count; i++ {
		sum := str1[i] - '0' + str2[i] - '0'
		if carry {
			sum++
		}
		if sum > 9 {
			sum = sum - 10
			carry = true
		} else {
			carry = false
		}
		nums[i] = sum + '0'
	}

	result := stringReverse(string(nums))
	if carry {
		result = "1" + result
	}
	return result

}

type Node struct {
	value int
	next  *Node
}

// 链表反转
func reverse(head *Node) *Node {

	var pre *Node = nil
	for head != nil {
		temp := head.next
		head.next = pre
		pre = head
		head = temp
	}
	return pre
}

func printNode(head *Node) {
	for head != nil {
		fmt.Println(head.value)
		head = head.next
	}
}

// 单链表是否存在环
func hasCycle(head *Node) bool {
	fast := head
	slow := head
	for fast != nil && fast.next != nil {
		slow = slow.next
		fast = fast.next.next
		if fast == slow {
			return true
		}
	}
	return false
}

func lengthOfNonRepeatingSubStr(s string) int {
	lastOccurred := make(map[rune]int)
	start := 0
	maxLength := 0
	for i, ch := range []int32(s) {

		if lastI, ok := lastOccurred[ch]; ok && lastI >= start {
			start = lastOccurred[ch] + 1
		}
		if i-start+1 > maxLength {
			maxLength = i - start + 1
		}
		lastOccurred[ch] = i
	}
	return maxLength
}

func ArrayIntersection(arr []int, arr1 []int) []int {

	var intersection []int
	arr = append(arr, arr1...)
	sameElem := make(map[int]int)

	for _, v := range arr {
		if _, ok := sameElem[v]; ok {
			intersection = append(intersection, v)
		} else {
			sameElem[v] = 1
		}
	}
	return intersection
}

func Test_Harbor(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   200 * time.Second,
				KeepAlive: 200 * time.Second,
			}).Dial,
			DisableKeepAlives: false,
			//MaxIdleConns:          MaxIdleConns,
			//MaxIdleConnsPerHost:   MaxIdleConnsPerHost,
			//IdleConnTimeout:       time.Duration(IdleConnTimeout) * time.Second,
			//ResponseHeaderTimeout: time.Duration(ResponseHeaderTimeout) * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	url := "https://101.37.169.150/api/v2.0/projects/test/repositories/nginx/artifacts?page=1&page_size=10&with_tag=true&with_label=false&with_scan_overview=false&with_signature=false&with_immutable_status=false"
	resp, err := client.Get(url)
	fmt.Println(err)
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(err)
	fmt.Println(string(data))
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func inits() {
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
