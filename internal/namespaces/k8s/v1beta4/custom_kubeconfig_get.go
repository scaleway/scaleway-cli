package k8s

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

func k8sKubeconfigCommand() *core.Command {
	return &core.Command{
		Short:     `Manage your Kubernetes Kapsule cluster's kubeconfig files`,
		Namespace: "k8s",
		Resource:  "kubeconfig",
	}
}
