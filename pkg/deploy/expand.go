package deploy

import (
	"encoding/json"
	"fmt"
	"github.com/lgy1027/kubemanage/inital"
	"github.com/lgy1027/kubemanage/utils"
	"github.com/pkg/errors"
	goyaml "gopkg.in/yaml.v2"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func int32Ptr(i int32) *int32 { return &i }

func intOrString(i int) *intstr.IntOrString {
	maxSurge := intstr.FromInt(i)
	return &maxSurge
}

func ExpandDeployment(req *DeployRequest) *appsv1.Deployment {
	if req.ImagePullPolicy == "" {
		req.ImagePullPolicy = "IfNotPresent"
	}
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       req.Kind,
			APIVersion: inital.GetGlobal().GetConfig().Deployment,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  req.PodName,
							Image: req.Image,
							//Resources:       *req.Resources,
							Command:         req.Command,
							Args:            req.Args,
							ImagePullPolicy: req.ImagePullPolicy, // PullAlways PullPolicy = "Always"  PullNever PullPolicy = "Never"  PullIfNotPresent PullPolicy = "IfNotPresent"
							Env:             req.Envs,            // 环境变量
							WorkingDir:      req.WorkingDir,      // 工作目录
						},
					},
				},
			},
		},
	}

	deployment.ObjectMeta.Namespace = req.Namespace
	if req.Annotations != nil {
		deployment.ObjectMeta.Annotations = req.Annotations
	}
	if req.ObjectMetaLabels != nil {
		deployment.ObjectMeta.Labels = req.ObjectMetaLabels
	}
	if req.MatchLabels != nil {
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: req.MatchLabels,
		}
	}
	if req.Resources != nil {
		deployment.Spec.Template.Spec.Containers[0].Resources = *req.Resources
	}
	if req.Replicas > 0 {
		deployment.Spec.Replicas = int32Ptr(req.Replicas)
	}
	if req.MaxSurge > 0 && req.MaxUnavailable > 0 {
		deployment.Spec.Strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{ // 由于replicas为3,则整个升级,pod个数在2-4个之间
			MaxSurge:       intOrString(req.MaxSurge),       // 滚动升级时会先启动1个pod
			MaxUnavailable: intOrString(req.MaxUnavailable), // 滚动升级时允许的最大Unavailable的pod个数

		}
	}
	if len(req.PodPort) > 0 {
		deployment.Spec.Template.Spec.Containers[0].Ports = req.PodPort
	}
	if req.TemplateLabels != nil {
		deployment.Spec.Template.ObjectMeta.Labels = req.TemplateLabels
	}
	if req.NodeSelector != nil {
		deployment.Spec.Template.Spec.NodeSelector = req.NodeSelector
	}
	if len(req.VolumeMounts) > 0 && len(req.Volumes) > 0 {
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = req.VolumeMounts
		deployment.Spec.Template.Spec.Volumes = req.Volumes
	}

	return deployment
}

func ExpandStatefulSets(req *DeployRequest) *appsv1.StatefulSet {

	if req.ImagePullPolicy == "" {
		req.ImagePullPolicy = "IfNotPresent"
	}
	statefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       req.Kind,
			APIVersion: inital.GetGlobal().GetConfig().StatefulSet,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: req.ServiceName,
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  req.PodName,
							Image: req.Image,
							//Resources:       *req.Resources,
							Command:         req.Command,
							Args:            req.Args,
							ImagePullPolicy: req.ImagePullPolicy, // PullAlways PullPolicy = "Always"  PullNever PullPolicy = "Never"  PullIfNotPresent PullPolicy = "IfNotPresent"
							Env:             req.Envs,            // 环境变量
							WorkingDir:      req.WorkingDir,      // 工作目录
						},
					},
				},
			},
		},
	}

	statefulSet.ObjectMeta.Namespace = req.Namespace
	if req.Annotations != nil {
		statefulSet.ObjectMeta.Annotations = req.Annotations
	}
	if req.ObjectMetaLabels != nil {
		statefulSet.ObjectMeta.Labels = req.ObjectMetaLabels
	}
	if req.MatchLabels != nil {
		statefulSet.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: req.MatchLabels,
		}
	}
	if req.Resources != nil {
		statefulSet.Spec.Template.Spec.Containers[0].Resources = *req.Resources
	}
	if req.Replicas > 0 {
		statefulSet.Spec.Replicas = int32Ptr(req.Replicas)
	}

	if req.StatefulSetUpdateStrategyType != "" {
		statefulSet.Spec.UpdateStrategy.Type = req.StatefulSetUpdateStrategyType
	}

	if req.Partition > 0 {
		statefulSet.Spec.UpdateStrategy.RollingUpdate = &appsv1.RollingUpdateStatefulSetStrategy{
			Partition: int32Ptr(req.Partition),
		}
	}

	if len(req.PodPort) > 0 {
		statefulSet.Spec.Template.Spec.Containers[0].Ports = req.PodPort
	}
	if req.TemplateLabels != nil {
		statefulSet.Spec.Template.ObjectMeta.Labels = req.TemplateLabels
	}
	if req.NodeSelector != nil {
		statefulSet.Spec.Template.Spec.NodeSelector = req.NodeSelector
	}
	if len(req.VolumeMounts) > 0 {
		statefulSet.Spec.Template.Spec.Containers[0].VolumeMounts = req.VolumeMounts
	}

	if len(req.Volumes) > 0 {
		statefulSet.Spec.Template.Spec.Volumes = req.Volumes
	}

	return statefulSet
}

