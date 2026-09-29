// Package manpage generates and installs roff man pages for the scw CLI.
//
// It follows the same pattern as internal/docgen: command metadata from core
// is rendered through an embedded text/template, producing one POSIX man page
// per command.
package manpage

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/core"
)

//go:embed man.tmpl
var manTemplate string

// BinaryName is the name of the CLI binary used in page names and synopses.
const BinaryName = "scw"

// SectionDir is the man section subdirectory pages are installed into.
const SectionDir = "man1"

// Option is a single flag or positional argument documented in the OPTIONS
// section of a man page.
type Option struct {
	// Term is the roff-safe term displayed for the option (e.g. "name=value"
	// or "<server-id>" for positional arguments).
	Term string

	// Description is the option description including annotations
	// (Required, Deprecated, Default, One of).
	Description string

	// Deprecated mirrors core.ArgSpec.Deprecated.
	Deprecated bool
}

// Example is a rendered usage example.
type Example struct {
	// Short is the example title.
	Short string

	// CmdLine is the rendered example command line.
	CmdLine string
}

// SeeAlso is a cross-reference to another command or resource.
type SeeAlso struct {
	// Line is the rendered roff line for the SEE ALSO section.
	Line string
}

// EnvVar is a documented environment variable.
type EnvVar struct {
	// Name is the environment variable name (e.g. SCW_ACCESS_KEY).
	Name string

	// Description of the environment variable.
	Description string
}

// ManPage holds the data used to render the man page of one command.
type ManPage struct {
	// PageName is the man page name (e.g. scw-instance-server-create).
	PageName string

	// FileName is the page file name (PageName + ".1").
	FileName string

	// Version is the CLI version displayed in the .TH header.
	Version string

	// Short is the command short description.
	Short string

	// UsageCommand is the bold part of the SYNOPSIS (e.g. "scw instance server create").
	UsageCommand string

	// UsageArgs is the remaining part of the SYNOPSIS (arguments, [arg=value ...]).
	UsageArgs string

	// Description is the command long description (ANSI stripped, roff escaped).
	Description string

	// Options are the command arguments and flags.
	Options []Option

	// Examples are the rendered usage examples.
	Examples []Example

	// SeeAlsos are the cross-references.
	SeeAlsos []SeeAlso

	// EnvVars are the documented global environment variables.
	EnvVars []EnvVar

	// Deprecated: indicates the command is deprecated.
	Deprecated bool
}

// NamespaceEntry is one line of the NAMESPACES section of the root page.
type NamespaceEntry struct {
	// Name is the namespace name.
	Name string

	// Short is the namespace short description.
	Short string
}

// RootPage holds the data used to render the top-level `scw` man page.
type RootPage struct {
	// Version is the CLI version displayed in the .TH header.
	Version string

	// SynopsisCommand is the bold part of the SYNOPSIS.
	SynopsisCommand string

	// SynopsisArgs is the remaining part of the SYNOPSIS.
	SynopsisArgs string

	// Description is the tool description.
	Description string

	// GlobalOptions are the global flags available on every command.
	GlobalOptions []Option

	// Namespaces are the available namespaces.
	Namespaces []NamespaceEntry

	// EnvVars are the documented global environment variables.
	EnvVars []EnvVar
}

// ansiRegexp matches ANSI escape sequences (same regex as internal/docgen).
var ansiRegexp = regexp.MustCompile(
	"[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))",
)

// PageName returns the man page name for a command path. An empty path
// refers to the top-level `scw` page.
func PageName(commandPath ...string) string {
	if len(commandPath) == 0 {
		return BinaryName
	}

	return BinaryName + "-" + strings.Join(commandPath, "-")
}

// FileName returns the man page file name for a page name.
func FileName(pageName string) string {
	return pageName + ".1"
}

// IsOwnedFileName reports whether a file in the man section directory is
// owned by the scw CLI (i.e. created by a previous `scw man install`).
func IsOwnedFileName(name string) bool {
	if name == FileName(PageName()) {
		return true
	}

	return strings.HasPrefix(name, BinaryName+"-") && strings.HasSuffix(name, ".1")
}

