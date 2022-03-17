package assemble

import (
	"github.com/lgy1027/kubemanage/model"
	"github.com/shopspring/decimal"
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func AssembleResourceList(nodeMetrics []v1beta1.NodeMetrics, nodes []v1.Node) []model.ResourceDetail {
	resourceList := make([]model.ResourceDetail, 0)
	for _, metric := range nodeMetrics {
		for _, node := range nodes {
			if metric.GetName() == node.GetName() {
				cpuNum := node.Status.Capacity.Cpu().Value()
				useCpuNum := metric.Usage.Cpu().Value()
				cpuNumValue := decimal.NewFromInt(cpuNum)
				useCpuNumValue := decimal.NewFromInt(useCpuNum)
				cpuNumValue = cpuNumValue.Div(decimal.NewFromInt(1000))
				useCpuNumValue = useCpuNumValue.Div(decimal.NewFromInt(1000))
				memNum := node.Status.Capacity.Memory().Value()
				userMemNum := metric.Usage.Memory().Value()
				memNumValue := decimal.NewFromInt(memNum)
				memNumValue = memNumValue.Div(decimal.NewFromInt(1024)).Div(decimal.NewFromInt(1024)).DivRound(decimal.NewFromInt(1024), 2)
				useMemNumValue := decimal.NewFromInt(userMemNum)
				useMemNumValue = useMemNumValue.Div(decimal.NewFromInt(1024)).Div(decimal.NewFromInt(1024)).DivRound(decimal.NewFromInt(1024), 2)
				resourceList = append(resourceList, model.ResourceDetail{
					CpuNum:         cpuNumValue.String(),
					CpuFree:        cpuNumValue.Sub(useCpuNumValue).String(),
					CpuUse:         useCpuNumValue.String(),
					CpuFreePercent: cpuNumValue.Sub(useCpuNumValue).DivRound(cpuNumValue, 2).String(),
					CpuUsePercent:  useCpuNumValue.DivRound(cpuNumValue, 2).String(),
					MemSize:        memNumValue.String(),
					MemFree:        memNumValue.Sub(useMemNumValue).String(),
					MemUse:         useMemNumValue.String(),
					MemFreePercent: memNumValue.Sub(useMemNumValue).DivRound(memNumValue, 2).String(),
					MemUsePercent:  useMemNumValue.DivRound(memNumValue, 2).String(),
				})
				break
			}
		}
	}
	return resourceList
}

func AssembleResource(metric v1beta1.NodeMetrics, node v1.Node) model.ResourceDetail {
	var resource model.ResourceDetail
	cpuNum := node.Status.Capacity.Cpu().MilliValue()
	useCpuNum := metric.Usage.Cpu().MilliValue()
	cpuNumValue := decimal.NewFromInt(cpuNum)
	useCpuNumValue := decimal.NewFromInt(useCpuNum)
	cpuNumValue = cpuNumValue.Div(decimal.NewFromInt(1000))
	useCpuNumValue = useCpuNumValue.Div(decimal.NewFromInt(1000))
	memNum := node.Status.Capacity.Memory().Value()
	userMemNum := metric.Usage.Memory().Value()
	memNumValue := decimal.NewFromInt(memNum)
	memNumValue = memNumValue.Div(decimal.NewFromInt(1024)).Div(decimal.NewFromInt(1024)).DivRound(decimal.NewFromInt(1024), 2)
	useMemNumValue := decimal.NewFromInt(userMemNum)
	useMemNumValue = useMemNumValue.Div(decimal.NewFromInt(1024)).Div(decimal.NewFromInt(1024)).DivRound(decimal.NewFromInt(1024), 2)
	resource = model.ResourceDetail{
		CpuNum:         cpuNumValue.String(),
		CpuFree:        cpuNumValue.Sub(useCpuNumValue).String(),
		CpuUse:         useCpuNumValue.String(),
		CpuFreePercent: cpuNumValue.Sub(useCpuNumValue).DivRound(cpuNumValue, 2).String(),
		CpuUsePercent:  useCpuNumValue.DivRound(cpuNumValue, 2).String(),
		MemSize:        memNumValue.String(),
		MemFree:        memNumValue.Sub(useMemNumValue).String(),
		MemUse:         useMemNumValue.String(),
		MemFreePercent: memNumValue.Sub(useMemNumValue).DivRound(memNumValue, 2).String(),
		MemUsePercent:  useMemNumValue.DivRound(memNumValue, 2).String(),
	}
	return resource
}
