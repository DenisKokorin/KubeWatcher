package monitor

import (
	"context"
	"k8s-mon/internal/models"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Monitor) GetDeployments(ctx context.Context) ([]models.Deployment, error) {
	if !m.enabled() {
		return []models.Deployment{}, nil
	}

	deployments, err := m.clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentInfos []models.Deployment
	for _, deployment := range deployments.Items {
		strategy := "RollingUpdate"
		if deployment.Spec.Strategy.Type == appsv1.RecreateDeploymentStrategyType {
			strategy = "Recreate"
		}

		deploymentInfo := models.Deployment{
			Name:              deployment.Name,
			Namespace:         deployment.Namespace,
			ReadyReplicas:     deployment.Status.ReadyReplicas,
			DesiredReplicas:   *deployment.Spec.Replicas,
			AvailableReplicas: deployment.Status.AvailableReplicas,
			Strategy:          strategy,
			CreatedAt:         deployment.CreationTimestamp.Time,
		}

		deploymentInfos = append(deploymentInfos, deploymentInfo)
	}

	return deploymentInfos, nil
}
