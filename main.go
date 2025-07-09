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
			{
				"name":        "delete_user",
				"description": "Pixelaでユーザーを削除します",
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
					},
					"required": []string{"username", "token"},
				},
			},
			{
				"name":        "update_user",
				"description": "Pixelaでユーザー情報を更新します",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "ユーザー名",
						},
						"token": map[string]interface{}{
							"type":        "string",
							"description": "現在の認証トークン",
						},
						"newToken": map[string]interface{}{
							"type":        "string",
							"description": "新しい認証トークン",
						},
						"thanksCode": map[string]interface{}{
							"type":        "string",
							"description": "サンクスコード（オプション）",
						},
					},
					"required": []string{"username", "token", "newToken"},
				},
			},
			{
				"name":        "update_user_profile",
				"description": "Pixelaでユーザープロフィールを更新します",
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
						"displayName": map[string]interface{}{
							"type":        "string",
							"description": "表示名（オプション）",
						},
						"profileURL": map[string]interface{}{
							"type":        "string",
							"description": "プロフィールURL（オプション）",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "プロフィール説明（オプション）",
						},
						"avatarURL": map[string]interface{}{
							"type":        "string",
							"description": "アバター画像URL（オプション）",
						},
						"twitter": map[string]interface{}{
							"type":        "string",
							"description": "Twitterユーザー名（オプション）",
						},
						"github": map[string]interface{}{
							"type":        "string",
							"description": "GitHubユーザー名（オプション）",
						},
						"website": map[string]interface{}{
							"type":        "string",
							"description": "ウェブサイトURL（オプション）",
						},
					},
					"required": []string{"username", "token"},
				},
			},
			{
				"name":        "get_graphs",
				"description": "Pixelaでユーザーのグラフ定義一覧を取得します",
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
					},
					"required": []string{"username", "token"},
				},
			},
			{
				"name":        "get_graph_definition",
				"description": "Pixelaで特定のグラフ定義を取得します",
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
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "update_graph",
				"description": "Pixelaでグラフ定義を更新します",
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
							"description": "グラフ名（オプション）",
						},
						"unit": map[string]interface{}{
							"type":        "string",
							"description": "単位（オプション）",
						},
						"color": map[string]interface{}{
							"type":        "string",
							"description": "グラフの色（オプション）",
						},
						"timezone": map[string]interface{}{
							"type":        "string",
							"description": "タイムゾーン（オプション）",
						},
						"selfSufficient": map[string]interface{}{
							"type":        "string",
							"description": "自己充足（yes/no、オプション）",
						},
						"isSecret": map[string]interface{}{
							"type":        "string",
							"description": "秘密グラフ（yes/no、オプション）",
						},
						"publishOptionalData": map[string]interface{}{
							"type":        "string",
							"description": "オプションデータ公開（yes/no、オプション）",
						},
					},
					"required": []string{"username", "token", "graphID"},
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Pixela MCP Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, server.router))
}
