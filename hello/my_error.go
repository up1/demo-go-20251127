package demo

import "fmt"

func HiError() string {
	r, _ := doSth2()
	return r
}

func doSth2() (string, error) {
	return "", fmt.Errorf("My Error")
}
