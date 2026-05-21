package models

import "time"

type Cluster struct {
	Nodes       NodeSummary       `json:"nodes"`
	Pods        PodSummary        `json:"pods"`
	Deployments DeploymentSummary `json:"deployments"`
	Services    ServiceSummary    `json:"services"`
	Timestamp   int64             `json:"timestamp"`
}

type NodeSummary struct {
	Total    int `json:"total"`
	Ready    int `json:"ready"`
	NotReady int `json:"not_ready"`
}

type PodSummary struct {
	Total   int `json:"total"`
	Running int `json:"running"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
}

type DeploymentSummary struct {
	Total    int `json:"total"`
	Ready    int `json:"ready"`
	NotReady int `json:"not_ready"`
}

type ServiceSummary struct {
	Total        int `json:"total"`
	ClusterIP    int `json:"cluster_ip"`
	LoadBalancer int `json:"load_balancer"`
	NodePort     int `json:"node_port"`
}

type Node struct {
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	OS             string    `json:"os"`
	KubeletVersion string    `json:"kubelet_version"`
	Architecture   string    `json:"architecture"`
	CPU            string    `json:"cpu"`
	Memory         string    `json:"memory"`
	Pods           string    `json:"pods"`
	CreatedAt      time.Time `json:"created_at"`
}

type Pod struct {
	Name       string          `json:"name"`
	Namespace  string          `json:"namespace"`
	Status     string          `json:"status"`
	Node       string          `json:"node"`
	Restarts   int32           `json:"restarts"`
	IP         string          `json:"ip"`
	CreatedAt  time.Time       `json:"created_at"`
	Containers []ContainerInfo `json:"containers"`
}

type ContainerInfo struct {
	Name  string `json:"name"`
	Ready bool   `json:"ready"`
	Image string `json:"image"`
}

type Service struct {
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	Type       string   `json:"type"`
	ClusterIP  string   `json:"cluster_ip"`
	ExternalIP string   `json:"external_ip"`
	Ports      []string `json:"ports"`
}

type Deployment struct {
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace"`
	ReadyReplicas     int32     `json:"ready_replicas"`
	DesiredReplicas   int32     `json:"desired_replicas"`
	AvailableReplicas int32     `json:"available_replicas"`
	Strategy          string    `json:"strategy"`
	CreatedAt         time.Time `json:"created_at"`
}

type Namespace struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
