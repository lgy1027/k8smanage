package client

import k8s2 "relaper.com/kubemanage/k8s"

var baseClient *BaseClient

func init() {
	baseClient = &BaseClient{
		Deployment: k8s2.NewDeploy(),
		Sv:         k8s2.NewSv(),
		Sf:         k8s2.NewStateFulSet(),
		Ns:         k8s2.NewNs(),
		Pod:        k8s2.NewPod(),
		Node:       k8s2.NewNode(),
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
