// @Author : liguoyu
// @Date: 2019/10/29 15:42
package deploy

import (
	"github.com/pkg/errors"
	apisv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"net/http"
	"strings"
)

type Validator interface {
	Validate() error
}

// 处理文件上传
type UploadRequest struct {
	*http.Request
}

func (r *UploadRequest) SetRequest(httpReq *http.Request) {
	r.Request = httpReq
}

type UploadResponse struct {
	// @description 文件名
	File string `json:"file"`
}

type DeployRequest struct {
	// @description  ========  Object  ===========
	// @description 资源类型 可选 Deployment | StatefulSet | Service
	Kind string `json:"kind"`
	// @description 命名空间
	Namespace string `json:"namespace"`
	// @description 服务名称
	Name string `json:"name"`
	// @description 资源
	ObjectMetaLabels map[string]string `json:"objectMetaLabels"`
	// @description 注释
	Annotations map[string]string `json:"annotations"`
	// @description =========   spec   ============
	// @description 副本数量
	Replicas int32 `json:"replicas"`
	// @description 上层spec标签
	MatchLabels map[string]string `json:"matchLabels"`
	// @description 滚动升级时候,会优先启动的pod数量
	MaxSurge int `json:"maxSurge" default:"1"`
	// @description 滚动升级时候,最大的unavailable数量
	MaxUnavailable int `json:"maxUnavailable" default:"1"`

	// @description =============  template  =============
	// @description template标签
	TemplateLabels map[string]string `json:"templateLabels"`
	// @description 节点选择 node:node1
	NodeSelector map[string]string `json:"nodeSelector"`

	// @description ============  Containers  ==================
	// @description pod名
	PodName string `json:"podName"`
	// @description 镜像名称
	Image string `json:"image"`
	// @description 容器暴露端口
	PodPort []apiv1.ContainerPort `json:"podPort"`
	// @description 资源限制
	Resources apiv1.ResourceRequirements `json:"resources"`

	// @description 容器启动执行的命令
	Command []string `json:"command"`
	// @description 初始参数
	Args []string `json:"args"`

	// @description 镜像拉去策略 Always | Never | IfNotPresent
	ImagePullPolicy apiv1.PullPolicy `json:"imagePullPolicy" default:"IfNotPresent"`

	// @description 环境变量数据
	Envs []apiv1.EnvVar `json:"envs"`

	// @description 工作目录
	WorkingDir string `json:"workingDir"`

	// @description 目录挂载
	VolumeMounts []apiv1.VolumeMount `json:"volumeMounts"`

	// @description 外部目录映射
	Volumes []apiv1.Volume `json:"volumes"`

	// @description ============= statefulSet  ================
	// @description 有状态服务专用 可选 RollingUpdate | OnDelete
	StatefulSetUpdateStrategyType apisv1.StatefulSetUpdateStrategyType `json:"statefulType"`

	// @description Partition
	Partition int32 `json:"partition"`

	// @description ============   service   ==================
	// @description 是否创建serice
	CreateService bool `json:"createService"`
	// @description service名
	ServiceName string `json:"serviceName"`
	// @description 可选 ClusterIP | NodePort | LoadBalancer
	ServiceType apiv1.ServiceType `json:"serviceType"`

	ClusterIP string `json:"clusterIp"`
	// @description 端口映射
	ServicePorts []apiv1.ServicePort `json:"servicePorts"`
}

func (r *DeployRequest) Validate() error {
	r.Kind = strings.TrimSpace(r.Kind)
	r.Name = strings.TrimSpace(r.Name)
	r.Namespace = strings.TrimSpace(r.Namespace)
	r.PodName = strings.TrimSpace(r.PodName)
	r.ServiceName = strings.TrimSpace(r.ServiceName)
	r.ClusterIP = strings.TrimSpace(r.ClusterIP)
	r.Name = strings.TrimSpace(r.Name)
	r.Namespace = strings.TrimSpace(r.Namespace)
	if r.Replicas < 1 {
		r.Replicas = 1
	}

	if r.Namespace == "" {
		r.Namespace = "default"
	}

	return nil
}

type DeploymentResponse struct{}

type DeleteRequest struct {
	// @description 资源类型 可选 Deployment | StatefulSet | Service
	Kind string `json:"kind"`
	// @description 命名空间
	Namespace string `json:"namespace"`
	// @description 资源名称
	Name string `json:"name"`
}

func (r *DeleteRequest) Validate() error {
	r.Kind = strings.TrimSpace(r.Kind)
	r.Name = strings.TrimSpace(r.Name)
	r.Namespace = strings.TrimSpace(r.Namespace)
	return nil
}

type DeleteResponse struct{}

type ExpansionRequest struct {
	// @description 资源类型 可选 Deployment | StatefulSet | Service
	Kind string `json:"kind"`
	// @description 命名空间
	Namespace string `json:"namespace"`
	// @description 资源名称
	Name string `json:"name"`
	// @description 扩容大小
	Replicas int32 `json:"replicas"`
	// @description 资源限制
	Resources apiv1.ResourceRequirements `json:"resources"`
	// @description 镜像
	ImagePull string `json:"imagePull"`
}

func (r *ExpansionRequest) Validate() error {
	r.Kind = strings.TrimSpace(r.Kind)
	r.Name = strings.TrimSpace(r.Name)
	r.Namespace = strings.TrimSpace(r.Namespace)
	r.ImagePull = strings.TrimSpace(r.ImagePull)
	return nil
}

type ExpansionResponse struct{}

type NamespaceRequest struct {
	// @description 命名空间
	Namespace string `json:"namespace"`
}

func (r *NamespaceRequest) Validate() error {
	r.Namespace = strings.TrimSpace(r.Namespace)
	if r.Namespace == "" {
		return errors.New("命名空间不能为空")
	}
	return nil
}

type NamespaceResponse struct{}
