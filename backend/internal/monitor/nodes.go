package monitor

import (
	"context"
	"k8s-mon/internal/models"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Monitor) GetNodes(ctx context.Context) ([]models.Node, error) {
	if !m.enabled() {
		return []models.Node{}, nil
	}

	nodes, err := m.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeInfos []models.Node
	for _, node := range nodes.Items {
		ready := "NotReady"
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady {
				if condition.Status == v1.ConditionTrue {
					ready = "Ready"
				}
				break
			}
		}

		cpu := node.Status.Capacity[v1.ResourceCPU]
		memory := node.Status.Capacity[v1.ResourceMemory]
		pods := node.Status.Capacity[v1.ResourcePods]

		nodeInfo := models.Node{
			Name:           node.Name,
			Status:         ready,
			OS:             node.Status.NodeInfo.OSImage,
			KubeletVersion: node.Status.NodeInfo.KubeletVersion,
			Architecture:   node.Status.NodeInfo.Architecture,
			CPU:            cpu.String(),
			Memory:         memory.String(),
			Pods:           pods.String(),
			CreatedAt:      node.CreationTimestamp.Time,
		}

		nodeInfos = append(nodeInfos, nodeInfo)
	}

	return nodeInfos, nil
}
