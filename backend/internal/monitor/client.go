package monitor

import (
	"k8s.io/client-go/kubernetes"
)

type Monitor struct {
	clientset *kubernetes.Clientset
}

func NewMonitor(clientset *kubernetes.Clientset) *Monitor {
	return &Monitor{
		clientset: clientset,
	}
}

func (m *Monitor) enabled() bool {
	return m != nil && m.clientset != nil
}
