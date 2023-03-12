package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/huweiATgithub/chatgpt-apiserver/apiserver"
	"log"
	"net/http"
)

func statusOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
func main() {
	port := flag.String("port", "8080", "port to listen on")
	openaiConfigFile := flag.String("openai_config_file", "", "path to the openai config file")
	flag.Parse()

	// Create controllers pool
	pool := apiserver.NewControllersPoolRandom([]apiserver.Controller{
		apiserver.NewOpenAIController(*apiserver.NewOpenAIConfig(*openaiConfigFile)),
		apiserver.NewOpenAIController(*apiserver.NewOpenAIConfig(*openaiConfigFile)),
	})

	r := gin.Default()
	r.GET("/status", statusOK)
	r.POST("/v1/chat/completions", apiserver.CompleteChat(pool))
	if err := r.Run(":" + *port); err != nil {
		log.Fatal(err)
	}
}
