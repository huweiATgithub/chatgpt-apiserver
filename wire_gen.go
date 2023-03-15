// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/huweiATgithub/chatgpt-apiserver/apiserver"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"strings"
)

// Injectors from wire.go:

func Initialize() *Server {
	config := ProviderConfig()
	controllersPoolRandom := NewControllerPool(config)
	server := NewServer(controllersPoolRandom, config)
	return server
}

// wire.go:

func NewControllerPool(config *Config) *apiserver.ControllersPoolRandom {
	controllers := make([]apiserver.Controller, 0)
	for _, controllerConfig := range config.Controllers {
		controller, err := ControllerFactories[controllerConfig.Type](&controllerConfig)
		if err != nil {
			log.Printf("In constructing controller: Config %v with error %s", controllerConfig, err)
			continue
		}
		controllers = append(controllers, controller)
	}
	log.Printf("Create controllers pool with %v controllers\n", len(controllers))
	pool := apiserver.NewControllersPoolRandom(controllers)
	return pool
}

func ProviderConfig() *Config {
	appName := "chatgpt-apiserver"
	viper.SetConfigName(appName)
	viper.SetEnvPrefix(appName)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
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
		viper.SetConfigFile(viper.GetString("config_file"))
	}
	if viper.GetString("config_path") != "" {
		viper.AddConfigPath(viper.GetString("config_path"))
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", appName))
	viper.SetDefault("port", "8080")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s", err))
	}
	viper.AutomaticEnv()
	log.Printf("Loading config file %s", viper.ConfigFileUsed())

	// Unmarshal Config
	var config Config
	if viper.Unmarshal(&config) != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
	return &config
}

func NewServer(pool *apiserver.ControllersPoolRandom, config *Config) *Server {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.Default())

	r.GET("/status", statusOK)
	r.POST("/v1/chat/completions", apiserver.CompleteChat(pool))
	return &Server{
		r:      r,
		config: config,
	}
}