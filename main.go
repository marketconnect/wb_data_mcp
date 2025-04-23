package main

import (
	"context"
	"fmt"
	"log"

	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/marketconnect/wb_data_mcp/config"
	"github.com/marketconnect/wb_data_mcp/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connections
	log.Println("Connecting to databases...")
	db, err := server.NewDatabase(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Println("Database connections established successfully")

	// Create MCP server
	mcpImpl := mcpserver.NewMCPServer("wb_data_mcp", "1.0.0")
	if mcpImpl == nil {
		log.Fatal("Failed to create MCP server")
	}

	// Register tools with database context
	server.RegisterTools(mcpImpl, db)

	// Create SSE server with options
	baseURL := fmt.Sprintf("http://%s:%s", cfg.Server.IP, cfg.Server.Port)
	sseServer := mcpserver.NewSSEServer(mcpImpl,
		mcpserver.WithBaseURL(baseURL),
		mcpserver.WithBasePath("/"),
		mcpserver.WithSSEEndpoint("/events"),
		mcpserver.WithMessageEndpoint("/messages"),
		mcpserver.WithKeepAlive(true),
	)
	if sseServer == nil {
		log.Fatal("Failed to create SSE server")
	}
	fmt.Printf("Starting SSE server at %s\n", baseURL)

	// Start the server
	if err := sseServer.Start(fmt.Sprintf("%s:%s", cfg.Server.IP, cfg.Server.Port)); err != nil {
		log.Fatal("Failed to start SSE server:", err)
	}
}
