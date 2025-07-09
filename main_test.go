package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMCPServer_Initialize(t *testing.T) {
	server := NewMCPServer()

	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params:  nil,
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleMCPRequest(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response MCPResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC 2.0, got %s", response.JSONRPC)
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}
}

func TestMCPServer_ToolsList(t *testing.T) {
	server := NewMCPServer()

	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
		Params:  nil,
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleMCPRequest(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response MCPResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}

	// ツールリストが返されることを確認
	result, ok := response.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}

	tools, ok := result["tools"].([]interface{})
	if !ok {
		t.Fatal("Expected tools to be an array")
	}

	if len(tools) == 0 {
		t.Error("Expected at least one tool")
	}
}

func TestMCPServer_UnknownMethod(t *testing.T) {
	server := NewMCPServer()

	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "unknown_method",
		Params:  nil,
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.handleMCPRequest(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response MCPResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Error == nil {
		t.Error("Expected error for unknown method")
	}

	if response.Error.Code != -32601 {
		t.Errorf("Expected error code -32601, got %d", response.Error.Code)
	}
}
