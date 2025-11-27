package demo_test

import (
	"demo"
	"testing"
)

func TestRamdomEqual5(t *testing.T) {
	data := demo.GenerateData()
	expected := "My Data 5"
	if data != expected {
		t.Errorf("Expected %s but got %s", expected, data)
	}
}
