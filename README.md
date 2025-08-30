# Pixela MCP Server

A Model Context Protocol (MCP) server for operating the Pixela API, implemented in Go.

## Features

This MCP server supports the following Pixela API operations as tools:

### User Management
- **create_user**: Create a user on Pixela
- **update_user**: Update user authentication token
- **update_user_profile**: Update user profile information
- **delete_user**: Delete a user

### Graph Management
- **create_graph**: Create a graph for a user
- **update_graph**: Update a graph definition
- **delete_graph**: Delete a specific graph
- **get_graphs**: Get all graph definitions for a user
- **get_graph_definition**: Get a specific graph definition

### Pixel Management
- **post_pixel**: Post a pixel to a graph
- **update_pixel**: Update a pixel
- **delete_pixel**: Delete a pixel
- **get_pixels**: Get a list of pixels
- **get_pixel**: Get a specific pixel
- **get_latest_pixel**: Get the latest pixel
- **get_today_pixel**: Get today's pixel
- **batch_post_pixels**: Batch post pixels
- **increment_pixel**: Increment today's pixel
- **decrement_pixel**: Decrement today's pixel
- **add_pixel**: Add a value to today's pixel (Pixela Instant recording `/add` endpoint)
- **subtract_pixel**: Subtract a value from today's pixel (Pixela Instant recording `/subtract` endpoint)
- **stopwatch**: Start or stop the stopwatch for a specific graph (Pixela Instant recording `/stopwatch` endpoint)

### Webhook Management
- **create_webhook**: Create a webhook
- **get_webhooks**: Get a list of webhooks
- **invoke_webhook**: Invoke a webhook
- **delete_webhook**: Delete a webhook

## Setup

### Prerequisites

- Go 1.21 or later
- Pixela account (https://pixe.la/)

### Installation

#### Method 1: Run directly

```bash
git clone https://github.com/a-know/pixela-mcp.git
cd pixela-mcp
go mod tidy
go run .
```

#### Method 2: Using Docker

```bash
git clone https://github.com/a-know/pixela-mcp.git
cd pixela-mcp
docker-compose up -d
```
or
```bash
docker build -t pixela-mcp .
docker run -it --rm pixela-mcp
```

> **Note:** The server communicates via standard input/output (MCP protocol). It does **not** listen on a TCP port.

## Usage

### MCP Client Configuration Example (for Cursor)

```json
{
  "mcpServers": {
    "pixela": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/path/to/pixela-mcp"
    }
  }
}
```
or for Docker:
```json
{
  "mcpServers": {
    "pixela-mcp": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "pixelaapi/pixela-mcp"]
    }
  }
}
```

### Available Tools & Parameters

（ツール名・説明・パラメータは `main.go` の `handleToolsList` 実装に完全準拠）

#### User Management

- **create_user**
  - `username` (string): User name
  - `token` (string): Authentication token
  - `agreeTermsOfService` (string): Agreement to the terms of service ("yes"/"no")
  - `notMinor` (string): Confirmation of not being a minor ("yes"/"no")

- **update_user**
  - `username` (string): User name
  - `token` (string): Current authentication token
  - `newToken` (string): New authentication token

- **update_user_profile**
  - `username` (string): User name
  - `token` (string): Authentication token
  - `displayName` (string, optional): Display name
  - `gravatarIconEmail` (string, optional): Gravatar icon email address
  - `title` (string, optional): Title
  - `about` (string, optional): About
  - `pixelaGraph` (string, optional): Pixela graph URL
  - `timezone` (string, optional): Timezone
  - `contributeURLs` (string, optional): Contribute URLs (comma-separated)

- **delete_user**
  - `username` (string): User name
  - `token` (string): Authentication token

#### Graph Management

- **create_graph**
  - `username`, `token`, `graphID`, `name`, `unit`, `type`, `color` (all string, required)

- **update_graph**
  - `username`, `token`, `graphID` (required)
  - `name`, `unit`, `color`, `purgeCacheURLs`, `selfSufficient`, `isSecret`, `publishOptionalData` (optional)

- **delete_graph**
  - `username`, `token`, `graphID` (all string, required)

- **get_graphs**
  - `username`, `token` (both string, required)

- **get_graph_definition**
  - `username`, `token`, `graphID` (all string, required)

#### Pixel Management

- **post_pixel**
  - `username`, `token`, `graphID`, `date`, `quantity` (all string, required)

- **update_pixel**
  - `username`, `token`, `graphID`, `date`, `quantity` (all string, required)
  - `optionalData` (string, optional)

- **delete_pixel**
  - `username`, `token`, `graphID`, `date` (all string, required)

- **get_pixels**
  - `username`, `token`, `graphID` (required)
  - `from`, `to`, `mode` (optional)

- **get_pixel**
  - `username`, `token`, `graphID`, `date` (all string, required)

- **get_latest_pixel**
  - `username`, `token`, `graphID` (all string, required)

- **get_today_pixel**
  - `username`, `token`, `graphID` (all string, required)

- **batch_post_pixels**
  - `username`, `token`, `graphID`, `pixels` (all string, required; `pixels` is JSON array)

- **increment_pixel / decrement_pixel**
  - `username`, `token`, `graphID` (all string, required)
- **add_pixel**
  - `username` (string, required): User name
  - `token` (string, required): Authentication token
  - `graphID` (string, required): Graph ID
  - `quantity` (string, required): Value to add (as string)
- **subtract_pixel**
  - `username` (string, required): User name
  - `token` (string, required): Authentication token
  - `graphID` (string, required): Graph ID
  - `quantity` (string, required): Value to subtract (as string)
- **stopwatch**
  - `username` (string, required): User name
  - `token` (string, required): Authentication token
  - `graphID` (string, required): Graph ID

#### Webhook Management

- **create_webhook**
  - `username`, `token`, `graphID`, `type` (all string, required)
  - `quantity` (string, optional)

- **get_webhooks**
  - `username`, `token` (both string, required)

- **invoke_webhook**
  - `username`, `webhookHash` (both string, required)

- **delete_webhook**
  - `username`, `token`, `webhookHash` (all string, required)

## Technical Notes

- Implements MCP protocol version `2024-11-05` (JSON-RPC 2.0 over stdio)
- All tool definitions and parameters are dynamically listed via `tools/list`
- Pixela API quirks (e.g., type inconsistencies) are handled internally
- Some Pixela API features require a supporter account or may be rate-limited

## Project Structure

```
pixela-mcp/
├── main.go              # MCP server entry point
├── tools.go             # MCP tool implementations
├── main_test.go         # Tests
├── pixela/
│   └── client.go        # Pixela API client
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── cursor_log.md
└── README.md
```

## Testing

```bash
go test -v
```

## License

MIT License

## Author

a-know (https://github.com/a-know) 
