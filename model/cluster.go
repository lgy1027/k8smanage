package model

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

type Cluster struct {
	// @description 节点数量
	NodeNum int `json:"nodeNum,omitempty"`
	// @description 正常运行的节点数量
	RunNodeNum int `json:"runNodeNum,omitempty"`
	// @description 节点详情
	Nodes []NodeDetail `json:"nodes,omitempty"`
	// @description pod上线数量
	PodNum int64 `json:"podNum,omitempty"`
	// @description 运行的pod数量
	ActivePodNum int64 `json:"activePodNum,omitempty"`
	// @description 命名空间相关信息
	NameSpaceNum int `json:"namespaceNum,omitempty"`
	//NamespaceDetail []NamespaceDetail `json:"namespaceDetail,omitempty"`
	// @description 集群总的指标情况
	Resource ResourceDetail `json:"resource,omitempty"`
}

type NodeDetail struct {
	Name string `json:"name,omitempty"`
	// @description 主机IP
	HostIp string `json:"hostIp,omitempty"`
	// @description 主机标签
	Label map[string]string `json:"label,omitempty"`
	// @description 注释
	Annotation map[string]string `json:"annotation,omitempty"`
	// @description 状态
	Status string `json:"status,omitempty"`
	// @description pod数量
	PodNum int64 `json:"podNum,omitempty"`
	// @description 部署的pod总量
	PodTotal int64 `json:"podTotal,omitempty"`
	// @description 部署的pod总量
	PodRun int64 `json:"podRun,omitempty"`
	// @description 是否有效
	IsValid string `json:"isValid,omitempty"`
	// @description nodeId
	NodeID string `json:"nodeID,omitempty"`
	// @description 创建时间
	CreateTime string `json:"createTime,omitempty"`
	// @description 镜像数量
	ImageNum int `json:"imageNum,omitempty"`
	// @description 最后一次心跳时间
	LastHeartbeatTime string `json:"lastHeartbeatTime,omitempty"`
	// @description Kubelet版本
	KuBeLetVersion string `json:"kuBeLetVersion,omitempty"`
	// @description Kubelet版本
	KuProxyVersion string `json:"kuProxyVersion,omitempty"`
	// @description 操作系统类型
	SystemType string `json:"systemType,omitempty"`
	// @description 操作系统
	SystemOs string `json:"systemOs,omitempty"`
	// @description docker版本
	DockVersion string `json:"dockVersion,omitempty"`
	// @description 内核版本
	KernlVersion string `json:"kernlVersion,omitempty"`
	// @description 角色
	Role string `json:"role,omitempty"`
	// @description pod列表
	Pods []PodDetail `json:"pods,omitempty"`
	// @description 集群名
	ClusterName string         `json:"clusterName,omitempty"`
	Resource    ResourceDetail `json:"resource,omitempty"`
}

type PodDetail struct {
	Name      string `json:"name,omitempty"`
	Id        string `json:"id,omitempty"`
	NodeName  string `json:"nodeName,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	// @description 状态 可选 Pending：正在启动 Running：运行中 Succeeded：部署成功未启动 Failed：失败 Unknown：未知
	Status string `json:"status,omitempty"`
	// @description 创建时间
	CreateTime string `json:"createTime,omitempty"`
	// @description 主机标签
	Label map[string]string `json:"label,omitempty"`
	// @description 重启次数
	RestartCount int32 `json:"restartCount,omitempty"`
	// @description 宿主机地址
	HostIp string `json:"hostIp,omitempty"`
	// @description 容器IP
	PodIp string `json:"podIp,omitempty"`
	// @description 注释
	Annotation map[string]string `json:"annotation,omitempty"`
	Resource   ResourceDetail    `json:"resource,omitempty"`
	EventData  []EventData       `json:"eventData,omitempty"`
}

type Label struct {
	// @description 主机标签
	Label map[string]string `json:"label,omitempty"`
	// @description 注释
	Annotation map[string]string `json:"annotation,omitempty"`
}

type ResourceDetail struct {
	// @description cpu数量
	CpuNum string `json:"cpuNum,omitempty"`
	// @description cpu剩余量
	CpuFree string `json:"cpuFree,omitempty"`
	// @description cpu使用量
	CpuUse string `json:"cpuUse,omitempty"`
	// @description cpu剩余百分比
	CpuFreePercent string `json:"cpuFreePercent,omitempty"`
	// @description cpu使用量百分比
	CpuUsePercent string `json:"cpuUsePercent,omitempty"`
	// @description 内存大小
	MemSize string `json:"memSize,omitempty"`
	// @description 内存剩余量
	MemFree string `json:"memFree,omitempty"`
	// @description 内存使用量
	MemUse string `json:"memUse,omitempty"`
	// @description 内使剩余量百分比
	MemFreePercent string `json:"memFreePercent,omitempty"`
	// @description 内使用百分比
	MemUsePercent string `json:"memUsePercent,omitempty"`
}

type NamespaceDetail struct {
	Name       string `json:"name,omitempty"`
	CreateTime string `json:"createTime,omitempty"`
	// @description 状态 可选  Active： 正常使用   Terminating：正在终止
	Status string `json:"status,omitempty"`
	// @description 无状态资源
	DeploymentList []DeploymentDetail `json:"deployments,omitempty"`
	// @description 有状态资源
	StatefulSetList []StatefulSetDetail `json:"statefulSets,omitempty"`
	// @description 服务资源
	ServiceList []ServiceDetail `json:"services,omitempty"`
}

type DeploymentDetail struct {
	// @description 资源类型
	Kind      string `json:"kind,omitempty" default:"Deployment"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	//Spec      appsv1.DeploymentSpec   `json:"spec,omitempty"`
	Status appsv1.DeploymentStatus `json:"status,omitempty"`
}

type StatefulSetDetail struct {
	// @description 资源类型
	Kind      string `json:"kind" default:"StatefulSet"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	//Spec      appsv1.StatefulSetSpec   `json:"spec,omitempty"`
	Status appsv1.StatefulSetStatus `json:"status,omitempty"`
}

type ServiceDetail struct {
	// @description 资源类型
	Kind      string `json:"kind" default:"Service"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	//Spec      v1.ServiceSpec   `json:"spec,omitempty"`
	Status v1.ServiceStatus `json:"status,omitempty"`
}

type EventData struct {
	// @description 事件时间
	EventTime string `json:"eventTime,omitempty"`
	// @description 信息
	Messages string `json:"messages,omitempty"`
	// @description 原因
	Reason string `json:"reason,omitempty"`
	// @description 主机Ip
	Host string `json:"host,omitempty"`
}
