package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/a-know/pixela-mcp/pixela"
)

type ToolCallParams struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"arguments"`
}

type ToolCallResult struct {
	Content []map[string]interface{} `json:"content"`
}

func (s *MCPServer) handleToolsCall(params interface{}) map[string]interface{} {
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return s.createErrorResult("パラメータの解析に失敗しました")
	}

	var toolCall ToolCallParams
	if err := json.Unmarshal(paramsBytes, &toolCall); err != nil {
		return s.createErrorResult("ツール呼び出しパラメータの解析に失敗しました")
	}

	client := pixela.NewClient()

	switch toolCall.Name {
	case "create_user":
		return s.handleCreateUser(client, toolCall.Args)
	case "create_graph":
		return s.handleCreateGraph(client, toolCall.Args)
	case "post_pixel":
		return s.handlePostPixel(client, toolCall.Args)
	default:
		return s.createErrorResult(fmt.Sprintf("未知のツール: %s", toolCall.Name))
	}
}

func (s *MCPServer) handleCreateUser(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("usernameパラメータが必要です")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("tokenパラメータが必要です")
	}

	agreeTermsOfService, ok := args["agreeTermsOfService"].(string)
	if !ok {
		return s.createErrorResult("agreeTermsOfServiceパラメータが必要です")
	}

	notMinor, ok := args["notMinor"].(string)
	if !ok {
		return s.createErrorResult("notMinorパラメータが必要です")
	}

	req := pixela.CreateUserRequest{
		Token:               token,
		Username:            username,
		AgreeTermsOfService: agreeTermsOfService,
		NotMinor:            notMinor,
	}

	resp, err := client.CreateUser(req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("ユーザー作成に失敗しました: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("ユーザー '%s' が正常に作成されました", username))
	} else {
		return s.createErrorResult(fmt.Sprintf("ユーザー作成に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) handleCreateGraph(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("usernameパラメータが必要です")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("tokenパラメータが必要です")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphIDパラメータが必要です")
	}

	name, ok := args["name"].(string)
	if !ok {
		return s.createErrorResult("nameパラメータが必要です")
	}

	unit, ok := args["unit"].(string)
	if !ok {
		return s.createErrorResult("unitパラメータが必要です")
	}

	graphType, ok := args["type"].(string)
	if !ok {
		return s.createErrorResult("typeパラメータが必要です")
	}

	color, ok := args["color"].(string)
	if !ok {
		return s.createErrorResult("colorパラメータが必要です")
	}

	req := pixela.CreateGraphRequest{
		ID:    graphID,
		Name:  name,
		Unit:  unit,
		Type:  graphType,
		Color: color,
	}

	resp, err := client.CreateGraph(username, token, req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("グラフ作成に失敗しました: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("グラフ '%s' が正常に作成されました", name))
	} else {
		return s.createErrorResult(fmt.Sprintf("グラフ作成に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) handlePostPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("usernameパラメータが必要です")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("tokenパラメータが必要です")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphIDパラメータが必要です")
	}

	date, ok := args["date"].(string)
	if !ok {
		// 日付が指定されていない場合は今日の日付を使用
		date = time.Now().Format("20060102")
	}

	quantity, ok := args["quantity"].(string)
	if !ok {
		return s.createErrorResult("quantityパラメータが必要です")
	}

	req := pixela.PostPixelRequest{
		Date:     date,
		Quantity: quantity,
	}

	resp, err := client.PostPixel(username, token, graphID, req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("ピクセル投稿に失敗しました: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("ピクセルが正常に投稿されました (日付: %s, 数量: %s)", date, quantity))
	} else {
		return s.createErrorResult(fmt.Sprintf("ピクセル投稿に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) createSuccessResult(message string) map[string]interface{} {
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": message,
			},
		},
	}
}

func (s *MCPServer) createErrorResult(message string) map[string]interface{} {
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "エラー: " + message,
			},
		},
	}
}
