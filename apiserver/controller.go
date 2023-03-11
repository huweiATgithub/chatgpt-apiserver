// Package apiserver handles chat completion requests.
package apiserver

import (
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
