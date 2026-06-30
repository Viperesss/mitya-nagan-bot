package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"the-mitya-nagan-bot/lib/e"
	"time"
)

const (
	apiUrl  = "https://openrouter.ai/api/v1/chat/completions" // POST endpoint
	aiModel = "qwen/qwen3-8b"
)

type Client struct {
	apiKey  string
	aiModel string
	client  http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		aiModel: aiModel,
		client: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Chat(userPrompt string) (string, error) {
	request := AiRequest{
		Model: c.aiModel,
		Messages: []Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return "", e.Wrap("chat Marshal fail:", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return "", e.Wrap("chat NewRequest fail:", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(req)
	if err != nil {
		return "", e.Wrap("cannot send request:", err)
	}
	defer response.Body.Close()

	var aiResponse AiResponse

	if err := json.NewDecoder(response.Body).Decode(&aiResponse); err != nil {
		return "", e.Wrap("cannot Decode response:", err)
	}

	if len(aiResponse.Choices) == 0 {
		return "", errors.New("empty response from AI")
	}

	log.Println(aiResponse.Choices[0].Message.Content)

	return aiResponse.Choices[0].Message.Content, nil
}
