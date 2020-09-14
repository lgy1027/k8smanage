package models

// 集群页面管理使用数据
type CloudClusterDetail struct {
	CloudCluster
	ClusterResources
	ClusterMem     int64
	ClusterCpu     int64
	ClusterNode    int64
	ClusterService int64
	ClusterPods    int
	Services       int
	Couters        int
	Health         string
}

type ClusterResources struct {
	UsedCpu       int64
	UsedMem       int64
	CpuUsePercent float64
	MemUsePercent float64
	MemFree       int64
	CpuFree       int64
	Cpu           int64
	Mmeory        int64
	Services      int
}

type CloudCluster struct {
	//docker安装路径
	DockerInstallDir string
	//内网网卡名称
	NetworkCart string
	//
	ClusterId int64
	//集群显示名称
	ClusterAlias string
	//最近修改时间
	LastModifyTime string
	//docker版本
	DockerVersion string
	//集群类型
	ClusterType string
	//集群名称,必须英文
	ClusterName string
	//创建时间
	CreateTime string
	//创建用户
	CreateUser string
	// ca证书公钥文件
	CaData string
	// node证书公钥内容
	CertData string
	// node证书私钥内容
	KeyData string
	// 主节点地址
	ApiAddress string
}

type ClusterStatus struct {
	ClusterId   int64
	ClusterType string
	NodeStatus
	ClusterAlias string
	ClusterName  string
	Nodes        int64
	Services     int
}

type NodeStatus struct {
	CloudClusterHosts
	Lables     []string
	K8sVersion string
	ErrorMsg   string
	MemSize    int64
}

type CloudClusterHosts struct {
	//主机IP
	HostIp string
	//主机标签
	HostLabel string
	//状态
	Status string
	//容器数量
	ContainerNum int64
	//cpu剩余量
	CpuFree string
	//pod数量
	PodNum int
	//内存大小
	MemSize string
	//创建用户
	CreateUser string
	//最近修改时间
	LastModifyUser string
	//内存剩余量
	MemFree string
	//主机类型
	HostType string
	//cpu使用百分比
	CpuPercent string
	//是否有效
	IsValid int64
	//内使用百分比
	MemPercent string
	//
	HostId int64
	//创建方法
	CreateMethod string
	//cpu数量
	CpuNum int64
	//创建时间
	CreateTime string
	//最近修改时间
	LastModifyTime string
	// 集群名称
	ClusterName string
	// k8sAPi端口,只需要master有就行了
	ApiPort string
	// 镜像数量
	ImageNum int
}
