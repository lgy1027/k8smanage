package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"relaper.com/kubemanage/example/exporter/kube"
	"sync"
)

// 指标结构体
type Metrics struct {
	metrics map[string]*prometheus.Desc
	mutex   sync.Mutex
}

/**
 * 函数：newGlobalMetric
 * 功能：创建指标描述符
 */
func newGlobalMetric(metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(metricName, docString, labels, nil)
}

/**
 * 工厂方法：NewMetrics
 * 功能：初始化指标信息，即Metrics结构体
 */
func NewMetrics() *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"pod_metric_cpu":        newGlobalMetric("pod_metric_cpu", "pod_metric_cpu", []string{"instance", "namespace", "name"}),
			"pod_metric_mem":        newGlobalMetric("pod_metric_mem", "pod_metric_mem", []string{"instance", "namespace", "name"}),
			"containers_metric_cpu": newGlobalMetric("containers_metric_cpu", "containers_metric_cpu", []string{"instance", "namespace", "podname", "name"}),
			"containers_metric_mem": newGlobalMetric("containers_metric_mem", "containers_metric_mem", []string{"instance", "namespace", "podname", "name"}),
		},
	}
}

/**
 * 接口：Describe
 * 功能：传递结构体中的指标描述符到channel
 */
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

/**
 * 接口：Collect
 * 功能：抓取最新的数据，传递给channel
 */
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // 加锁
	defer c.mutex.Unlock()

	//mockCounterMetricData, mockGaugeMetricData := c.GenerateMockData()
	podMetics := kube.GetPodMetrics()
	if podMetics != nil && len(podMetics) > 0 {
		total_cpu, total_mem := int64(0), int64(0)
		for _, pod := range podMetics {
			if pod.Containers != nil && len(pod.Containers) > 0 {
				for _, contain := range pod.Containers {
					total_cpu += contain.Usage.Cpu().MilliValue()
					total_mem += contain.Usage.Memory().MilliValue()
					ch <- prometheus.MustNewConstMetric(c.metrics["containers_metric_cpu"], prometheus.GaugeValue, float64(contain.Usage.Cpu().MilliValue()), pod.NodeName, pod.NameSpace, pod.Name, contain.Name)
					ch <- prometheus.MustNewConstMetric(c.metrics["containers_metric_mem"], prometheus.GaugeValue, float64(contain.Usage.Memory().MilliValue()), pod.NodeName, pod.NameSpace, pod.Name, contain.Name)
				}
				ch <- prometheus.MustNewConstMetric(c.metrics["pod_metric_cpu"], prometheus.GaugeValue, float64(total_cpu), pod.NodeName, pod.NameSpace, pod.Name)
				ch <- prometheus.MustNewConstMetric(c.metrics["pod_metric_mem"], prometheus.GaugeValue, float64(total_mem), pod.NodeName, pod.NameSpace, pod.Name)
			}
		}
	}

	//for host, currentValue := range mockCounterMetricData {
	//	ch <- prometheus.MustNewConstMetric(c.metrics["my_counter_metric"], prometheus.CounterValue, float64(currentValue), host)
	//}
	//for host, currentValue := range mockGaugeMetricData {
	//	ch <- prometheus.MustNewConstMetric(c.metrics["my_gauge_metric"], prometheus.GaugeValue, float64(currentValue), host)
	//}
}

/**
 * 函数：GenerateMockData
 * 功能：生成模拟数据
 */
func (c *Metrics) GenerateMockData() (mockCounterMetricData map[string]int, mockGaugeMetricData map[string]int) {
	mockCounterMetricData = map[string]int{
		"yahoo.com":  int(rand.Int31n(1000)),
		"google.com": int(rand.Int31n(1000)),
	}
	mockGaugeMetricData = map[string]int{
		"yahoo.com":  int(rand.Int31n(10)),
		"google.com": int(rand.Int31n(10)),
	}
	return
}
