package demo

import (
	"fmt"
	"math/rand"
)

type MyRandom struct{}

func (mr MyRandom) GetNumber() int {
	return rand.Intn(10) + 1
}

func GenerateData() string {
	randomNumber := MyRandom{}.GetNumber()
	return fmt.Sprintf("My Data %d", randomNumber)
}
