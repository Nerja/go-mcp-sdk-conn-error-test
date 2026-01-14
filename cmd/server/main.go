package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const ServerPort = ":8765"

type EchoInput struct {
	Message string `json:"message"`
}

type EchoOutput struct {
	Echo string `json:"echo"`
}

func main() {
	log.Println("Creating MCP server...")

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}, nil)

	// Add a simple tool for testing
	mcp.AddTool(server, &mcp.Tool{
		Name:        "echo",
		Description: "Echo back the input message",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input EchoInput) (*mcp.CallToolResult, EchoOutput, error) {
		return nil, EchoOutput{Echo: "Echo: " + input.Message}, nil
	})

	// Create HTTP handler with streamable transport
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		SessionTimeout: 5 * time.Minute,
	})

	log.Printf("Starting MCP server on %s\n", ServerPort)
	log.Println("Press Ctrl+C to stop the server")

	if err := http.ListenAndServe(ServerPort, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
