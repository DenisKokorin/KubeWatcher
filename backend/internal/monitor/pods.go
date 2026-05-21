package monitor

import (
	"context"
	"k8s-mon/internal/models"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Monitor) GetPods(ctx context.Context, namespace string) ([]models.Pod, error) {
	if !m.enabled() {
		return []models.Pod{}, nil
	}

	var pods *v1.PodList
	var err error

	if namespace != "" {
		pods, err = m.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	} else {
		pods, err = m.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	}

	if err != nil {
		return nil, err
	}

	var podInfos []models.Pod
	for _, pod := range pods.Items {
		status := string(pod.Status.Phase)

		var containers []models.ContainerInfo
		restarts := int32(0)

		for _, containerStatus := range pod.Status.ContainerStatuses {
			containers = append(containers, models.ContainerInfo{
				Name:  containerStatus.Name,
				Ready: containerStatus.Ready,
				Image: containerStatus.Image,
			})
			restarts += containerStatus.RestartCount
		}

		podInfo := models.Pod{
			Name:       pod.Name,
			Namespace:  pod.Namespace,
			Status:     status,
			Node:       pod.Spec.NodeName,
			Restarts:   restarts,
			IP:         pod.Status.PodIP,
			CreatedAt:  pod.CreationTimestamp.Time,
			Containers: containers,
		}

		podInfos = append(podInfos, podInfo)
	}

	return podInfos, nil
}
