package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/huweiATgithub/chatgpt-apiserver/apiserver"
	"io"
	"log"
	"net/http"
	"os"
)

type ControllerConfig struct {
	Type       string            `json:"type"`
	Config     map[string]string `json:"config,omitempty"`
	ConfigFile string            `json:"config_file,omitempty"`
}

type Config struct {
	Port        string             `json:"port,omitempty"`
	Controllers []ControllerConfig `json:"controllers,omitempty"`
}

func statusOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func NewDefaultConfig() *Config {
	return &Config{
		Port: "8080",
	}
}

func allowOriginAll(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Next()
}

func main() {
	port := flag.String("port", "", "port to listen on")
	openaiConfigFile := flag.String("openai_config_file", "", "path to the openai config file")
	configFile := flag.String("config_file", "", "path to the server config file")
	flag.Parse()

	config := NewDefaultConfig()
	controllers := make([]apiserver.Controller, 0)
	// Parse config file
	if *configFile != "" {
		jsonFile, err := os.Open(*configFile)
		defer jsonFile.Close()
		if err != nil {
			log.Fatal(err)
		}
		byteValue, _ := io.ReadAll(jsonFile)
		err = json.Unmarshal(byteValue, &config)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, controllerConfig := range config.Controllers {
		log.Printf("Add controller by config: %v\n", controllerConfig)
		switch controllerConfig.Type {
		case "openai":
			if controllerConfig.ConfigFile != "" {
				openaiConfig, err := apiserver.NewOpenAIConfigFromFile(controllerConfig.ConfigFile)
				if err != nil {
					log.Fatal(err)
				}
				controllers = append(controllers, apiserver.NewOpenAIController(*openaiConfig))
			} else {
				openaiConfig, err := apiserver.NewOpenAIConfigFromMap(controllerConfig.Config)
				if err != nil {
					log.Fatal(err)
				}
				controllers = append(controllers, apiserver.NewOpenAIController(*openaiConfig))
			}
		}
	}

	// handle command line arguments
	if *port != "" {
		config.Port = *port
	}
	if *openaiConfigFile != "" {
		openaiConfig, err := apiserver.NewOpenAIConfigFromFile(*openaiConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		controllers = append(controllers, apiserver.NewOpenAIController(*openaiConfig))
		log.Printf("Add OpenAI controller from file %v\n", *openaiConfigFile)
	}

	// handle environment variables
	openaiConfig, err := apiserver.NewOpenAIConfigFromEnv()
	if err == nil {
		controllers = append(controllers, apiserver.NewOpenAIController(*openaiConfig))
		log.Printf("Add OpenAI controller from environment variables %v\n", openaiConfig)
	}

	// Create controllers pool
	log.Printf("Create controllers pool with %v controllers\n", len(controllers))
	pool := apiserver.NewControllersPoolRandom(controllers)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(allowOriginAll)
	r.GET("/status", statusOK)
	r.POST("/v1/chat/completions", apiserver.CompleteChat(pool))
	if err := r.Run(":" + config.Port); err != nil {
		log.Fatal(err)
	}
}
