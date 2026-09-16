package docgen

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/core"
)

//go:embed main.md.tmpl
var mainTemplate string

type Data struct {
	Namespaces map[string]*Namespace
}

type Namespace struct {
	Cmd       *core.Command
	Commands  *core.Commands
	Resources map[string]*Resource
}

type Resource struct {
	Cmd   *core.Command
	Verbs map[string]*core.Command
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

// GenerateDocs generates markdown documentation for a given list of commands
func GenerateDocs(commands *core.Commands, outDir string) error {
	// Prepare data that will be sent to template engine
	data := &Data{
		Namespaces: map[string]*Namespace{},
	}

	for _, c := range commands.GetAll() {
		if c.Hidden {
			continue
		}

		if data.Namespaces[c.Namespace] == nil {
			data.Namespaces[c.Namespace] = &Namespace{
				Commands:  commands,
				Resources: map[string]*Resource{},
			}
		}
		namespace := data.Namespaces[c.Namespace]

		// If we have no resource command is the namespace command
		if c.Resource == "" {
			namespace.Cmd = c

			continue
		}

		if namespace.Resources[c.Resource] == nil {
			namespace.Resources[c.Resource] = &Resource{
				Verbs: map[string]*core.Command{},
			}
		}
		resource := namespace.Resources[c.Resource]

		// If we have no verb command is the resource command
		if c.Verb == "" {
			resource.Cmd = c

			continue
		}

		resource.Verbs[c.Verb] = c
	}

	// Fallback: if a resource has no Cmd defined, pick one verb as the fallback
	for _, ns := range data.Namespaces {
		for _, res := range ns.Resources {
			if res.Cmd == nil && len(res.Verbs) > 0 {
				for _, cmd := range res.Verbs {
					res.Cmd = cmd

					break
				}
			}
		}
	}

	for name, namespace := range data.Namespaces {
		fmt.Println("Generating namespace", name)
		namespaceDoc, err := renderNamespace(namespace)
		if err != nil {
			return err
		}
		err = os.WriteFile(path.Join(outDir, name+".md"), []byte(namespaceDoc), 0o600)
		if err != nil {
			return err
		}
	}

	return nil
}

func renderNamespace(data *Namespace) (string, error) {
	buffer := bytes.Buffer{}
	err := newTemplate().ExecuteTemplate(&buffer, "namespace", data)
	if err != nil {
		return "", err
	}
	str := buffer.String()

	return str, nil
}

func newTemplate() *template.Template {
	tpl := template.New("index")
	tpl = tpl.Funcs(map[string]any{
		"bq": func(_ ...int) string {
			return "`"
		},
		"bbq": func(_ ...int) string {
			return "```"
		},
		"map": func(args ...any) map[string]any {
			res := map[string]any{}
			for i := 0; i < len(args); i += 2 {
				res[args[i].(string)] = args[i+1]
			}

			return res
		},
		"anchor": func(short string) string {
			res := strings.ToLower(short)
			res = strings.ReplaceAll(res, " ", "-")
			res = strings.ReplaceAll(res, "/", "")

			return res
		},
		"remove_escape_sequence": func(s string) string {
			re := regexp.MustCompile(ansi)

			return re.ReplaceAllString(s, "")
		},
		"arg_spec_flag": func(arg *core.ArgSpec) string {
			parts := []string(nil)
			if arg.Deprecated {
				parts = append(parts, "Deprecated")
			}
			if arg.Required {
				parts = append(parts, "Required")
			}
			if arg.Default != nil {
				_, doc := arg.Default(core.GetDocGenContext())
				parts = append(parts, fmt.Sprintf("Default: `%s`", doc))
			}
			if len(arg.EnumValues) > 0 {
				parts = append(
					parts,
					fmt.Sprintf("One of: `%s`", strings.Join(arg.EnumValues, "`, `")),
				)
			}

			return strings.Join(parts, "<br />")
		},
		"arg_spec_name": func(arg *core.ArgSpec) string {
			res := arg.Name
			if arg.Deprecated {
				res = "~~" + arg.Name + "~~"
			}

			return res
		},
		"default": func(defaultValue string, value string) string {
			if value == "" {
				return defaultValue
			}

			return value
		},
	})
	tpl = template.Must(tpl.Parse(mainTemplate))

	return tpl
}
