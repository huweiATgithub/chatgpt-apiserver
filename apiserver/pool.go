package apiserver

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// ControllersPool The interface for the pool of controllers.
type ControllersPool interface {
	Get() Controller
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
		if request.Stream {
			stream, err := controller.CompleteChatStream(request)
			defer stream.Close()
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.Stream(func(w io.Writer) bool {
				resp, err := stream.Recv()
				if err == io.EOF {
					c.SSEvent("", getStreamFinishData())
					return false
				}
				if err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return false
				}
				c.SSEvent("", resp)
				return true
			})
		} else {
			response, err := controller.CompleteChat(request)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.JSON(http.StatusOK, response)
		}
	}
}
