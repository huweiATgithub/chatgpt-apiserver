package apiserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ControllersPool The interface for the pool of controllers.
type ControllersPool interface {
	Get() Controller
}

// SimpleControllersPool A simple implementation of ControllersPool.
type SimpleControllersPool struct {
	Controllers []Controller
}

// Get implements the interface. Currently, we just return the first controller.
func (s *SimpleControllersPool) Get() Controller {
	return s.Controllers[0]
}

// readRequest reads the request from the context and returns the request as the unified ChatCompletionRequest type.
func readRequest(c *gin.Context) (request *ChatCompletionRequest, err error) {
	err = c.BindJSON(&request)
	return
}

// CompleteChat handles the chat completion endpoint.
func CompleteChat(pool ControllersPool) gin.HandlerFunc {
	return func(c *gin.Context) {
		request, err := readRequest(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		controller := pool.Get()
		response, err := controller.CompleteChat(request)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}
