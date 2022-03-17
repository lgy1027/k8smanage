package assemble

import (
	"github.com/lgy1027/kubemanage/model"
	"github.com/lgy1027/kubemanage/utils"
	appsv1 "k8s.io/api/apps/v1"
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

func AssembleDeploymentSimple(deploy appsv1.Deployment) model.DeploymentDetail {
	version, ok := deploy.Annotations["deployment.kubernetes.io/revision"]
	if !ok {
	}
	return model.DeploymentDetail{
		Kind:        utils.DEPLOY_DEPLOYMENT,
		Namespace:   deploy.GetNamespace(),
		MatchLabels: deploy.Spec.Selector.MatchLabels,
		Version:     version,
		Name:        deploy.GetName(),
		Status:      deploy.Status,
	}
}
