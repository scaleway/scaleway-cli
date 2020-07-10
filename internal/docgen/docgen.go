package docgen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"text/template"

	"github.com/scaleway/scaleway-cli/internal/interactive"

	"github.com/scaleway/scaleway-cli/internal/core"
)

type tplData struct {
	Namespaces map[string]*tplNamespace
}

type tplNamespace struct {
	Cmd       *core.Command
	Commands  *core.Commands
	Resources map[string]*tplResource
}

type tplResource struct {
	Cmd   *core.Command
	Verbs map[string]*core.Command
}

// Generate markdown documentation for a given list of commands
func GenerateDocs(commands *core.Commands, outDir string) error {
	// Prepare data that will be sent to template engine
	data := &tplData{
		Namespaces: map[string]*tplNamespace{},
	}

	for _, c := range commands.GetAll() {
		if c.Hidden {
			continue
		}

		if data.Namespaces[c.Namespace] == nil {
			data.Namespaces[c.Namespace] = &tplNamespace{
				Commands:  commands,
				Resources: map[string]*tplResource{},
			}
		}
		namespace := data.Namespaces[c.Namespace]

		// If we have no resource command is the namespace command
		if c.Resource == "" {
			namespace.Cmd = c
			continue
		}

		if namespace.Resources[c.Resource] == nil {
			namespace.Resources[c.Resource] = &tplResource{
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

	for name, namespace := range data.Namespaces {
		fmt.Println("Generating namespace", name)
		namespaceDoc, err := renderNamespace(namespace)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path.Join(outDir, name+".md"), []byte(namespaceDoc), 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

func renderNamespace(data *tplNamespace) (string, error) {
	buffer := bytes.Buffer{}
	err := newTemplate().ExecuteTemplate(&buffer, "namespace", data)
	if err != nil {
		return "", err
	}
	str := buffer.String()
	//str = interactive.UnIndent(str, 2)
	return str, nil
}

func newTemplate() *template.Template {
	tpl := template.New("index")
	tpl = tpl.Funcs(map[string]interface{}{
		"unindent": func(indent int, str string) string {
			return interactive.UnIndent(str, indent)
		},
		"bq": func(count ...int) string {
			return "`"
		},
		"bbq": func(count ...int) string {
			return "```"
		},
		"map": func(args ...interface{}) map[string]interface{} {
			res := map[string]interface{}{}
			for i := 0; i < len(args); i += 2 {
				res[args[i].(string)] = args[i+1]
			}
			return res
		},
		"join": func(sep string, slice []string) string {
			return strings.Join(slice, sep)
		},
		"concat": func(strs ...string) string {
			s := ""
			for _, str := range strs {
				s += str
			}
			return s
		},
		"anchor": func(short string) string {
			res := strings.ToLower(short)
			res = strings.ReplaceAll(res, " ", "-")
			res = strings.ReplaceAll(res, "/", "")
			return res
		},
		"arg_spec_flag": func(arg *core.ArgSpec) string {
			parts := []string(nil)
			if arg.Required {
				parts = append(parts, "Required")
			}
			if arg.Default != nil {
				_, doc := arg.Default(core.GetDocGenContext())
				parts = append(parts, fmt.Sprintf("Default: `%s`", doc))
			}
			if len(arg.EnumValues) > 0 {
				parts = append(parts, fmt.Sprintf("One of: `%s`", strings.Join(arg.EnumValues, "`, `")))
			}
			return strings.Join(parts, "<br />")
		},
		"default": func(defaultValue string, value string) string {
			if value == "" {
				return defaultValue
			}
			return value
		},
	})
	tpl = template.Must(tpl.Parse(tplStr))
	return tpl
}
