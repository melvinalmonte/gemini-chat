package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gemini-chat/config"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelFunc()

	client, err := setupGenAIClient(ctx, cfg.APIKey)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}
	defer client.Close()

	filename := fmt.Sprintf("chat_history_%s.md", time.Now().Format("2006_01_02_150405"))
	if err := conductChatSession(ctx, client, filename); err != nil {
		log.Printf("Failed to conduct chat session: %v", err)
	}
}

func setupGenAIClient(ctx context.Context, apiKey string) (*genai.Client, error) {
	genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return genaiClient, nil
}

func conductChatSession(ctx context.Context, client *genai.Client, filename string) error {
	dirPath := ".history"
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fullPath := filepath.Join(dirPath, filename)

	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
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
		return fmt.Errorf("error reading from input: %w", err)
	}

	return nil
}

func handleUserInput(ctx context.Context, userInput string, chatSession *genai.ChatSession, file *os.File) error {
	if userInput == "" {
		fmt.Println("Please enter a valid input.")

		return nil
	}

	if _, err := file.WriteString(fmt.Sprintf("\n\n*Me:* %s", userInput)); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	messageStream := chatSession.SendMessageStream(ctx, genai.Text(userInput))

	fmt.Println("Gemini: ")

	_, err := file.WriteString("\n\n**Gemini:**\n\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	for response, err := messageStream.Next(); !errors.Is(err, iterator.Done); response, err = messageStream.Next() {
		if err != nil {
			return fmt.Errorf("error receiving response: %w", err)
		}

		geminiResponse := response.Candidates[0].Content
		if _, err := file.WriteString(fmt.Sprintln(geminiResponse.Parts[0])); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		fmt.Println(geminiResponse.Parts[0])
	}

	return nil
}
