# MCP Go SDK Connection Error Test

This project contains a simple MCP client and server using the official MCP Go SDK (v1.2.0) to test connection behavior, including network blocking scenarios.

**Server Port:** `8765`

## Prerequisites

- Go 1.23 or later
- `iptables` (typically pre-installed on Ubuntu)
- `sudo` access for traffic blocking

## Setup

First, download dependencies:

```bash
go mod tidy
```

## Running the Programs

### Start the Server

```bash
go run ./cmd/server
```

The server will listen on port `8765`.

### Stop the Server

Press `Ctrl+C` in the terminal running the server.

---

### Start the Client

In a separate terminal:

```bash
go run ./cmd/client
```

The client will connect to the server and ping every 2 seconds.

### Stop the Client

Press `Ctrl+C` in the terminal running the client.

---

## Traffic Blocking (using iptables)

These commands use `iptables`, which is typically pre-installed on Ubuntu. They require `sudo` access.

### Block Traffic to Port 8765

This will drop all incoming TCP packets to port 8765:

```bash
sudo iptables -A INPUT -p tcp --dport 8765 -j DROP
```

### Unblock Traffic to Port 8765

Remove the blocking rule:

```bash
sudo iptables -D INPUT -p tcp --dport 8765 -j DROP
```

### Check Current iptables Rules

```bash
sudo iptables -L INPUT -n --line-numbers
```

### Alternative: Block Outgoing Traffic (Client Side)

If you want to block the client from reaching the server (useful when client and server are on the same machine and INPUT rules don't apply to localhost):

**Block:**
```bash
sudo iptables -A OUTPUT -p tcp --dport 8765 -j DROP
```

**Unblock:**
```bash
sudo iptables -D OUTPUT -p tcp --dport 8765 -j DROP
```

---

## Quick Test Scenario

1. **Terminal 1:** Start the server
   ```bash
   go run ./cmd/server
   ```

2. **Terminal 2:** Start the client
   ```bash
   go run ./cmd/client
   ```
   
   You should see successful pings every 2 seconds.

3. **Terminal 3:** Block traffic
   ```bash
   sudo iptables -A OUTPUT -p tcp --dport 8765 -j DROP
   ```
   
   The client should start reporting ping failures.

4. **Terminal 3:** Unblock traffic
   ```bash
   sudo iptables -D OUTPUT -p tcp --dport 8765 -j DROP
   ```
   
   The client may need to be restarted to re-establish the connection.

---

## Cleanup

Make sure to remove any iptables rules when done:

```bash
sudo iptables -D INPUT -p tcp --dport 8765 -j DROP 2>/dev/null
sudo iptables -D OUTPUT -p tcp --dport 8765 -j DROP 2>/dev/null
```

Or flush all iptables rules (use with caution on production systems):

```bash
sudo iptables -F
```
