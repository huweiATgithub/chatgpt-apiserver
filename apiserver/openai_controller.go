package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

// OpenAIController The controller for OpenAI.
type OpenAIController struct {
	Config OpenAIConfig
}

type OpenAIConfig struct {
	ApiKey string `json:"api_key"`
	Proxy  string `json:"proxy"`
}

// ReadConfigFile reads the config from file
func (o *OpenAIConfig) ReadConfigFile(configPath string) (err error) {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, o)
	if err != nil {
		return
	}
	return
}

// ReadConfigEnv reads the config from environment variables
func (o *OpenAIConfig) ReadConfigEnv() (err error) {
	err = nil
	o.ApiKey = os.Getenv("OPENAI_API_KEY")
	// ApiKey is required
	if o.ApiKey == "" {
		err = errors.New("OPENAI_API_KEY is required")
		return
	}
	o.Proxy = os.Getenv("http_proxy")
	return
}

// convertRequest converts the request from the unified request to the request of the specific API.
func convertRequest(r *ChatCompletionRequest) (request *openai.ChatCompletionRequest, err error) {
	_request := *r
	request = &_request
	err = nil
	return
}

// convertResponse converts the response from the response of the specific API to the unified response.
func convertResponse(r *openai.ChatCompletionResponse) (response *ChatCompletionResponse, err error) {
	_response := *r
	response = &_response
	err = nil
	return
}

// CompleteChat implement the interface
func (o *OpenAIController) CompleteChat(r *ChatCompletionRequest) (response *ChatCompletionResponse, err error) {
	request, err := convertRequest(r)
	if err != nil {
		log.Println(err)
		return
	}
	// Process the input data, generate a completion, and package it in the response struct
	openaiConfig := openai.DefaultConfig(o.Config.ApiKey)
	if o.Config.Proxy != "" {
		proxyUrl, err := url.Parse(o.Config.Proxy)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		openaiConfig.HTTPClient = &http.Client{
			Transport: transport,
		}
	}
	openaiConfig.HTTPClient = &http.Client{}
	client := openai.NewClientWithConfig(openaiConfig)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		*request,
	)
	if err != nil {
		log.Println(err)
		return
	}
	response, err = convertResponse(&resp)
	return
}
