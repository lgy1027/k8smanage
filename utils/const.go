package utils

const (
	ZERO = iota
	ONE
	TWO
)

const (
	DEPLOYMENT  = "deployments"
	STATEFULSET = "statefulsets"
	SERVICE     = "services"

	DEPLOY_DEPLOYMENT  = "Deployment"
	DEPLOY_STATEFULSET = "StatefulSet"
	DEPLOY_Service     = "Service"

	DEFAULTNS = "default"
)

type PodPhase string

const (
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	PodUnknown PodPhase = "Unknown"
)

const (
	CLUSTER_PREFIX_KEY  = "cluster_key_"
	CLUSTER_DETAIL_TIME = 600

	NAMESPACE_PREFIX_KEY = "namespace_key_"
	NAMESPACE_TIME       = 600

	POD_PREFIX_KEY = "pod_key_"
	POD_TIME       = 600
)
