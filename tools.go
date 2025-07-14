package main

import (
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
	// Convert parameters to map
	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return s.createErrorResult("Invalid parameters")
	}

	// Get tool name
	toolName, ok := paramsMap["name"].(string)
	if !ok {
		return s.createErrorResult("Tool name not found")
	}

	// Get tool arguments
	arguments, ok := paramsMap["arguments"].(map[string]interface{})
	if !ok {
		return s.createErrorResult("Arguments not found")
	}

	client := pixela.NewClient()

	switch toolName {
	case "create_user":
		return s.handleCreateUser(client, arguments)
	case "create_graph":
		return s.handleCreateGraph(client, arguments)
	case "post_pixel":
		return s.handlePostPixel(client, arguments)
	case "delete_user":
		return s.handleDeleteUser(client, arguments)
	case "update_user":
		return s.handleUpdateUser(client, arguments)
	case "update_user_profile":
		return s.handleUpdateUserProfile(client, arguments)
	case "get_graphs":
		return s.handleGetGraphs(client, arguments)
	case "get_graph_definition":
		return s.handleGetGraphDefinition(client, arguments)
	case "update_graph":
		return s.handleUpdateGraph(client, arguments)
	case "delete_graph":
		return s.handleDeleteGraph(client, arguments)
	case "get_pixels":
		return s.handleGetPixels(client, arguments)
	case "get_graph_stats":
		return s.handleGetGraphStats(client, arguments)
	case "batch_post_pixels":
		return s.handleBatchPostPixels(client, arguments)
	case "get_pixel":
		return s.handleGetPixel(client, arguments)
	case "get_latest_pixel":
		return s.handleGetLatestPixel(client, arguments)
	case "get_today_pixel":
		return s.handleGetTodayPixel(client, arguments)
	case "update_pixel":
		return s.handleUpdatePixel(client, arguments)
	case "delete_pixel":
		return s.handleDeletePixel(client, arguments)
	case "increment_pixel":
		return s.handleIncrementPixel(client, arguments)
	case "decrement_pixel":
		return s.handleDecrementPixel(client, arguments)
	case "create_webhook":
		return s.handleCreateWebhook(client, arguments)
	case "get_webhooks":
		return s.handleGetWebhooks(client, arguments)
	case "invoke_webhook":
		return s.handleInvokeWebhook(client, arguments)
	case "delete_webhook":
		return s.handleDeleteWebhook(client, arguments)
	case "add_pixel":
		return s.handleAddPixel(client, arguments)
	default:
		return s.createErrorResult(fmt.Sprintf("Unknown tool: %s", toolName))
	}
}

func (s *MCPServer) handleCreateUser(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	agreeTermsOfService, ok := args["agreeTermsOfService"].(string)
	if !ok {
		return s.createErrorResult("agreeTermsOfService parameter is required")
	}

	notMinor, ok := args["notMinor"].(string)
	if !ok {
		return s.createErrorResult("notMinor parameter is required")
	}

	req := pixela.CreateUserRequest{
		Token:               token,
		Username:            username,
		AgreeTermsOfService: agreeTermsOfService,
		NotMinor:            notMinor,
	}

	resp, err := client.CreateUser(req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to create user: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("User '%s' was created successfully", username))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to create user: %s", resp.Message))
	}
}

func (s *MCPServer) handleCreateGraph(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	name, ok := args["name"].(string)
	if !ok {
		return s.createErrorResult("name parameter is required")
	}

	unit, ok := args["unit"].(string)
	if !ok {
		return s.createErrorResult("unit parameter is required")
	}

	graphType, ok := args["type"].(string)
	if !ok {
		return s.createErrorResult("type parameter is required")
	}

	color, ok := args["color"].(string)
	if !ok {
		return s.createErrorResult("color parameter is required")
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
		return s.createErrorResult(fmt.Sprintf("Failed to create graph: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("Graph '%s' was created successfully", name))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to create graph: %s", resp.Message))
	}
}

func (s *MCPServer) handlePostPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	date, ok := args["date"].(string)
	if !ok {
		// If date is not specified, use today's date
		date = time.Now().Format("20060102")
	}

	quantity, ok := args["quantity"].(string)
	if !ok {
		return s.createErrorResult("quantity parameter is required")
	}

	req := pixela.PostPixelRequest{
		Date:     date,
		Quantity: quantity,
	}

	resp, err := client.PostPixel(username, token, graphID, req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to post pixel: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("Pixel was posted successfully (date: %s, quantity: %s)", date, quantity))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to post pixel: %s", resp.Message))
	}
}

