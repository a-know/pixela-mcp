package main

import (
	"encoding/json"
	"fmt"
	"strings"
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
	case "delete_user":
		return s.handleDeleteUser(client, toolCall.Args)
	case "update_user":
		return s.handleUpdateUser(client, toolCall.Args)
	case "update_user_profile":
		return s.handleUpdateUserProfile(client, toolCall.Args)
	case "get_graphs":
		return s.handleGetGraphs(client, toolCall.Args)
	case "get_graph_definition":
		return s.handleGetGraphDefinition(client, toolCall.Args)
	case "update_graph":
		return s.handleUpdateGraph(client, toolCall.Args)
	case "delete_graph":
		return s.handleDeleteGraph(client, toolCall.Args)
	case "get_pixels":
		return s.handleGetPixels(client, toolCall.Args)
	case "get_graph_stats":
		return s.handleGetGraphStats(client, toolCall.Args)
	case "batch_post_pixels":
		return s.handleBatchPostPixels(client, toolCall.Args)
	case "get_pixel":
		return s.handleGetPixel(client, toolCall.Args)
	case "get_latest_pixel":
		return s.handleGetLatestPixel(client, toolCall.Args)
	case "get_today_pixel":
		return s.handleGetTodayPixel(client, toolCall.Args)
	case "update_pixel":
		return s.handleUpdatePixel(client, toolCall.Args)
	case "delete_pixel":
		return s.handleDeletePixel(client, toolCall.Args)
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

func (s *MCPServer) handleDeleteUser(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("usernameパラメータが必要です")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("tokenパラメータが必要です")
	}

	// デバッグログを追加
	fmt.Printf("DEBUG: Deleting user '%s' with token '%s'\n", username, token)

	resp, err := client.DeleteUser(username, token)
	if err != nil {
		fmt.Printf("DEBUG: Error deleting user: %v\n", err)
		return s.createErrorResult(fmt.Sprintf("ユーザー削除に失敗しました: %v", err))
	}

	fmt.Printf("DEBUG: Pixela API response: %+v\n", resp)

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("ユーザー '%s' が正常に削除されました", username))
	} else {
		return s.createErrorResult(fmt.Sprintf("ユーザー削除に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) handleUpdateUser(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("usernameパラメータが必要です")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("tokenパラメータが必要です")
	}

	newToken, ok := args["newToken"].(string)
	if !ok {
		return s.createErrorResult("newTokenパラメータが必要です")
	}

	// thanksCodeはオプショナル
	thanksCode, _ := args["thanksCode"].(string)

	req := pixela.UpdateUserRequest{
		NewToken:   newToken,
		ThanksCode: thanksCode,
	}

	resp, err := client.UpdateUser(username, token, req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("ユーザー更新に失敗しました: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("ユーザー '%s' の情報が正常に更新されました", username))
	} else {
		return s.createErrorResult(fmt.Sprintf("ユーザー更新に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) handleUpdateUserProfile(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("usernameパラメータが必要です")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("tokenパラメータが必要です")
	}

	// プロフィール更新のパラメータはすべてオプショナル
	displayName, _ := args["displayName"].(string)
	profileURL, _ := args["profileURL"].(string)
	description, _ := args["description"].(string)
	avatarURL, _ := args["avatarURL"].(string)
	twitter, _ := args["twitter"].(string)
	github, _ := args["github"].(string)
	website, _ := args["website"].(string)

	req := pixela.UpdateUserProfileRequest{
		DisplayName: displayName,
		ProfileURL:  profileURL,
		Description: description,
		AvatarURL:   avatarURL,
		Twitter:     twitter,
		GitHub:      github,
		Website:     website,
	}

	resp, err := client.UpdateUserProfile(username, token, req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("ユーザープロフィール更新に失敗しました: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("ユーザー '%s' のプロフィールが正常に更新されました", username))
	} else {
		return s.createErrorResult(fmt.Sprintf("ユーザープロフィール更新に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) handleGetGraphs(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("usernameパラメータが必要です")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("tokenパラメータが必要です")
	}

	resp, err := client.GetGraphs(username, token)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("グラフ定義一覧取得に失敗しました: %v", err))
	}

	if len(resp.Graphs) == 0 {
		return s.createSuccessResult(fmt.Sprintf("ユーザー '%s' のグラフは見つかりませんでした", username))
	}

	// グラフ一覧を整形して返す
	var graphList []string
	for _, graph := range resp.Graphs {
		graphInfo := fmt.Sprintf("ID: %s, 名前: %s, 単位: %s, タイプ: %s, 色: %s",
			graph.ID, graph.Name, graph.Unit, graph.Type, graph.Color)
		graphList = append(graphList, graphInfo)
	}

	message := fmt.Sprintf("ユーザー '%s' のグラフ一覧（%d件）:\n%s",
		username, len(resp.Graphs), strings.Join(graphList, "\n"))

	return s.createSuccessResult(message)
}

func (s *MCPServer) handleGetGraphDefinition(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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

	graph, err := client.GetGraphDefinition(username, token, graphID)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("グラフ定義の取得に失敗しました: %v", err))
	}

	graphData := map[string]interface{}{
		"id":                  graph.ID,
		"name":                graph.Name,
		"unit":                graph.Unit,
		"type":                graph.Type,
		"color":               graph.Color,
		"timezone":            graph.Timezone,
		"selfSufficient":      bool(graph.SelfSufficient),
		"isSecret":            bool(graph.IsSecret),
		"publishOptionalData": bool(graph.PublishOptionalData),
	}

	return s.createSuccessResult(fmt.Sprintf("グラフ定義を取得しました: %s", graph.Name), graphData)
}

func (s *MCPServer) handleUpdateGraph(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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

	req := pixela.UpdateGraphRequest{}

	// オプションパラメータを設定
	if name, ok := args["name"].(string); ok {
		req.Name = name
	}
	if unit, ok := args["unit"].(string); ok {
		req.Unit = unit
	}
	if color, ok := args["color"].(string); ok {
		req.Color = color
	}
	if timezone, ok := args["timezone"].(string); ok {
		req.Timezone = timezone
	}
	if selfSufficient, ok := args["selfSufficient"].(string); ok {
		req.SelfSufficient = selfSufficient
	}
	if isSecret, ok := args["isSecret"].(string); ok {
		req.IsSecret = isSecret
	}
	if publishOptionalData, ok := args["publishOptionalData"].(string); ok {
		req.PublishOptionalData = publishOptionalData
	}

	resp, err := client.UpdateGraph(username, token, graphID, req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("グラフ更新に失敗しました: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("グラフ '%s' が正常に更新されました", graphID))
	} else {
		return s.createErrorResult(fmt.Sprintf("グラフ更新に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) handleDeleteGraph(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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

	resp, err := client.DeleteGraph(username, token, graphID)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("グラフ削除に失敗しました: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("グラフ '%s' が正常に削除されました", graphID))
	} else {
		return s.createErrorResult(fmt.Sprintf("グラフ削除に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) handleGetPixels(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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

	var from, to, withBody *string
	if v, ok := args["from"].(string); ok {
		from = &v
	}
	if v, ok := args["to"].(string); ok {
		to = &v
	}
	if v, ok := args["withBody"].(string); ok {
		withBody = &v
	}

	pixels, err := client.GetPixels(username, token, graphID, from, to, withBody)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("ピクセル一覧の取得に失敗しました: %v", err))
	}

	// withBodyがtrueなら詳細配列、そうでなければ日付配列
	if withBody != nil && *withBody == "true" {
		if len(pixels.Pixels.Details) == 0 {
			return s.createSuccessResult(fmt.Sprintf("グラフ '%s' にはピクセルが登録されていません", graphID))
		}
		var pixelList []map[string]interface{}
		for _, detail := range pixels.Pixels.Details {
			pixelData := map[string]interface{}{
				"date":     detail.Date,
				"quantity": detail.Quantity,
			}
			if detail.OptionalData != "" {
				pixelData["optionalData"] = detail.OptionalData
			}
			pixelList = append(pixelList, pixelData)
		}
		return s.createSuccessResult(fmt.Sprintf("グラフ '%s' のピクセル詳細一覧（%d件）を取得しました", graphID, len(pixelList)), pixelList)
	} else {
		if len(pixels.Pixels.Dates) == 0 {
			return s.createSuccessResult(fmt.Sprintf("グラフ '%s' にはピクセルが登録されていません", graphID))
		}
		var pixelList []map[string]interface{}
		for _, date := range pixels.Pixels.Dates {
			pixelData := map[string]interface{}{
				"date": date,
			}
			pixelList = append(pixelList, pixelData)
		}
		return s.createSuccessResult(fmt.Sprintf("グラフ '%s' のピクセル一覧（%d件）を取得しました", graphID, len(pixelList)), pixelList)
	}
}

func (s *MCPServer) handleGetGraphStats(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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

	stats, err := client.GetGraphStats(username, token, graphID)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("グラフ統計情報の取得に失敗しました: %v", err))
	}

	statsData := map[string]interface{}{
		"totalPixelsCount":  stats.TotalPixelsCount,
		"maxQuantity":       stats.MaxQuantity.String(),
		"minQuantity":       stats.MinQuantity.String(),
		"maxDate":           stats.MaxDate,
		"minDate":           stats.MinDate,
		"totalQuantity":     stats.TotalQuantity.String(),
		"avgQuantity":       stats.AvgQuantity.String(),
		"todaysQuantity":    stats.TodaysQuantity.String(),
		"yesterdayQuantity": stats.YesterdayQuantity.String(),
	}

	return s.createSuccessResult(fmt.Sprintf("グラフ '%s' の統計情報を取得しました", graphID), statsData)
}

func (s *MCPServer) handleBatchPostPixels(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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
	pixelsRaw, ok := args["pixels"].([]interface{})
	if !ok || len(pixelsRaw) == 0 {
		return s.createErrorResult("pixels配列パラメータが必要です")
	}
	var pixels []pixela.PostPixelRequest
	for _, p := range pixelsRaw {
		m, ok := p.(map[string]interface{})
		if !ok {
			return s.createErrorResult("pixels配列の要素が不正です")
		}
		date, _ := m["date"].(string)
		quantity, _ := m["quantity"].(string)
		optionalData, _ := m["optionalData"].(string)
		if date == "" || quantity == "" {
			return s.createErrorResult("pixels配列の各要素にはdate, quantityが必要です")
		}
		pixels = append(pixels, pixela.PostPixelRequest{
			Date:         date,
			Quantity:     quantity,
			OptionalData: optionalData,
		})
	}
	resp, err := client.BatchPostPixels(username, token, graphID, pixels)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("複数Pixel登録に失敗しました: %v", err))
	}
	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("%d件のPixelが正常に登録されました", len(pixels)))
	} else {
		return s.createErrorResult(fmt.Sprintf("複数Pixel登録に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) handleGetPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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
		return s.createErrorResult("dateパラメータが必要です")
	}

	pixel, err := client.GetPixel(username, token, graphID, date)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Pixel取得に失敗しました: %v", err))
	}

	pixelData := map[string]interface{}{
		"date":     pixel.Date,
		"quantity": pixel.Quantity,
	}
	if pixel.OptionalData != "" {
		pixelData["optionalData"] = pixel.OptionalData
	}

	return s.createSuccessResult(fmt.Sprintf("日付 %s のPixelを取得しました", date), pixelData)
}

func (s *MCPServer) handleGetLatestPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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

	pixel, err := client.GetLatestPixel(username, token, graphID)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("最新Pixel取得に失敗しました: %v", err))
	}

	pixelData := map[string]interface{}{
		"date":     pixel.Date,
		"quantity": pixel.Quantity,
	}
	if pixel.OptionalData != "" {
		pixelData["optionalData"] = pixel.OptionalData
	}

	return s.createSuccessResult(fmt.Sprintf("グラフ '%s' の最新Pixel（日付: %s）を取得しました", graphID, pixel.Date), pixelData)
}

