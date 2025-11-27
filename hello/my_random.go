package demo

import (
	"fmt"
	"math/rand"
)

func GenerateData() string {
	randomNumber := rand.Intn(10) + 1
	return fmt.Sprintf("My Data %d", randomNumber)
}
