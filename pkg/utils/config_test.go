package utils

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	c := Config{}

	config := c.Get()

	fmt.Printf("%v", config)
}
