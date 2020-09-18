package k8s

import (
	log "github.com/cihub/seelog"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"relaper.com/kubemanage/inital"
	"relaper.com/kubemanage/model"
)

func parseEvents(items []v1.Event) []model.EventData {
	var data []model.EventData
	for _, v := range items {
		t := model.EventData{}
		t.Reason = v.Reason
		t.Messages = v.Message
		t.EventTime = v.LastTimestamp.String()
		t.Host = v.Source.Host
		data = append(data, t)
	}
	return data
}

func GetEvents(namespace string, podName string) []model.EventData {
	opt := metav1.ListOptions{}
	opt.FieldSelector = "involvedObject.name=" + podName + ",involvedObject.namespace=" + namespace
	events, err := inital.GetGlobal().GetClientSet().CoreV1().Events(namespace).List(opt)
	if err != nil {
		log.Error("获取Pods错误", err.Error())
		return nil
	} else {
		return parseEvents(events.Items)
	}
}