// DetectNameCollisions verifies that no two distinct command paths map to the
// same page file name. It returns the file names of all visible commands (plus
// the top-level page) or an error describing the collision.
func DetectNameCollisions(commands *core.Commands) (map[string]bool, error) {
	fileNames := map[string]bool{FileName(PageName()): true}
	commandByFileName := map[string]string{}

	for _, cmd := range uniqueVisibleCommands(commands) {
		fileName := FileName(PageName(commandPath(cmd)...))
		if commandByFileName[fileName] != "" {
			return nil, fmt.Errorf(
				"page name collision for %q and %q: %s",
				commandByFileName[fileName],
				cmd.GetCommandLine(BinaryName),
				fileName,
			)
		}
		commandByFileName[fileName] = cmd.GetCommandLine(BinaryName)
		fileNames[fileName] = true
	}

	return fileNames, nil
}

// uniqueVisibleCommands returns the visible commands, deduplicated by command
// path with the last registration winning. This mirrors the behavior of the
// command index used at runtime (core.Commands.Add overwrites earlier
// registrations of the same path).
func uniqueVisibleCommands(commands *core.Commands) []*core.Command {
	// path -> index in result, so the last registration of a path wins while
	// keeping the order of first appearance.
	indexByPath := map[string]int{}
	var result []*core.Command

	for _, cmd := range commands.GetAll() {
		if cmd.Hidden {
			continue
		}

		path := strings.Join(commandPath(cmd), ".")
		if i, found := indexByPath[path]; found {
			result[i] = cmd
		} else {
			indexByPath[path] = len(result)
			result = append(result, cmd)
		}
	}

	return result
}

// commandPath returns the non-empty path levels of a command.
func commandPath(cmd *core.Command) []string {
	path := make([]string, 0, 3)
	if cmd.Namespace != "" {
		path = append(path, cmd.Namespace)
	}
	if cmd.Resource != "" {
		path = append(path, cmd.Resource)
	}
	if cmd.Verb != "" {
		path = append(path, cmd.Verb)
	}

	return path
}

// BuildPages builds the man page data for every visible command plus the
// top-level `scw` page. It fails on page name collisions.
func BuildPages(ctx context.Context, commands *core.Commands) ([]*ManPage, *RootPage, error) {
	version := ""
	if buildInfo := core.ExtractBuildInfo(ctx); buildInfo != nil && buildInfo.Version != nil {
		version = buildInfo.Version.String()
	}

	_, err := DetectNameCollisions(commands)
	if err != nil {
		return nil, nil, err
	}

	pages := make([]*ManPage, 0, len(uniqueVisibleCommands(commands)))
	for _, cmd := range uniqueVisibleCommands(commands) {
		pages = append(pages, buildPage(ctx, commands, cmd, version))
	}

	root := buildRootPage(commands, version)

	return pages, root, nil
}

func buildPage(
	ctx context.Context,
	commands *core.Commands,
	cmd *core.Command,
	version string,
) *ManPage {
	path := commandPath(cmd)

	options := make([]Option, 0, len(cmd.ArgSpecs))
	for _, arg := range cmd.ArgSpecs {
		options = append(options, buildOption(ctx, arg))
	}

	examples := make([]Example, 0, len(cmd.Examples))
	for _, example := range cmd.Examples {
		examples = append(examples, Example{
			Short:   example.Short,
			CmdLine: roffEscape(example.GetCommandLine(BinaryName, cmd)),
		})
	}

	seeAlsos := make([]SeeAlso, 0, len(cmd.SeeAlsos))
	for _, seeAlso := range cmd.SeeAlsos {
		seeAlsos = append(seeAlsos, SeeAlso{Line: renderSeeAlso(seeAlso)})
	}

	usage := cmd.GetUsage(BinaryName, commands)
	usageCommand := cmd.GetCommandLine(BinaryName)
	usageArgs := strings.TrimPrefix(usage, usageCommand)

	return &ManPage{
		PageName:     PageName(path...),
		FileName:     FileName(PageName(path...)),
		Version:      version,
		Short:        roffEscape(cmd.Short),
		UsageCommand: usageCommand,
		UsageArgs:    usageArgs,
		Description:  roffEscape(longDescription(cmd)),
		Options:      options,
		Examples:     examples,
		SeeAlsos:     seeAlsos,
		EnvVars:      EnvVariables(),
		Deprecated:   cmd.Deprecated,
	}
}

