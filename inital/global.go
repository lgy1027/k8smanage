package inital

import (
	"github.com/garyburd/redigo/redis"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"net/http"
	"path/filepath"
	"relaper.com/kubemanage/utils"
	"relaper.com/kubemanage/utils/cache"
	"strings"
)

var global = &Global{}

var resource = map[int]string{
	1: "Deployment",
	2: "StatefulSet",
	3: "Service",
}

// 下面必须保持只读
// 不可以保证并发安全。所以必须率先初始化
func GetGlobal() *Global {
	return global
}

type Global struct {
	root          string
	opts          *Options
	cfg           *config
	k8sConfig     *rest.Config
	client        *http.Client
	clientSet     *kubernetes.Clientset
	restClient    *rest.RESTClient
	dynamicClient dynamic.Interface
	metrics       *versioned.Clientset
}

func (g *Global) GetOptions() *Options {
	return g.opts
}

func (g *Global) GetConfig() *config {
	return g.cfg
}

func (g *Global) GetRealPath(path string) string {
	if !filepath.IsAbs(path) && g.root != "" {
		path = filepath.Join(g.root, path)
	}
	return path
}

func (g *Global) GetHttpClient() *http.Client {
	return g.client
}

func (g *Global) GetClientSet() *kubernetes.Clientset {
	return g.clientSet
}

func (g *Global) GetDynamicClient() dynamic.Interface {
	return g.dynamicClient
}

func (g *Global) GetK8sConfig() *rest.Config {
	return g.k8sConfig
}

func (g *Global) GetMetricsClient() *versioned.Clientset {
	return g.metrics
}

func (g *Global) GetRedisConn() redis.Conn {
	return g.cfg.Redis.Pool().Get()
}

func (g *Global) GetCache() cache.Cache {
	return cache.NewRedisCache(g.cfg.Redis.Pool())
}

func (g *Global) GetRes(resource string) schema.GroupVersionResource {
	var gvs []string
	switch resource {
	case utils.DEPLOYMENT:
		gvs = strings.Split(g.GetConfig().Deployment, "/")
		if len(gvs) == 2 {
			return schema.GroupVersionResource{Group: gvs[0], Version: gvs[1], Resource: utils.DEPLOYMENT}
		}
	case utils.STATEFULSET:
		gvs = strings.Split(g.GetConfig().StatefulSet, "/")
		if len(gvs) == 2 {
			return schema.GroupVersionResource{Group: gvs[0], Version: gvs[1], Resource: utils.STATEFULSET}
		}
	case utils.SERVICE:
		gvs = strings.Split(g.GetConfig().Deployment, "/")
		if len(gvs) == 1 {
			return schema.GroupVersionResource{Group: "", Version: gvs[0], Resource: utils.SERVICE}
		}
	}
	return schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: utils.DEPLOYMENT}
}

func (g *Global) GetResForCustom(resource, apiVersion string) schema.GroupVersionResource {
	var gvs []string
	switch resource {
	case utils.DEPLOYMENT:
		gvs = strings.Split(apiVersion, "/")
		if len(gvs) == 2 {
			return schema.GroupVersionResource{Group: gvs[0], Version: gvs[1], Resource: utils.DEPLOYMENT}
		}
	case utils.STATEFULSET:
		gvs = strings.Split(apiVersion, "/")
		if len(gvs) == 2 {
			return schema.GroupVersionResource{Group: gvs[0], Version: gvs[1], Resource: utils.STATEFULSET}
		}
	case utils.SERVICE:
		gvs = strings.Split(apiVersion, "/")
		if len(gvs) == 1 {
			return schema.GroupVersionResource{Group: "", Version: gvs[0], Resource: utils.SERVICE}
		}
	}
	return schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: utils.DEPLOYMENT}
}

func (g *Global) GetResourceClient(resource string) dynamic.NamespaceableResourceInterface {
	return g.GetDynamicClient().Resource(g.GetRes(resource))
}

func (g *Global) GetResourceClientForCustom(resource, apiVersion string) dynamic.NamespaceableResourceInterface {
	return g.GetDynamicClient().Resource(g.GetResForCustom(resource, apiVersion))
}

func (g *Global) GetResourceToString(rType int) string {
	rs, ok := resource[rType]
	if ok {
		return rs
	}
	return ""
}
