package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultBaseUrl    = "http://localhost:11434"
	DefaultModel      = "qwen3:1.7b"
	DefaultNumPredict = 512
)

type Client struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewClient(baseURL string, model string) *Client {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = DefaultBaseUrl
	}

	if strings.TrimSpace(model) == "" {
		model = DefaultModel
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		httpClient: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c Client) Generate(ctx context.Context, prompt string) (string, error) {
	think := false

	body, err := json.Marshal(generateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
		Think:  &think,
		Options: generateOptions{
			NumPredict: DefaultNumPredict,
		},
	})
	if err != nil {
		return "", fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/generate",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("ollama returned %s: %s", resp.Status, strings.TrimSpace(string(msg)))
	}

	var decoded generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if strings.TrimSpace(decoded.Response) == "" {
		return "", fmt.Errorf("ollama returned an empty response")
	}

	return decoded.Response, nil
}

type generateRequest struct {
	Model   string          `json:"model"`
	Prompt  string          `json:"prompt"`
	Stream  bool            `json:"stream"`
	Think   *bool           `json:"think,omitempty"`
	Options generateOptions `json:"options,omitempty"`
}

type generateOptions struct {
	NumPredict int `json:"num_predict,omitempty"`
}

type generateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}
