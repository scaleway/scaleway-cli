package terraform

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"reflect"
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

func GetHCL(data interface{}) (string, error) {
	association, ok := getAssociation(data)
	if !ok {
		return "", fmt.Errorf("no terraform association found for this resource type (%s)", reflect.TypeOf(data).Name())
	}

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "scw-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	err = createImportFile(tmpDir, association, data)
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
	output, err := os.ReadFile(filepath.Join(tmpDir, "output.tf"))
	if err != nil {
		return "", err
	}

	return string(output), nil
}
