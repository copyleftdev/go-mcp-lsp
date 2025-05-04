# PROJECT:

An intent-aware Golang Language Server integrates with MCP to enforce governance, generate code, and validate developer actions.

# SUMMARY:

This project creates an enhanced Go language server by forking `gopls` and embedding a Meta-Configuration Protocol (MCP) client. The server enforces organization-defined intents through structural linting, template-driven scaffolding, and interactive MCP tools. Developers interact through LSP features enriched with governed behavior, embedded documentation, and real-time validation. A supporting CLI (`mcplsp`) and reference MCP server complete the ecosystem.

# STEPS:

1. Fork and modify `gopls` to support MCP integration.
2. Build a JSON-RPC 2.0 client to communicate with the MCP server.
3. Define intent-mechanism YAML structures for governance enforcement.
4. Implement AST/graph-based mechanisms to fulfill declared intents.
5. Add structural linting and boundary enforcement rules.
6. Embed scaffolding via `getPrompt` and `callTool` MCP endpoints.
7. Inject hover-based documentation using `getResource`.
8. Create save-time formatting, comment injection, and validation hooks.
9. Develop CLI tool `mcplsp` to test, audit, simulate, and verify.
10. Provide a reference MCP server exposing sample rules and templates.
11. Use snapshot and E2E tests to verify correct enforcement.
12. Integrate drift analysis and reporting for CI pipelines.
13. Organize code into composable modules using `pkg/` and `internal/`.
14. Script project setup, dependency fetching, and build tasks.
15. Deliver with README and developer onboarding documentation.
16. Ensure secure, versioned rules and deterministic test results.

# STRUCTURE:

```
go-mcp-lsp/
├── cmd/
│   └── mcplsp/
├── pkg/
│   ├── intent/
│   ├── mechanism/
│   ├── mcpclient/
│   └── lsp/
├── internal/
│   └── analyzers/
├── server/
│   └── mcpserver/
├── testdata/
├── scripts/
│   └── audit.sh
└── README.md
```

# DETAILED EXPLANATION:

1. `cmd/mcplsp/` – CLI entrypoint for running checks, simulations, audits.
2. `pkg/intent/` – Contains all YAML/JSON-defined intents and descriptions.
3. `pkg/mechanism/` – Implements AST/graph analysis, enforcement logic.
4. `pkg/mcpclient/` – Manages JSON-RPC communication with MCP server.
5. `pkg/lsp/` – Extends and overrides `gopls` behavior for MCP integration.
6. `internal/analyzers/` – Encapsulates rule-specific static analysis code.
7. `server/mcpserver/` – Lightweight MCP server exposing toolchain and rules.
8. `testdata/` – Real-world Go examples for testing drift and snapshot.
9. `scripts/audit.sh` – CI-compatible drift reporting and enforcement checker.
10. `README.md` – Project overview, setup, usage, and contributing instructions.

# CODE:

### File: `cmd/mcplsp/main.go`

```go
// CLI entrypoint for mcplsp. Supports commands like check, audit, simulate, test.
package main

import (
	"log"
	"os"

	"github.com/go-mcp-lsp/pkg/mcplsp"
)

func main() {
	if err := mcplsp.Execute(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
```

---

### File: `pkg/intent/intents.go`

```go
// Loads intent YAMLs into typed structs
package intent

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Intent struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
	Rationale   string `yaml:"rationale"`
	Category    string `yaml:"category"`
	Severity    string `yaml:"severity"`
}

func Load(file string) (*Intent, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var i Intent
	if err := yaml.Unmarshal(data, &i); err != nil {
		return nil, err
	}
	return &i, nil
}
```

---

### File: `pkg/mcpclient/client.go`

```go
// JSON-RPC 2.0 MCP client for calling remote endpoints
package mcpclient

import (
	"context"
	"encoding/json"
	"net/rpc/jsonrpc"
)

type Request struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func Call(endpoint string, method string, params interface{}) (json.RawMessage, error) {
	conn, err := jsonrpc.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var result json.RawMessage
	if err := conn.Call(method, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}
```

---

### File: `scripts/audit.sh`

```bash
#!/bin/bash
# Run drift check for declared vs enforced intents

set -e

echo "[audit] Checking for enforcement drift..."
mcplsp audit --intents ./pkg/intent --mechanisms ./pkg/mechanism --output drift-report.json

echo "[audit] Report written to drift-report.json"
```

---

### File: `README.md`

````markdown
# Go MCP Language Server

## Overview

A fork of `gopls` enhanced with MCP client integration for intent enforcement, code generation, and live documentation.

## Features

- Intent-aware autocompletion
- Rule-driven static analysis
- Template-based scaffolding
- Embedded org-style docs
- CLI tools for CI auditing

## Quick Start

```bash
git clone https://github.com/yourorg/go-mcp-lsp
cd go-mcp-lsp
make install
./scripts/audit.sh
````

## Directory Layout

See STRUCTURE section in project definition.

## Contributing

1. Fork the repo.
2. Run `make test`.
3. Submit PR with new intents or mechanisms.

````

# SETUP:
```bash
#!/bin/bash
# setup.sh - Initialize the go-mcp-lsp project

set -e

echo "[init] Cloning gopls and applying MCP patch..."
git clone https://github.com/golang/tools.git
cd tools/gopls
git checkout latest
patch -p1 < ../../../patches/mcp.diff

cd ../../
mkdir -p go-mcp-lsp/{cmd/pkg/internal/server/testdata/scripts}
touch go-mcp-lsp/scripts/audit.sh
chmod +x go-mcp-lsp/scripts/audit.sh

go mod init github.com/yourorg/go-mcp-lsp
go mod tidy

echo "[init] Setup complete."
````

# TAKEAWAYS:

1. Enforces development intent through centrally defined governance.
2. Enhances LSP behavior with real-time rule enforcement.
3. Treats coding as a verifiable expression of purpose.
4. Integrates Go analysis with domain-specific rules.
5. Enables proactive scaffolding and inline documentation.

# SUGGESTIONS:

1. Add WebSocket support for live MCP updates.
2. Enable fine-grained role-based filtering (e.g., intern vs lead).
3. Provide VSCode plugin with enhanced MCP UI.
4. Cache MCP resources locally for offline dev.
5. Auto-generate intent diff reports on PRs.

Would you like to generate the reference MCP server next?
