package k8s

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/core"
	api "github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1/types"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func k8sKubeconfigCommand() *core.Command {
	return &core.Command{
		Short:     `Manage your Kubernetes Kapsule cluster's kubeconfig files`,
		Namespace: "k8s",
		Resource:  "kubeconfig",
	}
}

type authMethods string

const (
	// Deprecated: authMethodLegacy uses the legacy admin token in the kubeconfig
	authMethodLegacy = authMethods("legacy")
	// authMethodCLI uses the scaleway-cli as credential provider in the kubeconfig
	authMethodCLI = authMethods("cli")
	// authMethodCopyCLIToken copies the scw_secret_key used for the command to the kubeconfig
	authMethodCopyCLIToken = authMethods("copy-cli-token")

	KubeconfigAPIVersion = "v1"
	KubeconfigKind       = "Config"
	InstallHint          = `This kubeconfig profile require scaleway-cli (scw) to authenticate.
Installation instruction: https://github.com/scaleway/scaleway-cli#installation`
)

var (
	// enumAuthMethods lists choices for auth method
	enumAuthMethods = []string{
		string(authMethodCLI),
		string(authMethodCopyCLIToken),
		string(authMethodLegacy),
	}

	// defaultAuthMethod is the default choice
	defaultAuthMethod = string(authMethodCLI)
)

// get the path to the wanted kubeconfig on disk
// either the file pointed by the KUBECONFIG env variable (first one in case of a list)
// or the $HOME/.kube/config
func getKubeconfigPath(ctx context.Context) string {
	var kubeconfigPath string
	kubeconfigEnv := core.ExtractEnv(ctx, "KUBECONFIG")
	if kubeconfigEnv != "" {
		if runtime.GOOS == "windows" {
			kubeconfigPath = strings.Split(kubeconfigEnv, ";")[0] // list is separated by ; on windows
		} else {
			kubeconfigPath = strings.Split(kubeconfigEnv, ":")[0] // list is separated by : on linux/macos
		}
	} else {
		kubeconfigPath = path.Join(core.ExtractUserHomeDir(ctx), kubeLocationDir, "config")
	}

	return kubeconfigPath
}

func openAndUnmarshalKubeconfig(kubeconfigPath string) (*api.Config, error) {
	// getting the existing file
	file, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	var kubeconfig api.Config
	err = yaml.Unmarshal(file, &kubeconfig)
	if err != nil {
		return nil, err
	}

	return &kubeconfig, nil
}

func marshalAndWriteKubeconfig(kubeconfig *api.Config, kubeconfigPath string) error {
	newKubeconfig, err := yaml.Marshal(*kubeconfig)
	if err != nil {
		return err
	}

	return os.WriteFile(kubeconfigPath, newKubeconfig, 0o600)
}

func generateNamedAuthInfo(ctx context.Context, method authMethods) (*api.NamedAuthInfo, error) {
	switch method {
	case authMethodCopyCLIToken:
		return generateTokenNamedAuthInfo(ctx)
	case authMethodCLI:
		return generateCLINamedAuthInfo(ctx)
	case authMethodLegacy:
		return nil, errors.New("legacy method is not handled by generateNamedAuthInfo")
	}

	return nil, fmt.Errorf("unknown auth method: %s", method)
}

func generateCLINamedAuthInfo(ctx context.Context) (*api.NamedAuthInfo, error) {
	// if SCW_SECRET_KEY is set, fallback to token method to avoid strange behavior
	if core.ExtractEnv(ctx, scw.ScwSecretKeyEnv) != "" {
		core.ExtractLogger(
			ctx,
		).Warningf("Environment variable %s is set, generated kubeconfig will hardcode the token.\n", scw.ScwSecretKeyEnv)

		return generateTokenNamedAuthInfo(ctx)
	}

	keys := []string{}
	args := []string{}
	if profileFlag := core.ExtractProfileFlag(ctx); profileFlag != "" {
		args = append(args, "--profile", profileFlag)
		keys = append(keys, "--profile", profileFlag)
	}

	if configPath := core.ExtractConfigPathFlag(ctx); configPath != "" {
		args = append(args, "--config", configPath)
		keys = append(keys, "--config", configPath)
	}

	execEnvVars := []api.ExecEnvVar{}
	if profileEnv := core.ExtractEnv(ctx, scw.ScwActiveProfileEnv); profileEnv != "" {
		execEnvVars = append(execEnvVars, api.ExecEnvVar{
			Name:  scw.ScwActiveProfileEnv,
			Value: profileEnv,
		})
		keys = append(keys, scw.ScwActiveProfileEnv, profileEnv)
	}
	if configEnv := core.ExtractEnv(ctx, scw.ScwConfigPathEnv); configEnv != "" {
		execEnvVars = append(execEnvVars, api.ExecEnvVar{
			Name:  scw.ScwConfigPathEnv,
			Value: configEnv,
		})
		keys = append(keys, scw.ScwConfigPathEnv, configEnv)
	}

	configPathSum := crc32.ChecksumIEEE([]byte(strings.Join(keys, ",")))

	return &api.NamedAuthInfo{
		Name: fmt.Sprintf("cli-config-%08x", configPathSum),
		AuthInfo: api.AuthInfo{
			Exec: &api.ExecConfig{
				APIVersion: "client.authentication.k8s.io/v1",
				Command:    core.ExtractBinaryName(ctx),
				Args: append(args,
					"k8s",
					"exec-credential",
				),
				Env:             execEnvVars,
				InstallHint:     InstallHint,
				InteractiveMode: api.NeverExecInteractiveMode,
			},
		},
	}, nil
}

