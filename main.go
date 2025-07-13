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
							"description": "認証トークン",
						},
						"newToken": map[string]interface{}{
							"type":        "string",
							"description": "新しい認証トークン",
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
							"description": "表示名",
						},
						"gravatarIconEmail": map[string]interface{}{
							"type":        "string",
							"description": "Gravatarアイコン用メールアドレス",
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "タイトル",
						},
						"about": map[string]interface{}{
							"type":        "string",
							"description": "自己紹介",
						},
						"pixelaGraph": map[string]interface{}{
							"type":        "string",
							"description": "PixelaグラフURL",
						},
						"timezone": map[string]interface{}{
							"type":        "string",
							"description": "タイムゾーン",
						},
						"contributeURLs": map[string]interface{}{
							"type":        "string",
							"description": "貢献URL（カンマ区切り）",
						},
					},
					"required": []string{"username", "token"},
				},
			},
			{
				"name":        "get_graphs",
				"description": "Pixelaでグラフ一覧を取得します",
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
				"description": "Pixelaでグラフ定義を取得します",
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
				"description": "Pixelaでグラフを更新します",
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
						"color": map[string]interface{}{
							"type":        "string",
							"description": "グラフの色",
						},
						"purgeCacheURLs": map[string]interface{}{
							"type":        "string",
							"description": "キャッシュ削除URL（カンマ区切り）",
						},
						"selfSufficient": map[string]interface{}{
							"type":        "string",
							"description": "自己充足（increment/decrement/none）",
						},
						"isSecret": map[string]interface{}{
							"type":        "string",
							"description": "秘密グラフ（true/false）",
						},
						"publishOptionalData": map[string]interface{}{
							"type":        "string",
							"description": "オプションデータ公開（true/false）",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "delete_graph",
				"description": "Pixelaでグラフを削除します",
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
				"name":        "get_pixels",
				"description": "Pixelaでピクセル一覧を取得します",
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
						"from": map[string]interface{}{
							"type":        "string",
							"description": "開始日（yyyyMMdd形式）",
						},
						"to": map[string]interface{}{
							"type":        "string",
							"description": "終了日（yyyyMMdd形式）",
						},
						"mode": map[string]interface{}{
							"type":        "string",
							"description": "モード（short/shortDetail）",
						},
					},
					"required": []string{"username", "token", "graphID"},
				},
			},
			{
				"name":        "get_graph_stats",
				"description": "Pixelaでグラフ統計を取得します",
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
				"name":        "batch_post_pixels",
				"description": "Pixelaでピクセルを一括投稿します",
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
						"pixels": map[string]interface{}{
							"type":        "string",
							"description": "ピクセルデータ（JSON形式）",
						},
					},
					"required": []string{"username", "token", "graphID", "pixels"},
				},
			},
			{
				"name":        "get_pixel",
				"description": "Pixelaで特定のピクセルを取得します",
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
					},
					"required": []string{"username", "token", "graphID", "date"},
				},
			},
			{
				"name":        "get_latest_pixel",
				"description": "Pixelaで最新のピクセルを取得します",
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
				"name":        "get_today_pixel",
				"description": "Pixelaで今日のピクセルを取得します",
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
				"name":        "update_pixel",
				"description": "Pixelaでピクセルを更新します",
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
						"optionalData": map[string]interface{}{
							"type":        "string",
							"description": "オプションデータ（任意）",
						},
					},
					"required": []string{"username", "token", "graphID", "date", "quantity"},
				},
			},
			{
				"name":        "delete_pixel",
				"description": "Pixelaで特定のグラフの特定日付のPixelを削除します",
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
					},
					"required": []string{"username", "token", "graphID", "date"},
				},
			},
			{
				"name":        "increment_pixel",
				"description": "Pixelaで特定のグラフの今日のPixelをインクリメントします（int型グラフなら+1、float型グラフなら+0.01）",
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
				"name":        "decrement_pixel",
				"description": "Pixelaで特定のグラフの今日のPixelをデクリメントします（int型グラフなら-1、float型グラフなら-0.01）",
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
				"name":        "create_webhook",
				"description": "PixelaでWebhookを新規作成します",
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
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Webhookタイプ（increment/decrement）",
						},
						"quantity": map[string]interface{}{
							"type":        "string",
							"description": "数量（オプション）",
						},
					},
					"required": []string{"username", "token", "graphID", "type"},
				},
			},
			{
				"name":        "get_webhooks",
				"description": "Pixelaで作成済みWebhook一覧を取得します",
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
				"name":        "invoke_webhook",
				"description": "Pixelaで特定のWebhookを実行します",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "ユーザー名",
						},
						"webhookHash": map[string]interface{}{
							"type":        "string",
							"description": "Webhookハッシュ",
						},
					},
					"required": []string{"username", "webhookHash"},
				},
			},
			{
				"name":        "delete_webhook",
				"description": "Pixelaで特定のWebhookを削除します",
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
						"webhookHash": map[string]interface{}{
							"type":        "string",
							"description": "Webhookハッシュ",
						},
					},
					"required": []string{"username", "token", "webhookHash"},
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
