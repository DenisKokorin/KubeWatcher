package monitor

import (
	"context"
	"k8s-mon/internal/models"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Monitor) GetNamespaces(ctx context.Context) ([]models.Namespace, error) {
	namespaces, err := m.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaceInfos []models.Namespace
	for _, ns := range namespaces.Items {
		namespaceInfo := models.Namespace{
			Name:      ns.Name,
			Status:    string(ns.Status.Phase),
			CreatedAt: ns.CreationTimestamp.Time,
		}
		namespaceInfos = append(namespaceInfos, namespaceInfo)
	}

	return namespaceInfos, nil
}
