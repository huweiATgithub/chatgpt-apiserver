package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huweiATgithub/chatgpt-apiserver/apiserver"
	"log"
	"net/http"
)

func main() {
	server := Initialize()
	server.Start()
}

type ControllerConfig struct {
	Type       string            `mapstructure:"type"`
	Config     map[string]string `mapstructure:"config"`
	ConfigFile string            `mapstructure:"config_file"`
}

type Config struct {
	Port        string             `mapstructure:"port"`
	Controllers []ControllerConfig `mapstructure:"controllers"`
}

type Server struct {
	r      *gin.Engine
	config *Config
}

func (s *Server) Start() {
	if err := s.r.Run(":" + s.config.Port); err != nil {
		log.Fatal(err)
	}
}

func statusOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

type ControllerFactory func(*ControllerConfig) (apiserver.Controller, error)

var ControllerFactories = map[string]ControllerFactory{
	"openai": NewOpenAIController,
}

func NewOpenAIController(controllerConfig *ControllerConfig) (controller apiserver.Controller, err error) {
	var openaiConfig *apiserver.OpenAIConfig
	if controllerConfig.ConfigFile != "" {
		openaiConfig, err = apiserver.NewOpenAIConfigFromFile(controllerConfig.ConfigFile)
		if err != nil {
			return
		}
		controller = apiserver.NewOpenAIController(*openaiConfig)
	} else {
		openaiConfig, err = apiserver.NewOpenAIConfigFromMap(controllerConfig.Config)
		if err != nil {
			return
		}
	}
	controller = apiserver.NewOpenAIController(*openaiConfig)
	return
}
