package docgen

import (
	"bytes"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/scaleway/scaleway-cli/internal/core"
)

type tplResource struct {
	Cmd   *core.Command
	Verbs map[string]*core.Command
}

type tplNamespace struct {
	Cmd       *core.Command
	Resources map[string]*tplResource
}

type tplData struct {
	Namespaces map[string]*tplNamespace
}

// Generate markdown documentation for a given list of commands
func GenerateDocs(commands *core.Commands, outDir string) error {
	data := &tplData{
		Namespaces: map[string]*tplNamespace{},
	}

	for _, c := range commands.GetAll() {
		if data.Namespaces[c.Namespace] == nil {
			data.Namespaces[c.Namespace] = &tplNamespace{
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

	indexDoc, err := renderIndex(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(outDir, "index.md"), []byte(indexDoc), 0600)
	if err != nil {
		return err
	}

	for name, namespace := range data.Namespaces {
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

func renderIndex(data *tplData) (string, error) {
	buffer := bytes.Buffer{}
	err := newTemplate().ExecuteTemplate(&buffer, "index", data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func renderNamespace(data *tplNamespace) (string, error) {
	buffer := bytes.Buffer{}
	err := newTemplate().ExecuteTemplate(&buffer, "namespace", data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func newTemplate() *template.Template {
	tpl := template.New("index")
	template.Must(tpl.Parse(tplStr))
	return tpl
}
