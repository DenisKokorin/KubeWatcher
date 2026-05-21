package monitor

import (
	"context"
	"fmt"
	"k8s-mon/internal/models"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Monitor) GetServices(ctx context.Context) ([]models.Service, error) {
	if !m.enabled() {
		return []models.Service{}, nil
	}

	services, err := m.clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var serviceInfos []models.Service
	for _, service := range services.Items {
		// Форматируем порты
		var ports []string
		for _, port := range service.Spec.Ports {
			portStr := fmt.Sprintf("%d/%s", port.Port, port.Protocol)
			if port.NodePort > 0 {
				portStr = fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol)
			}
			ports = append(ports, portStr)
		}

		externalIP := ""
		if service.Spec.Type == v1.ServiceTypeLoadBalancer {
			for _, ingress := range service.Status.LoadBalancer.Ingress {
				if ingress.IP != "" {
					externalIP = ingress.IP
					break
				} else if ingress.Hostname != "" {
					externalIP = ingress.Hostname
					break
				}
			}
		}

		serviceInfo := models.Service{
			Name:       service.Name,
			Namespace:  service.Namespace,
			Type:       string(service.Spec.Type),
			ClusterIP:  service.Spec.ClusterIP,
			ExternalIP: externalIP,
			Ports:      ports,
		}

		serviceInfos = append(serviceInfos, serviceInfo)
	}

	return serviceInfos, nil
}
