package inital

import (
	"crypto/tls"
	"github.com/BurntSushi/toml"
	log "github.com/cihub/seelog"
	gooptions "github.com/mreiferson/go-options"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func InitRoot(opts *Options) {
	if opts.Root != "" {
		global.root, _ = filepath.Abs(opts.Root)
	} else {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		global.root = filepath.Dir(path)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// init options, fatal on errors
func InitOptions(arguments []string) *Options {
	log.Info("action", "start init options")

	opts := NewOptions()

	flagSet := opts.FlagSet()
	flagSet.Parse(arguments)

	var o options
	configFile := flagSet.Lookup("config").Value.String()
	if configFile != "" {
		log.Info("action", "init options from configFile", "configFile", configFile)

		_, err := toml.DecodeFile(configFile, &o)
		if err != nil {
			log.Errorf("failed to load config file %s - %s - err:v%/n", "configFile", configFile, err)
		}
	}
	o.Validate()

	gooptions.Resolve(opts, flagSet, o)

	global.opts = opts

	log.Info("name", "success init options")
	return opts
}

func InitCfg(opts *Options) error {
	log.Info("action", "start init config")
	var cfg config
	var err error
	var cfgFile = global.GetRealPath(opts.CfgPath)
	_, err = toml.DecodeFile(cfgFile, &cfg)
	if err != nil {
		return log.Error("ERROR", " failed to load cfgFile ", " cfgFile ", cfgFile, " err: ", err)
	}

	if opts.K8sConfig == "" {
		if home := homeDir(); home != "" {
			opts.K8sConfig = filepath.Join(home, ".kube", "config")
		} else {
			return log.Error("ERROR:", " failed to load K8sConfig ", " K8sConfig: ", opts.K8sConfig, " err: ", err)
		}
	}

	k8sConfig, err := clientcmd.BuildConfigFromFlags("", opts.K8sConfig)
	if err != nil {
		return log.Error("ERROR", " error creating inClusterConfig, falling back to default config: ", err)
	}
	global.k8sConfig = k8sConfig

	global.dynamicClient, err = dynamic.NewForConfig(k8sConfig)
	if err != nil {
		return log.Error("ERROR", " create dynamic client err: ", err)
	}
	global.clientSet, err = kubernetes.NewForConfig(k8sConfig)

	if err != nil {
		return log.Error("ERROR", " create clientSet err: ", err)
	}

	global.metrics, err = metricsv.NewForConfig(k8sConfig)

	if err != nil {
		return log.Error("ERROR", " create metrics client err: ", err)
	}
	err = cfg.Init()
	if err != nil {
		return log.Error("action ", " failed to load cfgFile ", " cfgFile ", cfgFile, " err: ", err.Error())
	}
	global.cfg = &cfg
	initClient()
	log.Info("action", "success init config")
	return nil
}

func CloseCache() {
	global.GetCache().Close()
}

const (
	MaxIdleConns          int = 100
	MaxIdleConnsPerHost   int = 100
	ResponseHeaderTimeout int = 200
	IdleConnTimeout       int = 200
)

func initClient() {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   200 * time.Second,
				KeepAlive: 200 * time.Second,
			}).Dial,
			DisableKeepAlives:     false,
			MaxIdleConns:          MaxIdleConns,
			MaxIdleConnsPerHost:   MaxIdleConnsPerHost,
			IdleConnTimeout:       time.Duration(IdleConnTimeout) * time.Second,
			ResponseHeaderTimeout: time.Duration(ResponseHeaderTimeout) * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}
	global.client = client
}