func generateTokenNamedAuthInfo(ctx context.Context) (*api.NamedAuthInfo, error) {
	secretKey, ok := core.ExtractClient(ctx).GetSecretKey()
	if !ok {
		return nil, errors.New("could not get secret key")
	}

	configPathSum := crc32.ChecksumIEEE([]byte(secretKey))

	return &api.NamedAuthInfo{
		Name: fmt.Sprintf("cli-token-%08x", configPathSum),
		AuthInfo: api.AuthInfo{
			Token: secretKey,
		},
	}, nil
}

type KubeMapConfig struct {
	preferences    api.Preferences
	clusters       map[string]api.Cluster
	users          map[string]api.AuthInfo
	contexts       map[string]api.Context
	CurrentContext string
	extensions     map[string]any
}

func NewKubeMapConfig() *KubeMapConfig {
	return &KubeMapConfig{
		clusters:   map[string]api.Cluster{},
		users:      map[string]api.AuthInfo{},
		contexts:   map[string]api.Context{},
		extensions: map[string]any{},
	}
}

func LoadKubeMapConfig(ctx context.Context, kubeconfigPath string) (*KubeMapConfig, error) {
	file, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	var kubeconfig api.Config
	err = yaml.Unmarshal(file, &kubeconfig)
	if err != nil {
		return nil, err
	}

	kubeMapConfig := &KubeMapConfig{
		preferences:    kubeconfig.Preferences,
		clusters:       map[string]api.Cluster{},
		users:          map[string]api.AuthInfo{},
		contexts:       map[string]api.Context{},
		CurrentContext: kubeconfig.CurrentContext,
		extensions:     map[string]any{},
	}

	for _, namedCluster := range kubeconfig.Clusters {
		if _, ok := kubeMapConfig.clusters[namedCluster.Name]; ok {
			return nil, fmt.Errorf("duplicated cluster '%s' found in kubeconfig", namedCluster.Name)
		}
		kubeMapConfig.clusters[namedCluster.Name] = namedCluster.Cluster
	}

	for _, namedAuthInfo := range kubeconfig.AuthInfos {
		if _, ok := kubeMapConfig.users[namedAuthInfo.Name]; ok {
			return nil, fmt.Errorf("duplicated user '%s' found in kubeconfig", namedAuthInfo.Name)
		}
		kubeMapConfig.users[namedAuthInfo.Name] = namedAuthInfo.AuthInfo
	}

	for _, namedContext := range kubeconfig.Contexts {
		if _, ok := kubeMapConfig.contexts[namedContext.Name]; ok {
			return nil, fmt.Errorf("duplicated context '%s' found in kubeconfig", namedContext.Name)
		}

		// Warn the user about its invalid kubeconfig
		if _, ok := kubeMapConfig.clusters[namedContext.Context.Cluster]; !ok {
			core.ExtractLogger(ctx).
				Warningf("context '%s' refers to cluster '%s' that does not exist", namedContext.Name, namedContext.Context.Cluster)
		}

		if _, ok := kubeMapConfig.users[namedContext.Context.AuthInfo]; !ok {
			core.ExtractLogger(ctx).
				Warningf("context '%s' refers to user '%s' that does not exist", namedContext.Name, namedContext.Context.AuthInfo)
		}

		kubeMapConfig.contexts[namedContext.Name] = namedContext.Context
	}

	for _, namedExtension := range kubeconfig.Extensions {
		if _, ok := kubeMapConfig.extensions[namedExtension.Name]; ok {
			return nil, fmt.Errorf(
				"duplicated extension '%s' found in kubeconfig",
				namedExtension.Name,
			)
		}
		kubeMapConfig.extensions[namedExtension.Name] = namedExtension.Extension
	}

	return kubeMapConfig, nil
}

func (c *KubeMapConfig) GetCluster(name string) (*api.NamedCluster, error) {
	if cluster, ok := c.clusters[name]; ok {
		return &api.NamedCluster{Name: name, Cluster: cluster}, nil
	}

	return nil, fmt.Errorf("cluster '%s' not found", name)
}

