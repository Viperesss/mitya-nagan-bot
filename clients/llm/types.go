package llm

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AiRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type AiResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}
