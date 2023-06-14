package terraform

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

func getResourceID(format string, data interface{}) (string, error) {
	tmpl, err := template.New("terraform").Parse(format)
	if err != nil {
		return "", err
	}

	var resourceID bytes.Buffer
	err = tmpl.Execute(&resourceID, data)
	if err != nil {
		return "", err
	}

	return resourceID.String(), nil
}

type hclImportTemplateData struct {
	ResourceID   string
	ResourceName string
}

const hclImportTemplate = `
terraform {
	required_providers {
		scaleway = {
			source = "scaleway/scaleway"
		}
	}
	required_version = ">= 0.13"
}

import {
	# ID of the cloud resource
	# Check provider documentation for importable resources and format
	id = "{{ .ResourceID }}"

	# Resource address
	to = {{ .ResourceName }}.main
}
`

func createImportFile(directory string, association *association, data interface{}) error {
	importFile, err := os.CreateTemp(directory, "*.tf")
	if err != nil {
		return err
	}
	defer importFile.Close()

	resourceID, err := getResourceID(association.ImportFormat, data)
	if err != nil {
		return err
	}

	tmpl, err := template.New("").Parse(hclImportTemplate)
	if err != nil {
		return err
	}
	// Write the terraform file
	err = tmpl.Execute(importFile, hclImportTemplateData{
		ResourceID:   resourceID,
		ResourceName: association.ResourceName,
	})
	if err != nil {
		return err
	}

	return nil
}

var (
	resourceReferenceRe                  = regexp.MustCompile(`(?P<type>(data)|(resource)) "(?P<module>[a-z_]+)" "(?P<name>[a-z_]+)"`)
	resourceReferenceResourceTypeIndex   = resourceReferenceRe.SubexpIndex("type")
	resourceReferenceResourceModuleIndex = resourceReferenceRe.SubexpIndex("module")
	resourceReferenceResourceNameIndex   = resourceReferenceRe.SubexpIndex("name")
)

func getResourceReferenceFromOutput(output string) (resourceModule string, resourceName string) {
	matches := resourceReferenceRe.FindAllStringSubmatch(output, -1)
	if matches == nil {
		return "", ""
	}

	match := matches[len(matches)-1]

	resourceType := match[resourceReferenceResourceTypeIndex]
	resourceModule = match[resourceReferenceResourceModuleIndex]
	resourceName = match[resourceReferenceResourceNameIndex]

	if resourceType == "data" {
		resourceModule = fmt.Sprintf("data.%s", resourceModule)
	}

	return
}

type GetHCLConfig struct {
	Client *scw.Client
	Data   interface{}

	SkipParents  bool
	SkipChildren bool
}

func GetHCL(config *GetHCLConfig) (string, error) {
	association, ok := getAssociation(config.Data)
	if !ok {
		resourceType := "nil"
		if typeOf := reflect.TypeOf(config.Data); typeOf != nil {
			resourceType = typeOf.Name()

			if resourceType == "" {
				resourceType = typeOf.String()
			}
		}

		return "", fmt.Errorf("no terraform association found for this resource type (%s)", resourceType)
	}

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "scw-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	err = createImportFile(tmpDir, association, config.Data)
	if err != nil {
		return "", err
	}

	res, err := runInitCommand(tmpDir)
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("terraform init failed: %s", res.Stderr)
	}

	res, err = runGenerateConfigCommand(tmpDir, "output.tf")
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("terraform generate failed: %s", res.Stderr)
	}

	// Read the generated output
	outputRaw, err := os.ReadFile(filepath.Join(tmpDir, "output.tf"))
	if err != nil {
		return "", err
	}

	output := string(outputRaw)
	// Remove first 4 lines (terraform header)
	lines := strings.Split(output, "\n")
	output = strings.Join(lines[4:], "\n")

	if config.Client == nil {
		return output, nil
	}

	parents := make([]string, 0, len(association.Parents))
	children := make([]string, 0, len(association.Children))

	if !config.SkipParents {
		for attributeName, resource := range association.Parents {
			resourceData, err := resource.Fetcher(config.Client, config.Data)
			if err != nil {
				return "", err
			}

			resourceOutput, err := GetHCL(&GetHCLConfig{
				Client:       config.Client,
				Data:         resourceData,
				SkipChildren: true,
			})
			if err != nil {
				return "", err
			}

			resourceModule, resourceName := getResourceReferenceFromOutput(resourceOutput)

			parents = append(parents, resourceOutput)

			re := regexp.MustCompile(fmt.Sprintf(`%s([ \t]+)= .*`, attributeName))
			matches := re.FindAllStringSubmatch(output, -1)
			spaces := matches[len(matches)-1][1]

			output = re.ReplaceAllString(output, fmt.Sprintf("%s%s= %s.%s", attributeName, spaces, resourceModule, resourceName))
		}
	}

	if !config.SkipChildren {
		parentResourceModule, parentResourceName := getResourceReferenceFromOutput(output)

		for _, child := range association.Children {
			resourceData, err := child.Fetcher(config.Client, config.Data)
			if err != nil {
				return "", err
			}

			// resourceData SHOULD be a slice
			slice := reflect.ValueOf(resourceData)
			for i := 0; i < slice.Len(); i++ {
				resourceOutput, err := GetHCL(&GetHCLConfig{
					Client:      config.Client,
					Data:        slice.Index(i).Interface(),
					SkipParents: true,
				})
				if err != nil {
					return "", err
				}

				for childField, parentField := range child.ParentFieldMap {
					re := regexp.MustCompile(fmt.Sprintf(`%s([ \t]+)= .*`, childField))
					matches := re.FindAllStringSubmatch(resourceOutput, -1)
					spaces := matches[len(matches)-1][1]

					resourceOutput = re.ReplaceAllString(resourceOutput, fmt.Sprintf("%s%s= %s.%s.%s", childField, spaces, parentResourceModule, parentResourceName, parentField))
				}

				children = append(children, resourceOutput)
			}
		}
	}

	for _, parent := range parents {
		output = parent + "\n" + output
	}

	for _, child := range children {
		output = output + "\n" + child
	}

	return output, nil
}