func (s *MCPServer) handleGetTodayPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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

	var returnEmpty *bool
	if v, ok := args["returnEmpty"].(string); ok {
		if v == "true" {
			trueVal := true
			returnEmpty = &trueVal
		} else if v == "false" {
			falseVal := false
			returnEmpty = &falseVal
		}
	}

	pixel, err := client.GetTodayPixel(username, token, graphID, returnEmpty)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("今日のPixel取得に失敗しました: %v", err))
	}

	pixelData := map[string]interface{}{
		"date":     pixel.Date,
		"quantity": pixel.Quantity,
	}
	if pixel.OptionalData != "" {
		pixelData["optionalData"] = pixel.OptionalData
	}

	return s.createSuccessResult(fmt.Sprintf("グラフ '%s' の今日のPixel（日付: %s）を取得しました", graphID, pixel.Date), pixelData)
}

func (s *MCPServer) handleUpdatePixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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
		return s.createErrorResult("dateパラメータが必要です")
	}

	quantity, ok := args["quantity"].(string)
	if !ok {
		return s.createErrorResult("quantityパラメータが必要です")
	}

	req := pixela.UpdatePixelRequest{
		Quantity: quantity,
	}

	if optionalData, ok := args["optionalData"].(string); ok && optionalData != "" {
		req.OptionalData = optionalData
	}

	resp, err := client.UpdatePixel(username, token, graphID, date, req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("ピクセル更新に失敗しました: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("ピクセル（%s）が正常に更新されました", date))
	} else {
		return s.createErrorResult(fmt.Sprintf("ピクセル更新に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) handleDeletePixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
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
		return s.createErrorResult("dateパラメータが必要です")
	}

	resp, err := client.DeletePixel(username, token, graphID, date)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("ピクセル削除に失敗しました: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("ピクセル（%s）が正常に削除されました", date))
	} else {
		return s.createErrorResult(fmt.Sprintf("ピクセル削除に失敗しました: %s", resp.Message))
	}
}

func (s *MCPServer) createSuccessResult(message string, data ...interface{}) map[string]interface{} {
	content := []map[string]interface{}{
		{
			"type": "text",
			"text": message,
		},
	}

	if len(data) > 0 && data[0] != nil {
		content = append(content, map[string]interface{}{
			"type": "json",
			"json": data[0],
		})
	}

	return map[string]interface{}{
		"content": content,
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
