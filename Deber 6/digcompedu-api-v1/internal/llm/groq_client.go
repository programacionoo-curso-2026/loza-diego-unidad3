// internal/llm/groq_client.go

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

type GroqClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewGroqClient(apiKey, model string) *GroqClient {
	return &GroqClient{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type groqRequest struct {
	Model    string        `json:"model"`
	Messages []groqMessage `json:"messages"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *GroqClient) Analizar(ctx context.Context, prompt string) (string, error) {
	reqBody := groqRequest{
		Model: c.model,
		Messages: []groqMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error al serializar request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al llamar a la API de Groq: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer respuesta: %w", err)
	}

	var parsed groqResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", fmt.Errorf("error al parsear respuesta: %w", err)
	}

	if parsed.Error != nil {
		return "", fmt.Errorf("error de la API de Groq: %s", parsed.Error.Message)
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("respuesta vacía de la API de Groq")
	}

	return parsed.Choices[0].Message.Content, nil
}