func (c *KubeMapConfig) SetCluster(name string, cluster api.Cluster, overwrite bool) error {
	if _, ok := c.clusters[name]; ok && !overwrite {
		return fmt.Errorf("duplicated cluster '%s' found in kubeconfig", name)
	}
	c.clusters[name] = cluster

	return nil
}

func (c *KubeMapConfig) RemoveCluster(name string) error {
	if _, ok := c.clusters[name]; !ok {
		return errors.New("cluster not found")
	}

	for contextName, contextValue := range c.contexts {
		if contextValue.Cluster == name {
			return fmt.Errorf(
				"unable to remove cluster: cluster '%s' still referenced in context '%s'",
				name,
				contextName,
			)
		}
	}

	delete(c.clusters, name)

	return nil
}

func (c *KubeMapConfig) GetUser(name string) (*api.NamedAuthInfo, error) {
	if user, ok := c.users[name]; ok {
		return &api.NamedAuthInfo{Name: name, AuthInfo: user}, nil
	}

	return nil, fmt.Errorf("user '%s' not found", name)
}

func (c *KubeMapConfig) SetUser(name string, user api.AuthInfo, overwrite bool) error {
	if _, ok := c.users[name]; ok && !overwrite {
		return fmt.Errorf("duplicated user '%s' found in kubeconfig", name)
	}
	c.users[name] = user

	return nil
}

func (c *KubeMapConfig) RemoveUser(name string) error {
	if _, ok := c.users[name]; !ok {
		return errors.New("user not found")
	}

	for contextName, contextValue := range c.contexts {
		if contextValue.AuthInfo == name {
			return fmt.Errorf(
				"unable to remove user: user '%s' referenced in context '%s'",
				name,
				contextName,
			)
		}
	}

	delete(c.users, name)

	return nil
}

func (c *KubeMapConfig) GetContext(name string) (*api.NamedContext, error) {
	if kubeContext, ok := c.contexts[name]; ok {
		return &api.NamedContext{Name: name, Context: kubeContext}, nil
	}

	return nil, fmt.Errorf("context '%s' not found", name)
}

func (c *KubeMapConfig) SetContext(name string, context api.Context, overwrite bool) error {
	if _, ok := c.contexts[name]; ok && !overwrite {
		return fmt.Errorf("duplicated context '%s' found in kubeconfig", name)
	}

	if _, ok := c.clusters[context.Cluster]; !ok {
		return fmt.Errorf("cluster '%s' not found in kubeconfig", context.Cluster)
	}

	if _, ok := c.users[context.AuthInfo]; !ok {
		return fmt.Errorf("user '%s' not found in kubeconfig", context.AuthInfo)
	}

	c.contexts[name] = context

	return nil
}

func (c *KubeMapConfig) RemoveContext(name string) error {
	if _, ok := c.contexts[name]; ok {
		delete(c.contexts, name)

		return nil
	}

	return errors.New("context not found")
}

func (c *KubeMapConfig) Kubeconfig() api.Config {
	resultingKubeconfig := api.Config{
		APIVersion:     KubeconfigAPIVersion,
		Kind:           KubeconfigKind,
		Preferences:    c.preferences,
		CurrentContext: c.CurrentContext,
	}

	for name, kubeCluster := range c.clusters {
		resultingKubeconfig.Clusters = append(resultingKubeconfig.Clusters, api.NamedCluster{
			Name:    name,
			Cluster: kubeCluster,
		})
	}

	for name, kubeUser := range c.users {
		resultingKubeconfig.AuthInfos = append(resultingKubeconfig.AuthInfos, api.NamedAuthInfo{
			Name:     name,
			AuthInfo: kubeUser,
		})
	}

	for name, kubeContext := range c.contexts {
		resultingKubeconfig.Contexts = append(resultingKubeconfig.Contexts, api.NamedContext{
			Name:    name,
			Context: kubeContext,
		})
	}

	for name, kubeExtension := range c.extensions {
		resultingKubeconfig.Extensions = append(resultingKubeconfig.Extensions, api.NamedExtension{
			Name:      name,
			Extension: kubeExtension,
		})
	}

	return resultingKubeconfig
}

func (c *KubeMapConfig) Save(kubeconfigPath string) error {
	kubeconfigStruct := c.Kubeconfig()

	kubeconfigBytes, err := yaml.Marshal(kubeconfigStruct)
	if err != nil {
		return err
	}

	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		// make sure the directory exists
		err = os.MkdirAll(path.Dir(kubeconfigPath), 0o755)
		if err != nil {
			return err
		}
	}

	// create the file
	return os.WriteFile(kubeconfigPath, kubeconfigBytes, 0o600)
}
