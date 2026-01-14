.PHONY: build run stop server client block unblock status clean tidy

PORT := 8765
PID_DIR := .pids

# Build both binaries
build:
	go build -o bin/server ./cmd/server
	go build -o bin/client ./cmd/client

# Start both server and client in background
run: build
	@mkdir -p $(PID_DIR)
	@echo "Starting server..."
	@./bin/server > server.log 2>&1 & echo $$! > $(PID_DIR)/server.pid
	@sleep 1
	@echo "Starting client..."
	@./bin/client > client.log 2>&1 & echo $$! > $(PID_DIR)/client.pid
	@echo "Both programs started. Logs: server.log, client.log"
	@echo "Use 'make stop' to stop both programs"

# Stop both server and client
stop:
	@if [ -f $(PID_DIR)/client.pid ]; then \
		kill $$(cat $(PID_DIR)/client.pid) 2>/dev/null && echo "Client stopped" || echo "Client not running"; \
		rm -f $(PID_DIR)/client.pid; \
	fi
	@if [ -f $(PID_DIR)/server.pid ]; then \
		kill $$(cat $(PID_DIR)/server.pid) 2>/dev/null && echo "Server stopped" || echo "Server not running"; \
		rm -f $(PID_DIR)/server.pid; \
	fi

# Run the server (foreground)
server: build
	./bin/server

# Run the client (foreground)
client: build
	./bin/client

# Block traffic to the server port
block:
	sudo iptables -A OUTPUT -p tcp --dport $(PORT) -j DROP
	@echo "Traffic to port $(PORT) is now BLOCKED"

# Unblock traffic to the server port
unblock:
	sudo iptables -D OUTPUT -p tcp --dport $(PORT) -j DROP
	@echo "Traffic to port $(PORT) is now UNBLOCKED"

# Show current iptables rules
status:
	@echo "Current iptables OUTPUT rules:"
	@sudo iptables -L OUTPUT -n --line-numbers

# Clean build artifacts
clean:
	rm -rf bin/

# Download/update dependencies
tidy:
	go mod tidy
