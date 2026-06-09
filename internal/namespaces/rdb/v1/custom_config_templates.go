package rdb

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/scaleway/scaleway-cli/v2/core"
)

//go:embed templates/*.tmpl
var configTemplatesFS embed.FS

var configTemplates = template.Must(template.ParseFS(configTemplatesFS, "templates/*.tmpl"))

type configTemplateData struct {
	Host             string
	Port             uint32
	User             string
	Database         string
	Password         string
	PrivateNetworkID string
	PostgresDSN      string
	MySQLDSN         string
	MySQLGoDSN       string
}

func (info *ConnectionInfo) configTemplateData() configTemplateData {
	data := configTemplateData{
		Host:             info.Host,
		Port:             info.Port,
		User:             info.User,
		Database:         info.Database,
		Password:         rdbPasswordPlaceholder,
		PrivateNetworkID: info.PrivateNetworkID,
	}

	switch info.EngineFamily {
	case PostgreSQL:
		data.PostgresDSN = fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=require",
			info.User,
			rdbPasswordPlaceholder,
			info.hostPort(),
			info.Database,
		)
	case MySQL:
		hostPort := info.hostPort()
		data.MySQLDSN = fmt.Sprintf(
			"mysql://%s:%s@%s/%s",
			info.User,
			rdbPasswordPlaceholder,
			hostPort,
			info.Database,
		)
		data.MySQLGoDSN = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?tls=true",
			info.User,
			rdbPasswordPlaceholder,
			hostPort,
			info.Database,
		)
	}

	return data
}

func engineTemplateSuffix(family engineFamily) (string, error) {
	switch family {
	case PostgreSQL:
		return "postgresql", nil
	case MySQL:
		return "mysql", nil
	default:
		return "", fmt.Errorf("unsupported engine family %q", family)
	}
}

func renderConfigTemplate(configType rdbConfigType, info *ConnectionInfo) (core.RawResult, error) {
	engineSuffix, err := engineTemplateSuffix(info.EngineFamily)
	if err != nil {
		return core.RawResult(""), err
	}

	templateName := fmt.Sprintf("templates/%s-%s.tmpl", configType, engineSuffix)

	var buf bytes.Buffer
	if err := configTemplates.ExecuteTemplate(&buf, templateName, info.configTemplateData()); err != nil {
		return core.RawResult(""), fmt.Errorf("failed to render template %q: %w", templateName, err)
	}

	return core.RawResult(strings.TrimRight(buf.String(), "\n") + "\n"), nil
}
