package k8s

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// clusterKubectlCommand executes a one-shot kubectl command against a cluster.
// The arguments are passed through to kubectl as-is (raw args), so kubectl
// arguments containing "=" (e.g. -l app=web) are not interpreted by scw.
// Only "region" is interpreted by scw; the first other non-dash argument is
// the cluster ID, everything after is passed to kubectl.
//
// Autocompletion of the kubectl arguments is forwarded to kubectl's own
// completion mechanism (kubectl __complete).
//
// ponytail: kubectl __complete runs without KUBECONFIG, so verbs, flags and
// resource kinds complete but resource names (pods, namespaces) do not.
// Upgrade path: pass KUBECONFIG=<temp> from k8sKubeconfigGetRun, cached in
// core's autoCompleteCache. Also: no cluster ID completion at position 0,
// and no completion after a kubectl-level "--" (exec -it pod -- <TAB>).
func clusterKubectlCommand() *core.Command {
	return &core.Command{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "kubectl",
		Short:     `Execute a one-shot kubectl command against a cluster`,
		Long: `Execute a one-shot kubectl command against a Kubernetes Kapsule cluster.

The cluster kubeconfig is downloaded to a temporary file (removed on exit)
and passed to kubectl through the KUBECONFIG environment variable.

The first argument is the cluster ID, followed by the arguments to pass to kubectl.
Only the "region" argument is interpreted by scw, every other argument is passed
to kubectl as-is. If a kubectl argument starts with a dash, prefix the kubectl
arguments with "--":

    scw k8s cluster kubectl 1111... get nodes
    scw k8s cluster kubectl 1111... -- get pods -o wide
    scw k8s cluster kubectl 1111... -- get pods -l app=web`,
		ArgsType:         reflect.TypeFor[args.RawArgs](),
		ArgSpecs:         core.ArgSpecs{},
		Run:              clusterKubectlRun,
		AutoCompleteFunc: kubectlAutoCompleteFunc,
		Examples: []*core.Example{
			{
				Short: "List the nodes of a cluster",
				Raw:   `scw k8s cluster kubectl 11111111-1111-1111-1111-111111111111 get nodes`,
			},
			{
				Short: "List pods in wide format (note the -- separator before kubectl flags)",
				Raw:   `scw k8s cluster kubectl 11111111-1111-1111-1111-111111111111 -- get pods -o wide`,
			},
			{
				Short: "List nodes of a cluster in a specific region",
				Raw:   `scw k8s cluster kubectl 11111111-1111-1111-1111-111111111111 region=fr-par get nodes`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw k8s cluster connect",
				Short:   "Connect to a cluster with an interactive tool (k9s by default)",
			},
			{
				Command: "scw k8s kubeconfig install",
				Short:   "Install the cluster kubeconfig locally to use kubectl directly",
			},
		},
		ExcludeFromMCP: true,
	}
}

func clusterKubectlRun(ctx context.Context, argsI any) (any, error) {
	rawArgs := *argsI.(*args.RawArgs)

	client := core.ExtractClient(ctx)
	region, _ := client.GetDefaultRegion()

	var clusterID string
	kubectlArgs := []string(nil)
	for _, arg := range rawArgs {
		if value, found := strings.CutPrefix(arg, "region="); found {
			r, err := scw.ParseRegion(value)
			if err != nil {
				return nil, &core.CliError{
					Err:  fmt.Errorf("invalid region %s", value),
					Hint: "Region format should look like XX-XXX (e.g. fr-par)",
				}
			}
			region = r

			continue
		}
		if clusterID == "" && !strings.HasPrefix(arg, "-") {
			clusterID = arg

			continue
		}
		kubectlArgs = append(kubectlArgs, arg)
	}

	if clusterID == "" {
		return nil, &core.CliError{
			Err:  errors.New("a cluster ID is required for this command"),
			Hint: "Try running: scw k8s cluster kubectl <cluster-id> get pods",
		}
	}
	if len(kubectlArgs) == 0 {
		return nil, &core.CliError{
			Err:  errors.New("at least one kubectl argument is required after the cluster ID"),
			Hint: "Try running: scw k8s cluster kubectl <cluster-id> get pods",
		}
	}

	kubeconfigRes, err := k8sKubeconfigGetRun(ctx, &k8sKubeconfigGetRequest{
		ClusterID:  clusterID,
		Region:     region,
		AuthMethod: authMethods(defaultAuthMethod),
	})
	if err != nil {
		return nil, err
	}

	cmdArgs := append([]string{"kubectl"}, kubectlArgs...)
	err = runWithKubeconfig(ctx, kubeconfigRes.(string), func(string) ([]string, error) {
		return cmdArgs, nil
	})
	if err != nil {
		return nil, err
	}

	return &core.SuccessResult{Empty: true}, nil
}

// kubectlAutoCompleteFunc forwards completion of the kubectl arguments to
// kubectl's own completion mechanism: kubectl is cobra-based and exposes a
// hidden "__complete" command which prints the matching candidates (one per
// line, optionally followed by a tab and a description) and, last, a
// ":<directive>" line.
func kubectlAutoCompleteFunc(
	ctx context.Context,
	prefix string,
	_ any,
) core.AutocompleteSuggestions {
	typedArgs := core.ExtractTypedArgs(ctx)
	if len(typedArgs) == 0 {
		return nil
	}

	// Rebuild the kubectl command line from the typed words: before "--"
	// only plain words are kubectl args (dash and equals words are scw flags
	// and args, e.g. -o json, region=fr-par); after "--" everything is
	// kubectl. The first plain word is the cluster ID, not a kubectl arg.
	kubectlArgs := []string(nil)
	clusterIDSeen := false
	afterSeparator := false
	for _, word := range typedArgs {
		if afterSeparator {
			kubectlArgs = append(kubectlArgs, word)

			continue
		}
		switch {
		case word == "--":
			afterSeparator = true
		case strings.HasPrefix(word, "-") || strings.Contains(word, "="):
			// scw flag or arg
		case !clusterIDSeen:
			clusterIDSeen = true
		default:
			kubectlArgs = append(kubectlArgs, word)
		}
	}

	if !clusterIDSeen {
		return nil
	}

	execCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(
		execCtx,
		"kubectl",
		append(append([]string{"__complete"}, kubectlArgs...), prefix)...,
	).Output()
	if err != nil {
		// kubectl missing, timed out or errored: no suggestions.
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// The last line is the completion directive (":0", ":4", ...): drop it.
	if len(lines) > 0 && strings.HasPrefix(lines[len(lines)-1], ":") {
		lines = lines[:len(lines)-1]
	}

	suggestions := make(core.AutocompleteSuggestions, 0, len(lines))
	for _, line := range lines {
		// Candidates may carry a "\t<description>" suffix: strip it.
		if tab := strings.IndexByte(line, '\t'); tab >= 0 {
			line = line[:tab]
		}
		suggestions = append(suggestions, line)
	}

	return suggestions
}
