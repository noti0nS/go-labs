package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	systemPrompt string = "Answer every request about user's topic in a cute way. You're running inside a CLI program, so if you detected user has the intention to leave you just to reply '%s'. DO NOT SAY ANYTHING ELSE."
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
		}
		req, err := MakeRequest(userPrompt, apiConfig)
		if err != nil {
			panic(err)
		}
		response, err := func() (*AIResponse, error) {
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			var response AIResponse
			json.NewDecoder(resp.Body).Decode(&response)

			return &response, nil
		}()
		if err != nil {
			panic(err)
		}
		aiMessage := response.GetMessage()
		if aiMessage == nil {
			fmt.Println(response.Id)
			fmt.Println("Something went wrong! Try once more.")
		} else if userWannaLeave(*aiMessage) {
			fmt.Println("AI > OK! I'm leaving. Thank you for your time :-)")
			return
		} else {
			fmt.Println("AI > " + *aiMessage)
		}
	}
}

func MakeRequest(userPrompt string, apiConfig *APIConfig) (*http.Request, error) {
	payload := &AIRequest{
		Model: apiConfig.Model,
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
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, apiConfig.Url+"/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return req, err // return it early
	}
	req.Header.Add("Authorization", "Bearer "+apiConfig.Key)
	return req, err
}

func userWannaLeave(message string) bool {
	return strings.ToLower(strings.TrimSpace(message)) == exit
}
