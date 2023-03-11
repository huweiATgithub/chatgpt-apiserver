package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huweiATgithub/chatgpt-apiserver/apiserver"
	"log"
	"net/http"
)

func main() {

	r := gin.Default()
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	var config apiserver.OpenAIConfig
	config.ReadConfig("config/openai.json")
	controller := apiserver.OpenAIController{
		Config: config,
	}
	pool := apiserver.SimpleControllersPool{
		Controllers: []apiserver.Controller{&controller},
	}

	r.POST("/v1/chat/completions", apiserver.CompleteChat(&pool))
	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
