package main

import (
	"demo02"
	"fmt"
)

func main() {
	result := demo02.SayHi()
	println(result)
	fmt.Println(result)

	for i := 0; i < 5; i++ {
		fmt.Printf("Count: %d\n", i)
	}
}
