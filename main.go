package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type MCPServer struct {
	router *mux.Router
}

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

func NewMCPServer() *MCPServer {
	router := mux.NewRouter()
	server := &MCPServer{router: router}
	server.setupRoutes()
	return server
}

func (s *MCPServer) setupRoutes() {
	s.router.HandleFunc("/", s.handleMCPRequest).Methods("POST")
}

func (s *MCPServer) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
				"description": "Pixelaでユーザーを作成します",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "ユーザー名",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "認証トークン",
						},
						"agreeTermsOfService": map[string]interface{}{
							"type":        "string",
							"description": "利用規約への同意（yes/no）",
						},
						"notMinor": map[string]interface{}{
							"type":        "string",
							"description": "未成年でないことの確認（yes/no）",
						},
					},
					"required": []string{"username", "token", "agreeTermsOfService", "notMinor"},
				},
			},
			{
				"name":        "create_graph",
				"description": "Pixelaでグラフを作成します",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "ユーザー名",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "認証トークン",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "グラフID",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "グラフ名",
						},
						"unit": map[string]interface{}{
							"type":        "string",
							"description": "単位",
						},
						"type": map[string]interface{}{
							"type":        "string",
							"description": "グラフタイプ（int/float）",
						},
						"color": map[string]interface{}{
							"type":        "string",
							"description": "グラフの色",
						},
					},
					"required": []string{"username", "token", "graphID", "name", "unit", "type", "color"},
				},
			},
			{
				"name":        "post_pixel",
				"description": "Pixelaにピクセルを投稿します",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "ユーザー名",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "認証トークン",
						},
						"graphID": map[string]interface{}{
							"type":        "string",
							"description": "グラフID",
						},
						"date": map[string]interface{}{
							"type":        "string",
							"description": "日付（yyyyMMdd形式）",
						},
						"quantity": map[string]interface{}{
							"type":        "string",
							"description": "数量",
						},
					},
					"required": []string{"username", "token", "graphID", "date", "quantity"},
				},
			},
		},
	}
}

func (s *MCPServer) handleToolsCall(params interface{}) map[string]interface{} {
	// ツール呼び出しの実装は後で追加
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "ツール呼び出し機能は実装中です",
			},
		},
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	server := NewMCPServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Pixela MCP Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, server.router))
}
