package clients

import (
	"YourChatGptBackendService/database"
	models "YourChatGptBackendService/models/api"
	"YourChatGptBackendService/models/data"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

func OpenAIChatCompletionStream(request models.SendMessageRequest, c *gin.Context) (data.MessageThreadData, error) {
	fmt.Println("Starting openai chat completion...")

	// Set headers for SSE response
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Create the OpenAI client and initialize the chat stream
	client := openai.NewClient(" *** Add your OpenAI API key here *** ")
	ctx := context.Background()

	msgArray, err := RetrieveMessageArray(request)
	if err != nil {
		return data.MessageThreadData{}, err
	}

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 2000,
		Messages: msgArray,
		Stream: true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Printf("Error happened when try to get stream: %v\n", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return data.MessageThreadData{}, err
	}

	fmt.Println("Got the response stream ==")

	defer stream.Close()

	var fullContent string
	var msgData data.MessageThreadData

	// Send stream responses as SSE events
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			break
		}

		if err != nil {
			log.Printf("Stream error: %v\n", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return data.MessageThreadData{}, err
		}

		if response.Choices[0].Delta.Role != "" {
			msgData.Role = response.Choices[0].Delta.Role
		}
		msgData.FinishReason = string(response.Choices[0].FinishReason)

		// Append the SSE event data to the full content
		fullContent += response.Choices[0].Delta.Content

		// Format the SSE event data
		eventData := fmt.Sprintf(response.Choices[0].Delta.Content)

		// Write the SSE event to the client
		c.SSEvent("message", eventData)
		c.Writer.Flush()
	}

	msgData.Content = fullContent
	
	return msgData, nil
}

func RetrieveMessageArray(request models.SendMessageRequest) ([]openai.ChatCompletionMessage, error) {
	// Get the chat data in database
	container, err := database.ConnectDatabase("", "", "")
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %v", err)
	}

	chatData, err := database.GetChatData(request.UserID, request.ChatID, container)
	if err != nil {
		return nil, fmt.Errorf("Failed to get chat data: %v", err)
	}

	var msgArray []openai.ChatCompletionMessage
	msgArray = append(msgArray, openai.ChatCompletionMessage{
		Content: request.Message.Content,
		Role: request.Message.Role,
	})

	if chatData.ChatHistory.LastUpdatedMessage == "" {
		return msgArray, nil
	}

	msg := chatData.ChatHistory.MessagesMapping[chatData.ChatHistory.LastUpdatedMessage]
	msgArray = append([]openai.ChatCompletionMessage{{
		Content: msg.Content,
		Role: msg.Role,
	}}, msgArray...)

	// TODO: optimize the counter logic
	for i := 0; i < 19; i++ {
		if (msg.ParentID == "") {
			break;
		}

		msg = chatData.ChatHistory.MessagesMapping[msg.ParentID]
		msgArray = append([]openai.ChatCompletionMessage{{
			Content: msg.Content,
			Role: msg.Role,
		}}, msgArray...)
	}

	return msgArray, nil
}
