package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	systemPrompt string = "Answer every request about user's topic in a cute way. You're running inside a CLI program, so if you detected user has the intention to leave you just need to reply '%s'. DO NOT SAY ANYTHING ELSE."
	exit         string = "exit"
)

func main() {
	apiConfig, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Type something here > ")
		if !scanner.Scan() {
			panic(scanner.Err())
		}
		userPrompt := strings.TrimSpace(scanner.Text())
		if userPrompt == exit {
			return
		} else if userPrompt == "" {
			continue
		}
		req, err := makeRequest(userPrompt, apiConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			continue
		}
		fmt.Println()
		leave, err := handleResponse(client, req)
		fmt.Println()
		fmt.Println()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			continue
		} else if *leave {
			return
		}
	}
}

func makeRequest(userPrompt string, apiConfig *APIConfig) (*http.Request, error) {
	payload := &AIRequest{
		Model:  apiConfig.Model,
		Stream: true,
		Messages: []AIMessage{
			{
				Role:    "system",
				Content: fmt.Sprintf(systemPrompt, exit),
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, apiConfig.Url+"/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiConfig.Key)
	return req, nil
}

// handles AI response and returns whether should finish the app
func handleResponse(client *http.Client, req *http.Request) (*bool, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to process request: %w", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	leave := false

	for {
		line, err := reader.ReadString('\n')
		// fmt.Println("[DEBUG]: " + line)
		if err != nil {
			if err == io.EOF {
				break // sse ends
			}
			return nil, fmt.Errorf("failed to read sse event: %w", err)
		}
		if after, ok := strings.CutPrefix(line, "data: "); ok {
			data := strings.TrimSpace(after)
			if data == "[DONE]" {
				break
			}
			leave = processMessage(data)
		}
	}

	return &leave, nil
}

// returns true if user wanna leave from application
func processMessage(data string) bool {
	var chunkedResponse AIResponse
	if err := json.Unmarshal([]byte(data), &chunkedResponse); err != nil {
		panic(err)
	}
	aiMessage := chunkedResponse.GetMessage()
	if aiMessage == nil {
		fmt.Println(chunkedResponse.Id)
		fmt.Print("Something went wrong!")
	} else if userWannaLeave(*aiMessage) {
		return true
	} else {
		fmt.Print(*aiMessage)
	}
	return false
}

func userWannaLeave(message string) bool {
	return strings.ToLower(strings.TrimSpace(message)) == exit
}
