package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type MCPServer struct {
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

func NewMCPServer() *MCPServer {
	return &MCPServer{
		scanner: bufio.NewScanner(os.Stdin),
		writer:  bufio.NewWriter(os.Stdout),
	}
}

func (s *MCPServer) run() {
	for s.scanner.Scan() {
		line := strings.TrimSpace(s.scanner.Text())
		if line == "" {
			continue
		}

		var req MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			continue
		}

		response := s.handleRequest(req)
		s.sendResponse(response)
	}

	if err := s.scanner.Err(); err != nil {
		log.Printf("Error reading stdin: %v", err)
	}
}

func (s *MCPServer) handleRequest(req MCPRequest) MCPResponse {
	var response MCPResponse
	response.JSONRPC = "2.0"
	response.ID = req.ID

	switch req.Method {
	case "initialize":
		response.Result = s.handleInitialize(req.Params)
	case "tools/list":
		response.Result = s.handleToolsList()
	case "tools/call":
		response.Result = s.handleToolsCall(req.Params)
	default:
		response.Error = &MCPError{
			Code:    -32601,
			Message: "Method not found",
		}
	}

	return response
}

func (s *MCPServer) sendResponse(response MCPResponse) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	_, err = s.writer.WriteString(string(responseBytes) + "\n")
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}

	if err := s.writer.Flush(); err != nil {
		log.Printf("Error flushing response: %v", err)
	}
}

func (s *MCPServer) handleInitialize(params interface{}) map[string]interface{} {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": true,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    "pixela-mcp",
			"version": "1.0.0",
		},
	}
}

