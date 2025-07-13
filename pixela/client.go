package pixela

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL = "https://pixe.la"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type CreateUserRequest struct {
	Token               string `json:"token"`
	Username            string `json:"username"`
	AgreeTermsOfService string `json:"agreeTermsOfService"`
	NotMinor            string `json:"notMinor"`
}

type CreateGraphRequest struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Unit                string `json:"unit"`
	Type                string `json:"type"`
	Color               string `json:"color"`
	Timezone            string `json:"timezone,omitempty"`
	SelfSufficient      string `json:"selfSufficient,omitempty"`
	IsSecret            string `json:"isSecret,omitempty"`
	PublishOptionalData string `json:"publishOptionalData,omitempty"`
}

type PostPixelRequest struct {
	Date         string `json:"date"`
	Quantity     string `json:"quantity"`
	OptionalData string `json:"optionalData,omitempty"`
}

type UpdateUserRequest struct {
	NewToken   string `json:"newToken"`
	ThanksCode string `json:"thanksCode,omitempty"`
}

type UpdateUserProfileRequest struct {
	DisplayName string `json:"displayName,omitempty"`
	ProfileURL  string `json:"profileURL,omitempty"`
	Description string `json:"description,omitempty"`
	AvatarURL   string `json:"avatarURL,omitempty"`
	Twitter     string `json:"twitter,omitempty"`
	GitHub      string `json:"github,omitempty"`
	Website     string `json:"website,omitempty"`
}

type UpdateGraphRequest struct {
	Name                string `json:"name,omitempty"`
	Unit                string `json:"unit,omitempty"`
	Color               string `json:"color,omitempty"`
	Timezone            string `json:"timezone,omitempty"`
	SelfSufficient      string `json:"selfSufficient,omitempty"`
	IsSecret            string `json:"isSecret,omitempty"`
	PublishOptionalData string `json:"publishOptionalData,omitempty"`
}

type Pixel struct {
	Date         string `json:"date"`
	Quantity     string `json:"quantity"`
	OptionalData string `json:"optionalData,omitempty"`
}

type PixelDetail struct {
	Date         string `json:"date"`
	Quantity     string `json:"quantity"`
	OptionalData string `json:"optionalData,omitempty"`
}

type PixelList struct {
	Dates   []string
	Details []PixelDetail
}

func (p *PixelList) UnmarshalJSON(data []byte) error {
	// まずstring配列として試す
	var dates []string
	if err := json.Unmarshal(data, &dates); err == nil {
		p.Dates = dates
		return nil
	}
	// 次にPixelDetail配列として試す
	var details []PixelDetail
	if err := json.Unmarshal(data, &details); err == nil {
		p.Details = details
		return nil
	}
	return fmt.Errorf("pixelsフィールドの型が不正です: %s", string(data))
}

type GetPixelsResponse struct {
	Pixels PixelList `json:"pixels"`
}

type GraphStats struct {
	TotalPixelsCount  int         `json:"totalPixelsCount"`
	MaxQuantity       json.Number `json:"maxQuantity"`
	MinQuantity       json.Number `json:"minQuantity"`
	MaxDate           string      `json:"maxDate"`
	MinDate           string      `json:"minDate"`
	TotalQuantity     json.Number `json:"totalQuantity"`
	AvgQuantity       json.Number `json:"avgQuantity"`
	TodaysQuantity    json.Number `json:"todaysQuantity"`
	YesterdayQuantity json.Number `json:"yesterdayQuantity"`
}

type BoolString bool

func (b *BoolString) UnmarshalJSON(data []byte) error {
	if string(data) == "true" || string(data) == "false" {
		var boolVal bool
		if err := json.Unmarshal(data, &boolVal); err != nil {
			return err
		}
		*b = BoolString(boolVal)
		return nil
	}
	var strVal string
	if err := json.Unmarshal(data, &strVal); err != nil {
		return err
	}
	*b = BoolString(strVal == "true")
	return nil
}

