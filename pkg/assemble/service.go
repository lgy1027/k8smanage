package assemble

import (
	v1 "k8s.io/api/core/v1"
	"relaper.com/kubemanage/model"
	"relaper.com/kubemanage/utils"
)

func AssembleService(namespace string, svcs []v1.Service) []model.ServiceDetail {
	svcList := make([]model.ServiceDetail, 0)
	for _, svc := range svcs {
		if namespace != "" {
			if svc.GetNamespace() != namespace {
				continue
			}
		}
		svcList = append(svcList, model.ServiceDetail{
			Kind:       utils.DEPLOY_Service,
			Namespace:  svc.GetNamespace(),
			Name:       svc.GetName(),
			Spec:       svc.Spec,
			ObjectMeta: svc.ObjectMeta,
			Status:     svc.Status,
		})
	}
	return svcList
}
