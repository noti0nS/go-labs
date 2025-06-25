package main

type AIRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
}

type AIResponse struct {
	Id      string     `json:"id"`
	Choices []AIChoice `json:"choices"`
}

type AIChoice struct {
	FinishReason string    `json:"finish_reason"`
	Message      AIMessage `json:"message"`
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (resp *AIResponse) GetMessage() *string {
	if len(resp.Choices) == 0 {
		return nil
	}
	return &resp.Choices[0].Message.Content
}
