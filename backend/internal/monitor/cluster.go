package monitor

import (
	"context"
	"k8s-mon/internal/models"
	"time"
)

func (m *Monitor) GetClusterSummary(ctx context.Context) (*models.Cluster, error) {
	nodes, _ := m.GetNodes(ctx)
	pods, _ := m.GetPods(ctx, "")
	services, _ := m.GetServices(ctx)
	deployments, _ := m.GetDeployments(ctx)

	nodeSummary := models.NodeSummary{
		Total: len(nodes),
	}
	for _, node := range nodes {
		if node.Status == "Ready" {
			nodeSummary.Ready++
		} else {
			nodeSummary.NotReady++
		}
	}

	podSummary := models.PodSummary{
		Total: len(pods),
	}
	for _, pod := range pods {
		switch pod.Status {
		case "Running":
			podSummary.Running++
		case "Pending":
			podSummary.Pending++
		case "Failed":
			podSummary.Failed++
		}
	}

	deploymentSummary := models.DeploymentSummary{
		Total: len(deployments),
	}
	for _, deployment := range deployments {
		if deployment.ReadyReplicas == deployment.DesiredReplicas {
			deploymentSummary.Ready++
		} else {
			deploymentSummary.NotReady++
		}
	}

	serviceSummary := models.ServiceSummary{
		Total: len(services),
	}
	for _, service := range services {
		switch service.Type {
		case "ClusterIP":
			serviceSummary.ClusterIP++
		case "LoadBalancer":
			serviceSummary.LoadBalancer++
		case "NodePort":
			serviceSummary.NodePort++
		}
	}

	summary := &models.Cluster{
		Nodes:       nodeSummary,
		Pods:        podSummary,
		Deployments: deploymentSummary,
		Services:    serviceSummary,
		Timestamp:   time.Now().Unix(),
	}

	return summary, nil
}
