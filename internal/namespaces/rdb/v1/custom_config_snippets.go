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
