package config

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func LoadKubeconfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
			kubeconfigPath = kubeconfig
		} else {
			if home := homedir.HomeDir(); home != "" {
				kubeconfigPath = filepath.Join(home, ".kube", "config")
			}
		}
	}

	if _, err := os.Stat(kubeconfigPath); err == nil {
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}

	return rest.InClusterConfig()
}
