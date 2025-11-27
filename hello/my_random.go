package demo

import (
	"fmt"
	"math/rand"
)

type XXX interface {
	GetNumber() int
}

type MyRandom struct{}

func (mr MyRandom) GetNumber() int {
	return rand.Intn(10) + 1
}

func GenerateData(x XXX) string {
	randomNumber := x.GetNumber()
	return fmt.Sprintf("My Data %d", randomNumber)
}