func (s *MCPServer) handleDeleteUser(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	// Add debug log
	fmt.Printf("DEBUG: Deleting user '%s' with token '%s'\n", username, token)

	resp, err := client.DeleteUser(username, token)
	if err != nil {
		fmt.Printf("DEBUG: Error deleting user: %v\n", err)
		return s.createErrorResult(fmt.Sprintf("Failed to delete user: %v", err))
	}

	fmt.Printf("DEBUG: Pixela API response: %+v\n", resp)

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("User '%s' was deleted successfully", username))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to delete user: %s", resp.Message))
	}
}

func (s *MCPServer) handleUpdateUser(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	newToken, ok := args["newToken"].(string)
	if !ok {
		return s.createErrorResult("newToken parameter is required")
	}

	// thanksCode is optional
	thanksCode, _ := args["thanksCode"].(string)

	req := pixela.UpdateUserRequest{
		NewToken:   newToken,
		ThanksCode: thanksCode,
	}

	resp, err := client.UpdateUser(username, token, req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to update user: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("User '%s' information was updated successfully", username))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to update user: %s", resp.Message))
	}
}

func (s *MCPServer) handleUpdateUserProfile(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	// All profile update parameters are optional
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
		return s.createErrorResult(fmt.Sprintf("Failed to update user profile: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("User '%s' profile was updated successfully", username))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to update user profile: %s", resp.Message))
	}
}

func (s *MCPServer) handleGetGraphs(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	resp, err := client.GetGraphs(username, token)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to get graph definitions: %v", err))
	}

	if len(resp.Graphs) == 0 {
		return s.createSuccessResult(fmt.Sprintf("No graphs found for user '%s'", username))
	}

	// Format graph list for return
	var graphList []string
	for _, graph := range resp.Graphs {
		graphInfo := fmt.Sprintf("ID: %s, Name: %s, Unit: %s, Type: %s, Color: %s",
			graph.ID, graph.Name, graph.Unit, graph.Type, graph.Color)
		graphList = append(graphList, graphInfo)
	}

	message := fmt.Sprintf("Graph list for user '%s' (%d items):\n%s",
		username, len(resp.Graphs), strings.Join(graphList, "\n"))

	return s.createSuccessResult(message)
}

func (s *MCPServer) handleGetGraphDefinition(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	graph, err := client.GetGraphDefinition(username, token, graphID)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to get graph definition: %v", err))
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

	return s.createSuccessResult(fmt.Sprintf("Graph definition retrieved: %s", graph.Name), graphData)
}

func (s *MCPServer) handleUpdateGraph(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	req := pixela.UpdateGraphRequest{}

	// Set optional parameters
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
		return s.createErrorResult(fmt.Sprintf("Failed to update graph: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("Graph '%s' was updated successfully", graphID))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to update graph: %s", resp.Message))
	}
}

func (s *MCPServer) handleDeleteGraph(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	resp, err := client.DeleteGraph(username, token, graphID)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to delete graph: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("Graph '%s' was deleted successfully", graphID))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to delete graph: %s", resp.Message))
	}
}

func (s *MCPServer) handleGetPixels(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
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
		return s.createErrorResult(fmt.Sprintf("Failed to get pixel list: %v", err))
	}

	// If withBody is true, return detailed array, otherwise return date array
	if withBody != nil && *withBody == "true" {
		if len(pixels.Pixels.Details) == 0 {
			return s.createSuccessResult(fmt.Sprintf("No pixels found for graph '%s'", graphID))
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
		return s.createSuccessResult(fmt.Sprintf("Retrieved pixel details list for graph '%s' (%d items)", graphID, len(pixelList)), pixelList)
	} else {
		if len(pixels.Pixels.Dates) == 0 {
			return s.createSuccessResult(fmt.Sprintf("No pixels found for graph '%s'", graphID))
		}
		var pixelList []map[string]interface{}
		for _, date := range pixels.Pixels.Dates {
			pixelData := map[string]interface{}{
				"date": date,
			}
			pixelList = append(pixelList, pixelData)
		}
		return s.createSuccessResult(fmt.Sprintf("Retrieved pixel list for graph '%s' (%d items)", graphID, len(pixelList)), pixelList)
	}
}

