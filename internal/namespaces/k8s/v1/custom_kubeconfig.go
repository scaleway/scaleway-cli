package k8s

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func k8sKubeconfigCommand() *core.Command {
	return &core.Command{
		Short:     `Manage your Kubernetes Kapsule cluster's kubeconfig files`,
		Namespace: "k8s",
		Resource:  "kubeconfig",
	}
}
