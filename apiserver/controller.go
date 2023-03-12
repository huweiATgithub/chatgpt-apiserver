// Package apiserver handles chat completion requests.
package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// ChatCompletionRequest The unified request. Currently, we alias it to the one from openai package.
type ChatCompletionRequest = openai.ChatCompletionRequest

// ChatCompletionResponse The unified response. Currently, we alias it to
type ChatCompletionResponse = openai.ChatCompletionResponse

// Controller The interface for the controller.
type Controller interface {
	CompleteChat(r *ChatCompletionRequest) (response *ChatCompletionResponse, err error)
}

// readRequest reads the request from the gin.Context and returns the request as the unified ChatCompletionRequest type.
func readRequest(c *gin.Context) (request *ChatCompletionRequest, err error) {
	err = c.BindJSON(&request)
	return
}
