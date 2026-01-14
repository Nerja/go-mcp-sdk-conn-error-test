package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const ServerURL = "http://localhost:8765"

func main() {
	log.Println("Creating MCP client...")

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	// Create HTTP transport to connect to the server
	transport := &mcp.StreamableClientTransport{
		Endpoint: ServerURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	log.Printf("Connecting to server at %s\n", ServerURL)

	ctx := context.Background()
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer session.Close()

	log.Println("Connected successfully!")
	log.Println("Starting ListTools loop (every 2 seconds). Press Ctrl+C to stop.")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	callCount := 0
	for range ticker.C {
		callCount++
		log.Printf("Calling ListTools #%d...", callCount)

		listCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		result, err := session.ListTools(listCtx, nil)
		cancel()

		if err != nil {
			log.Printf("ListTools #%d FAILED: %v (type: %T) isConnectionClosed=%t", callCount, err, err, errors.Is(err, mcp.ErrConnectionClosed))
		} else {
			log.Printf("ListTools #%d SUCCESS: %d tools found", callCount, len(result.Tools))
		}
	}
}
