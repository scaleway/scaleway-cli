package k8s

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

func k8sExecCredentialCommand() *core.Command {
	return &core.Command{
		Hidden:    true,
		Short:     `exec-credential is a kubectl plugin to communicate credentials to HTTP transports.`,
		Namespace: "k8s",
		Resource:  "exec-credential",
		ArgsType:  reflect.TypeOf(struct{}{}),
		ArgSpecs:  core.ArgSpecs{},
		Run:       k8sExecCredentialRun,

		// avoid calling checkAPIKey (Check if API Key is about to expire)
		DisableAfterChecks: true,
	}
}

func k8sExecCredentialRun(ctx context.Context, _ any) (i any, e error) {
	secretKey, err := extractSecretKey(ctx)
	if err != nil {
		return nil, err
	}

	if !validation.IsSecretKey(secretKey) {
		return nil, core.InvalidSecretKeyError(secretKey)
	}

	execCreds := ExecCredential{
		APIVersion: "client.authentication.k8s.io/v1",
		Kind:       "ExecCredential",
		Status: &ExecCredentialStatus{
			Token: secretKey,
		},
	}
	response, err := json.MarshalIndent(execCreds, "", "    ")
	if err != nil {
		return nil, err
	}

	return string(response), nil
}

// ExecCredential is used by exec-based plugins to communicate credentials to HTTP transports.
type ExecCredential struct {
	// APIVersion defines the versioned schema of this representation of an object.
	// Servers should convert recognized schemas to the latest internal value, and
	// may reject unrecognized values.
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	Kind string `json:"kind,omitempty"`

	// Status is filled in by the plugin and holds the credentials that the transport
	// should use to contact the API.
	Status *ExecCredentialStatus `json:"status,omitempty"`
}

// ExecCredentialStatus holds credentials for the transport to use.
type ExecCredentialStatus struct {
	// Token is a bearer token used by the client for request authentication.
	Token string `json:"token,omitempty"`
}
