package k8s

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/lgy1027/kubemanage/inital"
	"github.com/lgy1027/kubemanage/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func GetEvents(namespace string, name string, kind int) []model.EventData {
	opt := metav1.ListOptions{}
	switch kind {
	case 0:
		opt.FieldSelector = fmt.Sprintf("involvedObject.kind=Node,involvedObject.name=%s", name)
	case 1:
		opt.FieldSelector = fmt.Sprintf("involvedObject.kind=Deployment,involvedObject.name=%s", name)
	case 2:
		opt.FieldSelector = fmt.Sprintf("involvedObject.kind=StatefulSet,involvedObject.name=%s", name)
	case 3:
		opt.FieldSelector = fmt.Sprintf("involvedObject.kind=Service,involvedObject.name=%s", name)
	case 4:
		opt.FieldSelector = fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s", name, namespace)
	}
	events, err := inital.GetGlobal().GetClientSet().CoreV1().Events(namespace).List(opt)
	if err != nil {
		log.Error("获取Pods错误", err.Error())
		return nil
	} else {
		return parseEvents(events.Items)
	}
}
