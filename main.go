package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Delta struct {
	Content string `json:"content"`
}

type Choice struct {
	Delta Delta `json:"delta"`
}

type StreamResponse struct {
	Choices []Choice `json:"choices"`
}

func SendDeepSeekRequest(apiKey string, messages []Message) (string, error) {
	url := "https://api.deepseek.com/chat/completions"
	body := RequestBody{
		Model:    "deepseek-chat",
		Messages: messages,
		Stream:   true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	assistantResponse := ""
	fmt.Println("\nAssistant:")
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data: ") {
			jsonString := strings.TrimPrefix(line, "data: ")
			var streamResponse StreamResponse
			err := json.Unmarshal([]byte(jsonString), &streamResponse)
			if err != nil {
				continue
			}

			if len(streamResponse.Choices) > 0 && streamResponse.Choices[0].Delta.Content != "" {
				content := streamResponse.Choices[0].Delta.Content
				fmt.Print(content)
				assistantResponse += content
			}
		}
	}
	fmt.Println()
	return assistantResponse, nil
}

func ReadMultiLineInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("You:")
	var userInput string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading input:", err)
			return ""
		}
		userInput += line
	}
	return strings.TrimRight(userInput, "\n")
}

func StartDeepSeekConversation(apiKey string) {
	fmt.Println("DeepSeek Multi-Turn Conversation (type 'exit' to quit)")

	messages := []Message{}

	for {
		userInput := ReadMultiLineInput()

		if userInput == "exit" {
			fmt.Println("Exiting conversation.")
			break
		}

		if userInput == "" {
			fmt.Println("Error: Empty input provided. Please enter valid input.")
			continue
		}

		messages = append(messages, Message{
			Role:    "user",
			Content: userInput,
		})

		assistantResponse, err := SendDeepSeekRequest(apiKey, messages)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}

		if assistantResponse != "" {
			fmt.Println()
			messages = append(messages, Message{
				Role:    "assistant",
				Content: assistantResponse,
			})
		}
	}
}

func main() {
	apiKey := flag.String("key", "", "API key for DeepSeek")
	flag.Parse()

	if *apiKey == "" {
		fmt.Println("API key is required. Please provide it using the -key flag.")
		os.Exit(1)
	}

	StartDeepSeekConversation(*apiKey)
}
