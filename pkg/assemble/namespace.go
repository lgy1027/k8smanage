package assemble

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"relaper.com/kubemanage/model"
	"relaper.com/kubemanage/utils"
)

func AssembleNamespace(kind string, resource interface{}, ns v1.Namespace) model.NamespaceDetail {
	namespaceDetail := model.NamespaceDetail{
		Name:       ns.GetName(),
		CreateTime: ns.GetCreationTimestamp().String(),
	}
	switch kind {
	case utils.DEPLOY_DEPLOYMENT:
		deploys := resource.([]appsv1.Deployment)
		deployList := AssembleDeployment(ns.GetName(), deploys)
		namespaceDetail.DeploymentList = deployList
	case utils.DEPLOY_STATEFULSET:
		states := resource.([]appsv1.StatefulSet)
		statList := AssembleStatefulSet(ns.GetName(), states)
		namespaceDetail.StatefulSetList = statList
	case utils.DEPLOY_Service:
		svcs := resource.([]v1.Service)
		svcList := AssembleService(ns.GetName(), svcs)
		namespaceDetail.ServiceList = svcList
	}
	return namespaceDetail
}
