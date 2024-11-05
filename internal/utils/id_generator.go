package utils

import (
	"fmt"
	"sync"
	"time"
)

type IDGenerator struct {
	mu        sync.Mutex
	runNumber int
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{
		runNumber: 0,
	}
}

func (gen *IDGenerator) GenerateID() string {
	gen.mu.Lock()
	defer gen.mu.Unlock()

	timestamp := time.Now().Unix()
	gen.runNumber++
	id := fmt.Sprintf("[]-[%d+%d]", timestamp, gen.runNumber)
	return id
}

func main() {
	gen := NewIDGenerator()

	for i := 0; i < 10; i++ {
		fmt.Println(gen.GenerateID())
	}
}
