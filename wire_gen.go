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
	poolConfig := ProviderPoolConfig(config)
	v := ProviderControllerConfigs(config)
	v2 := NewControllers(v)
	v3 := NewWeights(v)
	controllersPool := NewPool(poolConfig, v2, v3)
	server := NewServer(controllersPool, config)
	return server
}

// wire.go:

func NewControllers(configs []ControllerConfig) []apiserver.Controller {
	controllers := make([]apiserver.Controller, len(configs))
	for i, config := range configs {
		controller, err := ControllerFactories[config.Type](&config)
		if err != nil {
			log.Printf("In constructing controller: Config %v with error %s", config, err)
			continue
		}
		controllers[i] = controller
	}
	return controllers
}

func NewWeights(configs []ControllerConfig) []int {
	weights := make([]int, len(configs))
	for i, config := range configs {
		weights[i] = config.Weight
		if weights[i] == 0 {
			weights[i] = defaultWeight
		}
	}
	return weights
}

func NewPool(config *PoolConfig, controllers []apiserver.Controller, weights []int) apiserver.ControllersPool {
	validControllers := make([]apiserver.Controller, 0, len(controllers))
	validWeights := make([]int, 0, len(weights))
	for i := 0; i < len(controllers); i++ {
		if controllers[i] != nil {
			validControllers = append(validControllers, controllers[i])
			validWeights = append(validWeights, weights[i])
		}
	}
	log.Printf("%v controllers are available for the pool %s\n", len(validControllers), config.Type)
	return apiserver.PoolFactories[config.Type](validControllers, validWeights)
}

func defaultConfig() {
	viper.SetDefault("port", "8080")
	viper.SetDefault("pool.type", "balanced")
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

	defaultConfig()

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

func ProviderPoolConfig(config *Config) *PoolConfig {
	return &config.Pool
}

func ProviderControllerConfigs(config *Config) []ControllerConfig {
	return config.Controllers
}

func NewServer(pool apiserver.ControllersPool, config *Config) *Server {
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
