// internal/llm/client.go

package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LLMProvider interface {
	Analizar(ctx context.Context, prompt string) (string, error)
}

type AnthropicClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewAnthropicClient(apiKey, model string) *AnthropicClient {
	return &AnthropicClient{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *AnthropicClient) Analizar(ctx context.Context, prompt string) (string, error) {
	reqBody := anthropicRequest{
		Model:     c.model,
		MaxTokens: 1024,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error al serializar request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al llamar a la API de IA: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer respuesta: %w", err)
	}

	var parsed anthropicResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", fmt.Errorf("error al parsear respuesta: %w", err)
	}

	if parsed.Error != nil {
		return "", fmt.Errorf("error de la API de IA: %s", parsed.Error.Message)
	}

	if len(parsed.Content) == 0 {
		return "", fmt.Errorf("respuesta vacía de la API de IA")
	}

	return parsed.Content[0].Text, nil
}