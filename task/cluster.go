package task

import (
	log "github.com/cihub/seelog"
	app "github.com/lgy1027/kubemanage/cache"
	"github.com/lgy1027/kubemanage/inital"
	"github.com/robfig/cron"
)

func ClusterTask() {

	log.Info("任务计划参数->", inital.GetGlobal().GetConfig().Task)
	if !inital.GetGlobal().GetConfig().Task {
		return
	}
	log.Info("Enable Cluster Task.........")
	c := cron.New()
	spec := "0 0/5 * * * ?"
	c.AddFunc(spec, func() {
		app.Cache()
	})
	c.Start()
}