// buildOption builds the OPTIONS entry for an argument spec.
func buildOption(ctx context.Context, arg *core.ArgSpec) Option {
	term := arg.Name + "=value"
	if arg.Positional {
		term = "<" + arg.Name + ">"
	}

	description := arg.Short
	if arg.Deprecated {
		description += " [Deprecated]"
	}
	if arg.Required {
		description += " [Required]"
	}
	if arg.Default != nil {
		if _, doc := arg.Default(ctx); doc != "" {
			description += fmt.Sprintf(" [Default: %s]", doc)
		}
	}
	if len(arg.EnumValues) > 0 {
		description += fmt.Sprintf(" [One of: %s]", strings.Join(arg.EnumValues, ", "))
	}

	return Option{
		Term:        roffEscape(term),
		Description: roffEscape(description),
		Deprecated:  arg.Deprecated,
	}
}

// longDescription returns the description used in the DESCRIPTION section.
func longDescription(cmd *core.Command) string {
	if cmd.Long != "" {
		return strings.Trim(cmd.Long, "\n")
	}

	return cmd.Short
}

// renderSeeAlso renders a SEE ALSO line. scw commands are referenced as man
// page cross-references, other entries are rendered as-is.
func renderSeeAlso(seeAlso *core.SeeAlso) string {
	command := strings.TrimSpace(seeAlso.Command)
	if rest, found := strings.CutPrefix(command, BinaryName+" "); found {
		pageName := PageName(strings.Split(rest, " ")...)

		return fmt.Sprintf("\\fI%s\\fP (1) \\- %s", pageName, roffEscape(seeAlso.Short))
	}

	return roffEscape(command)
}

func buildRootPage(commands *core.Commands, version string) *RootPage {
	// Collect one entry per namespace, preferring the namespace container
	// command for the short description.
	namespaceShort := map[string]string{}
	var order []string
	addNamespace := func(name, short string) {
		if _, exists := namespaceShort[name]; !exists {
			order = append(order, name)
		}
		if namespaceShort[name] == "" && short != "" {
			namespaceShort[name] = short
		}
	}

	for _, cmd := range uniqueVisibleCommands(commands) {
		if cmd.Namespace == "" {
			continue
		}
		if cmd.Resource == "" && cmd.Verb == "" {
			addNamespace(cmd.Namespace, cmd.Short)
		} else if namespaceShort[cmd.Namespace] == "" {
			addNamespace(cmd.Namespace, cmd.Short)
		}
	}

	namespaces := make([]NamespaceEntry, 0, len(order))
	for _, name := range order {
		namespaces = append(
			namespaces,
			NamespaceEntry{Name: name, Short: roffEscape(namespaceShort[name])},
		)
	}

	return &RootPage{
		Version:         version,
		SynopsisCommand: BinaryName,
		SynopsisArgs:    " [global option ...] <command> [args ...]",
		Description: roffEscape(
			"scw is the Scaleway command line tool. It allows you to manage your Scaleway " +
				"infrastructure: compute, storage, networking, databases and more.",
		),
		GlobalOptions: GlobalOptions,
		Namespaces:    namespaces,
		EnvVars:       EnvVariables(),
	}
}

// GlobalOptions are the global flags available on every scw command.
var GlobalOptions = []Option{
	{Term: "--profile <value>, -p <value>", Description: "The config profile to use"},
	{Term: "--config <path>, -c <path>", Description: "The path of the config file"},
	{Term: "--output <format>, -o <format>", Description: "Output format: json or human"},
	{Term: "--debug, -D", Description: "Enable debug mode"},
	{Term: "--beta, -B", Description: "Enable beta mode"},
}

// roffEscape sanitizes text so it can be embedded in a roff man page:
// ANSI escape sequences are stripped, backslashes are escaped and lines
// starting with a roff macro character are neutralized.
func roffEscape(s string) string {
	s = ansiRegexp.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, `\`, `\E`)

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, ".") || strings.HasPrefix(line, `\`) {
			line = `\&` + line
		}
		lines[i] = line
	}

	return strings.Join(lines, "\n")
}

// Render renders a command man page to roff.
func Render(page *ManPage) (string, error) {
	return executeTemplate("command", page)
}

// RenderRoot renders the top-level `scw` man page to roff.
func RenderRoot(page *RootPage) (string, error) {
	return executeTemplate("root", page)
}

func executeTemplate(name string, data any) (string, error) {
	tpl, err := template.New("man").Parse(manTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing man template: %w", err)
	}

	buffer := bytes.Buffer{}
	if err := tpl.ExecuteTemplate(&buffer, name, data); err != nil {
		return "", fmt.Errorf("rendering man page: %w", err)
	}

	return buffer.String(), nil
}
