package rdb

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

const rdbPasswordPlaceholder = "YOUR_PASSWORD"

type rdbConfigType string

const (
	rdbConfigTypePHP        rdbConfigType = "php"
	rdbConfigTypeNode       rdbConfigType = "node"
	rdbConfigTypeTypeScript rdbConfigType = "typescript"
	rdbConfigTypePython     rdbConfigType = "python"
	rdbConfigTypeGo         rdbConfigType = "go"
	rdbConfigTypeRust       rdbConfigType = "rust"
)

func renderRDBConfig(configType rdbConfigType, info *ConnectionInfo) (core.RawResult, error) {
	switch configType {
	case rdbConfigTypePHP,
		rdbConfigTypeNode,
		rdbConfigTypeTypeScript,
		rdbConfigTypePython,
		rdbConfigTypeGo,
		rdbConfigTypeRust:
		return renderConfigTemplate(configType, info)
	default:
		return core.RawResult(""), fmt.Errorf("unsupported config type %q", configType)
	}
}

// RenderPHPConfig renders a PHP database connection snippet.
func RenderPHPConfig(info *ConnectionInfo) core.RawResult {
	result, err := renderConfigTemplate(rdbConfigTypePHP, info)
	if err != nil {
		return core.RawResult("")
	}

	return result
}

// RenderNodeConfig renders a Node.js database connection snippet.
func RenderNodeConfig(info *ConnectionInfo) core.RawResult {
	result, err := renderConfigTemplate(rdbConfigTypeNode, info)
	if err != nil {
		return core.RawResult("")
	}

	return result
}

// RenderTypeScriptConfig renders a TypeScript database connection snippet.
func RenderTypeScriptConfig(info *ConnectionInfo) core.RawResult {
	result, err := renderConfigTemplate(rdbConfigTypeTypeScript, info)
	if err != nil {
		return core.RawResult("")
	}

	return result
}

// RenderPythonConfig renders a Python database connection snippet.
func RenderPythonConfig(info *ConnectionInfo) core.RawResult {
	result, err := renderConfigTemplate(rdbConfigTypePython, info)
	if err != nil {
		return core.RawResult("")
	}

	return result
}

// RenderGoConfig renders a Go database connection snippet.
func RenderGoConfig(info *ConnectionInfo) core.RawResult {
	result, err := renderConfigTemplate(rdbConfigTypeGo, info)
	if err != nil {
		return core.RawResult("")
	}

	return result
}

// RenderRustConfig renders a Rust database connection snippet.
func RenderRustConfig(info *ConnectionInfo) core.RawResult {
	result, err := renderConfigTemplate(rdbConfigTypeRust, info)
	if err != nil {
		return core.RawResult("")
	}

	return result
}
