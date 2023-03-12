package apiserver

import (
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// ControllersPool The interface for the pool of controllers.
type ControllersPool interface {
	Get() Controller
}

// ControllersPoolRandom A simple implementation of ControllersPool.
type ControllersPoolRandom struct {
	Controllers []Controller
	n           int
	generator   *rand.Rand
}

// Get implements the interface. Currently, we just return the first controller.
func (s *ControllersPoolRandom) Get() Controller {
	return s.Controllers[s.generator.Intn(s.n)]
}

// NewControllersPoolRandom Creates a new ControllersPoolRandom.
func NewControllersPoolRandom(controllers []Controller) *ControllersPoolRandom {

	return &ControllersPoolRandom{
		Controllers: controllers,
		n:           len(controllers),
		generator:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
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
