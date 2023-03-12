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
	config OpenAIConfig
	client *openai.Client
}

type OpenAIConfig struct {
	ApiKey string `json:"api_key"`
	Proxy  string `json:"proxy"`
}

type OpenAIStream struct {
	stream *openai.ChatCompletionStream
}

// NewOpenAIController creates a new OpenAIController.
func NewOpenAIController(config OpenAIConfig) *OpenAIController {
	openaiConfig := openai.DefaultConfig(config.ApiKey)
	if config.Proxy != "" {
		proxyUrl, err := url.Parse(config.Proxy)
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
	controller := OpenAIController{config, client}
	return &controller
}

// NewOpenAIConfig creates a new OpenAIConfig.
func NewOpenAIConfig(configFilePath string) *OpenAIConfig {

	var config OpenAIConfig
	if configFilePath != "" {
		if err := config.readConfigFile(configFilePath); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := config.readConfigEnv(); err != nil {
			log.Fatal(err)
		}
	}
	return &config
}

// readConfigFile reads the config from file
func (o *OpenAIConfig) readConfigFile(configFilePath string) (err error) {
	jsonFile, err := os.Open(configFilePath)
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

// readConfigEnv reads the config from environment variables
func (o *OpenAIConfig) readConfigEnv() (err error) {
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

// concertStreamResponse converts the stream response from the stream response of the specific API to the unified stream response.
func convertStreamResponse(r *openai.ChatCompletionStreamResponse) (response *ChatCompletionStreamResponse, err error) {
	_response := *r
	response = &_response
	err = nil
	return
}

// Recv implement ChatCompletionStream interface for OpenAIStream
func (o *OpenAIStream) Recv() (response ChatCompletionStreamResponse, err error) {
	resp, err := o.stream.Recv()
	if err != nil {
		return
	}
	r, err := convertStreamResponse(&resp)
	response = *r
	return
}

// Close implement ChatCompletionStream interface for OpenAIStream
func (o *OpenAIStream) Close() {
	o.stream.Close()
}

// CompleteChat implement the interface
func (o *OpenAIController) CompleteChat(r *ChatCompletionRequest) (response *ChatCompletionResponse, err error) {
	request, err := convertRequest(r)
	if err != nil {
		log.Println(err)
		return
	}
	// Process the input data, generate a completion, and package it in the response struct
	resp, err := o.client.CreateChatCompletion(
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

// CompleteChatStream implement the interface
func (o *OpenAIController) CompleteChatStream(r *ChatCompletionRequest) (stream ChatCompletionStream, err error) {
	request, err := convertRequest(r)
	if err != nil {
		log.Println(err)
		return
	}
	// Process the input data, generate a completion, and package it in the response struct
	s, err := o.client.CreateChatCompletionStream(
		context.Background(),
		*request,
	)
	stream = &OpenAIStream{stream: s}
	return
}