func (s *MCPServer) handleGetGraphStats(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	stats, err := client.GetGraphStats(username, token, graphID)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to get graph statistics: %v", err))
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

	return s.createSuccessResult(fmt.Sprintf("Graph '%s' statistics retrieved", graphID), statsData)
}

func (s *MCPServer) handleBatchPostPixels(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}
	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}
	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}
	pixelsRaw, ok := args["pixels"].([]interface{})
	if !ok || len(pixelsRaw) == 0 {
		return s.createErrorResult("pixels array parameter is required")
	}
	var pixels []pixela.PostPixelRequest
	for _, p := range pixelsRaw {
		m, ok := p.(map[string]interface{})
		if !ok {
			return s.createErrorResult("elements in pixels array are invalid")
		}
		date, _ := m["date"].(string)
		quantity, _ := m["quantity"].(string)
		optionalData, _ := m["optionalData"].(string)
		if date == "" || quantity == "" {
			return s.createErrorResult("each element in pixels array requires date and quantity")
		}
		pixels = append(pixels, pixela.PostPixelRequest{
			Date:         date,
			Quantity:     quantity,
			OptionalData: optionalData,
		})
	}
	resp, err := client.BatchPostPixels(username, token, graphID, pixels)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to batch post pixels: %v", err))
	}
	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("%d pixels were successfully registered", len(pixels)))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to batch post pixels: %s", resp.Message))
	}
}

func (s *MCPServer) handleGetPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}
	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}
	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}
	date, ok := args["date"].(string)
	if !ok {
		return s.createErrorResult("date parameter is required")
	}

	pixel, err := client.GetPixel(username, token, graphID, date)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to get pixel: %v", err))
	}

	pixelData := map[string]interface{}{
		"date":     pixel.Date,
		"quantity": pixel.Quantity,
	}
	if pixel.OptionalData != "" {
		pixelData["optionalData"] = pixel.OptionalData
	}

	return s.createSuccessResult(fmt.Sprintf("Pixel for date %s retrieved", date), pixelData)
}

func (s *MCPServer) handleGetLatestPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}
	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}
	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	pixel, err := client.GetLatestPixel(username, token, graphID)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to get latest pixel: %v", err))
	}

	pixelData := map[string]interface{}{
		"date":     pixel.Date,
		"quantity": pixel.Quantity,
	}
	if pixel.OptionalData != "" {
		pixelData["optionalData"] = pixel.OptionalData
	}

	return s.createSuccessResult(fmt.Sprintf("Latest pixel for graph '%s' (date: %s) retrieved", graphID, pixel.Date), pixelData)
}

func (s *MCPServer) handleGetTodayPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}
	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}
	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
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
		return s.createErrorResult(fmt.Sprintf("Failed to get today's pixel: %v", err))
	}

	pixelData := map[string]interface{}{
		"date":     pixel.Date,
		"quantity": pixel.Quantity,
	}
	if pixel.OptionalData != "" {
		pixelData["optionalData"] = pixel.OptionalData
	}

	return s.createSuccessResult(fmt.Sprintf("Today's pixel for graph '%s' (date: %s) retrieved", graphID, pixel.Date), pixelData)
}

func (s *MCPServer) handleUpdatePixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	date, ok := args["date"].(string)
	if !ok {
		return s.createErrorResult("date parameter is required")
	}

	quantity, ok := args["quantity"].(string)
	if !ok {
		return s.createErrorResult("quantity parameter is required")
	}

	req := pixela.UpdatePixelRequest{
		Quantity: quantity,
	}

	if optionalData, ok := args["optionalData"].(string); ok && optionalData != "" {
		req.OptionalData = optionalData
	}

	resp, err := client.UpdatePixel(username, token, graphID, date, req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to update pixel: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("Pixel (%s) updated successfully", date))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to update pixel: %s", resp.Message))
	}
}

func (s *MCPServer) handleDeletePixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	date, ok := args["date"].(string)
	if !ok {
		return s.createErrorResult("date parameter is required")
	}

	resp, err := client.DeletePixel(username, token, graphID, date)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to delete pixel: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("Pixel (%s) deleted successfully", date))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to delete pixel: %s", resp.Message))
	}
}

