package main

import (
	"YourChatGptBackendService/controllers"
	"YourChatGptBackendService/middleware"
	"log"

	docs "YourChatGptBackendService/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	// Initialize the Gin router
	r := gin.Default()

	// Load the middlewares
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.SecurityMiddleware())
	} else {
		r.Use(middleware.CorsMiddleware())
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/chats", controllers.GetChats) // domain:port/chats?eeid={id}
		v1.POST("/chat", controllers.CreateChat)
		v1.GET("/chats/:user_id/:chat_id", controllers.GetChatHistory)
		v1.POST("/completion", controllers.ChatCompletion)
		v1.DELETE("/:chat_id", controllers.DeleteChat) // domain:port/{chat_id}?user_id={user_id}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
