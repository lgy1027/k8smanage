package assemble

import (
	appsv1 "k8s.io/api/apps/v1"
	"relaper.com/kubemanage/model"
	"relaper.com/kubemanage/utils"
)

func AssembleDeployment(namespace string, deploys []appsv1.Deployment) []model.DeploymentDetail {
	deployList := make([]model.DeploymentDetail, 0)
	for _, deploy := range deploys {
		if namespace != "" {
			if deploy.GetNamespace() != namespace {
				continue
			}
		}
		deployList = append(deployList, model.DeploymentDetail{
			Kind:      utils.DEPLOY_DEPLOYMENT,
			Namespace: deploy.GetNamespace(),
			Name:      deploy.GetName(),
			//Spec:      deploy.Spec,
			Status: deploy.Status,
		})
	}
	return deployList
}
