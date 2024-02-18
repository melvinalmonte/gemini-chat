package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	apiKey, err := getEnvVar("GEMINI_API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelFunc()

	client, err := setupGenAIClient(ctx, apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	filename := fmt.Sprintf("chat_history_%s.md", time.Now().Format("2006_01_02_150405"))
	if err := conductChatSession(ctx, client, filename); err != nil {
		log.Fatal(err)
	}
}
func getEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is not set", key)
	}
	return value, nil
}
func setupGenAIClient(ctx context.Context, apiKey string) (*genai.Client, error) {
	return genai.NewClient(ctx, option.WithAPIKey(apiKey))
}
func conductChatSession(ctx context.Context, client *genai.Client, filename string) error {
	dirPath := ".history"
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	fullPath := filepath.Join(dirPath, filename)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	generativeModel := client.GenerativeModel("gemini-pro")
	chatSession := generativeModel.StartChat()

	fmt.Println("---------------------------------------------------------------------------------------")
	fmt.Println("Terminal chat started. Type something and press enter (CTRL+C to exit).")
	fmt.Println("---------------------------------------------------------------------------------------")
	fmt.Print("Me: ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := handleUserInput(ctx, scanner.Text(), chatSession, file); err != nil {
			return err
		}
		fmt.Print("Me: ")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from input: %v", err)
	}

	return nil
}
func handleUserInput(ctx context.Context, userInput string, chatSession *genai.ChatSession, file *os.File) error {
	if userInput == "" {
		fmt.Println("Please enter a valid input.")
		return nil
	}

	if _, err := file.WriteString(fmt.Sprintf("\n\n*Me:* %s", userInput)); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	messageStream := chatSession.SendMessageStream(ctx, genai.Text(userInput))
	fmt.Println("Gemini: ")
	_, err := file.WriteString("\n\n**Gemini:**\n\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	for response, err := messageStream.Next(); !errors.Is(err, iterator.Done); response, err = messageStream.Next() {
		if err != nil {
			return fmt.Errorf("error receiving response: %v", err)
		}

		geminiResponse := response.Candidates[0].Content
		if _, err := file.WriteString(fmt.Sprintln(geminiResponse.Parts[0])); err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
		fmt.Println(geminiResponse.Parts[0])
	}

	return nil
}
