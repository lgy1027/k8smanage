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
	Name string `query:"name",description:"节点名"`
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

type NsRequest struct{}

type NsResponse struct {
	Num        int                     `json:"num,omitempty"`
	Namespaces []model.NamespaceDetail `json:"namespaces,omitempty"`
}

type NameSpaceRequest struct {
	NameSpace string `query:"namespace"`
}

func (r *NameSpaceRequest) Validate() error {
	r.NameSpace = strings.TrimSpace(r.NameSpace)
	if r.NameSpace == "" {
		return errors.New("命名空间不能为空")
	}
	return nil
}

type NameSpaceResponse struct {
	Exist      bool                  `json:"exist,omitempty"`
	Namespaces model.NamespaceDetail `json:"namespace,omitempty"`
}

type PodInfoRequest struct {
	// @description 命名空间
	NameSpace string `query:"namespace"`
	// @description pod
	PodName string `query:"podName"`
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
	// @description 节点
	NodeName string `json:"nodeName"`
}

func (r *PodsRequest) Validate() error {
	r.NameSpace = strings.TrimSpace(r.NameSpace)
	return nil
}

type PodsResponse struct {
	// @description pods列表
	Pods []model.PodDetail `json:"pods,omitempty"`
}

type PodLogResponse struct {
	// @description pod日志
	Log string `json:"log,omitempty"`
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

type GetYamlRequest struct {
	// @description 资源类型
	Kind string `query:"kind"`
	// @description 命名空间
	Namespace string `query:"namespace"`
	// @description 资源名
	Name string `query:"name"`
}

type GetYamlResponse struct {
	Yaml map[string]interface{} `json:"yaml,omitempty"`
}

type EventRequest struct {
	// @description 事件类型   0 node 1 Deployment 2 StatefulSet 3 Service 4 pod
	Kind      int    `query:"kind"`
	Namespace string `query:"namespace"`
	Name      string `query:"name"`
}

type EventResponse struct {
	Event []model.EventData `json:"event,omitempty"`
}

type VersionRequest struct {
	Namespace string `query:"namespace"`
	Name      string `query:"name"`
	// @description 标签 app=demo
	Label string `query:"label"`
}

type VersionResponse struct {
	VersionList []model.Versions `json:"versions,omitempty"`
}
