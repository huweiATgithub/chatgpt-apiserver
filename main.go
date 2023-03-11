package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/huweiATgithub/chatgpt-apiserver/apiserver"
	"log"
	"net/http"
)

func main() {

	port := flag.String("port", "8080", "port to listen on")
	openaiConfigFile := flag.String("openai_config_file", "", "path to the openai config file")
	flag.Parse()

	// Create openai controller
	var config apiserver.OpenAIConfig
	if *openaiConfigFile != "" {
		if err := config.ReadConfigFile(*openaiConfigFile); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := config.ReadConfigEnv(); err != nil {
			log.Fatal(err)
		}
	}
	controller := apiserver.OpenAIController{
		Config: config,
	}

	// Create controllers pool
	pool := apiserver.SimpleControllersPool{
		Controllers: []apiserver.Controller{&controller},
	}

	r := gin.Default()
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	r.POST("/v1/chat/completions", apiserver.CompleteChat(&pool))
	if err := r.Run(":" + *port); err != nil {
		log.Fatal(err)
	}
}
