package main

type AIRequest struct {
	Model    string      `json:"model"`
	Stream   bool        `json:"stream"`
	Messages []AIMessage `json:"messages"`
}

type AIResponse struct {
	Id      string     `json:"id"`
	Choices []AIChoice `json:"choices"`
}

type AIChoice struct {
	FinishReason string    `json:"finish_reason"`
	Delta        AIMessage `json:"delta"`
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (resp *AIResponse) GetMessage() *string {
	if len(resp.Choices) == 0 {
		return nil
	}
	return &resp.Choices[0].Delta.Content
}