func ExpandService(req *DeployRequest) *apiv1.Service {
	service := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       utils.DEPLOY_Service,
			APIVersion: inital.GetGlobal().GetConfig().Service,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: req.ServiceName,
		},
	}
	if req.TemplateLabels != nil {
		service.Spec.Selector = req.TemplateLabels
	}
	if req.ClusterIP != "" {
		service.Spec.ClusterIP = req.ClusterIP
	}

	if req.ServiceType != "" {
		service.Spec.Type = req.ServiceType
	}

	if len(req.ServicePorts) > 0 {
		service.Spec.Ports = req.ServicePorts
	}

	return service
}

func ExpandObject(req *DeployRequest) *unstructured.Unstructured {

	if req.ImagePullPolicy == "" {
		req.ImagePullPolicy = "IfNotPresent"
	}
	kind := ""
	apiVersion := ""
	switch req.Kind {
	case utils.DEPLOY_DEPLOYMENT:
		apiVersion = inital.GetGlobal().GetConfig().Deployment
		kind = utils.DEPLOY_DEPLOYMENT
	case utils.DEPLOY_STATEFULSET:
		apiVersion = inital.GetGlobal().GetConfig().StatefulSet
		kind = utils.DEPLOY_STATEFULSET
	case utils.SERVICE:
		apiVersion = inital.GetGlobal().GetConfig().Service
		kind = utils.SERVICE
	}

	obj := &unstructured.Unstructured{}
	conf := map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":        req.Name,
			"labels":      req.ObjectMetaLabels,
			"annotations": req.Annotations,
			"namespace":   req.Namespace,
		},
		"spec": map[string]interface{}{
			"selector": map[string]interface{}{
				"matchLabels": req.MatchLabels,
			},
			"replicas":    req.Replicas,
			"serviceName": req.ServiceName,
			"template": map[string]interface{}{
				//"metadata": map[string]interface{}{
				//	"labels":req.TemplateLabels,
				//},
				"spec": map[string]interface{}{
					"nodeSelector": req.NodeSelector,
					"containers": []map[string]interface{}{
						map[string]interface{}{
							"args":            req.Args,
							"env":             ExpandEnv(req.Envs),
							"image":           req.Image,
							"imagePullPolicy": req.ImagePullPolicy,
							"name":            req.PodName,
							"ports":           ExpandPorts(req.PodPort),
							"volumeMounts":    ExpandVolumeMounts(req.VolumeMounts),
							"command":         req.Command,
							"workingDir":      req.WorkingDir,
							//"resources": map[string]interface{}{
							//
							//},
						},
					},
					//"volumes":				   ExpandVolumes(req.Volumes),
				},
			},
		},
	}

	if req.TemplateLabels != nil {
		conf["spec"].(map[string]interface{})["template"].(map[string]interface{})["metadata"] = map[string]interface{}{
			"labels": req.TemplateLabels,
		}
	}

	if len(req.Volumes) > 0 {
		conf["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["volumes"] = ExpandVolumes(req.Volumes)
	}

	if len(req.VolumeMounts) > 0 {
		conf["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]map[string]interface{})[0]["volumeMounts"] = ExpandVolumeMounts(req.VolumeMounts)
	}

	if req.Resources != nil {
		conf["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]map[string]interface{})[0]["resources"] = map[string]interface{}{
			"limits":   req.Resources.Limits,
			"requests": req.Resources.Requests,
		}
	}

	if req.StatefulSetUpdateStrategyType != "" {
		//statefulSet.Spec.UpdateStrategy.Type = req.StatefulSetUpdateStrategyType
		conf["spec"].(map[string]interface{})["updateStrategy"] = map[string]interface{}{
			"type": req.StatefulSetUpdateStrategyType,
		}
	}

	if req.Partition > 0 {
		if req.StatefulSetUpdateStrategyType != "" {
			conf["spec"].(map[string]interface{})["updateStrategy"].(map[string]interface{})["rollingUpdate"] = map[string]interface{}{
				"partition": req.Partition,
			}
		} else {
			conf["spec"].(map[string]interface{})["updateStrategy"] = map[string]interface{}{
				"rollingUpdate": map[string]interface{}{
					"partition": req.Partition,
				},
			}
		}
	}

	obj.Object = conf
	return obj
}

func ExpandEnv(envs []apiv1.EnvVar) []map[string]interface{} {
	envMap := make([]map[string]interface{}, 0)
	for _, env := range envs {
		envMap = append(envMap, map[string]interface{}{
			"name":  env.Name,
			"value": env.Value,
		})
	}
	return envMap
}

func ExpandVolumeMounts(volumes []apiv1.VolumeMount) []map[string]interface{} {
	volumeMap := make([]map[string]interface{}, 0)
	for _, volume := range volumes {
		volumeMap = append(volumeMap, map[string]interface{}{
			"mountPath": volume.MountPath,
			"name":      volume.Name,
		})
	}
	return volumeMap
}

func ExpandVolumes(volumes []apiv1.Volume) []map[string]interface{} {
	volumeMap := make([]map[string]interface{}, 0)
	for _, volume := range volumes {
		volumeMap = append(volumeMap, map[string]interface{}{
			"name": volume.Name,
			"hostPath": map[string]interface{}{
				"path": volume.HostPath.Path,
			},
		})
	}
	return volumeMap
}

func ExpandPorts(ports []apiv1.ContainerPort) []map[string]interface{} {

	portMap := make([]map[string]interface{}, 0)
	for _, port := range ports {
		portMap = append(portMap, map[string]interface{}{
			"name":          port.Name,
			"containerPort": port.ContainerPort,
			"protocol":      "TCP",
		})
	}
	return portMap
}

// 单文档读取
func ExpandSimpleYamlFileToObject(yaml []byte) (*unstructured.Unstructured, error) {
	temp := make(map[string]interface{})
	var data interface{}
	err := goyaml.Unmarshal(yaml, &data)
	if err != nil {
		return nil, err
	}
	yaml = nil
	err = transformData(&data)
	if err != nil {
		return nil, err
	}
	output, err := json.Marshal(data)
	fmt.Println(string(output))
	if err != nil {
		return nil, err
	}
	data = nil
	err = json.Unmarshal(output, &temp)
	if err != nil {
		return nil, errors.New("json 转换失败")
	}
	return &unstructured.Unstructured{Object: temp}, nil
}

// 多文档读取
func ExpandMultiYamlFileToObject(input io.Reader) ([]*unstructured.Unstructured, error) {
	yaml := goyaml.NewDecoder(input)
	var obj []*unstructured.Unstructured
	var data interface{}
	for yaml.Decode(&data) == nil {
		temp := make(map[string]interface{})
		err := transformData(&data)
		if err != nil {
			return nil, err
		}
		output, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		data = nil
		err = json.Unmarshal(output, &temp)
		if err != nil {
			return nil, errors.New("yaml文件转换失败")
		}
		obj = append(obj, &unstructured.Unstructured{Object: temp})
	}
	return obj, nil
}

func transformData(data *interface{}) (err error) {
	switch in := (*data).(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(in))
		for k, v := range in {
			if err = transformData(&v); err != nil {
				return err
			}
			var sk string
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			default:
				return fmt.Errorf("类型不匹配: 期望映射字符串或int类型; 当前类型: %T", k)
			}
			m[sk] = v
		}
		*data = m
	case []interface{}:
		for i := len(in) - 1; i >= 0; i-- {
			if err = transformData(&in[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
