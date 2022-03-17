package assemble

import (
	"github.com/lgy1027/kubemanage/model"
	"github.com/lgy1027/kubemanage/utils"
	v1 "k8s.io/api/core/v1"
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
