// @Author : liguoyu
// @Date: 2019/10/29 15:42
package cluster

import (
	"github.com/pkg/errors"
	"relaper.com/kubemanage/model"
	"strings"
)

type Validator interface {
	Validate() error
}

type NodesRequest struct{}

type NodesResponse struct {
	NodeList []model.NodeDetail `json:"nodeList"`
}

type ClusterRequest struct{}

type ClusterResponse struct {
	Cluster model.Cluster `json:"cluster"`
}

type NodeRequest struct {
	// @description 节点名
	Name string `json:"name"`
}

func (r *NodeRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return errors.New("节点名称不能为空")
	}
	return nil
}

type NodeResponse struct {
	// @description 是否存在
	Exist bool             `json:"exist"`
	Node  model.NodeDetail `json:"node"`
}

type NameSpacesRequest struct {
	// @description 命名空间 非必填
	Name string `json:"name"`
}

func (r *NameSpacesRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	return nil
}

type NameSpacesResponse struct {
	Namespaces []model.NamespaceDetail `json:"namespaces,omitempty"`
}

type PodInfoRequest struct {
	// @description 命名空间
	NameSpace string `json:"namespace"`
	// @description pod
	PodName string `json:"podName"`
}

func (r *PodInfoRequest) Validate() error {
	r.NameSpace = strings.TrimSpace(r.NameSpace)
	r.PodName = strings.TrimSpace(r.PodName)
	if r.PodName == "" {
		return errors.New("pod名称不能为空")
	}
	if r.NameSpace == "" {
		r.NameSpace = "default"
	}
	return nil
}

type PodInfoResponse struct {
	// @description pod信息
	Pod model.PodDetail `json:"pod,omitempty"`
	// @description 是否存在
	Exist bool `json:"exist,omitempty"`
}

type PodsRequest struct {
	// @description 命名空间
	NameSpace string `json:"namespace"`
}

func (r *PodsRequest) Validate() error {
	r.NameSpace = strings.TrimSpace(r.NameSpace)
	return nil
}

type PodsResponse struct {
	// @description pods列表
	Pods []model.PodDetail `json:"pods,omitempty"`
}

type ResourceRequest struct {
	// @description 命名空间
	NameSpace string `json:"namespace"`
	// @description 资源名称
	Name string `json:"name"`
}

func (r *ResourceRequest) Validate() error {
	r.NameSpace = strings.TrimSpace(r.NameSpace)
	r.Name = strings.TrimSpace(r.Name)
	return nil
}

type DeploymentsResponse struct {
	// @description 无状态资源列表
	Items []model.DeploymentDetail `json:"items"`
}

type StatefulSetsResponse struct {
	// @description 有状态资源列表
	Items []model.StatefulSetDetail `json:"items"`
}

type ServiceResponse struct {
	// @description 有状态资源列表
	Items []model.ServiceDetail `json:"items"`
}
