package assemble

import (
	"github.com/lgy1027/kubemanage/model"
	"github.com/lgy1027/kubemanage/utils"
	appsv1 "k8s.io/api/apps/v1"
)

func AssembleStatefulSet(namespace string, stats []appsv1.StatefulSet) []model.StatefulSetDetail {
	statList := make([]model.StatefulSetDetail, 0)
	for _, stat := range stats {
		if namespace != "" {
			if stat.GetNamespace() != namespace {
				continue
			}
		}
		statList = append(statList, model.StatefulSetDetail{
			Kind:      utils.DEPLOY_Service,
			Namespace: stat.GetNamespace(),
			Name:      stat.GetName(),
			//Spec:      stat.Spec,
			Status: stat.Status,
		})
	}
	return statList
}

func AssembleStatefulSetSimple(stat appsv1.StatefulSet) model.StatefulSetDetail {

	stateful := model.StatefulSetDetail{
		Kind:      utils.DEPLOY_Service,
		Namespace: stat.GetNamespace(),
		Name:      stat.GetName(),
		//Spec:      stat.Spec,
		Status: stat.Status,
	}
	return stateful
}
