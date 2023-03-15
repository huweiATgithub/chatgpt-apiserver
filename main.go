package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/huweiATgithub/chatgpt-apiserver/apiserver"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
)

type ControllerConfig struct {
	Type       string            `mapstructure:"type"`
	Config     map[string]string `mapstructure:"config"`
	ConfigFile string            `mapstructure:"config_file"`
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

type ControllerFactory func(*ControllerConfig) (apiserver.Controller, error)

var ControllerFactories = map[string]ControllerFactory{
	"openai": NewOpenAIController,
}

func main() {
	config := LoadConfig()

	controllers := make([]apiserver.Controller, 0)
	for _, controllerConfig := range config.Controllers {
		controller, err := ControllerFactories[controllerConfig.Type](&controllerConfig)
		if err != nil {
			log.Printf("In constructing controller: Config %v with error %s", controllerConfig, err)
			continue
		}
		controllers = append(controllers, controller)
	}

	// Create controllers pool
	log.Printf("Create controllers pool with %v controllers\n", len(controllers))
	pool := apiserver.NewControllersPoolRandom(controllers)
	StartServer(pool, config)
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

func LoadConfig() *Config {
	appName := "chatgpt-apiserver"
	viper.SetConfigName(appName)
	viper.SetEnvPrefix(appName)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// Command-line or Env provided config file / path
	pflag.String("config_file", "", "path to the server config file")
	pflag.String("config_path", "", "path to the server config file")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(fmt.Errorf("fatal error binding pflags: %s", err))
	}
	viper.MustBindEnv("config_file", "CONFIG_FILE")
	viper.MustBindEnv("config_path", "CONFIG_PATH")
	viper.AutomaticEnv()
	if viper.GetString("config_file") != "" {
		viper.SetConfigFile(viper.GetString("config_file")) // it has the highest priority if set
	}
	if viper.GetString("config_path") != "" {
		viper.AddConfigPath(viper.GetString("config_path")) // earlier path has higher priority
	}

	// Other Paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", appName))

	// Default Values
	viper.SetDefault("port", "8080")

	// Read Config
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s", err))
	}
	viper.AutomaticEnv() // Environment variables have higher priority
	log.Printf("Loading config file %s", viper.ConfigFileUsed())

	// Unmarshal Config
	var config Config
	if viper.Unmarshal(&config) != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
	return &config
}

func StartServer(pool *apiserver.ControllersPoolRandom, config *Config) {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.Default())

	r.GET("/status", statusOK)
	r.POST("/v1/chat/completions", apiserver.CompleteChat(pool))
	if err := r.Run(":" + config.Port); err != nil {
		log.Fatal(err)
	}
}
