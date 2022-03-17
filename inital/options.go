package inital

import (
	"flag"
	log "github.com/cihub/seelog"
	"github.com/lgy1027/kubemanage/protocol"
)

type Options struct {
	Dev       bool   `flag:"dev"`
	Root      string `flag:"root"`
	CfgPath   string `flag:"cfg-path"`
	LogPath   string `flag:"log-path"`
	Address   string `flag:"address"`
	K8sConfig string `flag:"k8s-config"`
}

func NewOptions() *Options {
	return &Options{}
}

func (opts *Options) FlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("kubemanage", flag.ExitOnError)

	// basic options
	flagSet.String("config", "config/config.toml", "path to config file")
	flagSet.String("root", "", "project root dir")
	flagSet.Bool("dev", false, "development mode")
	flagSet.String("cfg-path", "", "cfg path")
	flagSet.String("log-path", "", "log path")
	flagSet.String("address", ":80", "service address")
	flagSet.String("k8s-config", "", "k8s cluster config file path")

	return flagSet
}

type options map[string]interface{}

// Validate settings in the config file, and fatal on errors
func (o *options) Validate() {

}

type config struct {
	Deployment  string
	StatefulSet string
	Service     string
	Task        bool
	Redis       protocol.Redis
}

func (cfg *config) Init() (err error) {
	if err = cfg.Redis.Init(); err != nil {
		log.Debug("init redis failed:" + err.Error())
	}
	return
}
