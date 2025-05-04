# MCP Server Reference Implementation

A lightweight Meta-Configuration Protocol (MCP) server that exposes governance rules, templates, and tools to the Go Language Server.

## Overview

This reference implementation provides the JSON-RPC 2.0 endpoints required by the MCP specification. It enables:

- Intent-driven governance through rule validation
- Template-based code scaffolding
- Documentation resources for hover tooltips
- Toolchain integration for validation and generation

## Endpoints

The server exposes these primary endpoints:

- `GetResource` - Retrieve rules and documentation
- `GetPrompt` - Access templates for code generation
- `CallTool` - Execute tools like validation and scaffolding
- `ValidateIntent` - Validate code against rules

## Usage

Start the server with:

```bash
go run cmd/main.go --address=localhost:9000 --rules=./rules --templates=./templates
```

Or use a configuration file:

```bash
go run cmd/main.go --config=config.json
```

## Directory Structure

- `rules/` - YAML files defining coding standards and constraints
- `templates/` - Go templates for scaffolding and generation
- `endpoints/` - Implementation of MCP protocol handlers
- `cmd/` - Server entry point and CLI

## Integration

The Go Language Server integrates with this MCP server through the mcpclient module, which provides a JSON-RPC client for calling the exposed endpoints.
