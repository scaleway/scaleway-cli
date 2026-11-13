package rdb

import "github.com/scaleway/scaleway-cli/v2/core"

// ResetAutocompleteCaches clears package-level autocomplete caches between tests.
func ResetAutocompleteCaches() {
	completeListEngineCache = nil
	completeListNodeTypeCache = nil
}

// Exported for black-box tests of config snippet rendering.

func RenderPHPConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypePHP, info)
}

func RenderNodeConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypeNode, info)
}

func RenderTypeScriptConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypeTypeScript, info)
}

func RenderPythonConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypePython, info)
}

func RenderGoConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypeGo, info)
}

func RenderRustConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypeRust, info)
}
