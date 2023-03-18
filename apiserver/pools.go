package apiserver

import (
	"github.com/thedevsaddam/traffic"
	"math/rand"
	"time"
)

// ControllersPoolRandom A simple implementation of ControllersPool.
type ControllersPoolRandom struct {
	controllers []Controller
	n           int
	generator   *rand.Rand
}

// Get implements the interface. Currently, we randomly choose one.
func (s *ControllersPoolRandom) Get() Controller {
	return s.controllers[s.generator.Intn(s.n)]
}

type PoolFactory func(controllers []Controller, weights []int) ControllersPool

// NewControllersPoolRandom Creates a new ControllersPoolRandom.
func NewControllersPoolRandom(controllers []Controller, _ []int) ControllersPool {

	return &ControllersPoolRandom{
		controllers: controllers,
		n:           len(controllers),
		generator:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type ControllersPoolLoadBalance struct {
	controllers []Controller
	n           int
	t           traffic.Traffic
}

// Get implements the interface. Currently, we just return the first controller.
func (s *ControllersPoolLoadBalance) Get() Controller {
	return s.controllers[s.t.Next().(int)]
}

// NewControllersPoolSmoothWeightedRoundRobin Creates a new ControllersPoolSmoothWeightedRoundRobin.
func NewControllersPoolSmoothWeightedRoundRobin(controllers []Controller, weights []int) ControllersPool {
	t := traffic.NewSmoothWeightedRoundRobin()
	for i, w := range weights {
		err := t.Add(i, w)
		if err != nil {
			panic(err)
		}
	}

	return &ControllersPoolLoadBalance{
		controllers: controllers,
		n:           len(controllers),
		t:           t,
	}
}

var PoolFactories = map[string]PoolFactory{
	"random":   NewControllersPoolRandom,
	"balanced": NewControllersPoolSmoothWeightedRoundRobin,
}
