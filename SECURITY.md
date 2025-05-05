# Security Policy

## Supported Versions

Currently supported versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

The Go-MCP-LSP team takes all security vulnerabilities seriously. We appreciate your efforts to responsibly disclose your findings.

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please follow these steps:

1. Email [security@example.com](mailto:security@example.com) with a detailed description of the vulnerability
2. Include steps to reproduce, if possible
3. Indicate the affected version(s)
4. If known, include potential fixes or mitigation steps

You should receive a response within 48 hours. We will work with you to understand and address the issue promptly.

## Security Best Practices

When using Go-MCP-LSP:

1. Keep the software updated to the latest version
2. Be cautious when implementing custom rules that parse or execute user input
3. Validate all YAML configurations before loading them into production
4. Apply the principle of least privilege when deploying the MCP server

## Security Features

The Go-MCP-LSP includes several security-focused governance rules:

- Detection of weak cryptographic algorithms
- Identification of SQL injection vulnerabilities
- Detection of hardcoded credentials
- Validation of secure API design patterns

## Acknowledgments

We would like to thank the following individuals for their responsible vulnerability disclosures (alphabetically):

*No acknowledgments yet.*
