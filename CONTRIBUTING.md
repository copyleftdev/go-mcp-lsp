# Contributing to Go-MCP-LSP

Thank you for your interest in contributing to Go-MCP-LSP! This document provides guidelines and instructions for contributing to this project.

## Development Philosophy

- **Minimal Over-Engineering**: Keep implementations simple and straightforward.
- **Test-Driven Development**: Write tests that clearly reflect the purpose of the code.
- **Documentation**: Focus on concise, relevant documentation.
- **Modularity**: Each package should have a single, well-defined responsibility.

## Getting Started

1. Fork the repository
2. Clone your fork locally
3. Set up your development environment
4. Create a new branch for your changes

## Development Environment

Ensure you have:

- Go 1.19 or higher
- Basic understanding of Go's AST package

## Coding Standards

- Use Go's standard formatting tools (`gofmt`, `golint`)
- Follow Go's standard naming conventions
- Avoid global variables; use dependency injection instead
- Keep code comments minimal and only when necessary to explain complex logic

## Pull Request Process

1. Create a focused PR that addresses a single concern
2. Include tests for any new functionality
3. Update documentation if necessary
4. Ensure all tests pass
5. Request review from a maintainer

## Adding New Rules

1. Create a YAML rule definition in the appropriate category directory
2. Implement the AST-based analysis method
3. Add test cases demonstrating both compliant and non-compliant code
4. Update the analyzer engine to recognize the new rule

## Commit Messages

Follow conventional commits format:

```
feat: add new error handling rule
fix: resolve issue with AST analyzer for nested blocks
docs: update README with new rule information
test: add test cases for concurrent map access
```

## License

By contributing to this project, you agree that your contributions will be licensed under the project's MIT License.

## Code Review

All submissions require review. We use GitHub pull requests for this purpose.

## Questions?

If you have questions about the project or the contribution process, open an issue with the "question" label.
