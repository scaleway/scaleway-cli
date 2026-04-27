# MCP Server for Scaleway CLI

This namespace provides MCP (Model Context Protocol) server functionality, exposing Scaleway CLI commands as AI tools.

## Usage

### Starting the MCP Server

```bash
scw mcp server serve
```

This starts the MCP server over stdio, which can be connected to by MCP-compatible AI assistants.

### Configuration

The MCP server uses the same configuration as the CLI:
- Config file: `~/.config/scw/config.yaml` or `$SCW_CONFIG_PATH`
- Environment variables: `SCW_ACCESS_KEY`, `SCW_SECRET_KEY`, etc.

### Available Tools

All CLI commands are exposed as MCP tools with the naming convention:
- `namespace_resource_verb` (e.g., `instance_server_list`)
- `namespace_resource` (e.g., `config_get`)

Examples:
- `config_get` - Get a value from the config file
- `instance_server_list` - List servers
- `iam_user_list` - List IAM users

## Architecture

- `server/` - Core MCP server implementation
  - `server.go` - MCP server wrapper and tool registration
  - `tool.go` - Command-to-tool conversion and execution
  - `schema.go` - JSON Schema generation from ArgSpecs

## Development

### Testing

```bash
go test ./internal/namespaces/mcp/... -v
```

### Debug Mode

Enable debug logging:
```bash
SCW_DEBUG=true scw mcp server serve
```

## Example: Using with an MCP Client

Configure your MCP client to connect to the Scaleway CLI server:

```json
{
  "mcpServers": {
    "scaleway": {
      "command": "scw",
      "args": ["mcp", "server", "serve"]
    }
  }
}
```

The client will then have access to all Scaleway CLI commands as tools.
