package demo_test

import (
	"demo"
	"testing"
)

func TestHelloSuccess(t *testing.T) {
	actual := demo.SayHi()
	expected := "Hello Go 2024"

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
