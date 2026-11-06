package docs

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
)

type CommandHierarchy struct {
	Namespaces map[string]*NamespaceNode
}

type NamespaceNode struct {
	Command  *core.Command
	Children map[string]*ResourceNode
}

type ResourceNode struct {
	Command *core.Command
	Verbs   map[string]*core.Command
}

func BuildHierarchy(commands *core.Commands) *CommandHierarchy {
	hierarchy := &CommandHierarchy{
		Namespaces: make(map[string]*NamespaceNode),
	}

	for _, cmd := range commands.GetAll() {
		if cmd.Hidden {
			continue
		}

		ns := cmd.Namespace
		if ns == "" {
			continue
		}

		if hierarchy.Namespaces[ns] == nil {
			hierarchy.Namespaces[ns] = &NamespaceNode{
				Children: make(map[string]*ResourceNode),
			}
		}

		if cmd.Resource == "" {
			hierarchy.Namespaces[ns].Command = cmd

			continue
		}

		if hierarchy.Namespaces[ns].Children[cmd.Resource] == nil {
			hierarchy.Namespaces[ns].Children[cmd.Resource] = &ResourceNode{}
		}

		if cmd.Verb == "" {
			hierarchy.Namespaces[ns].Children[cmd.Resource].Command = cmd

			continue
		}

		if hierarchy.Namespaces[ns].Children[cmd.Resource].Verbs == nil {
			hierarchy.Namespaces[ns].Children[cmd.Resource].Verbs = make(map[string]*core.Command)
		}
		hierarchy.Namespaces[ns].Children[cmd.Resource].Verbs[cmd.Verb] = cmd
	}

	return hierarchy
}

func RenderAllNamespaces(commands *core.Commands) string {
	hierarchy := BuildHierarchy(commands)

	namespaces := make([]string, 0, len(hierarchy.Namespaces))
	for ns := range hierarchy.Namespaces {
		namespaces = append(namespaces, ns)
	}
	sort.Strings(namespaces)

	var sb strings.Builder
	sb.WriteString("Available namespaces:\n\n")
	for _, ns := range namespaces {
		node := hierarchy.Namespaces[ns]
		short := ""
		if node.Command != nil {
			short = node.Command.Short
		}
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", ns, short))
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

func RenderNamespacePage(commands *core.Commands, namespace string) string {
	hierarchy := BuildHierarchy(commands)
	node := hierarchy.Namespaces[namespace]
	if node == nil {
		return "Namespace \"" + namespace + "\" not found."
	}

	var sb strings.Builder

	if node.Command != nil {
		if node.Command.Long != "" {
			sb.WriteString(node.Command.Long + "\n\n")
		} else if node.Command.Short != "" {
			sb.WriteString(node.Command.Short + "\n\n")
		}
	}

	sb.WriteString("Resources:\n\n")

	resources := make([]string, 0, len(node.Children))
	for res := range node.Children {
		resources = append(resources, res)
	}
	sort.Strings(resources)

	for _, res := range resources {
		resNode := node.Children[res]
		short := ""
		if resNode.Command != nil {
			short = resNode.Command.Short
		}
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", res, short))
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

func RenderResourcePage(commands *core.Commands, namespace, resource string) string {
	hierarchy := BuildHierarchy(commands)
	node := hierarchy.Namespaces[namespace]
	if node == nil {
		return "Namespace \"" + namespace + "\" not found."
	}

	resNode := node.Children[resource]
	if resNode == nil {
		return "Resource \"" + resource + "\" not found in namespace \"" + namespace + "\"."
	}

	var sb strings.Builder

	if resNode.Command != nil {
		if resNode.Command.Long != "" {
			sb.WriteString(resNode.Command.Long + "\n\n")
		} else if resNode.Command.Short != "" {
			sb.WriteString(resNode.Command.Short + "\n\n")
		}
	}

	if len(resNode.Verbs) > 0 {
		sb.WriteString("Commands:\n\n")

		verbs := make([]string, 0, len(resNode.Verbs))
		for v := range resNode.Verbs {
			verbs = append(verbs, v)
		}
		sort.Strings(verbs)

		for _, v := range verbs {
			cmd := resNode.Verbs[v]
			sb.WriteString(fmt.Sprintf("  %-20s %s\n", v, cmd.Short))
		}
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

func RenderCommandPage(ctx context.Context, cmd *core.Command, commands *core.Commands) string {
	var sb strings.Builder

	if cmd.Short != "" {
		sb.WriteString(cmd.Short + "\n\n")
	}

	if cmd.Long != "" {
		sb.WriteString(cmd.Long + "\n\n")
	}

	usage := cmd.GetUsage("scw", commands)
	if usage != "" {
		sb.WriteString("Usage:\n")
		sb.WriteString("  " + usage + "\n\n")
	}

	args := core.BuildUsageArgs(ctx, cmd, false)
	if args != "" {
		sb.WriteString("Arguments:\n")
		sb.WriteString(args + "\n\n")
	}

	if len(cmd.Examples) > 0 {
		sb.WriteString("Examples:\n")
		for _, example := range cmd.Examples {
			sb.WriteString("  " + example.Short + "\n")
			cmdLine := example.GetCommandLine("scw", cmd)
			cmdLine = indent(cmdLine, 4)
			sb.WriteString(cmdLine + "\n\n")
		}
	}

	if len(cmd.SeeAlsos) > 0 {
		sb.WriteString("See Also:\n")
		for _, sa := range cmd.SeeAlsos {
			sb.WriteString(fmt.Sprintf("  %-20s %s\n", sa.Command, sa.Short))
		}
		sb.WriteString("\n")
	}

	tips := GetTips(cmd.DebugString())
	if len(tips) > 0 {
		sb.WriteString("Tips:\n")
		for _, tip := range tips {
			sb.WriteString("  - " + tip + "\n")
		}
		sb.WriteString("\n")
	}

	pitfalls := GetPitfalls(cmd.DebugString())
	if len(pitfalls) > 0 {
		sb.WriteString("Common Pitfalls:\n")
		for _, p := range pitfalls {
			sb.WriteString("  - " + p + "\n")
		}
		sb.WriteString("\n")
	}

	relatedTutorials := GetRelatedTutorials(cmd.DebugString())
	if len(relatedTutorials) > 0 {
		sb.WriteString("Related Tutorials:\n")
		for _, t := range relatedTutorials {
			sb.WriteString("  - scw tutorial " + t + "\n")
		}
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

func RenderHelpTopicPage(cmd *core.Command) string {
	var sb strings.Builder

	if cmd.Short != "" {
		sb.WriteString(cmd.Short + "\n\n")
	}

	if cmd.Long != "" {
		sb.WriteString(cmd.Long)
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

func indent(str string, n int) string {
	padding := strings.Repeat(" ", n)
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = padding + line
		}
	}

	return strings.Join(lines, "\n")
}
