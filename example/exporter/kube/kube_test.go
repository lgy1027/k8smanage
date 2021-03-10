package kube_test

import (
	"fmt"
	"relaper.com/kubemanage/example/exporter/kube"
	"testing"
)

func Test_GetPodMetrics(t *testing.T) {
	pod := kube.GetPodMetrics()
	fmt.Println(pod)
}
