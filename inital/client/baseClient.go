package client

import (
	k8s2 "github.com/lgy1027/kubemanage/k8s"
	"k8s.io/client-go/kubernetes"
)

var baseClient *BaseClient

func Init(clientSet kubernetes.Interface) {
	baseClient = &BaseClient{
		Deployment: k8s2.NewDeploy(clientSet),
		Sv:         k8s2.NewSv(clientSet),
		Sf:         k8s2.NewStateFulSet(clientSet),
		Ns:         k8s2.NewNs(clientSet),
		Pod:        k8s2.NewPod(clientSet),
		Node:       k8s2.NewNode(clientSet),
	}
}

type BaseClient struct {
	*k8s2.Deployment
	*k8s2.Sv
	*k8s2.Sf
	*k8s2.Ns
	*k8s2.Pod
	*k8s2.Node
}

func GetBaseClient() *BaseClient {
	return baseClient
}
