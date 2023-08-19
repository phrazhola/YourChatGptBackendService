package controllers

import (
	"YourChatGptBackendService/clients"
	"YourChatGptBackendService/database"
	"YourChatGptBackendService/models/api"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetChats retrieves all conversations for a given user eeid.
func GetChats(c *gin.Context) {
	email := c.Query("eeid")

	response, err := database.GetUserChats(email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// CreateChat creates a new conversation for a given user eeid.
func CreateChat(c *gin.Context) {
	var request models.CreateChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, chatID, err := database.CreateChat(request.Email, request.Title)
	var response models.CreateChatResponse

	if err != nil {
		response = models.CreateChatResponse{
			UserID: userID,
			ChatID: chatID,
			Status:  "Failed",
			Message: fmt.Sprintf("Error: %v", err),
		}

		c.JSON(http.StatusInternalServerError, response)
	} else {
		response = models.CreateChatResponse{
			UserID: userID,
			ChatID: chatID,
			Status:  "Succeeded",
			Message: "null",
		}
	
		c.JSON(http.StatusOK, response)
	}
}

// GetChatHistory retrieves the chat history for given chat id and user id.
func GetChatHistory(c *gin.Context) {
	userID := c.Param("user_id")
	chatID := c.Param("chat_id")

	response, err := database.GetChatHistory(userID, chatID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// ChatCompletion returns the completion response (stream) from openai chat api
func ChatCompletion(c *gin.Context) {
	var request models.SendMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return SSE events to the client
	completionMsg, err := clients.OpenAIChatCompletionStream(request, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		database.UpdateChatHistory(completionMsg, request)
	}
}

func DeleteChat(c *gin.Context) {
	// Get chat_id and user_id from the request URL parameters
	chatID := c.Param("chat_id")
	userID := c.Query("user_id")

	err := database.DeleteChat(userID, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Chat deleted successfully",
			"chat_id": chatID,
			"user_id": userID,
		})
	}
}