type GraphDefinition struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	Unit                string     `json:"unit"`
	Type                string     `json:"type"`
	Color               string     `json:"color"`
	Timezone            string     `json:"timezone,omitempty"`
	SelfSufficient      BoolString `json:"selfSufficient"`
	IsSecret            BoolString `json:"isSecret"`
	PublishOptionalData BoolString `json:"publishOptionalData"`
}

type GetGraphsResponse struct {
	Graphs []GraphDefinition `json:"graphs"`
}

type PixelaResponse struct {
	Message   string `json:"message"`
	IsSuccess bool   `json:"isSuccess"`
}

func NewClient() *Client {
	return &Client{
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) CreateUser(req CreateUserRequest) (*PixelaResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/v1/users", c.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) CreateGraph(username, token string, req CreateGraphRequest) (*PixelaResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/users/%s/graphs", c.BaseURL, username),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create graph: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) PostPixel(username, token, graphID string, req PostPixelRequest) (*PixelaResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/users/%s/graphs/%s", c.BaseURL, username, graphID),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to post pixel: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) DeleteUser(username, token string) (*PixelaResponse, error) {
	httpReq, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/v1/users/%s", c.BaseURL, username),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) UpdateUser(username, token string, req UpdateUserRequest) (*PixelaResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/v1/users/%s", c.BaseURL, username),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) UpdateUserProfile(username, token string, req UpdateUserProfileRequest) (*PixelaResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/@%s", c.BaseURL, username),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) UpdateGraph(username, token, graphID string, req UpdateGraphRequest) (*PixelaResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/v1/users/%s/graphs/%s", c.BaseURL, username, graphID),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update graph: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) DeleteGraph(username, token, graphID string) (*PixelaResponse, error) {
	httpReq, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/v1/users/%s/graphs/%s", c.BaseURL, username, graphID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to delete graph: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) GetPixels(username, token, graphID string, from, to, withBody *string) (*GetPixelsResponse, error) {
	baseURL := fmt.Sprintf("%s/v1/users/%s/graphs/%s/pixels", c.BaseURL, username, graphID)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	q := u.Query()
	if from != nil && *from != "" {
		q.Set("from", *from)
	}
	if to != nil && *to != "" {
		q.Set("to", *to)
	}
	if withBody != nil && *withBody != "" {
		q.Set("withBody", *withBody)
	}
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequest(
		"GET",
		u.String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get pixels: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get pixels: status %d, body: %s", resp.StatusCode, string(body))
	}

	var pixelsResp GetPixelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&pixelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &pixelsResp, nil
}

func (c *Client) GetGraphStats(username, token, graphID string) (*GraphStats, error) {
	httpReq, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/v1/users/%s/graphs/%s/stats", c.BaseURL, username, graphID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get graph stats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get graph stats: status %d, body: %s", resp.StatusCode, string(body))
	}

	var stats GraphStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &stats, nil
}

func (c *Client) GetGraphs(username, token string) (*GetGraphsResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/v1/users/%s/graphs", c.BaseURL, username),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get graphs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var graphsResp GetGraphsResponse
	if err := json.Unmarshal(body, &graphsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &graphsResp, nil
}

func (c *Client) GetGraphDefinition(username, token, graphID string) (*GraphDefinition, error) {
	httpReq, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/v1/users/%s/graphs/%s/graph-def", c.BaseURL, username, graphID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-USER-TOKEN", token)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get graph definition: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var graphDef GraphDefinition
	if err := json.Unmarshal(body, &graphDef); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &graphDef, nil
}

func (c *Client) GetGraph(username, graphID string) (string, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/v1/users/%s/graphs/%s", c.BaseURL, username, graphID))
	if err != nil {
		return "", fmt.Errorf("failed to get graph: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

func (c *Client) parseResponse(resp *http.Response) (*PixelaResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var pixelaResp PixelaResponse
	if err := json.Unmarshal(body, &pixelaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &pixelaResp, nil
}
