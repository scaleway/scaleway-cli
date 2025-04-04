/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

// NOTE: This Config type definition is copied from k8s.io/client-go/tools/clientcmd/api/v1/types.go
//       for parsing the kube config yaml. The "k8s.io/apimachinery/pkg/runtime" dependency has
//       been removed.

// Where possible, json tags match the cli argument names.
// Top level config objects and all values required for proper functioning are not "omitempty".  Any truly optional piece of config is allowed to be omitted.

// Config holds the information needed to build connect to remote kubernetes clusters as a given user
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Config struct {
	// Legacy field from pkg/api/types.go TypeMeta.
	// TODO(jlowdermilk): remove this after eliminating downstream dependencies.
	// +k8s:conversion-gen=false
	// +optional
	Kind string `json:"kind,omitempty"`
	// Legacy field from pkg/api/types.go TypeMeta.
	// TODO(jlowdermilk): remove this after eliminating downstream dependencies.
	// +k8s:conversion-gen=false
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`
	// Preferences holds general information to be use for cli interactions
	Preferences Preferences `json:"preferences"`
	// Clusters is a map of referencable names to cluster configs
	Clusters []NamedCluster `json:"clusters"`
	// AuthInfos is a map of referencable names to user configs
	AuthInfos []NamedAuthInfo `json:"users"`
	// Contexts is a map of referencable names to context configs
	Contexts []NamedContext `json:"contexts"`
	// CurrentContext is the name of the context that you would like to use by default
	CurrentContext string `json:"current-context"`
	// Extensions holds additional information. This is useful for extenders so that reads and writes don't clobber unknown fields
	// +optional
	Extensions []NamedExtension `json:"extensions,omitempty"`
}

type Preferences struct {
	Extensions []NamedExtension `json:"extensions,omitempty"`
	Colors     bool             `json:"colors,omitempty"`
}

// Cluster contains information about how to communicate with a kubernetes cluster
type Cluster struct {
	Server                   string           `json:"server"`
	TLSServerName            string           `json:"tls-server-name,omitempty"`
	CertificateAuthority     string           `json:"certificate-authority,omitempty"`
	ProxyURL                 string           `json:"proxy-url,omitempty"`
	CertificateAuthorityData []byte           `json:"certificate-authority-data,omitempty"`
	Extensions               []NamedExtension `json:"extensions,omitempty"`
	InsecureSkipTLSVerify    bool             `json:"insecure-skip-tls-verify,omitempty"`
	DisableCompression       bool             `json:"disable-compression,omitempty"`
}

// AuthInfo contains information that describes identity information.  This is use to tell the kubernetes cluster who you are.
type AuthInfo struct {
	// ClientCertificate is the path to a client cert file for TLS.
	// +optional
	ClientCertificate string `json:"client-certificate,omitempty"`
	// ClientCertificateData contains PEM-encoded data from a client cert file for TLS. Overrides ClientCertificate
	// +optional
	ClientCertificateData []byte `json:"client-certificate-data,omitempty"`
	// ClientKey is the path to a client key file for TLS.
	// +optional
	ClientKey string `json:"client-key,omitempty"`
	// ClientKeyData contains PEM-encoded data from a client key file for TLS. Overrides ClientKey
	// +optional
	ClientKeyData []byte `json:"client-key-data,omitempty" datapolicy:"security-key"`
	// Token is the bearer token for authentication to the kubernetes cluster.
	// +optional
	Token string `json:"token,omitempty" datapolicy:"token"`
	// TokenFile is a pointer to a file that contains a bearer token (as described above).  If both Token and TokenFile are present, Token takes precedence.
	// +optional
	TokenFile string `json:"tokenFile,omitempty"`
	// Impersonate is the username to impersonate.  The name matches the flag.
	// +optional
	Impersonate string `json:"as,omitempty"`
	// ImpersonateUID is the uid to impersonate.
	// +optional
	ImpersonateUID string `json:"as-uid,omitempty"`
	// ImpersonateGroups is the groups to impersonate.
	// +optional
	ImpersonateGroups []string `json:"as-groups,omitempty"`
	// ImpersonateUserExtra contains additional information for impersonated user.
	// +optional
	ImpersonateUserExtra map[string][]string `json:"as-user-extra,omitempty"`
	// Username is the username for basic authentication to the kubernetes cluster.
	// +optional
	Username string `json:"username,omitempty"`
	// Password is the password for basic authentication to the kubernetes cluster.
	// +optional
	Password string `json:"password,omitempty" datapolicy:"password"`
	// AuthProvider specifies a custom authentication plugin for the kubernetes cluster.
	// +optional
	AuthProvider *AuthProviderConfig `json:"auth-provider,omitempty"`
	// Exec specifies a custom exec-based authentication plugin for the kubernetes cluster.
	// +optional
	Exec *ExecConfig `json:"exec,omitempty"`
	// Extensions holds additional information. This is useful for extenders so that reads and writes don't clobber unknown fields
	// +optional
	Extensions []NamedExtension `json:"extensions,omitempty"`
}

// Context is a tuple of references to a cluster (how do I communicate with a kubernetes cluster), a user (how do I identify myself), and a namespace (what subset of resources do I want to work with)
type Context struct {
	// Cluster is the name of the cluster for this context
	Cluster string `json:"cluster"`
	// AuthInfo is the name of the authInfo for this context
	AuthInfo string `json:"user"`
	// Namespace is the default namespace to use on unspecified requests
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// Extensions holds additional information. This is useful for extenders so that reads and writes don't clobber unknown fields
	// +optional
	Extensions []NamedExtension `json:"extensions,omitempty"`
}

// NamedCluster relates nicknames to cluster information
type NamedCluster struct {
	// Name is the nickname for this Cluster
	Name string `json:"name"`
	// Cluster holds the cluster information
	Cluster Cluster `json:"cluster"`
}

// NamedContext relates nicknames to context information
type NamedContext struct {
	// Name is the nickname for this Context
	Name string `json:"name"`
	// Context holds the context information
	Context Context `json:"context"`
}

// NamedAuthInfo relates nicknames to auth information
type NamedAuthInfo struct {
	// Name is the nickname for this AuthInfo
	Name string `json:"name"`
	// AuthInfo holds the auth information
	AuthInfo AuthInfo `json:"user"`
}

// NamedExtension relates nicknames to extension information
type NamedExtension struct {
	Extension interface{} `json:"extension"`
	Name      string      `json:"name"`
}

// AuthProviderConfig holds the configuration for a specified auth provider.
type AuthProviderConfig struct {
	Config map[string]string `json:"config"`
	Name   string            `json:"name"`
}

// ExecConfig specifies a command to provide client credentials. The command is exec'd
// and outputs structured stdout holding credentials.
//
// See the client.authentication.k8s.io API group for specifications of the exact input
// and output format
type ExecConfig struct {
	Command            string              `json:"command"`
	APIVersion         string              `json:"apiVersion,omitempty"`
	InstallHint        string              `json:"installHint,omitempty"`
	InteractiveMode    ExecInteractiveMode `json:"interactiveMode,omitempty"`
	Args               []string            `json:"args"`
	Env                []ExecEnvVar        `json:"env"`
	ProvideClusterInfo bool                `json:"provideClusterInfo"`
}

// ExecEnvVar is used for setting environment variables when executing an exec-based
// credential plugin.
type ExecEnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ExecInteractiveMode is a string that describes an exec plugin's relationship with standard input.
type ExecInteractiveMode string

const (
	// NeverExecInteractiveMode declares that this exec plugin never needs to use standard
	// input, and therefore the exec plugin will be run regardless of whether standard input is
	// available for user input.
	NeverExecInteractiveMode ExecInteractiveMode = "Never"
	// IfAvailableExecInteractiveMode declares that this exec plugin would like to use standard input
	// if it is available, but can still operate if standard input is not available. Therefore, the
	// exec plugin will be run regardless of whether stdin is available for user input. If standard
	// input is available for user input, then it will be provided to this exec plugin.
	IfAvailableExecInteractiveMode ExecInteractiveMode = "IfAvailable"
	// AlwaysExecInteractiveMode declares that this exec plugin requires standard input in order to
	// run, and therefore the exec plugin will only be run if standard input is available for user
	// input. If standard input is not available for user input, then the exec plugin will not be run
	// and an error will be returned by the exec plugin runner.
	AlwaysExecInteractiveMode ExecInteractiveMode = "Always"
)
