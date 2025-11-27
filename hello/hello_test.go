package demo_test

import (
	"demo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloSuccess(t *testing.T) {
	actual := demo.SayHi()
	expected := "Hello Go 2025"

	assert.Equal(t, expected, actual)
}
