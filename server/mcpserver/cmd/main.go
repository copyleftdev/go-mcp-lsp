package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/yourorg/go-mcp-lsp/server/mcpserver"
)

func main() {
	var (
		address     = flag.String("address", "localhost:9000", "Address to listen on")
		rulesDir    = flag.String("rules", "./rules", "Path to rules directory")
		templatesDir = flag.String("templates", "./templates", "Path to templates directory")
	)

	flag.Parse()

	absRulesDir, err := filepath.Abs(*rulesDir)
	if err != nil {
		log.Fatalf("Failed to resolve rules directory: %v", err)
	}

	absTemplatesDir, err := filepath.Abs(*templatesDir)
	if err != nil {
		log.Fatalf("Failed to resolve templates directory: %v", err)
	}

	server, err := mcpserver.NewMCPServer(absRulesDir, absTemplatesDir)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down MCP server...")
		if err := server.Stop(); err != nil {
			log.Printf("Error stopping server: %v", err)
		}
		os.Exit(0)
	}()

	log.Printf("Starting MCP server on %s", *address)
	if err := server.Start(*address); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