func (s *MCPServer) handleToolsList() map[string]interface{} {
	return map[string]interface{}{
		"tools": []map[string]interface{}{
			{
				"name":        "create_user",
				"description": "Create a user on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"agreeTermsOfService": map[string]interface{}{
							"type":        "string",
							"description": "Agreement to the terms of service (yes/no)",
						},
						"notMinor": map[string]interface{}{
							"type":        "string",
							"description": "Confirmation of not being a minor (yes/no)",
						},
					},
					"required": []string{"username", "token", "agreeTermsOfService", "notMinor"},
				},
			},
			{
				"name":        "create_graph",
				"description": "Create a graph on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Graph name",
						},
						"unit": map[string]interface{}{
							"type":        "string",
							"description": "Unit",
						},
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Graph type (int/float)",
						},
						"color": map[string]interface{}{
							"type":        "string",
							"description": "Graph color",
						},
					},
					"required": []string{"username", "token", "graphID", "name", "unit", "type", "color"},
				},
			},
			{
				"name":        "post_pixel",
				"description": "Post a pixel to Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
						"date": map[string]interface{}{
							"type":        "string",
							"description": "Date (yyyyMMdd format)",
						},
						"quantity": map[string]interface{}{
							"type":        "string",
							"description": "Quantity",
						},
					},
					"required": []string{"username", "token", "graphID", "date", "quantity"},
				},
			},
			{
				"name":        "delete_user",
				"description": "Delete a user on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
					},
					"required": []string{"username", "token"},
				},
			},
			{
				"name":        "update_user",
				"description": "Update user information on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"newToken": map[string]interface{}{
							"type":        "string",
							"description": "New authentication token",
						},
					},
					"required": []string{"username", "token", "newToken"},
				},
			},
			{
				"name":        "update_user_profile",
				"description": "Update user profile on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"displayName": map[string]interface{}{
							"type":        "string",
							"description": "Display name",
						},
						"gravatarIconEmail": map[string]interface{}{
							"type":        "string",
							"description": "Gravatar icon email address",
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Title",
						},
						"about": map[string]interface{}{
							"type":        "string",
							"description": "About",
						},
						"pixelaGraph": map[string]interface{}{
							"type":        "string",
							"description": "Pixela graph URL",
						},
						"timezone": map[string]interface{}{
							"type":        "string",
							"description": "Timezone",
						},
						"contributeURLs": map[string]interface{}{
							"type":        "string",
							"description": "Contribute URLs (comma-separated)",
						},
					},
					"required": []string{"username", "token"},
				},
			},
			{
				"name":        "get_graphs",
				"description": "Get a list of graphs on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
					},
					"required": []string{"username", "token"},
				},
			},
			{
				"name":        "get_graph_definition",
				"description": "Get graph definition on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "update_graph",
				"description": "Update a graph on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Graph name",
						},
						"unit": map[string]interface{}{
							"type":        "string",
							"description": "Unit",
						},
						"color": map[string]interface{}{
							"type":        "string",
							"description": "Graph color",
						},
						"purgeCacheURLs": map[string]interface{}{
							"type":        "string",
							"description": "Purge cache URLs (comma-separated)",
						},
						"selfSufficient": map[string]interface{}{
							"type":        "string",
							"description": "Self-sufficient (increment/decrement/none)",
						},
						"isSecret": map[string]interface{}{
							"type":        "string",
							"description": "Is secret graph (true/false)",
						},
						"publishOptionalData": map[string]interface{}{
							"type":        "string",
							"description": "Publish optional data (true/false)",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "delete_graph",
				"description": "Delete a graph on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "get_pixels",
				"description": "Get a list of pixels on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
						"from": map[string]interface{}{
							"type":        "string",
							"description": "Start date (yyyyMMdd format)",
						},
						"to": map[string]interface{}{
							"type":        "string",
							"description": "End date (yyyyMMdd format)",
						},
						"mode": map[string]interface{}{
							"type":        "string",
							"description": "Mode (short/shortDetail)",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "get_graph_stats",
				"description": "Get graph statistics on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "batch_post_pixels",
				"description": "Batch post pixels to Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
						"pixels": map[string]interface{}{
							"type":        "string",
							"description": "Pixel data (JSON format)",
						},
					},
					"required": []string{"username", "token", "graphID", "pixels"},
				},
			},
			{
				"name":        "get_pixel",
				"description": "Get a specific pixel on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
						"date": map[string]interface{}{
							"type":        "string",
							"description": "Date (yyyyMMdd format)",
						},
					},
					"required": []string{"username", "token", "graphID", "date"},
				},
			},
			{
				"name":        "get_latest_pixel",
				"description": "Get the latest pixel on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "get_today_pixel",
				"description": "Get today's pixel on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "update_pixel",
				"description": "Update a pixel on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
						"date": map[string]interface{}{
							"type":        "string",
							"description": "Date (yyyyMMdd format)",
						},
						"quantity": map[string]interface{}{
							"type":        "string",
							"description": "Quantity",
						},
						"optionalData": map[string]interface{}{
							"type":        "string",
							"description": "Optional data (optional)",
						},
					},
					"required": []string{"username", "token", "graphID", "date", "quantity"},
				},
			},
			{
				"name":        "delete_pixel",
				"description": "Delete a specific pixel on a specific graph on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
						"date": map[string]interface{}{
							"type":        "string",
							"description": "Date (yyyyMMdd format)",
						},
					},
					"required": []string{"username", "token", "graphID", "date"},
				},
			},
			{
				"name":        "increment_pixel",
				"description": "Increment the today's pixel on a specific graph on Pixela (for int graphs +1, for float graphs +0.01)",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "decrement_pixel",
				"description": "Decrement the today's pixel on a specific graph on Pixela (for int graphs -1, for float graphs -0.01)",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "create_webhook",
				"description": "Create a new webhook on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Webhook type (increment/decrement)",
						},
						"quantity": map[string]interface{}{
							"type":        "string",
							"description": "Quantity (optional)",
						},
					},
					"required": []string{"username", "token", "graphID", "type"},
				},
			},
			{
				"name":        "get_webhooks",
				"description": "Get a list of existing webhooks on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
					},
					"required": []string{"username", "token"},
				},
			},
			{
				"name":        "invoke_webhook",
				"description": "Invoke a specific webhook on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"webhookHash": map[string]interface{}{
							"type":        "string",
							"description": "Webhook hash",
						},
					},
					"required": []string{"username", "webhookHash"},
				},
			},
			{
				"name":        "delete_webhook",
				"description": "Delete a specific webhook on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"webhookHash": map[string]interface{}{
							"type":        "string",
							"description": "Webhook hash",
						},
					},
					"required": []string{"username", "token", "webhookHash"},
				},
			},
			{
				"name":        "add_pixel",
				"description": "Add a value to today's pixel on a specific graph on Pixela",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "User name",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "Authentication token",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "Graph ID",
						},
						"quantity": map[string]interface{}{
							"type":        "string",
							"description": "Value to add (string, required)",
						},
					},
					"required": []string{"username", "token", "graphID", "quantity"},
				},
			},
		},
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	server := NewMCPServer()
	server.run()
}