func (s *MCPServer) handleIncrementPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	resp, err := client.IncrementPixel(username, token, graphID)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to increment pixel: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult("Today's pixel incremented successfully")
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to increment pixel: %s", resp.Message))
	}
}

func (s *MCPServer) handleDecrementPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	resp, err := client.DecrementPixel(username, token, graphID)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to decrement pixel: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult("Today's pixel decremented successfully")
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to decrement pixel: %s", resp.Message))
	}
}

func (s *MCPServer) handleCreateWebhook(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}

	webhookType, ok := args["type"].(string)
	if !ok {
		return s.createErrorResult("type parameter is required")
	}

	req := pixela.CreateWebhookRequest{
		GraphID: graphID,
		Type:    webhookType,
	}

	if quantity, ok := args["quantity"].(string); ok && quantity != "" {
		req.Quantity = quantity
	}

	webhook, err := client.CreateWebhook(username, token, req)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to create webhook: %v", err))
	}

	webhookData := map[string]interface{}{
		"webhookHash": webhook.WebhookHash,
		"graphID":     webhook.GraphID,
		"type":        webhook.Type,
	}
	if webhook.Quantity != "" {
		webhookData["quantity"] = webhook.Quantity
	}

	return s.createSuccessResult(fmt.Sprintf("Webhook created successfully (webhookHash: %s)", webhook.WebhookHash), webhookData)
}

func (s *MCPServer) handleGetWebhooks(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}

	webhooksResponse, err := client.GetWebhooks(username, token)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to get webhook list: %v", err))
	}

	var webhooksData []map[string]interface{}
	for _, webhook := range webhooksResponse.Webhooks {
		webhookData := map[string]interface{}{
			"webhookHash": webhook.WebhookHash,
			"graphID":     webhook.GraphID,
			"type":        webhook.Type,
		}
		if webhook.Quantity != "" {
			webhookData["quantity"] = webhook.Quantity
		}
		webhooksData = append(webhooksData, webhookData)
	}

	return s.createSuccessResult(fmt.Sprintf("%d webhooks retrieved", len(webhooksResponse.Webhooks)), webhooksData)
}

func (s *MCPServer) handleInvokeWebhook(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}

	webhookHash, ok := args["webhookHash"].(string)
	if !ok {
		return s.createErrorResult("webhookHash parameter is required")
	}

	resp, err := client.InvokeWebhook(username, webhookHash)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to invoke webhook: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("Webhook '%s' executed successfully", webhookHash))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to invoke webhook: %s", resp.Message))
	}
}

func (s *MCPServer) handleDeleteWebhook(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}
	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}
	webhookHash, ok := args["webhookHash"].(string)
	if !ok {
		return s.createErrorResult("webhookHash parameter is required")
	}
	resp, err := client.DeleteWebhook(username, token, webhookHash)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to delete webhook: %v", err))
	}
	if resp.IsSuccess {
		return s.createSuccessResult(fmt.Sprintf("Webhook '%s' deleted successfully", webhookHash))
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to delete webhook: %s", resp.Message))
	}
}

func (s *MCPServer) handleAddPixel(client *pixela.Client, args map[string]interface{}) map[string]interface{} {
	username, ok := args["username"].(string)
	if !ok {
		return s.createErrorResult("username parameter is required")
	}
	token, ok := args["token"].(string)
	if !ok {
		return s.createErrorResult("token parameter is required")
	}
	graphID, ok := args["graphID"].(string)
	if !ok {
		return s.createErrorResult("graphID parameter is required")
	}
	quantity, ok := args["quantity"].(string)
	if !ok {
		return s.createErrorResult("quantity parameter is required")
	}

	resp, err := client.AddPixel(username, token, graphID, quantity)
	if err != nil {
		return s.createErrorResult(fmt.Sprintf("Failed to add pixel: %v", err))
	}

	if resp.IsSuccess {
		return s.createSuccessResult("Today's pixel added successfully")
	} else {
		return s.createErrorResult(fmt.Sprintf("Failed to add pixel: %s", resp.Message))
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
				"text": "Error: " + message,
			},
		},
	}
}
