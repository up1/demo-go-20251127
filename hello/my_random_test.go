package demo_test

import (
	"demo"
	"testing"
)

type MyStubRandom struct{}

func (mr MyStubRandom) GetNumber() int {
	return 5
}

func TestRamdomEqual5(t *testing.T) {
	data := demo.GenerateData(MyStubRandom{})
	expected := "My Data 5"
	if data != expected {
		t.Errorf("Expected %s but got %s", expected, data)
	}
}
